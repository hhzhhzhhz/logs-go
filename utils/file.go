package utils

import (
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
