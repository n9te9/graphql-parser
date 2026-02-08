package parser_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/lexer"
	"github.com/n9te9/graphql-parser/parser"
	"github.com/n9te9/graphql-parser/token"
)

func TestParseObjectTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name: "Simple Type Definition",
			input: `
            type User {
                id: ID!
                name: String
            }
            `,
			wantErr: "Unexpected token at top level: type",
			expect:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()

			if tt.wantErr != "" {
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
					return
				}
				found := false
				for _, err := range errors {
					if strings.Contains(err, tt.wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, errors)
				}
				return
			} else {
				if len(errors) > 0 {
					t.Fatalf("unexpected parser errors: %v", errors)
				}
			}

			opts := []cmp.Option{
				cmpopts.IgnoreTypes(token.Token{}),
			}

			if diff := cmp.Diff(tt.expect, got, opts...); diff != "" {
				t.Errorf("ParseDocument() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
