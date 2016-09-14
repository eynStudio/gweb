package gweb

import (
	"reflect"
	"strings"
)

type Router struct {
}

func (p *Router) Route(n INode, c *Ctx) {
	if c.IsErr() {
		return
	}
	n.RunInterceptors(c)

	if c.hasNextPart() {
		p.RouteSubNodes(n, c)
	}

	if !c.Handled {
		p.autoHandle(n, c)
	}
	if !c.Handled {
		n.Handle(c)
	}
}

func (p *Router) autoHandle(n INode, c *Ctx) bool {
	method := strings.ToLower(c.Req.Method)
	in := []reflect.Value{reflect.ValueOf(c)}
	if m, ok := n.Actions()[method]; ok {
		reflect.ValueOf(n).MethodByName(m.Name).Call(in)
		c.Handled = true
		return true
	}
	return false
}

func (p *Router) RouteSubNodes(n INode, c *Ctx) {
	for _, it := range n.GetNodes() {
		if it.CanRoute(c.NextPart().path, c) {
			if it.NeedAuth() {
				if s := c.ValidAuth(c); s.IsErr() {
					c.HandleStatusJson(s)
					break
				}
			}
			c.moveNextPart()
			p.Route(it, c)
			break
		}
	}
}
