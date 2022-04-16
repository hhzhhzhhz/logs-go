# LOGS-GO
    logs-go is a log component based on uber-go and native log 
    support output console file and rsyslog.

## Features

* [Blazing fast](#benchmarks)
* Support to rsyslog/file/sdout 
* Support json format
* Support custom format
* graceful shutdown

## Installation

```bash
go get -u github.com/hhzhhzhhz/logs-go
```

## Benchmarks
```text
output file
logs-go    	   768952	      1547 ns/op	     461 B/op	       7 allocs/op
lumberjack    	   307975	      3960 ns/op	     344 B/op	       4 allocs/op
rotatelogs    	   31593	     37487 ns/op	     802 B/op	       9 allocs/op
```
## Getting Started

### format Logging Example

For simple logging, output rsyslog

```go
t.Run("rsyslog", func(t *testing.T) {
    cfg := NewLogfConfig()
    cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
    cfg.Stdout = true
    l, err := cfg.BuildLogf()
    if err != nil {
        t.Error(err)
    }
    l.Info("rsyslog %s", "rsyslog")
    l.Close()
})
// output: 2022/04/16 14:31:37 log_test.go:185: [INFO] rsyslog rsyslog

t.Run("rsyslog", func(t *testing.T) {
    fileds := map[string]interface{}{}
    fileds["@rsyslog_tag"] = "rsyslog_tag"
    cfg := NewLogJconfig()
    cfg.InitialFields = fileds
    cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
    cfg.Stdout = true
    l, err := cfg.BuildLogJ()
    if err != nil {
        t.Error(err)
    }
    l.Info("rsyslog", zap.String("out", "rsyslog"))
    l.Close()
})
// output: {"level":"info","timestamp":"2022-04-16T14:31:57.338+08:00","caller":"logs-go/log.go:49","tag":"rsyslog","@rsyslog_tag":"rsyslog_tag","out":"rsyslog"}
```

For simple logging, output disk
```go
t.Run("disk", func(t *testing.T) {
    cfg := NewLogfConfig()
    cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
    cfg.Stdout = true
    l, err := cfg.BuildLogf()
    if err != nil {
        t.Error(err)
    }
    l.Info("disk %s", "file")
    l.Close()
})
// 2022/04/16 14:28:23 log_test.go:215: [INFO] disk file

t.Run("disk", func(t *testing.T) {
    cfg := NewLogJconfig()
    cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
    cfg.Stdout = true
    l, err := cfg.BuildLogJ()
        if err != nil {
    t.Error(err)
    }
    l.Info("disk", zap.String("out", "file"))
    l.Close()
})
// output: {"level":"info","timestamp":"2022-04-16T14:28:34.688+08:00","caller":"logs-go/log.go:49","tag":"disk","out":"file"}
```

For simple logging, output stdout
```go
t.Run("stdout", func(t *testing.T) {
    cfg := NewLogfConfig()
    cfg.Stdout = true
    l, err := cfg.BuildLogf()
    if err != nil {
        t.Error(err)
    }
    l.Info("stdout %s", "stdout")
    l.Close()
})
// output: 2022/04/16 14:26:33 log_test.go:240: [INFO] stdout stdout

t.Run("stdout", func(t *testing.T) {
    cfg := NewLogJconfig()
    cfg.Stdout = true
    l, err := cfg.BuildLogJ()
    if err != nil {
        t.Error(err)
    }
    l.Info("stdout", zap.String("out", "stdout"))
	l.Close()
})
// output: {"level":"info","timestamp":"2022-04-16T14:26:53.911+08:00","caller":"logs-go/log.go:49","tag":"stdout","out":"stdout"}
```