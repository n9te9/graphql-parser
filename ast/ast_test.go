package ast_test

import (
	"testing"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func TestDocument_String(t *testing.T) {
	tests := []struct {
		name     string
		node     *ast.Document
		expected string
	}{
		{
			name: "Document with Single Query",
			node: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Token:     token.Token{Literal: "query"},
						Operation: ast.Query,
						Name:      &ast.Name{Token: token.Token{Literal: "MyQuery"}, Value: "MyQuery"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Token: token.Token{Literal: "user"},
								Name:  &ast.Name{Token: token.Token{Literal: "user"}, Value: "user"},
							},
						},
					},
				},
			},
			expected: "query MyQuery { user }",
		},
		{
			name: "Document with Multiple Definitions",
			node: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Token:     token.Token{Literal: "query"},
						Operation: ast.Query,
						Name:      &ast.Name{Token: token.Token{Literal: "GetU"}, Value: "GetU"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Token: token.Token{Literal: "u"},
								Name:  &ast.Name{Token: token.Token{Literal: "u"}, Value: "u"},
							},
						},
					},
					&ast.OperationDefinition{
						Token:     token.Token{Literal: "mutation"},
						Operation: ast.Mutation,
						Name:      &ast.Name{Token: token.Token{Literal: "SaveU"}, Value: "SaveU"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Token: token.Token{Literal: "s"},
								Name:  &ast.Name{Token: token.Token{Literal: "s"}, Value: "s"},
							},
						},
					},
				},
			},
			expected: "query GetU { u }mutation SaveU { s }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.String(); got != tt.expected {
				t.Errorf("String() mismatch.\n got:      %s\n expected: %s", got, tt.expected)
			}
		})
	}
}

func TestDocument_TokenLiteral(t *testing.T) {
	tests := []struct {
		name     string
		node     *ast.Document
		expected string
	}{
		{
			name: "Document with Definitions",
			node: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Token: token.Token{Literal: "query"},
					},
				},
			},
			expected: "query",
		},
		{
			name: "Empty Document",
			node: &ast.Document{
				Definitions: []ast.Definition{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.TokenLiteral(); got != tt.expected {
				t.Errorf("TokenLiteral() mismatch.\n got:      %s\n expected: %s", got, tt.expected)
			}
		})
	}
}
