package main

import (
	"github.com/elliotchance/gedcom"
	"github.com/elliotchance/gedcom/html"
)

type sourceProperty struct {
	document *gedcom.Document
	node     gedcom.Node
}

func newSourceProperty(document *gedcom.Document, node gedcom.Node) *sourceProperty {
	return &sourceProperty{
		document: document,
		node:     node,
	}
}

func (c *sourceProperty) String() string {
	tag := c.node.Tag().String()
	value := c.node.Value()

	s := html.Sprintf(`
		<tr>
			<th nowrap="nowrap">%s</th>
			<td>%s</td>
		</tr>`, tag, value)

	for _, node := range c.node.Nodes() {
		s += newSourceProperty(c.document, node).String()
	}

	return s
}
