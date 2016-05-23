package gweb

import (
	"regexp"
)

type RegexNode struct {
	*Node
	regex   string
	curPath string
}

func NewRegexNode(path, regex string, auth bool) *RegexNode {
	return &RegexNode{Node: NewNode(path, auth), regex: regex}
}

func (p *RegexNode) CanRoute(test string) bool {
	match, _ := regexp.MatchString(p.regex, test)
	if match {
		p.curPath = test
	}
	return match
}

func (p *RegexNode) Handle(c *Ctx) {
	c.Scope[p.Path] = p.curPath
}
