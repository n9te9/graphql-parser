package token

import "strconv"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INT
	FLOAT
	STRING
	BLOCK_STRING

	BANG      // !
	DOLLAR    // $
	PAREN_L   // (
	PAREN_R   // )
	SPREAD    // ...
	COLON     // :
	EQUALS    // =
	AT        // @
	BRACKET_L // [
	BRACKET_R // ]
	BRACE_L   // {
	PIPE      // |
	BRACE_R   // }
	AMP       // &
	COMMA     // ,

	NAME
	QUERY        // query
	MUTATION     // mutation
	SUBSCRIPTION // subscription
	FRAGMENT     // fragment
	TYPE         // type
	INPUT        // input
	ENUM         // enum
	UNION        // union
	INTERFACE    // interface
	SCALAR       // scalar
	DIRECTIVE    // directive
	EXTEND       // extend
	SCHEMA       // schema
	IMPLEMENTS   // implements
	ON           // on
	TRUE         // true
	FALSE        // false
	NULL         // null
)

var tokens = map[TokenType]string{
	ILLEGAL:      "ILLEGAL",
	EOF:          "EOF",
	IDENT:        "IDENT",
	INT:          "INT",
	FLOAT:        "FLOAT",
	STRING:       "STRING",
	BLOCK_STRING: "BLOCK_STRING",

	BANG:      "!",
	DOLLAR:    "$",
	PAREN_L:   "(",
	PAREN_R:   ")",
	SPREAD:    "...",
	COLON:     ":",
	EQUALS:    "=",
	AT:        "@",
	BRACKET_L: "[",
	BRACKET_R: "]",
	BRACE_L:   "{",
	PIPE:      "|",
	BRACE_R:   "}",
	AMP:       "&",

	NAME:         "NAME",
	QUERY:        "query",
	MUTATION:     "mutation",
	SUBSCRIPTION: "subscription",
	FRAGMENT:     "fragment",
	TYPE:         "type",
	INPUT:        "input",
	ENUM:         "enum",
	UNION:        "union",
	INTERFACE:    "interface",
	SCALAR:       "scalar",
	DIRECTIVE:    "directive",
	EXTEND:       "extend",
	SCHEMA:       "schema",
	IMPLEMENTS:   "implements",
	ON:           "on",
	TRUE:         "true",
	FALSE:        "false",
	NULL:         "null",
}

var keywords = map[string]TokenType{
	"query":        QUERY,
	"mutation":     MUTATION,
	"subscription": SUBSCRIPTION,
	"fragment":     FRAGMENT,
	"type":         TYPE,
	"input":        INPUT,
	"enum":         ENUM,
	"union":        UNION,
	"interface":    INTERFACE,
	"scalar":       SCALAR,
	"directive":    DIRECTIVE,
	"extend":       EXTEND,
	"schema":       SCHEMA,
	"implements":   IMPLEMENTS,
	"on":           ON,
	"true":         TRUE,
	"false":        FALSE,
	"null":         NULL,
}

func (t TokenType) String() string {
	if s, ok := tokens[t]; ok {
		return s
	}
	return "token(" + strconv.Itoa(int(t)) + ")"
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Start   int
	End     int
}

func New(tokenType TokenType, ch byte, line int) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
	}
}
