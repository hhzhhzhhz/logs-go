package fileout

import (
	"fmt"
	"go.uber.org/multierr"
	"logs-go/utils"
	"path/filepath"
	"testing"
)

func Test_Split(t *testing.T) {
	l, err := NewFileout("log")
	if err != nil {
		t.Error(err)
	}
	name := ".tmpxxxx.tmp"
	if l.rename(name) != ".tmpxxxx" {
		t.Error("error")
	}
	name = "xxxx.tmp"
	if l.rename(name) != "xxxx" {
		t.Error("error")
	}
	name = "x.tmpx.tmpxx.tmp"
	if l.rename(name) != "x.tmpx.tmpxx" {
		t.Error("error")
	}
	name = ".tmpx.tmpx.tmpxx.tmp"
	if l.rename(name) != ".tmpx.tmpx.tmpxx" {
		t.Error("error")
	}
	name = "test.tm.mp.tlp"
	if l.rename(name) != "" {
		t.Error("error")
	}
}

func Test_log(t *testing.T) {
	//t.Skip()
	l, err := NewFileout("log-")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := l.Write([]byte("hello world!\n")); err != nil {
			t.Error(err)
		}
		//time.Sleep(500 * time.Millisecond)
	}
	l.Close()
}

func Test_Rang_Dir(t *testing.T) {
	t.Skip()
	l, err := NewFileout("./log/%Y/log-%H")
	if err != nil {
		t.Error(err)
	}
	list, err := filepath.Glob(l.match)
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func Test_compressLogFile(t *testing.T) {
	//t.Skip()
	if err := utils.GzipFile("log-.7", "log-.7.gz"); err != nil {
		t.Error(err.Error())
	}
}

func Test_errors(t *testing.T) {
	var errs error
	errs = multierr.Append(errs, fmt.Errorf("xxx"))
	errs = multierr.Append(errs, fmt.Errorf("xxx1"))
	errs = multierr.Append(errs, fmt.Errorf("xxx2"))
	errs = multierr.Append(errs, fmt.Errorf("xxx3"))
	t.Log(errs.Error())

}
