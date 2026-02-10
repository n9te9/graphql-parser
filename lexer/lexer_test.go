package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/graphql-parser/token"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "Punctuators",
			input: "! $ ( ) ... : = @ [ ] { | } &",
			expected: []token.Token{
				{Type: token.BANG, Literal: "!", Line: 1, Start: 0, End: 1},
				{Type: token.DOLLAR, Literal: "$", Line: 1, Start: 2, End: 3},
				{Type: token.PAREN_L, Literal: "(", Line: 1, Start: 4, End: 5},
				{Type: token.PAREN_R, Literal: ")", Line: 1, Start: 6, End: 7},
				{Type: token.SPREAD, Literal: "...", Line: 1, Start: 8, End: 11},
				{Type: token.COLON, Literal: ":", Line: 1, Start: 12, End: 13},
				{Type: token.EQUALS, Literal: "=", Line: 1, Start: 14, End: 15},
				{Type: token.AT, Literal: "@", Line: 1, Start: 16, End: 17},
				{Type: token.BRACKET_L, Literal: "[", Line: 1, Start: 18, End: 19},
				{Type: token.BRACKET_R, Literal: "]", Line: 1, Start: 20, End: 21},
				{Type: token.BRACE_L, Literal: "{", Line: 1, Start: 22, End: 23},
				{Type: token.PIPE, Literal: "|", Line: 1, Start: 24, End: 25},
				{Type: token.BRACE_R, Literal: "}", Line: 1, Start: 26, End: 27},
				{Type: token.AMP, Literal: "&", Line: 1, Start: 28, End: 29},
				{Type: token.EOF, Literal: "", Line: 1, Start: 29, End: 29},
			},
		},
		{
			name: "Keywords",
			input: `query mutation subscription fragment type input enum union 
interface scalar directive extend schema implements on true false null`,
			expected: []token.Token{
				{Type: token.QUERY, Literal: "query", Line: 1, Start: 0, End: 5},
				{Type: token.MUTATION, Literal: "mutation", Line: 1, Start: 6, End: 14},
				{Type: token.SUBSCRIPTION, Literal: "subscription", Line: 1, Start: 15, End: 27},
				{Type: token.FRAGMENT, Literal: "fragment", Line: 1, Start: 28, End: 36},
				{Type: token.TYPE, Literal: "type", Line: 1, Start: 37, End: 41},
				{Type: token.INPUT, Literal: "input", Line: 1, Start: 42, End: 47},
				{Type: token.ENUM, Literal: "enum", Line: 1, Start: 48, End: 52},
				{Type: token.UNION, Literal: "union", Line: 1, Start: 53, End: 58},
				{Type: token.INTERFACE, Literal: "interface", Line: 2, Start: 60, End: 69},
				{Type: token.SCALAR, Literal: "scalar", Line: 2, Start: 70, End: 76},
				{Type: token.DIRECTIVE, Literal: "directive", Line: 2, Start: 77, End: 86},
				{Type: token.EXTEND, Literal: "extend", Line: 2, Start: 87, End: 93},
				{Type: token.SCHEMA, Literal: "schema", Line: 2, Start: 94, End: 100},
				{Type: token.IMPLEMENTS, Literal: "implements", Line: 2, Start: 101, End: 111},
				{Type: token.ON, Literal: "on", Line: 2, Start: 112, End: 114},
				{Type: token.TRUE, Literal: "true", Line: 2, Start: 115, End: 119},
				{Type: token.FALSE, Literal: "false", Line: 2, Start: 120, End: 125},
				{Type: token.NULL, Literal: "null", Line: 2, Start: 126, End: 130},
				{Type: token.EOF, Literal: "", Line: 2, Start: 130, End: 130},
			},
		},
		{
			name:  "Line Terminators (CRLF and CR)",
			input: "a\r\nb\rc",
			expected: []token.Token{
				{Type: token.IDENT, Literal: "a", Line: 1, Start: 0, End: 1},
				{Type: token.IDENT, Literal: "b", Line: 2, Start: 3, End: 4},
				{Type: token.IDENT, Literal: "c", Line: 3, Start: 5, End: 6},
				{Type: token.EOF, Literal: "", Line: 3, Start: 6, End: 6},
			},
		},
		{
			name:  "Multi-byte Characters (Valid)",
			input: "\"あ\"\n# 🍺\n\"\"\"\nあ\n\"\"\"",
			expected: []token.Token{
				{Type: token.STRING, Literal: "あ", Line: 1, Start: 0, End: 5},
				{Type: token.BLOCK_STRING, Literal: "あ", Line: 3, Start: 13, End: 24},
				{Type: token.EOF, Literal: "", Line: 5, Start: 24, End: 24},
			},
		},
		{
			name:  "Multi-byte Characters (Illegal Identifier)",
			input: "queryあ",
			expected: []token.Token{
				{Type: token.QUERY, Literal: "query", Line: 1, Start: 0, End: 5},
				{Type: token.ILLEGAL, Literal: "\xe3", Line: 1, Start: 5, End: 6},
				{Type: token.ILLEGAL, Literal: "\x81", Line: 1, Start: 6, End: 7},
				{Type: token.ILLEGAL, Literal: "\x82", Line: 1, Start: 7, End: 8},
				{Type: token.EOF, Literal: "", Line: 1, Start: 8, End: 8},
			},
		},
		{
			name: "Integers and Floats",
			input: `0
1234
-56
1.23
-1.23
1e5
1.23e+4
-1.23E-5`,
			expected: []token.Token{
				{Type: token.INT, Literal: "0", Line: 1, Start: 0, End: 1},
				{Type: token.INT, Literal: "1234", Line: 2, Start: 2, End: 6},
				{Type: token.INT, Literal: "-56", Line: 3, Start: 7, End: 10},
				{Type: token.FLOAT, Literal: "1.23", Line: 4, Start: 11, End: 15},
				{Type: token.FLOAT, Literal: "-1.23", Line: 5, Start: 16, End: 21},
				{Type: token.FLOAT, Literal: "1e5", Line: 6, Start: 22, End: 25},
				{Type: token.FLOAT, Literal: "1.23e+4", Line: 7, Start: 26, End: 33},
				{Type: token.FLOAT, Literal: "-1.23E-5", Line: 8, Start: 34, End: 42},
				{Type: token.EOF, Literal: "", Line: 8, Start: 42, End: 42},
			},
		},
		{
			name: "Strings",
			input: `"simple"
"with \" escaped quote"
"with unicode \u1234"
"""block string"""
"""
multi
line
"""`,
			expected: []token.Token{
				{Type: token.STRING, Literal: "simple", Line: 1, Start: 0, End: 8},
				{Type: token.STRING, Literal: "with \" escaped quote", Line: 2, Start: 9, End: 32},
				{Type: token.STRING, Literal: "with unicode ሴ", Line: 3, Start: 33, End: 54},
				{Type: token.BLOCK_STRING, Literal: "block string", Line: 4, Start: 55, End: 73},
				{Type: token.BLOCK_STRING, Literal: "multi\nline", Line: 5, Start: 74, End: 92},
				{Type: token.EOF, Literal: "", Line: 8, Start: 92, End: 92},
			},
		},
		{
			name:  "Comments and Commas (Ignored Tokens)",
			input: "query, # This is a comment\n{ id }",
			expected: []token.Token{
				{Type: token.QUERY, Literal: "query", Line: 1, Start: 0, End: 5},
				{Type: token.BRACE_L, Literal: "{", Line: 2, Start: 27, End: 28},
				{Type: token.IDENT, Literal: "id", Line: 2, Start: 29, End: 31},
				{Type: token.BRACE_R, Literal: "}", Line: 2, Start: 32, End: 33},
				{Type: token.EOF, Literal: "", Line: 2, Start: 33, End: 33},
			},
		},
		{
			name:  "Variables and Directives",
			input: "query ($var: String) @directive",
			expected: []token.Token{
				{Type: token.QUERY, Literal: "query", Line: 1, Start: 0, End: 5},
				{Type: token.PAREN_L, Literal: "(", Line: 1, Start: 6, End: 7},
				{Type: token.DOLLAR, Literal: "$", Line: 1, Start: 7, End: 8},
				{Type: token.IDENT, Literal: "var", Line: 1, Start: 8, End: 11},
				{Type: token.COLON, Literal: ":", Line: 1, Start: 11, End: 12},
				{Type: token.IDENT, Literal: "String", Line: 1, Start: 13, End: 19},
				{Type: token.PAREN_R, Literal: ")", Line: 1, Start: 19, End: 20},
				{Type: token.AT, Literal: "@", Line: 1, Start: 21, End: 22},
				{Type: token.DIRECTIVE, Literal: "directive", Line: 1, Start: 22, End: 31},
				{Type: token.EOF, Literal: "", Line: 1, Start: 31, End: 31},
			},
		},
		{
			name: "Complex Schema Definition",
			input: `type User implements Node {
  id: ID!
  name: String
}`,
			expected: []token.Token{
				{Type: token.TYPE, Literal: "type", Line: 1, Start: 0, End: 4},
				{Type: token.IDENT, Literal: "User", Line: 1, Start: 5, End: 9},
				{Type: token.IMPLEMENTS, Literal: "implements", Line: 1, Start: 10, End: 20},
				{Type: token.IDENT, Literal: "Node", Line: 1, Start: 21, End: 25},
				{Type: token.BRACE_L, Literal: "{", Line: 1, Start: 26, End: 27},
				{Type: token.IDENT, Literal: "id", Line: 2, Start: 30, End: 32},
				{Type: token.COLON, Literal: ":", Line: 2, Start: 32, End: 33},
				{Type: token.IDENT, Literal: "ID", Line: 2, Start: 34, End: 36},
				{Type: token.BANG, Literal: "!", Line: 2, Start: 36, End: 37},
				{Type: token.IDENT, Literal: "name", Line: 3, Start: 40, End: 44},
				{Type: token.COLON, Literal: ":", Line: 3, Start: 44, End: 45},
				{Type: token.IDENT, Literal: "String", Line: 3, Start: 46, End: 52},
				{Type: token.BRACE_R, Literal: "}", Line: 4, Start: 53, End: 54},
				{Type: token.EOF, Literal: "", Line: 4, Start: 54, End: 54},
			},
		},
		{
			name: "BOM (Byte Order Mark)",
			// 1. Start of file: \xef\xbb\xbf
			// 2. Concatenated file (start of line): ...\n\xef\xbb\xbf...
			input: "\xef\xbb\xbfquery\n\xef\xbb\xbf{ id }",
			expected: []token.Token{
				{Type: token.QUERY, Literal: "query", Line: 1, Start: 3, End: 8}, // BOM(3bytes) skipped
				{Type: token.BRACE_L, Literal: "{", Line: 2, Start: 12, End: 13}, // \n(1) + BOM(3) skipped
				{Type: token.IDENT, Literal: "id", Line: 2, Start: 14, End: 16},
				{Type: token.BRACE_R, Literal: "}", Line: 2, Start: 17, End: 18},
				{Type: token.EOF, Literal: "", Line: 2, Start: 18, End: 18},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			got := l.Tokens()

			if len(got) != len(tt.expected) {
				t.Fatalf("Tokens length mismatch. want=%d, got=%d", len(tt.expected), len(got))
			}

			if d := cmp.Diff(tt.expected, got); d != "" {
				t.Errorf("Tokens mismatch (-want +got):\n%s", d)
			}
		})
	}
}
