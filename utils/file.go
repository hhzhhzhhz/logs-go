package utils

import (
    "logs-go/strftime"
    "os"
    "time"
)

// FileExist determine if the file exists.
func FileExist(file string) bool {
    _, err := os.Stat(file)
    return err == nil || os.IsExist(err)
}

// GenRolaFileName
func GenRolaFileName(pattern *strftime.Strftime, now time.Time, rotationTime time.Duration) string {
    var base time.Time
    if now.Location() != time.UTC {
        base = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
        base =    base.Truncate(rotationTime)
        base = time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
    } else {
        base = now.Truncate(rotationTime)
    }

    return pattern.FormatString(base)
}

