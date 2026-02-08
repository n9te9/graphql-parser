package ast

import (
	"bytes"
	"strings"

	"github.com/n9te9/graphql-parser/token"
)

type OperationType string

const (
	Query        OperationType = "query"
	Mutation     OperationType = "mutation"
	Subscription OperationType = "subscription"
)

type OperationDefinition struct {
	Token               token.Token
	Operation           OperationType
	Name                *Name
	VariableDefinitions []*VariableDefinition
	Directives          []*Directive
	SelectionSet        []Selection
}

func (o *OperationDefinition) TokenLiteral() string { return o.Token.Literal }
func (o *OperationDefinition) String() string {
	var out bytes.Buffer

	if o.Operation != "" {
		out.WriteString(string(o.Operation))
		out.WriteString(" ")
	}

	if o.Name != nil {
		out.WriteString(o.Name.String())
	}

	if len(o.VariableDefinitions) > 0 {
		out.WriteString("(")
		defs := []string{}
		for _, d := range o.VariableDefinitions {
			defs = append(defs, d.String())
		}
		out.WriteString(strings.Join(defs, ", "))
		out.WriteString(")")
	}

	if len(o.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range o.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	if len(o.SelectionSet) > 0 {
		out.WriteString(" ")
		out.WriteString(selectionSetToString(o.SelectionSet))
	}

	return out.String()
}
func (o *OperationDefinition) definitionNode() {}

type VariableDefinition struct {
	Token        token.Token
	Variable     *Variable
	Type         Type
	DefaultValue Value
	Directives   []*Directive
}

func (v *VariableDefinition) TokenLiteral() string { return v.Token.Literal }
func (v *VariableDefinition) String() string {
	var out bytes.Buffer
	out.WriteString(v.Variable.String())
	out.WriteString(": ")
	out.WriteString(v.Type.String())

	if v.DefaultValue != nil {
		out.WriteString(" = ")
		out.WriteString(v.DefaultValue.String())
	}

	if len(v.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range v.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	return out.String()
}

type Field struct {
	Token        token.Token
	Alias        *Name
	Name         *Name
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet []Selection
}

func (f *Field) TokenLiteral() string { return f.Token.Literal }
func (f *Field) String() string {
	var out bytes.Buffer

	if f.Alias != nil {
		out.WriteString(f.Alias.String() + ": ")
	}
	out.WriteString(f.Name.String())

	if len(f.Arguments) > 0 {
		out.WriteString("(")
		args := []string{}
		for _, a := range f.Arguments {
			args = append(args, a.String())
		}
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(")")
	}

	if len(f.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range f.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	if len(f.SelectionSet) > 0 {
		out.WriteString(" ")
		out.WriteString(selectionSetToString(f.SelectionSet))
	}

	return out.String()
}
func (f *Field) selectionNode() {}

type FragmentSpread struct {
	Token      token.Token
	Name       *Name
	Directives []*Directive
}

func (fs *FragmentSpread) TokenLiteral() string { return fs.Token.Literal }
func (fs *FragmentSpread) String() string {
	var out bytes.Buffer
	out.WriteString("...")
	out.WriteString(fs.Name.String())

	if len(fs.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range fs.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	return out.String()
}
func (fs *FragmentSpread) selectionNode() {}

type InlineFragment struct {
	Token         token.Token
	TypeCondition *NamedType
	Directives    []*Directive
	SelectionSet  []Selection
}

func (i *InlineFragment) TokenLiteral() string { return i.Token.Literal }
func (i *InlineFragment) String() string {
	var out bytes.Buffer
	out.WriteString("...")

	if i.TypeCondition != nil {
		out.WriteString(" on ")
		out.WriteString(i.TypeCondition.String())
	}

	if len(i.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range i.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	out.WriteString(" ")
	out.WriteString(selectionSetToString(i.SelectionSet))
	return out.String()
}
func (i *InlineFragment) selectionNode() {}

func selectionSetToString(selections []Selection) string {
	var out bytes.Buffer
	out.WriteString("{ ")
	sels := []string{}
	for _, s := range selections {
		sels = append(sels, s.String())
	}
	out.WriteString(strings.Join(sels, " "))
	out.WriteString(" }")
	return out.String()
}

type FragmentDefinition struct {
	Token         token.Token
	Name          *Name
	TypeCondition *NamedType
	Directives    []*Directive
	SelectionSet  []Selection
}

func (fd *FragmentDefinition) TokenLiteral() string { return fd.Token.Literal }
func (fd *FragmentDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("fragment ")
	out.WriteString(fd.Name.String())
	out.WriteString(" on ")
	out.WriteString(fd.TypeCondition.String())

	if len(fd.Directives) > 0 {
		out.WriteString(" ")
		dirs := []string{}
		for _, d := range fd.Directives {
			dirs = append(dirs, d.String())
		}
		out.WriteString(strings.Join(dirs, " "))
	}

	out.WriteString(" ")
	out.WriteString(selectionSetToString(fd.SelectionSet))
	return out.String()
}

func (fd *FragmentDefinition) definitionNode() {}
