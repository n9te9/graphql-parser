package ast

import (
	"bytes"
	"strings"

	"github.com/n9te9/graphql-parser/token"
)

type Directive struct {
	Token     token.Token
	Name      string
	Arguments []*Argument
}

func (d *Directive) TokenLiteral() string { return d.Token.Literal }
func (d *Directive) String() string {
	var out bytes.Buffer
	out.WriteString("@")
	out.WriteString(d.Name)

	if len(d.Arguments) > 0 {
		out.WriteString("(")
		args := []string{}
		for _, a := range d.Arguments {
			args = append(args, a.String())
		}
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(")")
	}
	return out.String()
}

type Argument struct {
	Token token.Token
	Name  *Name
	Value Value
}

func (a *Argument) TokenLiteral() string { return a.Token.Literal }
func (a *Argument) String() string {
	return a.Name.String() + ": " + a.Value.String()
}
