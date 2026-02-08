package ast

import (
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Definition interface {
	Node
	definitionNode()
}

type Value interface {
	Node
	valueNode()
}

type Type interface {
	Node
	typeNode()
}

type Selection interface {
	Node
	selectionNode()
}

type Document struct {
	Definitions []Definition
}

func (d *Document) TokenLiteral() string {
	if len(d.Definitions) > 0 {
		return d.Definitions[0].TokenLiteral()
	}
	return ""
}

func (d *Document) String() string {
	var out bytes.Buffer
	for _, s := range d.Definitions {
		out.WriteString(s.String())
	}
	return out.String()
}
