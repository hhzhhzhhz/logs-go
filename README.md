# LGOS-GO

logs-go is based on uber zap package, which can quickly write to the file system and connect to rsyslog.

## Features

* [Blazing fast](#benchmarks)
* Support file output according to regular format
* Support for rsyslog output
* graceful shutdown

## Installation

```bash
go get -u https://github.com/hhzhhzhhz/logs-go
```

#Benchmarks
See ..logs-go\log-go_test.go
```text
BenchmarkForFile-12    	  874861	      1358 ns/op	     373 B/op	       6 allocs/op
BenchmarkOneForFile-12    	  422335	      2414 ns/op	     370 B/op	       6 allocs/op
```
## Getting Started

### Simple Logging Example

For simple logging, output rsyslog

```go
func Test_log_rsyslog(t *testing.T) {
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "rsyslog_tag"
	cfg := logs_go.NewDefaultConfig()
	cfg.InitialFields = fileds
	cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
	cfg.Stdout = true
	l, err := cfg.Build()
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
func Test_log_file(t *testing.T) {
	cfg := logs_go.NewDefaultConfig()
	cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
	cfg.Stdout = true
	l, err := cfg.Build()
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
func Test_logs_stdout(t *testing.T) {
	cfg := logs_go.NewDefaultConfig()
	cfg.Stdout = true
	l, err := cfg.Build()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_logs_stdout", zap.String("out", "stdout"))
	l.Close()
}
stdout: {"level":"info","timestamp":"2022-04-07T00:15:53.759+08:00","caller":"logs-go/logs-go.go:24","tag":"Test_logs_stdout","out":"stdout"}
```
In this case, many consumers will take the last value, but this is not guaranteed; check yours if in doubt.
