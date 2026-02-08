package ast_test

import (
	"testing"

	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/token"
)

func TestDirective_String(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected string
	}{
		{
			name: "Directive without Arguments",
			node: &ast.Directive{
				Token: token.Token{Literal: "@"},
				Name:  "deprecated",
			},
			expected: "@deprecated",
		},
		{
			name: "Directive with Single Argument",
			node: &ast.Directive{
				Token: token.Token{Literal: "@"},
				Name:  "key",
				Arguments: []*ast.Argument{
					{
						Token: token.Token{Literal: "fields"},
						Name:  &ast.Name{Token: token.Token{Literal: "fields"}, Value: "fields"},
						Value: &ast.StringValue{Token: token.Token{Literal: "id"}, Value: "id"},
					},
				},
			},
			expected: "@key(fields: \"id\")",
		},
		{
			name: "Directive with Multiple Arguments",
			node: &ast.Directive{
				Token: token.Token{Literal: "@"},
				Name:  "myDirective",
				Arguments: []*ast.Argument{
					{
						Token: token.Token{Literal: "arg1"},
						Name:  &ast.Name{Token: token.Token{Literal: "arg1"}, Value: "arg1"},
						Value: &ast.IntValue{Token: token.Token{Literal: "1"}, Value: 1},
					},
					{
						Token: token.Token{Literal: "arg2"},
						Name:  &ast.Name{Token: token.Token{Literal: "arg2"}, Value: "arg2"},
						Value: &ast.BooleanValue{Token: token.Token{Literal: "true"}, Value: true},
					},
				},
			},
			// Note: Argument order depends on implementation logic, assuming preserving order here
			expected: "@myDirective(arg1: 1, arg2: true)",
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
