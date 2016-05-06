package gweb

import (
	"time"

	"github.com/eynstudio/gobreak/log/datelog"
)

type Log struct {
	dl *datelog.DateLog
}

func NewLog(path string) *Log {
	return &Log{dl: datelog.New(path)}
}

func (p *Log) Log(r *Req) {
	p.dl.Logf("%s %s %s %s %s %s %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		r.RemoteAddr,
		r.Proto,
		r.Method,
		r.URL.RequestURI(),
		r.Referer(),
		r.UserAgent())
}
