package ast_test

import (
	"testing"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected string
	}{
		{
			name: "Named Type",
			node: &ast.NamedType{
				Token: token.Token{Literal: "String"},
				Name:  &ast.Name{Token: token.Token{Literal: "String"}, Value: "String"},
			},
			expected: "String",
		},
		{
			name: "List Type",
			node: &ast.ListType{
				Token: token.Token{Literal: "["},
				Type: &ast.NamedType{
					Token: token.Token{Literal: "Int"},
					Name:  &ast.Name{Token: token.Token{Literal: "Int"}, Value: "Int"},
				},
			},
			expected: "[Int]",
		},
		{
			name: "Non-Null Type (Named)",
			node: &ast.NonNullType{
				Token: token.Token{Literal: "!"},
				Type: &ast.NamedType{
					Token: token.Token{Literal: "ID"},
					Name:  &ast.Name{Token: token.Token{Literal: "ID"}, Value: "ID"},
				},
			},
			expected: "ID!",
		},
		{
			name: "Non-Null Type (List)",
			node: &ast.NonNullType{
				Token: token.Token{Literal: "!"},
				Type: &ast.ListType{
					Token: token.Token{Literal: "["},
					Type: &ast.NamedType{
						Token: token.Token{Literal: "User"},
						Name:  &ast.Name{Token: token.Token{Literal: "User"}, Value: "User"},
					},
				},
			},
			expected: "[User]!",
		},
		{
			name: "List of Non-Null Type",
			node: &ast.ListType{
				Token: token.Token{Literal: "["},
				Type: &ast.NonNullType{
					Token: token.Token{Literal: "!"},
					Type: &ast.NamedType{
						Token: token.Token{Literal: "String"},
						Name:  &ast.Name{Token: token.Token{Literal: "String"}, Value: "String"},
					},
				},
			},
			expected: "[String!]",
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
