package jinja

import "testing"

func Test_DigitToken(t *testing.T) {
	input := "{{ 5 }}"
	tests := []Token{
		{Token: START_EXPRESSION, Value: "{{"},
		{Token: INT, Value: "5"},
		{Token: END_EXPRESSION, Value: "}}"},
	}

	runTests(input, tests, t)
}

func Test_NextToken(t *testing.T) {
	input := "{{ =+-(); }}"
	tests := []Token{
		{Token: START_EXPRESSION, Value: "{{"},
		{Token: ASSIGN, Value: "="},
		{Token: PLUS, Value: "+"},
		{Token: MINUS, Value: "-"},
		{Token: LEFT_BRACKET, Value: "("},
		{Token: RIGHT_BRACKET, Value: ")"},
		{Token: SEMI_COLON, Value: ";"},
		{Token: END_EXPRESSION, Value: "}}"},
		{Token: EOF, Value: ""},
	}

	runTests(input, tests, t)
}

func Test_NextTextToken(t *testing.T) {
	input := "hello my name is joro {{ =+-(); }}"
	tests := []Token{
		{Token: TEXT, Value: "hello my name is joro "},
		{Token: START_EXPRESSION, Value: "{{"},
		{Token: ASSIGN, Value: "="},
		{Token: PLUS, Value: "+"},
		{Token: MINUS, Value: "-"},
		{Token: LEFT_BRACKET, Value: "("},
		{Token: RIGHT_BRACKET, Value: ")"},
		{Token: SEMI_COLON, Value: ";"},
		{Token: END_EXPRESSION, Value: "}}"},
		{Token: EOF, Value: ""},
	}

	runTests(input, tests, t)
}

func Test_Identifiers(t *testing.T) {
	input := "{{ set result = [\"thing\"] }}"
	tests := []Token{
		{Token: START_EXPRESSION, Value: "{{"},
		{Token: SET, Value: "set"},
		{Token: IDENT, Value: "result"},
		{Token: ASSIGN, Value: "="},
		{Token: START_COLLECTION, Value: "["},
		{Token: QUOTE, Value: "\""},
		{Token: IDENT, Value: "thing"},
		{Token: QUOTE, Value: "\""},
		{Token: END_COLLECTION, Value: "]"},
		{Token: END_EXPRESSION, Value: "}}"},
		{Token: EOF, Value: ""},
	}
	runTests(input, tests, t)
}

func runTests(input string, tokens []Token, t *testing.T) {
	lexer := NewJinjaLexer(input)
	for i, tt := range tokens {
		tok := lexer.NextToken()

		if tok.Value != tt.Value {
			t.Fatalf("test[%d] - token type wrong expected %v, got=%v", i, tt.Value, tok.Value)
		}

		if tok.Token != tt.Token {
			t.Fatalf("test[%d] - token type wrong expected %v, got=%v", i, tt.Token, tok.Token)
		}
	}
}
