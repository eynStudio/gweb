package gweb

type HandleNode struct {
	*Node
	handle func(c *Ctx)
}

func NewHandleNode(path string, handle func(*Ctx), auth bool) *HandleNode {
	return &HandleNode{Node: NewNode(path, auth), handle: handle}
}

func (p *HandleNode) Handle(c *Ctx) {
	p.handle(c)
	c.Handled = true
}
