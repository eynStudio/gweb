package gweb

type ParamNode struct {
	*Node
}

func NewParamNode(path string, auth bool) *ParamNode { return &ParamNode{Node: NewNode(path, auth)} }

func (p *ParamNode) CanRoute(test string, c *Ctx) bool {
	c.Scope[p.Path] = test
	return true
}
