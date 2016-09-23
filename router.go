package gweb

import (
	"encoding/json"
	"log"
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

func GetFuncArgs(t reflect.Type) []reflect.Type {
	l := t.NumIn()
	in := make([]reflect.Type, l)
	for i := 0; i < l; i++ {
		in[i] = t.In(i)
	}
	return in
}

func (p *Router) autoHandle(n INode, c *Ctx) bool {
	actions := p.findActions(c)
	for _, it := range actions {
		if act, ok := n.Actions()[it]; ok {
			c.Handled = true
			in := []reflect.Value{reflect.ValueOf(c)}
			m := reflect.ValueOf(n).MethodByName(act.Name)
			args := GetFuncArgs(m.Type())
			if len(args) == 2 {
				obj := reflect.New(args[1].Elem()).Interface()
				if err := json.NewDecoder(c.Body).Decode(obj); err != nil {
					log.Println(err)
				}
				defer c.Body.Close()
				in = append(in, reflect.ValueOf(obj))
			}
			m.Call(in)
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
