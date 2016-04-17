package gweb

import (
	"fmt"
	"io"
	"net/http"
	"time"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/conf"
	"github.com/eynstudio/gobreak/di"
)

type App struct {
	Root INode
	di.Container
	Name string
	*Cfg
	Server *http.Server
	*Router
	*Tmpl
	*Sessions
}

func NewApp(name string) *App {
	var cfg Cfg
	conf.MustLoadJsonCfg(&cfg, "conf/"+name+".json")
	return NewAppWithCfg(&cfg)
}

func NewAppWithCfg(c *Cfg) *App {
	app := &App{
		Root:      NewNode("", false),
		Container: di.Root,
		Name:      "",
		Cfg:       c,
		Router:    &Router{},
		Tmpl:      &Tmpl{},
		Sessions:  &Sessions{},
	}

	if c.UseTmpl {
		app.Tmpl.Load()
	}
	app.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		Handler:      http.HandlerFunc(app.handler),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	app.Map(c)
	return app
}

func (P *App) Start() {
	if P.Cfg.Tls {
		err := P.Server.ListenAndServeTLS(P.Cfg.CertFile, P.Cfg.KeyFile)
		if err != nil {
			panic(err)
		}
	} else {
		err := P.Server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}
}

func (p *App) NewCtx(r *http.Request, rw http.ResponseWriter) *Ctx {
	req := NewReq(r)
	resp := &Resp{rw}
	c := &Ctx{App: p, Container: di.New(), Req: req, Resp: resp, Scope: M{}}
	c.Map(c) //需要吗？
	c.Map(resp)
	c.Map(req)
	c.SetParent(p)
	return c
}

func (p *App) handler(w http.ResponseWriter, r *http.Request) {
	ctx := p.NewCtx(r, w)
	if !ctx.ServeFile() {
		p.Route(p.Root, ctx)
		if !ctx.Handled {
			ctx.NotFound()
		}
	}

	if w, ok := ctx.Resp.ResponseWriter.(io.Closer); ok {
		w.Close()
	}
}
