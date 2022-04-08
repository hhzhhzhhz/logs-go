package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"logs-go/strftime"
	"os"
	"strconv"
	"time"
)

const (
	flag = "."
)

// FileExist determine if the file exists.
func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

// GenRolaFileName
func GenRolaFileName(pattern *strftime.Strftime, now time.Time, rotationTime time.Duration, generation int, requriedTimezone bool, compensate string) string {
	var base time.Time
	if requriedTimezone {
		base = now.Truncate(rotationTime)
		return pattern.FormatString(base) + flag + strconv.Itoa(generation)
	}

	if now.Location() != time.UTC {
		base = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
		base = base.Truncate(rotationTime)
		base = time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
	} else {
		base = now.Truncate(rotationTime)
	}

	return pattern.FormatString(base) + flag + strconv.Itoa(generation) + compensate
}

// compressLogFile compresses the given log file, removing the
// uncompressed log file if successful.
func GzipFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	// If this file already exists, we presume it was created by
	// a previous attempt to compress the log file.
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)

	defer func() {
		if err != nil {
			os.Remove(dst)
			fmt.Println(src, dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Remove(src)
}
