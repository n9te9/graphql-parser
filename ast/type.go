package ast

import (
	"bytes"

	"github.com/n9te9/graphql-parser/token"
)

type NamedType struct {
	Token token.Token
	Name  *Name
}

func (n *NamedType) TokenLiteral() string { return n.Token.Literal }
func (n *NamedType) String() string       { return n.Name.String() }
func (n *NamedType) typeNode()            {}

type ListType struct {
	Token token.Token
	Type  Type
}

func (l *ListType) TokenLiteral() string { return l.Token.Literal }
func (l *ListType) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	out.WriteString(l.Type.String())
	out.WriteString("]")
	return out.String()
}
func (l *ListType) typeNode() {}

type NonNullType struct {
	Token token.Token
	Type  Type
}

func (n *NonNullType) TokenLiteral() string { return n.Token.Literal }
func (n *NonNullType) String() string {
	return n.Type.String() + "!"
}
func (n *NonNullType) typeNode() {}
