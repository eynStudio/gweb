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
	Sessions ISessions
	*Log
	NeedAuth       func(*Ctx)
	SessionExpired func(*Ctx)
	NotFoundHandle func(*Ctx)
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
		Log:       NewLog("./logs/serv"),
		Router:    &Router{},
		Tmpl:      &Tmpl{},
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

func (p *App) Start() {
	p.injectNodes(p.Root)

	if p.Cfg.Tls {
		err := p.Server.ListenAndServeTLS(p.Cfg.CertFile, p.Cfg.KeyFile)
		if err != nil {
			panic(err)
		}
	} else {
		err := p.Server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}
}

func (p *App) injectNodes(n INode) {
	p.Apply(n)
	nodes := n.GetNodes()
	for i, l := 0, len(nodes); i < l; i++ {
		p.injectNodes(nodes[i])
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
			if p.NotFoundHandle != nil {
				p.NotFoundHandle(ctx)
			} else {
				ctx.NotFound()
			}
		}
	}

	p.Log.Log(ctx.Req)
	if w, ok := ctx.Resp.ResponseWriter.(io.Closer); ok {
		w.Close()
	}
}
