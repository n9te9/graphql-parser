package parser

import (
	"fmt"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/lexer"
	"github.com/n9te9/graphql-parser/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseDocument() *ast.Document {
	doc := &ast.Document{}
	doc.Definitions = []ast.Definition{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseDefinition()
		if stmt != nil {
			doc.Definitions = append(doc.Definitions, stmt)
		}

		if len(p.errors) > 0 {
			return nil
		}
	}

	return doc
}

func (p *Parser) parseDefinition() ast.Definition {
	switch p.curToken.Type {
	case token.QUERY, token.MUTATION, token.SUBSCRIPTION, token.BRACE_L:
		return p.parseOperationDefinition()
	case token.FRAGMENT:
		return p.parseFragmentDefinition()
	default:
		p.errors = append(p.errors, fmt.Sprintf("Unexpected token at top level: %s", p.curToken.Literal))
		return nil
	}
}

// curTokenIs checks if the current token type matches t.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if the next token type matches t.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the next token is t. If so, it advances the tokens.
// If not, it records an error.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseType() ast.Type {
	var t ast.Type

	if p.curTokenIs(token.BRACKET_L) {
		listType := &ast.ListType{Token: p.curToken}
		p.nextToken()
		listType.Type = p.parseType()
		if !p.curTokenIs(token.BRACKET_R) {
			p.peekError(token.BRACKET_R)
			return nil
		}
		t = listType
		p.nextToken()
	} else if p.curTokenIs(token.IDENT) {
		t = &ast.NamedType{
			Token: p.curToken,
			Name:  &ast.Name{Token: p.curToken, Value: p.curToken.Literal},
		}
		p.nextToken()
	}

	if p.curTokenIs(token.BANG) {
		nonNull := &ast.NonNullType{Token: p.curToken, Type: t}
		p.nextToken()
		return nonNull
	}
	return t
}
