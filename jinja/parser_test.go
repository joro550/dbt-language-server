package jinja

import (
	"testing"
)

func TestSetStatement(t *testing.T) {
	input := "{% set name = 5 %}"

	lexer := NewJinjaLexer(input)
	parser := NewParser(lexer)

	file := parser.Parse()
	if file == nil {
		t.Fatalf("parse returned nil")
	}

	expectedIdentifier := "name"
	t.Logf("length %v", len(file.Statements))
	stmt := file.Statements[0].(*SetStatment)
	testSetStatement(t, stmt, expectedIdentifier)
}

func testSetStatement(t *testing.T, s *SetStatment, name string) bool {
	if s.TokenLiteral() != "set" {
		t.Errorf("s.TokenLiteral not let go %s", s.TokenLiteral())
		return false
	}

	if s.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %v got%v", name, s.Name.TokenLiteral())
	}

	return true
}
