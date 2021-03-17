package gweb

import (
	"reflect"
	"strings"
)

type INode interface {
	AddNode(n INode) INode
	NewParamNode(path string, auth bool) INode
	NewRegexNode(path, regex string, auth bool) INode
	NewSubNode(path string, auth bool) INode
	NewHandleNode(path string, handle func(*Ctx), auth bool) INode
	CanRoute(test string, c *Ctx) bool
	Handle(c *Ctx)
	CheckHttpMethod(n INode)
	RunInterceptors(c *Ctx) INode
	GetNodes() []INode
	NeedAuth() bool
	Actions() map[string]NodeAction
}

type NodeAction struct {
	Name string
	Type reflect.Type
}

type Node struct {
	Path         string
	Interceptors []*Interceptor
	Nodes        []INode
	needAuth     bool
	actions      map[string]NodeAction
}

func NewNode(path string, auth bool) (pn *Node) {
	paths := strings.Split(path, "/")
	var cn INode
	for i, it := range paths {
		if i == 0 {
			pn = newNode(it, auth)
			cn = pn
		} else if strings.HasPrefix(it, "{") && strings.HasSuffix(it, "}") {
			cn = cn.NewParamNode(strings.Trim(it, "{}"), auth)
		} else {
			cn = cn.NewSubNode(it, auth)
		}
	}
	return
}

func newNode(path string, auth bool) *Node {
	return &Node{
		Path:         path,
		Interceptors: []*Interceptor{},
		Nodes:        []INode{},
		needAuth:     auth,
		actions:      make(map[string]NodeAction),
	}
}

func (p *Node) NewSubNode(path string, auth bool) INode {
	return p.addNode(newNode(path, auth))
}

func (p *Node) NewParamNode(path string, auth bool) INode {
	return p.addNode(NewParamNode(path, auth))
}

func (p *Node) NewRegexNode(path, regex string, auth bool) INode {
	return p.addNode(NewRegexNode(path, regex, auth))
}

func (p *Node) NewHandleNode(path string, handle func(*Ctx), auth bool) INode {
	return p.addNode(NewHandleNode(path, handle, auth))
}

func (p *Node) AddNode(n INode) INode {
	p.addNode(n)
	return p
}
func (p *Node) CheckHttpMethod(n INode) {
	nodeType := reflect.TypeOf(n)
	for i, j := 0, nodeType.NumMethod(); i < j; i++ {
		m := nodeType.Method(i)
		if isHttpMethod(m) {
			p.actions[strings.ToLower(m.Name)] = NodeAction{m.Name, m.Type}
		}
	}
}
func (p *Node) addNode(n INode) INode {
	n.CheckHttpMethod(n)
	p.Nodes = append(p.Nodes, n)
	return n
}

func (p *Node) Handle(c *Ctx) {}

func (p *Node) Actions() map[string]NodeAction { return p.actions }

func (p *Node) GetNodes() []INode                 { return p.Nodes }
func (p *Node) NeedAuth() bool                    { return p.needAuth }
func (p *Node) CanRoute(test string, c *Ctx) bool { return p.Path == test }

func (p *Node) Interceptor(m *Interceptor) *Node {
	p.Interceptors = append(p.Interceptors, m)
	return p
}

func (p *Node) RunInterceptors(c *Ctx) INode {
	if c.IsErr() {
		return p
	}

	for _, i := range p.Interceptors {
		if nil != i.After {
			c.afters = append(c.afters, i.After)
		}

		if nil != i.Before {
			i.Before(c)
			if c.IsErr() {
				break
			}
		}
	}

	return p
}

func isHttpMethod(m reflect.Method) bool {
	lst := []string{"Get", "Post", "Del", "Put"}
	for _, it := range lst {
		if strings.HasPrefix(m.Name, it) && m.Type.NumIn() > 1 && m.Type.In(1) == ctxType {
			return true
		}
	}
	return false
}
