package gweb

type INode interface {
	AddNode(n INode) INode
	NewParamNode(path string, auth bool) INode
	NewRegexNode(path, regex string, auth bool) INode
	CanRouter(test string) bool
	Handler(c *Ctx)
	RunInterceptors(c *Ctx) INode
	GetNodes() []INode
	NeedAuth() bool
}

type Node struct {
	Path         string
	Interceptors []*Interceptor
	Nodes        []INode
	needAuth     bool
}

func NewNode(path string, auth bool) *Node {
	return &Node{
		Path:         path,
		Interceptors: []*Interceptor{},
		Nodes:        []INode{},
		needAuth:     auth,
	}
}

func (p *Node) NewParamNode(path string, auth bool) INode {
	n := NewParamNode(path, auth)
	p.Nodes = append(p.Nodes, n)
	return n
}
func (p *Node) NewRegexNode(path, regex string, auth bool) INode {
	n := NewRegexNode(path, regex, auth)
	p.Nodes = append(p.Nodes, n)
	return n
}

func (p *Node) AddNode(n INode) INode {
	p.Nodes = append(p.Nodes, n)
	return p
}

func (p *Node) Handler(c *Ctx)             {}
func (p *Node) GetNodes() []INode          { return p.Nodes }
func (p *Node) NeedAuth() bool             { return p.needAuth }
func (p *Node) CanRouter(test string) bool { return p.Path == test }

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
