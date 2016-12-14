package gweb

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	http2 "github.com/eynstudio/gobreak/net/http"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/di"
)

var ctxType = reflect.TypeOf(&Ctx{})

type Ctx struct {
	*App
	di.Container
	*Req
	*Resp
	Scope   M
	isErr   bool
	afters  []Handler
	Handled bool
	Uid     GUID
}

func (p *Ctx) Error(code int) *Ctx {
	p.WriteHeader(code)
	p.isErr = true
	return p
}

func (p *Ctx) ErrorJson(code int, m T) *Ctx {
	p.Handled = true
	if code != http.StatusOK {
		p.WriteHeader(code)
		p.isErr = true
	}
	if m == nil {
		return p
	}
	if b, err := json.Marshal(m); err != nil {
		p.Error(http.StatusInternalServerError)
	} else {
		p.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		p.Resp.Write(b)
	}
	return p
}

func (p *Ctx) Set(k string, v T)   { p.Scope[k] = v }
func (p *Ctx) IsErr() bool         { return p.isErr }
func (p *Ctx) Get(k string) string { return p.Scope.GetStr(k) }

func (p *Ctx) IsGET() bool  { return p.Method == GET }
func (p *Ctx) IsPOST() bool { return p.Method == POST }
func (p *Ctx) IsPUT() bool  { return p.Method == PUT }
func (p *Ctx) IsDEL() bool  { return p.Method == DEL }

func (p *Ctx) SetHandled()         { p.Handled = true }
func (p *Ctx) OK()                 { p.WriteHeader(http.StatusOK) }
func (p *Ctx) NotFound()           { p.Error(http.StatusNotFound) }
func (p *Ctx) Forbidden()          { p.Error(http.StatusForbidden) }
func (p *Ctx) Redirect(url string) { http.Redirect(p.Resp, p.Request, url, http.StatusFound) }
func (p *Ctx) HandleStatusJson(s IStatus) {
	p.Json(s)
	p.Handled = true
}

func (p *Ctx) Json(m T) *Ctx {
	if p.IsErr() {
		return p
	}
	if b, err := json.Marshal(m); err != nil {
		p.Error(http.StatusInternalServerError)
	} else {
		p.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		p.Resp.Write(b)
	}
	return p
}

func (p *Ctx) Text(str string) {
	if p.IsErr() {
		return
	}
	p.Resp.Write([]byte(str))
}

func (p *Ctx) SetCookie(c http.Cookie) { http.SetCookie(p.Resp, &c) }

func (p *Ctx) Tmpl(tpl string, o T) {
	p.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := p.App.Tmpl.Execute(p.Resp, tpl, o); err != nil {
		log.Println(err)
		p.Error(http.StatusInternalServerError)
	}
}

func (p *Ctx) ServeFile() bool {
	url := p.Url()
	for _, path := range p.Cfg.ServeFiles {
		if strings.HasPrefix(url, path) {
			if fi, err := os.Stat(url[1:]); err != nil || fi.IsDir() {
				p.NotFound()
			} else {
				http.ServeFile(p.Resp, p.Request, url[1:])
			}
			return true
		}
	}
	return false
}

func (p Ctx) IsAuth() bool {
	return !p.Uid.IsEmpty()
}

func (p *Ctx) UserId() GUID {
	if !p.HasToken() {
		return GUID("")
	}
	if p.Uid != "" {
		return p.Uid
	}
	uid, _ := p.Sessions.GetSessUid(p.Token)
	p.Uid = GUID(uid)
	return p.Uid
}

func (p *Ctx) ReqIp() string { return http2.GetReqIp(p.Request) }
