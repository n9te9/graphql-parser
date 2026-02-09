package parser

import (
	"fmt"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func (p *Parser) parseObjectTypeDefinition(description string) ast.Definition {
	def := &ast.ObjectTypeDefinition{
		Description: description,
		Token:       p.curToken,
	}
	p.nextToken()

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	def.Name = name

	if p.curTokenIs(token.IMPLEMENTS) {
		def.Interfaces = p.parseImplementsInterfaces()
	}

	def.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		def.Fields = p.parseFieldDefinitions()
	}

	return def
}

func (p *Parser) parseImplementsInterfaces() []*ast.NamedType {
	var interfaces []*ast.NamedType
	p.nextToken()

	if p.curTokenIs(token.AMP) {
		p.nextToken()
	}

	for {

		name, err := p.parseName()
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}

		interfaces = append(interfaces, &ast.NamedType{
			Token: name.Token,
			Name:  name,
		})

		if !p.curTokenIs(token.AMP) {
			break
		}

		p.nextToken()
	}

	return interfaces
}

func (p *Parser) parseFieldDefinitions() []*ast.FieldDefinition {
	p.nextToken()

	var fields []*ast.FieldDefinition
	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		canParse := p.curTokenIs(token.IDENT) ||
			p.isKeywordToken() ||
			p.curTokenIs(token.STRING) ||
			p.curTokenIs(token.BLOCK_STRING)
		if canParse {
			fields = append(fields, p.parseFieldDefinition())
		} else {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.curToken.Literal))
			p.nextToken()
		}
	}

	if !p.curTokenIs(token.BRACE_R) {
		p.errors = append(p.errors, fmt.Sprintf("expected token } but, got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()

	return fields
}

func (p *Parser) parseFieldDefinition() *ast.FieldDefinition {
	description := p.parseDescription()
	def := &ast.FieldDefinition{
		Description: description,
		Token:       p.curToken,
	}

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	def.Name = name

	if p.curTokenIs(token.PAREN_L) {
		def.Arguments = p.parseInputValueDefinitions()
	}

	if !p.curTokenIs(token.COLON) {
		p.errors = append(p.errors, fmt.Sprintf("expected token : but, got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()

	def.Type = p.parseType()
	def.Directives = p.parseDirectives()

	return def
}

func (p *Parser) parseInputValueDefinitions() []*ast.InputValueDefinition {
	if !p.curTokenIs(token.PAREN_L) {
		p.errors = append(p.errors, fmt.Sprintf("expected token ( but, got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()

	var args []*ast.InputValueDefinition

	for !p.curTokenIs(token.PAREN_R) && !p.curTokenIs(token.EOF) {
		description := p.parseDescription()
		arg := &ast.InputValueDefinition{
			Description: description,
			Token:       p.curToken,
		}

		name, err := p.parseName()
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}
		arg.Name = name
		if !p.curTokenIs(token.COLON) {
			p.errors = append(p.errors, fmt.Sprintf("expected token : but, got %s", p.curToken.Literal))
			return nil
		}
		p.nextToken()

		arg.Type = p.parseType()
		if p.curTokenIs(token.EQUALS) {
			p.nextToken()
			arg.DefaultValue = p.parseValue()
		}
		arg.Directives = p.parseDirectives()

		args = append(args, arg)
	}

	if !p.curTokenIs(token.PAREN_R) {
		p.errors = append(p.errors, fmt.Sprintf("expected token ) but, got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()

	return args
}

func (p *Parser) parseInterfaceTypeDefinition(description string) ast.Definition {
	def := &ast.InterfaceTypeDefinition{
		Token:       p.curToken,
		Description: description,
	}
	p.nextToken()

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	def.Name = name
	if p.curTokenIs(token.IMPLEMENTS) {
		def.Interfaces = p.parseImplementsInterfaces()
	}

	def.Directives = p.parseDirectives()
	if p.curTokenIs(token.BRACE_L) {
		def.Fields = p.parseFieldDefinitions()
	}

	return def
}

func (p *Parser) parseDescription() string {
	if p.curTokenIs(token.STRING) {
		val := unquoteGeneric(p.curToken.Literal)
		p.nextToken()
		return val
	}
	if p.curTokenIs(token.BLOCK_STRING) {
		val := dedentBlockStringValue(p.curToken.Literal)
		p.nextToken()
		return val
	}
	return ""
}

func (p *Parser) parseUnionTypeDefinition(description string) ast.Definition {
	def := &ast.UnionTypeDefinition{
		Token:       p.curToken,
		Description: description,
	}
	p.nextToken()

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	def.Name = name
	def.Directives = p.parseDirectives()

	if p.curTokenIs(token.EQUALS) {
		p.nextToken()
		def.Types = p.parseUnionMemberTypes()
	}

	return def
}

func (p *Parser) parseUnionMemberTypes() []*ast.NamedType {
	if p.curTokenIs(token.PIPE) {
		p.nextToken()
	}

	var types []*ast.NamedType
	for {
		name, err := p.parseName()
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("expected type name in union, got %s", p.curToken.Literal))
			p.nextToken()
			return nil
		}

		types = append(types, &ast.NamedType{
			Token: name.Token,
			Name:  name,
		})

		if !p.curTokenIs(token.PIPE) {
			break
		}
		p.nextToken()
	}
	return types
}

func (p *Parser) parseEnumTypeDefinition(description string) ast.Definition {
	def := &ast.EnumTypeDefinition{
		Token:       p.curToken,
		Description: description,
	}
	p.nextToken()

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	def.Name = name
	def.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		def.Values = p.parseEnumValueDefinitions()
	}

	return def
}

func (p *Parser) parseEnumValueDefinitions() []*ast.EnumValueDefinition {
	p.nextToken()

	var values []*ast.EnumValueDefinition
	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		canEnumValidToken := p.curTokenIs(token.IDENT) ||
			p.isKeywordToken() ||
			p.curTokenIs(token.STRING) ||
			p.curTokenIs(token.BLOCK_STRING)
		if canEnumValidToken {
			values = append(values, p.parseEnumValueDefinition())
		} else {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.curToken.Literal))
			p.nextToken()
		}
	}

	if !p.curTokenIs(token.BRACE_R) {
		p.errors = append(p.errors, fmt.Sprintf("expected token }, but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken()

	return values
}

func (p *Parser) parseEnumValueDefinition() *ast.EnumValueDefinition {
	description := p.parseDescription()
	def := &ast.EnumValueDefinition{
		Token:       p.curToken,
		Description: description,
	}

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	def.Name = name
	def.Directives = p.parseDirectives()

	return def
}

func (p *Parser) parseInputObjectTypeDefinition(description string) ast.Definition {
	def := &ast.InputObjectTypeDefinition{
		Description: description,
		Token:       p.curToken,
	}
	p.nextToken()

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	def.Name = name
	def.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		def.Fields = p.parseInputObjectFields()
	}

	return def
}

func (p *Parser) parseInputObjectFields() []*ast.InputValueDefinition {
	p.nextToken() // skip '{'

	var fields []*ast.InputValueDefinition
	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		validInputObjectToken := p.curTokenIs(token.IDENT) ||
			p.isKeywordToken() ||
			p.curTokenIs(token.STRING) ||
			p.curTokenIs(token.BLOCK_STRING)
		if validInputObjectToken {
			fields = append(fields, p.parseInputObjectField())
		} else {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token in input object: %s", p.curToken.Literal))
			p.nextToken()
		}
	}

	if !p.curTokenIs(token.BRACE_R) {
		p.errors = append(p.errors, fmt.Sprintf("expected token } but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip '}'

	return fields
}

func (p *Parser) parseInputObjectField() *ast.InputValueDefinition {
	description := p.parseDescription()

	def := &ast.InputValueDefinition{
		Description: description,
		Token:       p.curToken,
	}

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	def.Name = name

	if !p.curTokenIs(token.COLON) {
		p.errors = append(p.errors, fmt.Sprintf("expected token :, but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip ':'

	def.Type = p.parseType()

	if p.curTokenIs(token.EQUALS) {
		p.nextToken()
		def.DefaultValue = p.parseValue()
	}

	def.Directives = p.parseDirectives()

	return def
}

func (p *Parser) parseScalarTypeDefinition(description string) ast.Definition {
	def := &ast.ScalarTypeDefinition{
		Description: description,
		Token:       p.curToken,
	}
	p.nextToken() // skip 'scalar'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	def.Name = name
	def.Directives = p.parseDirectives()

	return def
}

func (p *Parser) parseDirectiveDefinition(description string) ast.Definition {
	def := &ast.DirectiveDefinition{
		Description: description,
		Token:       p.curToken,
	}
	p.nextToken() // skip 'directive'

	if !p.curTokenIs(token.AT) {
		p.errors = append(p.errors, fmt.Sprintf("expected token @, but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip '@'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	def.Name = name

	if p.curTokenIs(token.PAREN_L) {
		def.Arguments = p.parseInputValueDefinitions()
	}

	if p.curTokenIs(token.IDENT) && p.curToken.Literal == "repeatable" {
		def.Repeatable = true
		p.nextToken()
	}

	if !p.curTokenIs(token.ON) {
		p.errors = append(p.errors, fmt.Sprintf("expected token on, but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip 'on'

	// Locations (FIELD | OBJECT | ...)
	def.Locations = p.parseDirectiveLocations()

	return def
}

func (p *Parser) parseDirectiveLocations() []*ast.Name {
	var locations []*ast.Name

	if p.curTokenIs(token.PIPE) {
		p.nextToken()
	}

	for {
		name, err := p.parseName()
		if err != nil {
			p.errors = append(p.errors, err.Error())
		}
		locations = append(locations, name)

		if !p.curTokenIs(token.PIPE) {
			break
		}
		p.nextToken() // skip '|'
	}

	return locations
}

func (p *Parser) parseSchemaDefinition(description string) ast.Definition {
	def := &ast.SchemaDefinition{
		Token:       p.curToken,
		Description: description,
	}
	p.nextToken() // skip 'schema'

	def.Directives = p.parseDirectives()

	if !p.curTokenIs(token.BRACE_L) {
		p.errors = append(p.errors, fmt.Sprintf("expected token { but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip '{'

	// Root Operation Types (query: Query, etc.)
	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) || p.isKeywordToken() {
			def.OperationTypes = append(def.OperationTypes, p.parseOperationTypeDefinition())
		}
	}

	if !p.curTokenIs(token.BRACE_R) {
		p.errors = append(p.errors, fmt.Sprintf("expected token } but got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // skip '}'

	return def
}

func (p *Parser) parseOperationTypeDefinition() *ast.OperationTypeDefinition {
	def := &ast.OperationTypeDefinition{}

	if p.curTokenIs(token.QUERY) || p.curTokenIs(token.MUTATION) || p.curTokenIs(token.SUBSCRIPTION) {
		def.Operation = p.curToken.Type
		p.nextToken()
	} else if p.curTokenIs(token.IDENT) {
		switch p.curToken.Literal {
		case "query":
			def.Operation = token.QUERY
		case "mutation":
			def.Operation = token.MUTATION
		case "subscription":
			def.Operation = token.SUBSCRIPTION
		default:
			p.errors = append(p.errors, fmt.Sprintf("unexpected operation type: %s", p.curToken.Literal))
			return nil
		}
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected operation type, got %s", p.curToken.Literal))
		return nil
	}

	if !p.curTokenIs(token.COLON) {
		return nil
	}
	p.nextToken()

	t, ok := p.parseType().(*ast.NamedType)
	if !ok {
		p.errors = append(p.errors, "failed type assertion *ast.NamedType")
		return nil
	}

	def.Type = t

	return def
}

func (p *Parser) parseExtendDefinition(description string) ast.Definition {
	if description != "" {
		p.errors = append(p.errors, "unexpected description before 'extend': extensions cannot have descriptions")
		return nil
	}

	t := p.curToken
	p.nextToken()

	switch p.curToken.Type {
	case token.SCHEMA:
		return p.parseSchemaExtension(t)
	case token.SCALAR:
		return p.parseScalarTypeExtension(t)
	case token.TYPE:
		return p.parseObjectTypeExtension(t)
	case token.INPUT:
		return p.parseInputObjectTypeExtension(t)
	case token.INTERFACE:
		return p.parseInterfaceTypeExtension(t)
	case token.UNION:
		return p.parseUnionTypeExtension(t)
	case token.ENUM:
		return p.parseEnumTypeExtension(t)
	}

	p.errors = append(p.errors, fmt.Sprintf("unexpected token after extend: %s", p.curToken.Literal))
	return nil
}

func (p *Parser) parseSchemaExtension(startToken token.Token) ast.Definition {
	ext := &ast.SchemaExtension{Token: startToken}
	p.nextToken() // skip 'schema'

	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		p.nextToken() // skip '{'
		for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
			if p.curTokenIs(token.IDENT) || p.isKeywordToken() {
				ext.OperationTypes = append(ext.OperationTypes, p.parseOperationTypeDefinition())
			} else {
				p.nextToken()
			}
		}
		if !p.curTokenIs(token.BRACE_R) {
			return nil
		}

		p.nextToken() // skip '}'
	}

	if len(ext.Directives) == 0 && len(ext.OperationTypes) == 0 {
		p.errors = append(p.errors, "unexpected token: extend schema must have directives or operation types")
		return nil
	}

	return ext
}

func (p *Parser) parseScalarTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.ScalarTypeExtension{Token: startToken}
	p.nextToken() // skip 'scalar'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name
	ext.Directives = p.parseDirectives()

	if len(ext.Directives) == 0 {
		p.errors = append(p.errors, "unexpected token: extend scalar must have directives")
		return nil
	}

	return ext
}

func (p *Parser) parseObjectTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.ObjectTypeExtension{Token: startToken}
	p.nextToken() // skip 'type'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name

	if p.curTokenIs(token.IMPLEMENTS) {
		ext.Interfaces = p.parseImplementsInterfaces()
	}

	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		ext.Fields = p.parseFieldDefinitions()
	}

	if len(ext.Interfaces) == 0 && len(ext.Directives) == 0 && len(ext.Fields) == 0 {
		p.errors = append(p.errors, "unexpected token: extend type must have implements, directives or fields")
		return nil
	}

	return ext
}

func (p *Parser) parseInterfaceTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.InterfaceTypeExtension{Token: startToken}
	p.nextToken() // skip 'interface'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name

	if p.curTokenIs(token.IMPLEMENTS) {
		ext.Interfaces = p.parseImplementsInterfaces()
	}

	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		ext.Fields = p.parseFieldDefinitions()
	}

	if len(ext.Interfaces) == 0 && len(ext.Directives) == 0 && len(ext.Fields) == 0 {
		p.errors = append(p.errors, "unexpected token: extend interface must have implements, directives or fields")
		return nil
	}

	return ext
}

// -------------------------------------------------------------------------
// 5. Union Extension
// extend union Name @directive = ...
// -------------------------------------------------------------------------
func (p *Parser) parseUnionTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.UnionTypeExtension{Token: startToken}
	p.nextToken() // skip 'union'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name
	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.EQUALS) {
		p.nextToken()
		ext.Types = p.parseUnionMemberTypes()
	}

	if len(ext.Directives) == 0 && len(ext.Types) == 0 {
		p.errors = append(p.errors, "unexpected token: extend union must have directives or member types")
		return nil
	}

	return ext
}

func (p *Parser) parseEnumTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.EnumTypeExtension{Token: startToken}
	p.nextToken() // skip 'enum'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name
	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		ext.Values = p.parseEnumValueDefinitions()
	}

	if len(ext.Directives) == 0 && len(ext.Values) == 0 {
		p.errors = append(p.errors, "unexpected token: extend enum must have directives or values")
		return nil
	}

	return ext
}

func (p *Parser) parseInputObjectTypeExtension(startToken token.Token) ast.Definition {
	ext := &ast.InputObjectTypeExtension{Token: startToken}
	p.nextToken() // skip 'input'

	name, err := p.parseName()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}

	ext.Name = name
	ext.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		ext.Fields = p.parseInputObjectFields()
	}

	if len(ext.Directives) == 0 && len(ext.Fields) == 0 {
		p.errors = append(p.errors, "unexpected token: extend input must have directives or fields")
		return nil
	}

	return ext
}
