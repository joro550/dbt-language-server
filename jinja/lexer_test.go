package jinja

import "testing"

func Test_NextToken(t *testing.T) {
	input := "=+-();"
	tests := []Token{
		{Token: ASSIGN, Value: "="},
		{Token: PLUS, Value: "+"},
		{Token: MINUS, Value: "-"},
		{Token: LEFT_BRACKET, Value: "("},
		{Token: RIGHT_BRACKET, Value: ")"},
		{Token: SEMI_COLON, Value: ";"},
		{Token: EOF, Value: ""},
	}

	lexer := NewJinjaLexer(input)
	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Token != tt.Token {
			t.Fatalf("test[%d] - token type wrong expected %q, got=%q", i, tt.Token, tok.Token)
		}

		if tok.Value != tt.Value {
			t.Fatalf("test[%d] - token type wrong expected %q, got=%q", i, tt.Value, tok.Value)
		}
	}
}

func Test_Identifiers(t *testing.T) {
	input := "set result = [\"thing\"]"
	tests := []Token{
		{Token: SET, Value: "set"},
		{Token: IDENT, Value: "result"},
		{Token: ASSIGN, Value: "="},
		{Token: START_COLLECTION, Value: "["},
		{Token: QUOTE, Value: "\""},
		{Token: IDENT, Value: "thing"},
		{Token: QUOTE, Value: "\""},
		{Token: END_COLLECTION, Value: "]"},
	}

	lexer := NewJinjaLexer(input)
	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Value != tt.Value {
			t.Fatalf("test[%d] - token type wrong expected %v, got=%v", i, tt.Value, tok.Value)
		}

		if tok.Token != tt.Token {
			t.Fatalf("test[%d] - token type wrong expected %v, got=%v", i, tt.Token, tok.Token)
		}
	}
}
