package lexer

import (
	"github.com/n9te9/graphql-parser/token"
)

type Lexer struct {
	input        []byte
	position     int
	readPosition int
	ch           byte
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: []byte(input),
		line:  1,
	}
	l.readChar()
	return l
}

func (l *Lexer) Tokens() []token.Token {
	tokens := make([]token.Token, 0)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token
	startPos := l.position
	startLine := l.line

	switch l.ch {
	case '$':
		tok = l.newToken(token.DOLLAR, l.ch)
	case '!':
		tok = l.newToken(token.BANG, l.ch)
	case '(':
		tok = l.newToken(token.PAREN_L, l.ch)
	case ')':
		tok = l.newToken(token.PAREN_R, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case '=':
		tok = l.newToken(token.EQUALS, l.ch)
	case '@':
		tok = l.newToken(token.AT, l.ch)
	case '[':
		tok = l.newToken(token.BRACKET_L, l.ch)
	case ']':
		tok = l.newToken(token.BRACKET_R, l.ch)
	case '{':
		tok = l.newToken(token.BRACE_L, l.ch)
	case '}':
		tok = l.newToken(token.BRACE_R, l.ch)
	case '|':
		tok = l.newToken(token.PIPE, l.ch)
	case '&':
		tok = l.newToken(token.AMP, l.ch)
	case '.':
		if l.peekChar() == '.' && l.peekChar2() == '.' {
			l.readChar()
			l.readChar()
			tok = token.Token{Type: token.SPREAD, Literal: "...", Line: startLine}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case '"':
		if l.peekChar() == '"' && l.peekChar2() == '"' {
			l.readChar()
			l.readChar()
			tok.Literal = l.readBlockString()
			tok.Type = token.BLOCK_STRING
		} else {
			tok.Literal = l.readString()
			tok.Type = token.STRING
		}
		tok.Line = startLine
		tok.Start = startPos
		tok.End = l.position
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Start = startPos
		tok.End = startPos
		return tok
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = startLine
			tok.Start = startPos
			tok.End = l.position
			return tok
		} else if isDigit(l.ch) || l.ch == '-' {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = startLine
			tok.Start = startPos
			tok.End = l.position
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()

	tok.Start = startPos
	tok.End = l.position
	tok.Line = startLine

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekChar2() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string([]byte{ch}), Line: l.line}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	tokenType := token.INT

	if l.ch == '-' {
		l.readChar()
	}

	if l.ch == '0' {
		l.readChar()
		if isDigit(l.ch) {
			return string(l.input[position:l.position]), token.ILLEGAL
		}
	} else {
		l.readDigits()
	}

	if l.ch == '.' {
		tokenType = token.FLOAT
		l.readChar()
		l.readDigits()
	}

	if l.ch == 'e' || l.ch == 'E' {
		tokenType = token.FLOAT
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		l.readDigits()
	}

	return string(l.input[position:l.position]), tokenType
}

func (l *Lexer) readDigits() {
	for isDigit(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	l.readChar()
	position := l.position
	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
		}
		l.readChar()
	}
	str := string(l.input[position:l.position])

	if l.ch == '"' {
		l.readChar()
	}

	return str
}

func (l *Lexer) readBlockString() string {
	l.readChar()
	position := l.position
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '"' && l.peekChar() == '"' && l.peekChar2() == '"' {
			break
		}
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
	str := string(l.input[position:l.position])

	// skip the ending """
	l.readChar()
	l.readChar()
	l.readChar()

	return str
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case ' ', '\t', ',':
			l.readChar()

		case '\n':
			l.line++
			l.readChar()

		case '\r':
			l.line++
			if l.peekChar() == '\n' {
				l.readChar()
			}
			l.readChar()

		case '#':
			for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
				l.readChar()
			}

		case 0xEF:
			// --- BOM (Byte Order Mark) check ---
			if l.peekChar() == 0xBB && l.peekChar2() == 0xBF {
				l.readChar()
				l.readChar()
				l.readChar()
			} else {
				return
			}

		default:
			return
		}
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
