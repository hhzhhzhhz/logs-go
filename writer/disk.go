package writer

import (
	"bufio"
	"fmt"
	"github.com/hhzhhzhhz/logs-go/strftime"
	"go.uber.org/multierr"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxSize = 100
	defaultBufSize = 4
	randnum        = 10
	templog        = ".tmp"
	compressSuffix = ".gz"
)

var (
	// log start time
	currentTime = time.Now

	// os_Stat exists so it can be mocked out by tests.
	os_Stat = os.Stat
	// to file_log mb.
	megabyte = 1024 * 1024
	// avoid duplicate files
	dittotime = int64(5 * time.Minute)

	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type Option func(*Options)

// WithGenerationRule
func WithGenerationRule(gtr string) Option {
	return func(o *Options) {
		o.gtr = gtr
	}
}

// WithNorequriedTimezone
func WithNorequriedTimezone() Option {
	return func(o *Options) {
		o.requriedTimezone = false
	}
}

// WithBufSize
func WithBufSize(mb int) Option {
	return func(o *Options) {
		o.bufSize = mb
	}
}

// WithMaxSize
func WithMaxSize(mb int) Option {
	return func(o *Options) {
		o.maxSize = mb
	}
}

// WithMaxAge
func WithMaxAge(day int) Option {
	return func(o *Options) {
		o.maxAge = day
	}
}

// WithRotationTime
func WithRotationTime(t time.Duration) Option {
	return func(o *Options) {
		o.rotationTime = t
	}
}

// WithCompression
func WithCompression(open bool) Option {
	return func(o *Options) {
		o.compression = open
	}
}

type Options struct {
	// filename generation rule
	gtr string
	// maxSize bufsize default Mb
	bufSize int
	// maxSize file size default Mb
	maxSize int
	// log survival days default day
	maxAge int
	// rotationTime
	rotationTime time.Duration
	// requriedTimezone
	requriedTimezone bool
	// compression
	compression bool
}

func NewFileout(name string, opts ...Option) (*fileout, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	strf, err := strftime.New(name)
	if err != nil {
		return nil, fmt.Errorf("time format invalid %s", err)
	}

	match := name
	for _, re := range patternConversionRegexps {
		match = re.ReplaceAllString(match, "*") + "*"
	}

	return &fileout{
		opt:   o,
		strf:  strf,
		match: match,
	}, nil
}

type fileout struct {
	opt *Options

	strf *strftime.Strftime

	currTime time.Time

	match string

	w *bufio.Writer

	fr *os.File

	size int64

	mu sync.Mutex
	// avoid duplicate files
	generation int
	// handler age and log file
	oldStuff chan *Event

	sw sync.WaitGroup

	startMill sync.Once
}

type Event struct {
	// tp =0-> handler oldfiles >0 -> gzip/rename
	tp  int
	msg string
}

func (l *fileout) test() int {
	return l.opt.maxSize
}

// maxSize returns the maximum size in bytes of log files
func (l *fileout) maxSize() int64 {
	if l.opt.maxSize <= 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.opt.maxSize) * int64(megabyte)
}

// bufsize
func (l *fileout) bufsize() int {
	if l.opt.bufSize < 0 {
		return 1024
	}
	if l.opt.bufSize == 0 {
		return defaultBufSize * megabyte
	}
	return l.opt.bufSize * megabyte
}

// rotationTime
func (l *fileout) rotationTime() time.Duration {
	if l.opt.rotationTime < 0 {
		return 0
	}
	if l.opt.rotationTime < time.Minute {
		return time.Minute
	}
	return l.opt.rotationTime * time.Minute
}

// maxAge
func (l *fileout) maxAge() time.Duration {
	if l.opt.maxAge > 0 {
		return time.Duration(l.opt.maxAge) * 24 * time.Hour
	}
	return 24 * 365 * time.Hour
}

// Sync
func (d *fileout) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.w.Flush()
}

// Close
func (d *fileout) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.close()
}

// close
func (d *fileout) close() error {
	var errs error
	if d.w != nil {
		if err := d.w.Flush(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	if d.fr != nil {
		if err := d.fr.Close(); err != nil {
			errs = multierr.Append(errs, err)
		}
		select {
		case d.oldStuff <- &Event{tp: 1, msg: d.fr.Name()}:
			d.sw.Add(1)
		case <-time.After(time.Millisecond * 10):
		}
	}
	d.sw.Wait()

	return errs
}

// getWriter
func (d *fileout) getWriter(b []byte, createFile bool) (io.Writer, error) {
	var filename string
	var gentime bool

	rotationtime := d.rotationTime()
	// use rotationtime gen files
	if rotationtime > 0 {
		gentime = false
		d.currTime = time.Now()
		filename = GenRolaFileName(d.strf, d.currTime, rotationtime, d.generation, d.opt.requriedTimezone, templog)
	}

	writeLen := int64(len(b))

	var forceNewFile bool
	// create new file
	if d.fr == nil || (d.size+writeLen) > d.maxSize() || (rotationtime > 0 && filename != d.fr.Name()) {
		if (d.size + writeLen) > d.maxSize() {
			// avoid duplicate files
			d.generation++
		}
		forceNewFile = true
	}
	var errs error
	if forceNewFile {
		if gentime {
			d.currTime = time.Now()
		}
		d.startMill.Do(func() {
			d.oldStuff = make(chan *Event, 5)
			go d.stduffRun()
		})
		select {
		case d.oldStuff <- &Event{tp: 0, msg: d.match}:
			d.sw.Add(1)
		case <-time.After(time.Millisecond * 10):
		}
		// Prevent duplicate filenames after restart
		if d.currTime.Sub(currentTime()).Milliseconds() < dittotime && d.generation < 1 {
			d.generation = globalRand.Intn(randnum)
		}
		for {
			filename = GenRolaFileName(d.strf, d.currTime, rotationtime, d.generation, d.opt.requriedTimezone, templog)
			_, err := os_Stat(d.rename(filename))
			if err != nil {
				break
			}
			d.generation++
		}
		if !createFile {
			if err := d.close(); err != nil {
				errs = multierr.Append(errs, err)
			}
			d.size = 0
			d.fr = nil
			return d.w, errs
		}
		if err := d.close(); err != nil {
			errs = multierr.Append(errs, err)
		}
		nf, err := d.createFile(filename)
		if err != nil {
			return nil, multierr.Append(errs, err)
		}
		if d.w != nil {
			d.w.Reset(nf)
		} else {
			d.w = bufio.NewWriterSize(nf, d.bufsize())
		}
		d.size = 0
		d.fr = nf
	}
	return d.w, errs
}

// olderStduffRun runs in a goroutine to manage post-rotation compression and removal
// of old log files.
func (d *fileout) stduffRun() {
	tick := time.Tick(d.rotationTime())
	for {
		select {
		case stduff := <-d.oldStuff:
			_ = d.stduffHandler(stduff)
			d.sw.Done()
		case <-tick:
			if len(d.oldStuff) == 0 {
				d.mu.Lock()
				d.getWriter(nil, false)
				d.mu.Unlock()
			}
		}
	}
}

// stduffHandler rename/remove/callback old files
func (d *fileout) stduffHandler(stduff *Event) error {
	if stduff.tp > 0 {
		if d.opt.compression {
			gzName := d.rename(d.fr.Name()) + compressSuffix
			return GzipFile(d.fr.Name(), gzName)
		}
		return d.renameFile(d.fr.Name())
	}
	matches, err := filepath.Glob(stduff.msg)
	if err != nil {
		return err
	}
	for _, fullName := range matches {
		f, err := os.Stat(fullName)
		if err != nil {
			continue
		}
		if d.currTime.Sub(f.ModTime()) > d.maxAge() {
			os.Remove(fullName)
			continue
		}
		if strings.HasSuffix(fullName, templog) {
			if d.currTime.Sub(f.ModTime()) >= time.Duration(float64(d.rotationTime())*0.1) {
				if d.opt.compression {
					cpname := d.rename(fullName) + compressSuffix
					return GzipFile(fullName, cpname)
				} else {
					d.renameFile(f.Name())
				}
			}
		}
	}
	return nil
}

// oldLogFiles
func (d *fileout) oldLogFiles(stdff string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(stdff)
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	return files, nil
}

// renameFile
func (d *fileout) renameFile(fullName string) error {
	if na := d.rename(fullName); na != "" {
		return os.Rename(fullName, na)
	}
	return nil
}

// rename
func (d *fileout) rename(old string) string {
	var name string
	filenames := strings.Split(old, templog)
	lenth := len(filenames)
	switch {
	case lenth == 2:
		name = filenames[0]
	case lenth > 2:
		for n, spe := range filenames {
			if n < lenth-2 {
				name += spe + templog
			}
			if n == lenth-2 && spe != "" {
				name += spe
			}
		}
	}
	return name
}

// Write
func (d *fileout) Write(b []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var errs error
	w, err := d.getWriter(b, true)
	errs = multierr.Append(errs, err)
	if err != nil && w == nil {
		return 0, fmt.Errorf("failed to acquite target io.Writer, cause: %s", errs.Error())
	}
	n, err = w.Write(b)
	if err != nil {
		return 0, multierr.Append(errs, err)
	}
	d.size += int64(n)
	return n, errs
}

// CreateFile creates a new file in the given path, creating parent directories
func (d *fileout) createFile(filename string) (*os.File, error) {
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s", dirname)
	}
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", filename, err)
	}

	return fh, nil
}

func (d *fileout) RuningInfo() map[string]interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()
	runing := make(map[string]interface{}, 0)
	runing["bufSize"] = fmt.Sprintf("%.4fMB", float64(d.bufsize())/float64(megabyte))
	runing["maxSize"] = fmt.Sprintf("%dMB", d.maxSize()/int64(megabyte))
	runing["rotationTime"] = fmt.Sprintf("%ds", d.rotationTime())
	runing["currentSize"] = d.size
	if d.fr != nil {
		runing["currentName"] = d.fr.Name()
	}
	return runing
}
