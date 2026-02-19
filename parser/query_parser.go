package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func (p *Parser) parseOperationDefinition() ast.Definition {
	if p.curTokenIs(token.BRACE_L) {
		return &ast.OperationDefinition{
			Operation:    ast.Query,
			SelectionSet: p.parseSelectionSet(),
		}
	}

	if !p.isOperationType() {
		return nil
	}

	stmt := &ast.OperationDefinition{
		Token:     p.curToken,
		Operation: ast.OperationType(p.curToken.Literal),
	}
	p.nextToken()

	if p.curTokenIs(token.IDENT) {
		stmt.Name = &ast.Name{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		p.nextToken()
	}

	stmt.VariableDefinitions = p.parseVariableDefinitions()
	stmt.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		stmt.SelectionSet = p.parseSelectionSet()
	} else {
		p.peekError(token.BRACE_L)
		return nil
	}

	return stmt
}

func (p *Parser) parseFragmentDefinition() ast.Definition {
	stmt := &ast.FragmentDefinition{Token: p.curToken}
	p.nextToken()

	if !p.isKeywordToken() && !p.curTokenIs(token.IDENT) {
		p.errors = append(p.errors, "expected fragment name")
		return nil
	}

	stmt.Name = &ast.Name{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ON) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.TypeCondition = &ast.NamedType{
		Token: p.curToken,
		Name:  &ast.Name{Token: p.curToken, Value: p.curToken.Literal},
	}

	p.nextToken()
	stmt.Directives = p.parseDirectives()

	if !p.curTokenIs(token.BRACE_L) {
		p.peekError(token.BRACE_L)
		return nil
	}

	stmt.SelectionSet = p.parseSelectionSet()

	return stmt
}

func (p *Parser) parseSelectionSet() []ast.Selection {
	selections := []ast.Selection{}

	if p.curTokenIs(token.BRACE_L) {
		p.nextToken()
	}

	if p.curTokenIs(token.BRACE_R) {
		p.errors = append(p.errors, "empty selection set")
		return nil
	}

	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		selection := p.parseSelection()
		if selection != nil {
			selections = append(selections, selection)
		}
	}

	if p.curTokenIs(token.BRACE_R) {
		p.nextToken()
	} else {
		p.peekError(token.BRACE_R)
		return nil
	}

	return selections
}

func (p *Parser) parseSelection() ast.Selection {
	if p.curTokenIs(token.IDENT) || p.isKeywordToken() {
		return p.parseField()
	}

	if p.curTokenIs(token.SPREAD) {
		return p.parseFragment()
	}

	p.nextToken()
	return nil
}

func (p *Parser) isKeywordToken() bool {
	return p.curTokenIs(token.QUERY) ||
		p.curTokenIs(token.MUTATION) ||
		p.curTokenIs(token.SUBSCRIPTION) ||
		p.curTokenIs(token.FRAGMENT) ||
		p.curTokenIs(token.ON) ||
		p.curTokenIs(token.TYPE) ||
		p.curTokenIs(token.INPUT) ||
		p.curTokenIs(token.ENUM) ||
		p.curTokenIs(token.UNION) ||
		p.curTokenIs(token.INTERFACE) ||
		p.curTokenIs(token.SCALAR) ||
		p.curTokenIs(token.DIRECTIVE)
}

func (p *Parser) parseFragment() ast.Selection {
	startToken := p.curToken
	p.nextToken()

	if p.curTokenIs(token.ON) {
		p.nextToken()

		typeCondition := &ast.NamedType{
			Token: p.curToken,
			Name:  &ast.Name{Token: p.curToken, Value: p.curToken.Literal},
		}
		p.nextToken()

		dirs := p.parseDirectives()

		return &ast.InlineFragment{
			Token:         startToken,
			TypeCondition: typeCondition,
			Directives:    dirs,
			SelectionSet:  p.parseSelectionSet(),
		}
	}

	if p.curTokenIs(token.AT) || p.curTokenIs(token.BRACE_L) {
		dirs := p.parseDirectives()

		return &ast.InlineFragment{
			Token:        startToken,
			Directives:   dirs,
			SelectionSet: p.parseSelectionSet(),
		}
	}

	spread := &ast.FragmentSpread{
		Token: startToken,
		Name:  &ast.Name{Token: p.curToken, Value: p.curToken.Literal},
	}
	p.nextToken()
	spread.Directives = p.parseDirectives()

	return spread
}

func (p *Parser) parseField() *ast.Field {
	field := &ast.Field{Token: p.curToken}

	if p.peekTokenIs(token.COLON) {
		field.Alias = &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		p.nextToken()
	}

	field.Name = &ast.Name{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()

	if p.curTokenIs(token.PAREN_L) {
		field.Arguments = p.parseArguments()
	}
	field.Directives = p.parseDirectives()

	if p.curTokenIs(token.BRACE_L) {
		field.SelectionSet = p.parseSelectionSet()
	}

	return field
}

func (p *Parser) isOperationType() bool {
	return p.curTokenIs(token.QUERY) || p.curTokenIs(token.MUTATION) || p.curTokenIs(token.SUBSCRIPTION)
}

// parseValue parses a value literal or variable
func (p *Parser) parseValue() ast.Value {
	switch p.curToken.Type {
	case token.INT:
		val, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
		lit := &ast.IntValue{Token: p.curToken, Value: val}
		p.nextToken()
		return lit
	case token.FLOAT:
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		lit := &ast.FloatValue{Token: p.curToken, Value: val}
		p.nextToken()
		return lit
	case token.STRING:
		lit := &ast.StringValue{Token: p.curToken, Value: unquoteGeneric(p.curToken.Literal)}
		p.nextToken()
		return lit
	case token.TRUE, token.FALSE:
		val, _ := strconv.ParseBool(p.curToken.Literal)
		lit := &ast.BooleanValue{Token: p.curToken, Value: val}
		p.nextToken()
		return lit
	case token.NULL:
		lit := &ast.NullValue{Token: p.curToken}
		p.nextToken()
		return lit
	case token.IDENT:
		// Enum Value (e.g. ACTIVE)
		lit := &ast.EnumValue{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		return lit
	case token.DOLLAR:
		v := p.parseVariable()
		p.nextToken()
		return v
	case token.BRACKET_L:
		v := p.parseListValue()
		return v
	case token.BRACE_L:
		v := p.parseObjectValue()
		return v
	case token.BLOCK_STRING:
		lit := &ast.StringValue{Token: p.curToken, Value: dedentBlockStringValue(p.curToken.Literal)}
		p.nextToken()
		return lit
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token for value: %s", p.curToken.Type))
		return nil
	}
}

func unquoteGeneric(s string) string {
	quoted := `"` + s + `"`

	if unquoted, err := strconv.Unquote(quoted); err == nil {
		return unquoted
	}
	return s
}

func dedentBlockStringValue(rawString string) string {
	lines := strings.Split(rawString, "\n")

	commonIndent := math.MaxInt
	foundCommonIndent := false

	for i, line := range lines {
		if i == 0 {
			continue
		}

		indent := countLeadingWhitespace(line)

		if indent < len(line) {
			if indent < commonIndent {
				commonIndent = indent
				foundCommonIndent = true
			}
		}
	}

	if !foundCommonIndent {
		commonIndent = 0
	}

	if commonIndent > 0 {
		for i, line := range lines {
			if i == 0 {
				continue
			}
			if len(line) >= commonIndent {
				lines[i] = line[commonIndent:]
			}
		}
	}

	for len(lines) > 0 {
		if isBlank(lines[0]) {
			lines = lines[1:]
		} else {
			break
		}
	}

	for len(lines) > 0 {
		if isBlank(lines[len(lines)-1]) {
			lines = lines[:len(lines)-1]
		} else {
			break
		}
	}

	return strings.Join(lines, "\n")
}

func countLeadingWhitespace(s string) int {
	count := 0
	for _, r := range s {
		if r == ' ' || r == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}

func isBlank(s string) bool {
	return countLeadingWhitespace(s) == len(s)
}

// parseVariable parses $name
func (p *Parser) parseVariable() *ast.Variable {
	v := &ast.Variable{Token: p.curToken}
	if !p.peekTokenIs(token.IDENT) && !p.peekTokenIs(token.INPUT) {
		return nil
	}
	p.nextToken()
	v.Name = p.curToken.Literal

	return v
}

// parseListValue parses [val1, val2]
func (p *Parser) parseListValue() *ast.ListValue {
	lit := &ast.ListValue{Token: p.curToken}
	p.nextToken() // skip '['

	for !p.curTokenIs(token.BRACKET_R) && !p.curTokenIs(token.EOF) {
		val := p.parseValue()
		if val != nil {
			lit.Values = append(lit.Values, val)
		}
	}

	if !p.curTokenIs(token.BRACKET_R) {
		return nil
	}
	p.nextToken()

	return lit
}

// parseObjectValue parses { key: value }
func (p *Parser) parseObjectValue() *ast.ObjectValue {
	lit := &ast.ObjectValue{Token: p.curToken}
	p.nextToken() // skip '{'

	for !p.curTokenIs(token.BRACE_R) && !p.curTokenIs(token.EOF) {
		field := p.parseObjectField()
		if field != nil {
			lit.Fields = append(lit.Fields, field)
		}
	}

	if !p.curTokenIs(token.BRACE_R) {
		return nil
	}
	p.nextToken()

	return lit
}

func (p *Parser) parseObjectField() *ast.ObjectField {
	field := &ast.ObjectField{Token: p.curToken}
	field.Name = &ast.Name{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextToken()
	field.Value = p.parseValue()
	return field
}

// parseArguments parses (arg1: val1, arg2: val2)
func (p *Parser) parseArguments() []*ast.Argument {
	p.nextToken()

	var args []*ast.Argument
	for !p.curTokenIs(token.PAREN_R) && !p.curTokenIs(token.EOF) {
		arg := &ast.Argument{Token: p.curToken}
		arg.Name = &ast.Name{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(token.COLON) {
			p.peekError(token.COLON)
			return nil
		}
		p.nextToken()

		arg.Value = p.parseValue()
		args = append(args, arg)

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	p.nextToken() // skip ')'
	return args
}

// parseDirectives parses @dir(args) @dir2
func (p *Parser) parseDirectives() []*ast.Directive {
	var directives []*ast.Directive

	for p.curTokenIs(token.AT) {
		d := &ast.Directive{Token: p.curToken}
		p.nextToken() // skip '@'

		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, "expected directive name")
			return nil
		}
		d.Name = p.curToken.Literal
		p.nextToken()

		if p.curTokenIs(token.PAREN_L) {
			d.Arguments = p.parseArguments()
		}

		directives = append(directives, d)
	}

	return directives
}

func (p *Parser) parseVariableDefinitions() []*ast.VariableDefinition {
	if !p.curTokenIs(token.PAREN_L) {
		return nil
	}
	p.nextToken()

	var vars []*ast.VariableDefinition
	for !p.curTokenIs(token.PAREN_R) && !p.curTokenIs(token.EOF) {
		def := &ast.VariableDefinition{Token: p.curToken}

		def.Variable = p.parseVariable()

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()

		def.Type = p.parseType()

		if p.curTokenIs(token.EQUALS) {
			p.nextToken()
			def.DefaultValue = p.parseValue()
		}

		def.Directives = p.parseDirectives()
		vars = append(vars, def)

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.PAREN_R) {
		p.nextToken()
	}
	return vars
}
