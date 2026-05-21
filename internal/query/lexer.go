package query

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a lexer token.
type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_IDENT
	TOKEN_STRING
	TOKEN_NUMBER
	TOKEN_EQ    // ==
	TOKEN_NEQ   // !=
	TOKEN_AND   // AND
	TOKEN_OR    // OR
	TOKEN_LPAREN
	TOKEN_RPAREN
)

// Token holds a lexer token with its type and literal value.
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer tokenizes a query string.
type Lexer struct {
	input []rune
	pos   int
}

// NewLexer creates a new Lexer for the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0}
}

// Next returns the next token from the input.
func (l *Lexer) Next() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TOKEN_EOF}
	}

	ch := l.input[l.pos]

	switch {
	case ch == '(':
		l.pos++
		return Token{Type: TOKEN_LPAREN, Literal: "("}
	case ch == ')':
		l.pos++
		return Token{Type: TOKEN_RPAREN, Literal: ")"}
	case ch == '=' && l.peek() == '=':
		l.pos += 2
		return Token{Type: TOKEN_EQ, Literal: "=="}
	case ch == '!' && l.peek() == '=':
		l.pos += 2
		return Token{Type: TOKEN_NEQ, Literal: "!="}
	case ch == '"' || ch == '\'':
		return l.readString(ch)
	case unicode.IsDigit(ch) || (ch == '-' && unicode.IsDigit(l.peek())):
		return l.readNumber()
	case unicode.IsLetter(ch) || ch == '_':
		return l.readIdent()
	}

	l.pos++
	return Token{Type: TOKEN_IDENT, Literal: string(ch)}
}

func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *Lexer) readString(quote rune) Token {
	l.pos++ // skip opening quote
	var sb strings.Builder
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
			l.pos++
		}
		sb.WriteRune(l.input[l.pos])
		l.pos++
	}
	l.pos++ // skip closing quote
	return Token{Type: TOKEN_STRING, Literal: sb.String()}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	if l.input[l.pos] == '-' {
		l.pos++
	}
	for l.pos < len(l.input) && (unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		l.pos++
	}
	return Token{Type: TOKEN_NUMBER, Literal: string(l.input[start:l.pos])}
}

func (l *Lexer) readIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || l.input[l.pos] == '_' || unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		l.pos++
	}
	lit := string(l.input[start:l.pos])
	switch strings.ToUpper(lit) {
	case "AND":
		return Token{Type: TOKEN_AND, Literal: lit}
	case "OR":
		return Token{Type: TOKEN_OR, Literal: lit}
	}
	return Token{Type: TOKEN_IDENT, Literal: lit}
}
