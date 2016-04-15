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
		n.Handler(c)
	}
}

func (p *Router) RouteSubNodes(n INode, c *Ctx) {
	for _, it := range n.GetNodes() {
		if it.NeedAuth() && !c.hasUserId() {
			continue
		}
		if it.CanRouter(c.NextPart().path) {
			c.moveNextPart()
			p.Route(it, c)
			break
		}
	}
}
