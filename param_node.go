package gweb

type ParamNode struct {
	*Node
}

func NewParamNode(path string, auth bool) *ParamNode { return &ParamNode{NewNode(path, auth)} }

func (p *ParamNode) CanRoute(test string) bool { return true }

func (p *ParamNode) Handle(c *Ctx) {
	c.Scope[p.Path] = c.CurPart().path
}
