package ast_test

import (
	"testing"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func TestValue_String(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected string
	}{
		{
			name: "Integer Value",
			node: &ast.IntValue{
				Token: token.Token{Literal: "123"},
				Value: 123,
			},
			expected: "123",
		},
		{
			name: "Float Value",
			node: &ast.FloatValue{
				Token: token.Token{Literal: "12.34"},
				Value: 12.34,
			},
			expected: "12.34",
		},
		{
			name: "String Value",
			node: &ast.StringValue{
				Token: token.Token{Literal: "hello"},
				Value: "hello",
			},
			expected: "\"hello\"",
		},
		{
			name: "Boolean Value True",
			node: &ast.BooleanValue{
				Token: token.Token{Literal: "true"},
				Value: true,
			},
			expected: "true",
		},
		{
			name: "Boolean Value False",
			node: &ast.BooleanValue{
				Token: token.Token{Literal: "false"},
				Value: false,
			},
			expected: "false",
		},
		{
			name: "Null Value",
			node: &ast.NullValue{
				Token: token.Token{Literal: "null"},
			},
			expected: "null",
		},
		{
			name: "Enum Value",
			node: &ast.EnumValue{
				Token: token.Token{Literal: "USER_ROLE"},
				Value: "USER_ROLE",
			},
			expected: "USER_ROLE",
		},
		{
			name: "Variable",
			node: &ast.Variable{
				Token: token.Token{Literal: "$"},
				Name:  "userId",
			},
			expected: "$userId",
		},
		{
			name: "List Value",
			node: &ast.ListValue{
				Token: token.Token{Literal: "["},
				Values: []ast.Value{
					&ast.IntValue{Token: token.Token{Literal: "1"}, Value: 1},
					&ast.IntValue{Token: token.Token{Literal: "2"}, Value: 2},
				},
			},
			expected: "[1, 2]",
		},
		{
			name: "Object Value",
			node: &ast.ObjectValue{
				Token: token.Token{Literal: "{"},
				Fields: []*ast.ObjectField{
					{
						Token: token.Token{Literal: "id"},
						Name:  &ast.Name{Token: token.Token{Literal: "id"}, Value: "id"},
						Value: &ast.IntValue{Token: token.Token{Literal: "1"}, Value: 1},
					},
					{
						Token: token.Token{Literal: "name"},
						Name:  &ast.Name{Token: token.Token{Literal: "name"}, Value: "name"},
						Value: &ast.StringValue{Token: token.Token{Literal: "test"}, Value: "test"},
					},
				},
			},
			expected: "{id: 1, name: \"test\"}",
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
