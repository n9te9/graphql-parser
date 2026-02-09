package ast

import (
	"bytes"

	"github.com/n9te9/graphql-parser/token"
)

type ObjectTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Interfaces  []*NamedType
	Directives  []*Directive
	Fields      []*FieldDefinition
}

func (def *ObjectTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *ObjectTypeDefinition) String() string {
	var out bytes.Buffer

	out.WriteString("type ")
	out.WriteString(def.Name.String())

	if len(def.Interfaces) > 0 {
		out.WriteString(" implements ")
		for i, iface := range def.Interfaces {
			if i > 0 {
				out.WriteString(" & ")
			}
			out.WriteString(iface.String())
		}
	}

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	if len(def.Fields) > 0 {
		out.WriteString(" {")
		for _, f := range def.Fields {
			out.WriteString(" ")
			out.WriteString(f.String())
		}
		out.WriteString(" }")
	}

	return out.String()
}

func (def *ObjectTypeDefinition) definitionNode() {}

type FieldDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Arguments   []*InputValueDefinition
	Type        Type
	Directives  []*Directive
}

func (fd *FieldDefinition) TokenLiteral() string { return fd.Token.Literal }
func (fd *FieldDefinition) String() string {
	var out bytes.Buffer

	out.WriteString(fd.Name.String())

	if len(fd.Arguments) > 0 {
		out.WriteString("(")
		for i, arg := range fd.Arguments {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(arg.String())
		}
		out.WriteString(")")
	}

	out.WriteString(": ")
	out.WriteString(fd.Type.String())

	if len(fd.Directives) > 0 {
		for _, d := range fd.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	return out.String()
}

type InputValueDefinition struct {
	Description  string
	Token        token.Token
	Name         *Name
	Type         Type
	DefaultValue Value
	Directives   []*Directive
}

func (ivd *InputValueDefinition) TokenLiteral() string { return ivd.Token.Literal }
func (ivd *InputValueDefinition) String() string {
	var out bytes.Buffer

	out.WriteString(ivd.Name.String())
	out.WriteString(": ")
	out.WriteString(ivd.Type.String())

	if ivd.DefaultValue != nil {
		out.WriteString(" = ")
		out.WriteString(ivd.DefaultValue.String())
	}

	if len(ivd.Directives) > 0 {
		for _, d := range ivd.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	return out.String()
}

type InterfaceTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Interfaces  []*NamedType
	Directives  []*Directive
	Fields      []*FieldDefinition
}

func (def *InterfaceTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *InterfaceTypeDefinition) String() string {
	var out bytes.Buffer

	out.WriteString("interface ")
	out.WriteString(def.Name.String())

	if len(def.Interfaces) > 0 {
		out.WriteString(" implements ")
		for i, iface := range def.Interfaces {
			if i > 0 {
				out.WriteString(" & ")
			}
			out.WriteString(iface.String())
		}
	}

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	if len(def.Fields) > 0 {
		out.WriteString(" {")
		for _, f := range def.Fields {
			out.WriteString(" ")
			out.WriteString(f.String())
		}
		out.WriteString(" }")
	}

	return out.String()
}

func (def *InterfaceTypeDefinition) definitionNode() {}

type UnionTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Directives  []*Directive
	Types       []*NamedType
}

func (def *UnionTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *UnionTypeDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("union ")
	out.WriteString(def.Name.String())

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	if len(def.Types) > 0 {
		out.WriteString(" = ")
		for i, t := range def.Types {
			if i > 0 {
				out.WriteString(" | ")
			}
			out.WriteString(t.String())
		}
	}
	return out.String()
}

func (def *UnionTypeDefinition) definitionNode() {}

type EnumTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Directives  []*Directive
	Values      []*EnumValueDefinition
}

func (def *EnumTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *EnumTypeDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("enum ")
	out.WriteString(def.Name.String())

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	if len(def.Values) > 0 {
		out.WriteString(" {")
		for _, v := range def.Values {
			out.WriteString(" ")
			out.WriteString(v.String())
		}
		out.WriteString(" }")
	}
	return out.String()
}

type EnumValueDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Directives  []*Directive
}

func (v *EnumValueDefinition) TokenLiteral() string { return v.Token.Literal }
func (v *EnumValueDefinition) String() string {
	var out bytes.Buffer
	out.WriteString(v.Name.String())
	if len(v.Directives) > 0 {
		for _, d := range v.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}
	return out.String()
}

func (def *EnumTypeDefinition) definitionNode() {}

type ScalarTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Directives  []*Directive
}

func (def *ScalarTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *ScalarTypeDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("scalar ")
	out.WriteString(def.Name.String())
	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}
	return out.String()
}

func (def *ScalarTypeDefinition) definitionNode() {}

type InputObjectTypeDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Directives  []*Directive
	Fields      []*InputValueDefinition
}

func (def *InputObjectTypeDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *InputObjectTypeDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("input ")
	out.WriteString(def.Name.String())

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	if len(def.Fields) > 0 {
		out.WriteString(" {")
		for _, f := range def.Fields {
			out.WriteString(" ")
			out.WriteString(f.String())
		}
		out.WriteString(" }")
	}
	return out.String()
}

func (def *InputObjectTypeDefinition) definitionNode() {}

type DirectiveDefinition struct {
	Description string
	Token       token.Token
	Name        *Name
	Arguments   []*InputValueDefinition
	Repeatable  bool
	Locations   []*Name
}

func (def *DirectiveDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *DirectiveDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("directive @")
	out.WriteString(def.Name.String())

	if len(def.Arguments) > 0 {
		out.WriteString("(")
		for i, arg := range def.Arguments {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(arg.String())
		}
		out.WriteString(")")
	}

	if def.Repeatable {
		out.WriteString(" repeatable")
	}

	out.WriteString(" on ")
	for i, loc := range def.Locations {
		if i > 0 {
			out.WriteString(" | ")
		}
		out.WriteString(loc.String())
	}
	return out.String()
}

func (def *DirectiveDefinition) definitionNode() {}

type SchemaDefinition struct {
	Token          token.Token // 'schema'
	Description    string
	Directives     []*Directive
	OperationTypes []*OperationTypeDefinition
}

func (def *SchemaDefinition) TokenLiteral() string { return def.Token.Literal }
func (def *SchemaDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("schema")

	if len(def.Directives) > 0 {
		for _, d := range def.Directives {
			out.WriteString(" ")
			out.WriteString(d.String())
		}
	}

	out.WriteString(" {")
	for _, op := range def.OperationTypes {
		out.WriteString(" ")
		out.WriteString(op.String())
	}
	out.WriteString(" }")
	return out.String()
}

func (def *SchemaDefinition) definitionNode() {}

type OperationTypeDefinition struct {
	Operation token.TokenType
	Type      *NamedType
}

func (op *OperationTypeDefinition) String() string {
	return op.Operation.String() + ": " + op.Type.String()
}

type SchemaExtension struct {
	Token          token.Token
	Directives     []*Directive
	OperationTypes []*OperationTypeDefinition
}

func (e *SchemaExtension) TokenLiteral() string { return e.Token.Literal }
func (e *SchemaExtension) String() string {
	return ""
}

func (e *SchemaExtension) definitionNode() {}

type ScalarTypeExtension struct {
	Token      token.Token
	Name       *Name
	Directives []*Directive
}

func (e *ScalarTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *ScalarTypeExtension) String() string {
	return ""
}

func (e *ScalarTypeExtension) definitionNode() {}

type ObjectTypeExtension struct {
	Token      token.Token
	Name       *Name
	Interfaces []*NamedType
	Directives []*Directive
	Fields     []*FieldDefinition
}

func (e *ObjectTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *ObjectTypeExtension) String() string {
	return "extend type " + e.Name.String() + " ..."
}

func (e *ObjectTypeExtension) definitionNode() {}

type InterfaceTypeExtension struct {
	Token      token.Token
	Name       *Name
	Interfaces []*NamedType
	Directives []*Directive
	Fields     []*FieldDefinition
}

func (e *InterfaceTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *InterfaceTypeExtension) String() string {
	return "extend interface " + e.Name.String() + " ..."
}

func (e *InterfaceTypeExtension) definitionNode() {}

type UnionTypeExtension struct {
	Token      token.Token
	Name       *Name
	Directives []*Directive
	Types      []*NamedType
}

func (e *UnionTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *UnionTypeExtension) String() string {
	return "extend union " + e.Name.String() + " ..."
}

func (e *UnionTypeExtension) definitionNode() {}

type EnumTypeExtension struct {
	Token      token.Token
	Name       *Name
	Directives []*Directive
	Values     []*EnumValueDefinition
}

func (e *EnumTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *EnumTypeExtension) String() string {
	return "extend enum " + e.Name.String() + " ..."
}

func (e *EnumTypeExtension) definitionNode() {}

type InputObjectTypeExtension struct {
	Token      token.Token
	Name       *Name
	Directives []*Directive
	Fields     []*InputValueDefinition
}

func (e *InputObjectTypeExtension) TokenLiteral() string { return e.Token.Literal }
func (e *InputObjectTypeExtension) String() string {
	return "extend input " + e.Name.String() + " ..."
}

func (e *InputObjectTypeExtension) definitionNode() {}
