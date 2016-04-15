package main

import (
	"log"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gweb"
)

type Home struct {
	*gweb.Node
}

func NewHome() *Home {
	h := &Home{gweb.NewNode("", false)}
	h.NewParamNode("id", false)
	return h
}

func (p *Home) Handler(c *gweb.Ctx) {
	handled := true
	switch c.Method {
	case "GET":
		p.get(c)
	case "POST":
		p.post(c)
	default:
		handled = false
	}
	c.Handled = handled
}

func (p *Home) get(c *gweb.Ctx) {
	c.Json(M{"get": "aa"})
}

func (p *Home) post(c *gweb.Ctx) {
	var h H
	c.JsonBody(&h)
	log.Println(h)
	c.Json(M{"post": "aa"})
}

type H struct {
	Id int
}
