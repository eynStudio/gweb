package gweb

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
		n.Handle(c)
	}
}

func (p *Router) RouteSubNodes(n INode, c *Ctx) {
	for _, it := range n.GetNodes() {
		if it.CanRoute(c.NextPart().path) {
			if it.NeedAuth() {
				if !c.HasToken() || !c.HasSession() {
					c.Forbidden()
					break
				}
			}
			c.moveNextPart()
			p.Route(it, c)
			break
		}
	}
}
