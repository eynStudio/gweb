package gweb

type ParamNode struct {
	*Node
}

func NewParamNode(path string, auth bool) *ParamNode { return &ParamNode{NewNode(path, auth)} }

func (p *ParamNode) CanRouter(test string) bool { return true }

func (p *ParamNode) Handler(c *Ctx) {
	c.Scope[p.Path] = c.CurPart().path
}
