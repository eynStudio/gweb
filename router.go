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
	in := []reflect.Value{reflect.ValueOf(c)}
	actions := p.findActions(c)
	for _, it := range actions {
		if m, ok := n.Actions()[it]; ok {
			c.Handled = true
			reflect.ValueOf(n).MethodByName(m.Name).Call(in)
			return c.Handled
		}
	}
	return false
}
func (p *Router) findActions(c *Ctx) (actions []string) {
	method := strings.ToLower(c.Req.Method)
	jBreakMethod := c.Req.Header.Get("jBreak-Method")
	if jBreakMethod != "" {
		jBreakMethod = strings.ToLower(jBreakMethod)
	}

	hasId := c.Scope.HasKey("id")
	appendActions := func(act string) {
		if jBreakMethod != "" {
			actions = append(actions, method+act+jBreakMethod)
		}
		if hasId {
			actions = append(actions, method+act+"id")
		}
		actions = append(actions, method+act)
	}

	if c.Scope.HasKeys("act", "act1") {
		appendActions(strings.ToLower(c.Get("act") + c.Get("act1")))
	}
	if c.Scope.HasKey("act1") {
		appendActions(strings.ToLower(c.Get("act1")))
	}
	if c.Scope.HasKey("act") {
		appendActions(c.Get("act"))
	}
	appendActions("")
	return
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
