package main

import (
	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gweb"
)

type Api struct {
	*gweb.Node
}

func NewApi() *Api {
	h := &Api{gweb.NewNode("api", false)}
	h.NewParamNode("id", false)
	return h
}

func (p *Api) Handler(c *gweb.Ctx) {
	handled := true
	switch c.Method {
	case "GET":
		p.Get(c)
	default:
		handled = false
	}
	c.Handled = handled
}

func (p *Api) Get(c *gweb.Ctx) {
	c.Json(M{"haah": "aaffd"})
}
