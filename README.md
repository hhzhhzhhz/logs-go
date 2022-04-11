# LOGS-GO
    logs-go is a log component based on uber-go and native log 
    support output console file and rsyslog.

## Features

* [Blazing fast](#benchmarks)
* Support output to rsyslog/file/sdout 
* Support json format
* Support custom format
* graceful shutdown

## Installation

```bash
go get -u https://github.com/hhzhhzhhz/logs-go
```

## Benchmarks
```text
output file
logs-go    	   768952	      1547 ns/op	     461 B/op	       7 allocs/op
lumberjack    	   307975	      3960 ns/op	     344 B/op	       4 allocs/op
rotatelogs    	   31593	     37487 ns/op	     802 B/op	       9 allocs/op
```
## Getting Started

### Simple Logging Example

For simple logging, output rsyslog

```go
t.Run("rsyslog", func(t *testing.T) {
    cfg := logs_go.NewSimpleConfig()
    cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
    cfg.Stdout = true
    l, err := cfg.BuildSimpleLog()
    if err != nil {
        t.Error(err)
    }
    l.Info("Test_log_rsyslog", "rsyslog")
    l.Close()
})
output: 2022/04/07 23:32:39 log_example_test.go:68: [INFO] Test_log_rsyslog

func Test_log_rsyslog(t *testing.T) {
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "rsyslog_tag"
	cfg := logs_go.NewJsonConfig()
	cfg.InitialFields = fileds
	cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
	cfg.Stdout = true
	l, err := cfg.BuildJsonLog()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_log_rsyslog", zap.String("out", "rsyslog"))
	l.Close()
}
output: {"level":"info","timestamp":"2022-04-07T00:10:30.953+08:00","caller":"logs-go/logs-go.go:24","tag":"Test_log_rsyslog","@rsyslog_tag":"rsyslog_tag","out":"rsyslog"}
```

For simple logging, output files
```go
t.Run("file", func(t *testing.T) {
    cfg := logs_go.NewSimpleConfig()
    cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
    cfg.Stdout = true
    l, err := cfg.BuildSimpleLog()
    if err != nil { 
        t.Error(err)
    }
    l.Info("Test_log_file %s", "file")
    l.Close()
})
output: 2022/04/07 23:34:14 log_example_test.go:57: [INFO] Test_log_file file

func Test_log_file(t *testing.T) {
	cfg := logs_go.NewJsonConfig()
	cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
	cfg.Stdout = true
	l, err := cfg.BuildJsonLog()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_log_file", zap.String("out", "file"))
	l.Close()
}
output: {"level":"info","timestamp":"2022-04-07T00:15:24.368+08:00","caller":"logs-go/logs-go.go:24","tag":"Test_log_file","out":"file"}
```

For simple logging, output stdout
```go
t.Run("stdout", func(t *testing.T) {
    cfg := logs_go.NewSimpleConfig()
    cfg.Stdout = true
    l, err := cfg.BuildSimpleLog()
    if err != nil {
        t.Error(err)
    }
    l.Info("Test_logs_stdout %s", "stdout")
    l.Close()
})
output: 2022/04/07 23:29:57 log_example_test.go:78: [INFO] Test_logs_stdout stdout

func Test_logs_stdout(t *testing.T) {
	cfg := logs_go.NewJsonConfig()
	cfg.Stdout = true
	l, err := cfg.BuildJsonLog()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_logs_stdout", zap.String("out", "stdout"))
	l.Close()
}
stdout: {"level":"info","timestamp":"2022-04-07T00:15:53.759+08:00","caller":"logs-go/logs-go.go:24","tag":"Test_logs_stdout","out":"stdout"}
```
In this case, many consumers will take the last value, but this is not guaranteed; check yours if in doubt.
