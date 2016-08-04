package gweb

import (
	"regexp"
)

type RegexNode struct {
	*Node
	regex string
}

func NewRegexNode(path, regex string, auth bool) *RegexNode {
	return &RegexNode{Node: NewNode(path, auth), regex: regex}
}

func (p *RegexNode) CanRoute(test string, c *Ctx) bool {
	match, _ := regexp.MatchString(p.regex, test)
	if match {
		c.Scope[p.Path] = test
	}
	return match
}
