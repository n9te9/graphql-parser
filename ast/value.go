package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/n9te9/graphql-parser/token"
)

type IntValue struct {
	Token token.Token
	Value int64
}

func (i *IntValue) TokenLiteral() string { return i.Token.Literal }
func (i *IntValue) String() string       { return i.Token.Literal }
func (i *IntValue) valueNode()           {}

type FloatValue struct {
	Token token.Token
	Value float64
}

func (f *FloatValue) TokenLiteral() string { return f.Token.Literal }
func (f *FloatValue) String() string       { return f.Token.Literal }
func (f *FloatValue) valueNode()           {}

type StringValue struct {
	Token token.Token
	Value string
}

func (s *StringValue) TokenLiteral() string { return s.Token.Literal }
func (s *StringValue) String() string       { return fmt.Sprintf("%q", s.Value) } // 引用符付きで出力
func (s *StringValue) valueNode()           {}

type BooleanValue struct {
	Token token.Token
	Value bool
}

func (b *BooleanValue) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanValue) String() string       { return b.Token.Literal }
func (b *BooleanValue) valueNode()           {}

type NullValue struct {
	Token token.Token
}

func (n *NullValue) TokenLiteral() string { return n.Token.Literal }
func (n *NullValue) String() string       { return "null" }
func (n *NullValue) valueNode()           {}

type EnumValue struct {
	Token token.Token
	Value string
}

func (e *EnumValue) TokenLiteral() string { return e.Token.Literal }
func (e *EnumValue) String() string       { return e.Value }
func (e *EnumValue) valueNode()           {}

type ListValue struct {
	Token  token.Token
	Values []Value
}

func (l *ListValue) TokenLiteral() string { return l.Token.Literal }
func (l *ListValue) String() string {
	var out bytes.Buffer
	vals := []string{}
	for _, v := range l.Values {
		vals = append(vals, v.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(vals, ", "))
	out.WriteString("]")
	return out.String()
}
func (l *ListValue) valueNode() {}

type ObjectValue struct {
	Token  token.Token
	Fields []*ObjectField
}

func (o *ObjectValue) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectValue) String() string {
	var out bytes.Buffer
	fields := []string{}
	for _, f := range o.Fields {
		fields = append(fields, f.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString("}")
	return out.String()
}
func (o *ObjectValue) valueNode() {}

type ObjectField struct {
	Token token.Token
	Name  *Name
	Value Value
}

func (o *ObjectField) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectField) String() string {
	return o.Name.String() + ": " + o.Value.String()
}

type Variable struct {
	Token token.Token
	Name  string
}

func (v *Variable) TokenLiteral() string { return v.Token.Literal }
func (v *Variable) String() string       { return "$" + v.Name }
func (v *Variable) valueNode()           {}

type Name struct {
	Token token.Token
	Value string
}

func (n *Name) TokenLiteral() string { return n.Token.Literal }
func (n *Name) String() string       { return n.Value }
