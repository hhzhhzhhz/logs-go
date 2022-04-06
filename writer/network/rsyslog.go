package network

import (
	"bytes"
)

type Priority int

const (
	// https://github.com/rsyslog/rsyslog/blob/master/runtime/rsyslog.h
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
)

func NewRsyslogCoder(prefix string) *RsyslogCoder {
	return &RsyslogCoder{
		prefix: []byte(prefix),
	}
}

type RsyslogCoder struct {
	prefix []byte
}

func (r *RsyslogCoder) Encoder(b []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(r.prefix)
	buf.Write(b)
	return buf.Bytes()
}