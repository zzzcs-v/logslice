package query

import (
	"testing"
)

func tokens(input string) []Token {
	l := NewLexer(input)
	var result []Token
	for {
		tok := l.Next()
		result = append(result, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return result
}

// assertTokenTypes checks that the token types match the expected sequence.
func assertTokenTypes(t *testing.T, toks []Token, expected []TokenType) {
	t.Helper()
	if len(toks) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %+v", len(expected), len(toks), toks)
	}
	for i, tt := range expected {
		if toks[i].Type != tt {
			t.Errorf("token[%d]: expected type %d, got %d (literal=%q)", i, tt, toks[i].Type, toks[i].Literal)
		}
	}
}

func TestLexerEmpty(t *testing.T) {
	toks := tokens("")
	if len(toks) != 1 || toks[0].Type != TOKEN_EOF {
		t.Fatalf("expected EOF, got %+v", toks)
	}
}

func TestLexerEqOperator(t *testing.T) {
	toks := tokens("level == \"error\"")
	expected := []TokenType{TOKEN_IDENT, TOKEN_EQ, TOKEN_STRING, TOKEN_EOF}
	assertTokenTypes(t, toks, expected)
	if toks[0].Literal != "level" {
		t.Errorf("expected ident 'level', got %q", toks[0].Literal)
	}
	if toks[2].Literal != "error" {
		t.Errorf("expected string 'error', got %q", toks[2].Literal)
	}
}

func TestLexerNeqOperator(t *testing.T) {
	toks := tokens("status != 200")
	expected := []TokenType{TOKEN_IDENT, TOKEN_NEQ, TOKEN_NUMBER, TOKEN_EOF}
	assertTokenTypes(t, toks, expected)
	if toks[2].Literal != "200" {
		t.Errorf("expected number '200', got %q", toks[2].Literal)
	}
}

func TestLexerAndOrKeywords(t *testing.T) {
	toks := tokens("a == 'x' AND b != 'y' OR c == 'z'")
	expected := []TokenType{
		TOKEN_IDENT, TOKEN_EQ, TOKEN_STRING,
		TOKEN_AND,
		TOKEN_IDENT, TOKEN_NEQ, TOKEN_STRING,
		TOKEN_OR,
		TOKEN_IDENT, TOKEN_EQ, TOKEN_STRING,
		TOKEN_EOF,
	}
	assertTokenTypes(t, toks, expected)
}

func TestLexerParens(t *testing.T) {
	toks := tokens("(field == \"val\")")
	if toks[0].Type != TOKEN_LPAREN {
		t.Errorf("expected LPAREN, got %+v", toks[0])
	}
	if toks[len(toks)-2].Type != TOKEN_RPAREN {
		t.Errorf("expected RPAREN, got %+v", toks[len(toks)-2])
	}
}

func TestLexerNegativeNumber(t *testing.T) {
	toks := tokens("code == -1")
	if toks[2].Type != TOKEN_NUMBER || toks[2].Literal != "-1" {
		t.Errorf("expected number '-1', got %+v", toks[2])
	}
}

func TestLexerDottedIdent(t *testing.T) {
	toks := tokens("request.method == \"GET\"")
	if toks[0].Type != TOKEN_IDENT || toks[0].Literal != "request.method" {
		t.Errorf("expected dotted ident 'request.method', got %+v", toks[0])
	}
}
