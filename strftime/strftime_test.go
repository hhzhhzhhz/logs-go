package strftime

import (
	"testing"
	"time"
)

func Test_Gen_Name(t *testing.T) {
	s, err := New("./log/%Y%m%d")
	if err != nil {
		t.Error(err)
	}
	n := time.Now()
	t.Log(s.FormatString(n.Truncate(time.Minute)))
}
