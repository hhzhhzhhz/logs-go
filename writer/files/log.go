package fileout

import (
    "bufio"
    "fmt"
    "io"
    "io/fs"
    "io/ioutil"
    "logs-go/strftime"
    "logs-go/utils"
    "math/rand"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

const (
    defaultMaxSize   = 100
    defaultBufSize   = 4
    defaultPath = "/var/logs"
    randnum = 50
    templog = ".temp"
)

var (
	// currentTime exists so it can be mocked out by tests.
	currentTime = time.Now

	// os_Stat exists so it can be mocked out by tests.
	os_Stat = os.Stat
	// megabyte is the conversion factor between MaxSize and bytes.  It is a
	// variable so tests can mock it out and not need to write megabytes of data
	// to file_log.
	megabyte = 1024 * 1024
    // avoid duplicate files
    dittotime = int64(10 * time.Minute)
    
    globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Option func(*Options)
// GenerationRule
func GenerationRule(gtr string) Option {
    return func(o *Options) {
        o.gtr = gtr
    }
}
// BufSize
func BufSize(mb int) Option {
    return func(o *Options) {
        o.bufSize = mb
    }
}
// MaxSize
func MaxSize(mb int) Option {
    return func(o *Options) {
        o.maxSize = mb
    }
}
// MaxAge
func MaxAge(day int) Option {
    return func(o *Options) {
        o.maxAge = day
    }
}
// RotationTime
func RotationTime(t time.Duration) Option {
    return func(o *Options) {
        o.rotationTime = t
    }
}

type Options struct {
    // filename generation rule
    gtr string
	// max bufsize default Mb
	bufSize int
	// max file size default Mb
	maxSize int
    // save time
    maxAge int
	// rotationTime
    rotationTime time.Duration
}

func NewFileout(name string, opts ...Option) (*fileout, error) {
    o := &Options{}
    for _, opt := range opts {
        opt(o)
    }
    strf, err :=  strftime.New(name)
    if err != nil {
        return nil, fmt.Errorf("time format invalid", err)
    }

	return &fileout{
        opt: o,
        strf: strf,
    }, nil
}

type fileout struct {
	opt *Options

    strf *strftime.Strftime

    currTime time.Time

    matchName string
    
	w *bufio.Writer

    fr *os.File

	size int64

	mu sync.Mutex
    // avoid duplicate files
    generation int
    // handler age and log file
    millCh chan string
    
    startMill sync.Once
}

func (d *fileout) updateConfig(opts ...Option) {
    o := &Options{}
    for _, opt := range opts {
        opt(o)
    }
    d.mu.Lock()
    defer d.mu.Unlock()
    d.opt = o
}

// max returns the maximum size in bytes of log files before rolling.
func (l *fileout) max() int64 {
    if l.opt.maxSize == 0 {
        return int64(defaultMaxSize * megabyte)
    }
    return int64(l.opt.maxSize) * int64(megabyte)
}

// bufsize
func (l *fileout) bufsize() int  {
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
    var err error
    if d.w != nil {
        d.w.Flush()
    }
    if d.fr != nil {
        err =  d.fr.Close()
        d.renameFile(d.fr.Name())
    }
    
    return err
}

// getWriter
func (d *fileout) getWriter(b []byte) (io.Writer, error) {
    // 轮训文件名
    var filename string
    if d.rotationTime() > 0 {
        filename = utils.GenRolaFileName(d.strf, d.currTime, d.rotationTime())
    }
    
    writeLen := int64(len(b))

    var forceNewFile bool
    // d.w空, 轮训条件, 文件最大值
    if d.w == nil || (d.rotationTime() > 0 && filename != d.fr.Name()) || (d.size + writeLen) > d.max() {
        if (d.size + writeLen) > d.max() {
            // 防止文件重名
            d.generation++
        }
        forceNewFile = true
    }
    
    if forceNewFile {
        d.startMill.Do(func() {
            d.millCh = make(chan string, 1)
            go d.millRun()
        })
        select {
        case d.millCh <- d.opt.gtr:
        case <-time.After(time.Millisecond * 10):
        }
        // Prevent duplicate filenames after restart
        if time.Since(currentTime()).Milliseconds() < dittotime && d.generation < 1 {
            d.generation = globalRand.Intn(randnum)
        }
        var name string
        for {
            name = utils.GenRolaFileName(d.strf, d.currTime, d.rotationTime())
            _, err := os_Stat(name)
            if err != nil {
                break 
            }
            d.generation++
        }

        nf, err := d.createFile(name)
        if err != nil {
            return nil, err
        }
        d.close()
        d.w = bufio.NewWriterSize(nf, d.bufsize())
        d.size = 0
        d.fr = nf
    }

    return d.w, nil
}

// millRun runs in a goroutine to manage post-rotation compression and removal
// of old log files.
func (d *fileout) millRun() {
    for pdir := range d.millCh {
        _ = d.millRunOnce(pdir)
    }
}

// millRunOnce
func (d *fileout) millRunOnce(projectdir string) error {
    files, err := d.oldLogFiles(projectdir)
    if err != nil {
        return err
    }
    for _, f := range files {
        if !f.IsDir() {
            if !strings.HasSuffix(f.Name(), templog) {
                if time.Now().Sub(f.ModTime()) >= d.rotationTime() * 2 || (f.Size() + 2048) > d.max() {
                    return d.renameFile(f.Name())
                }
            }
        }
    }
    return nil
}
// oldLogFiles
func (d *fileout) oldLogFiles(projectdir string) ([]fs.FileInfo, error) {
    files, err := ioutil.ReadDir(projectdir)
    if err != nil {
        return nil, fmt.Errorf("can't read log file directory: %s", err)
    }
    return files, nil
}

// renameFile
func (d *fileout) renameFile(old string) error {
    return os.Rename(old, old)
}

// Write
func (d *fileout) Write(b []byte) (n int, err error) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    w, err := d.getWriter(b)
    if err != nil {
        return 0, fmt.Errorf("failed to acquite target io.Writer, cause: %s", err.Error())
    }
    n, err =  w.Write(b)
    d.size += int64(n)
    return n, err
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
    runing["maxSize"] = fmt.Sprintf("%dMB", d.max()/int64(megabyte))
    runing["rotationTime"] = fmt.Sprintf("%ds", d.rotationTime())
    runing["currentSize"] = d.size
    if d.fr != nil {
        runing["currentName"] = d.fr.Name()
    }
    return runing
}
