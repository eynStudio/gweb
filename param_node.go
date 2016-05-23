package gweb

type ParamNode struct {
	*Node
	curPath string
}

func NewParamNode(path string, auth bool) *ParamNode { return &ParamNode{Node: NewNode(path, auth)} }

func (p *ParamNode) CanRoute(test string) bool {
	p.curPath = test
	return true
}

func (p *ParamNode) Handle(c *Ctx) {
	c.Scope[p.Path] = p.curPath
}
