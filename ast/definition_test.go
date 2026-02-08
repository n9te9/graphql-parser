package ast_test

import (
	"testing"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func TestDefinition_String(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected string
	}{
		{
			name: "Simple Field",
			node: &ast.Field{
				Token: token.Token{Literal: "name"},
				Name:  &ast.Name{Token: token.Token{Literal: "name"}, Value: "name"},
			},
			expected: "name",
		},
		{
			name: "Field with Alias and Arguments",
			node: &ast.Field{
				Token: token.Token{Literal: "user"},
				Alias: &ast.Name{Token: token.Token{Literal: "me"}, Value: "me"},
				Name:  &ast.Name{Token: token.Token{Literal: "user"}, Value: "user"},
				Arguments: []*ast.Argument{
					{
						Token: token.Token{Literal: "id"},
						Name:  &ast.Name{Token: token.Token{Literal: "id"}, Value: "id"},
						Value: &ast.IntValue{Token: token.Token{Literal: "1"}, Value: 1},
					},
				},
			},
			expected: "me: user(id: 1)",
		},
		{
			name: "Field with Directives and SelectionSet",
			node: &ast.Field{
				Token: token.Token{Literal: "hero"},
				Name:  &ast.Name{Token: token.Token{Literal: "hero"}, Value: "hero"},
				Directives: []*ast.Directive{
					{
						Token: token.Token{Literal: "@"},
						Name:  "skip",
						Arguments: []*ast.Argument{
							{
								Token: token.Token{Literal: "if"},
								Name:  &ast.Name{Token: token.Token{Literal: "if"}, Value: "if"},
								Value: &ast.BooleanValue{Token: token.Token{Literal: "true"}, Value: true},
							},
						},
					},
				},
				SelectionSet: []ast.Selection{
					&ast.Field{
						Token: token.Token{Literal: "name"},
						Name:  &ast.Name{Token: token.Token{Literal: "name"}, Value: "name"},
					},
				},
			},
			expected: "hero @skip(if: true) { name }",
		},
		{
			name: "Fragment Spread",
			node: &ast.FragmentSpread{
				Token: token.Token{Literal: "..."},
				Name:  &ast.Name{Token: token.Token{Literal: "MyFragment"}, Value: "MyFragment"},
			},
			expected: "...MyFragment",
		},
		{
			name: "Inline Fragment",
			node: &ast.InlineFragment{
				Token: token.Token{Literal: "..."},
				TypeCondition: &ast.NamedType{
					Token: token.Token{Literal: "User"},
					Name:  &ast.Name{Token: token.Token{Literal: "User"}, Value: "User"},
				},
				SelectionSet: []ast.Selection{
					&ast.Field{
						Token: token.Token{Literal: "id"},
						Name:  &ast.Name{Token: token.Token{Literal: "id"}, Value: "id"},
					},
				},
			},
			expected: "... on User { id }",
		},
		{
			name: "Operation Definition (Query)",
			node: &ast.OperationDefinition{
				Token:     token.Token{Literal: "query"},
				Operation: ast.Query,
				Name:      &ast.Name{Token: token.Token{Literal: "MyQuery"}, Value: "MyQuery"},
				VariableDefinitions: []*ast.VariableDefinition{
					{
						Token:    token.Token{Literal: "$"},
						Variable: &ast.Variable{Token: token.Token{Literal: "$"}, Name: "id"},
						Type: &ast.NonNullType{
							Token: token.Token{Literal: "!"},
							Type: &ast.NamedType{
								Token: token.Token{Literal: "ID"},
								Name:  &ast.Name{Token: token.Token{Literal: "ID"}, Value: "ID"},
							},
						},
					},
				},
				SelectionSet: []ast.Selection{
					&ast.Field{
						Token: token.Token{Literal: "user"},
						Name:  &ast.Name{Token: token.Token{Literal: "user"}, Value: "user"},
					},
				},
			},
			expected: "query MyQuery($id: ID!) { user }",
		},
		{
			name: "Operation Definition (Mutation)",
			node: &ast.OperationDefinition{
				Token:     token.Token{Literal: "mutation"},
				Operation: ast.Mutation,
				SelectionSet: []ast.Selection{
					&ast.Field{
						Token: token.Token{Literal: "createUser"},
						Name:  &ast.Name{Token: token.Token{Literal: "createUser"}, Value: "createUser"},
					},
				},
			},
			expected: "mutation  { createUser }",
		},
		{
			name: "Fragment Definition",
			node: &ast.FragmentDefinition{
				Token: token.Token{Literal: "fragment"},
				Name:  &ast.Name{Token: token.Token{Literal: "UserParts"}, Value: "UserParts"},
				TypeCondition: &ast.NamedType{
					Token: token.Token{Literal: "User"},
					Name:  &ast.Name{Token: token.Token{Literal: "User"}, Value: "User"},
				},
				SelectionSet: []ast.Selection{
					&ast.Field{
						Token: token.Token{Literal: "id"},
						Name:  &ast.Name{Token: token.Token{Literal: "id"}, Value: "id"},
					},
					&ast.Field{
						Token: token.Token{Literal: "name"},
						Name:  &ast.Name{Token: token.Token{Literal: "name"}, Value: "name"},
					},
				},
			},
			expected: "fragment UserParts on User { id name }",
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
