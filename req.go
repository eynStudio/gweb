package gweb

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db/filter"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
	DEL  = "DELETE"
)

type Req struct {
	*http.Request
	*urlParts
	Token string
}

func NewReq(r *http.Request) *Req {
	rr := &Req{Request: r}
	rr.parseToken()
	rr.urlParts = newUrlParts(rr.Url())
	return rr
}

func (p *Req) parseToken() {
	if authHeader := p.Header.Get("Authorization"); authHeader != "" {
		if tokens := strings.Split(authHeader, " "); len(tokens) == 2 {
			p.Token = tokens[1]
		}
		return
	}
	if t := p.URL.Query().Get("token"); t != "" {
		p.Token = t
	}
}

func (p *Req) Url() string     { return p.URL.Path }
func (p *Req) HasToken() bool  { return len(p.Token) > 0 }
func (p *Req) JMethod() string { return p.Header.Get("jBreak-Method") }
func (p *Req) IsJPage() bool   { return p.JMethod() == "Page" }
func (p *Req) IsJList() bool   { return p.JMethod() == "List" }
func (p *Req) JPage() *filter.PageFilter {
	var pf filter.PageFilter
	if p.IsJsonContent() && p.IsJPage() && p.JsonBody(&pf) {
		return &pf
	}
	return nil
}

func (p *Req) JsonBody(m T) bool {
	if p.IsJsonContent() && p.Body != nil {
		defer p.Body.Close()
		if err := json.NewDecoder(p.Body).Decode(&m); err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false
}

func (p *Req) IsJsonContent() bool {
	return strings.Contains(p.Header.Get("Content-Type"), "application/json")
}
func (p *Req) IsAcceptJson() bool { return strings.Contains(p.Header.Get("Accept"), "application/json") }

type urlPart struct {
	path string
}

type urlParts struct {
	curIdx int
	parts  []*urlPart
}

func newUrlParts(path string) *urlParts {
	m := &urlParts{}
	m.parseParts(path)
	return m
}

func (p *urlParts) parseParts(path string) {
	parts := strings.Split(path, "/")
	for _, it := range parts {
		p.parts = append(p.parts, &urlPart{it})
	}
}
func (p *urlParts) moveNextPart()     { p.curIdx += 1 }
func (p *urlParts) hasNextPart() bool { return p.curIdx < len(p.parts)-1 }
func (p *urlParts) CurPart() *urlPart { return p.parts[p.curIdx] }
func (p *urlParts) NextPart() *urlPart {
	if p.hasNextPart() {
		return p.parts[p.curIdx+1]
	}
	return nil
}
