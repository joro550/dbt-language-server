package jinja

import (
	"bytes"
	"fmt"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Node
type File struct {
	Statements []Statement
}

func (f *File) TokenLiteral() string {
	if len(f.Statements) > 0 {
		return f.Statements[0].TokenLiteral()
	}
	return ""
}

func (f *File) String() string {
	var out bytes.Buffer

	for _, s := range f.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Statement
type SetStatment struct {
	Value Expression
	Name  *Identifier
	Token Token
}

func (ss *SetStatment) statementNode()       {}
func (ss *SetStatment) TokenLiteral() string { return ss.Token.Value }

func (ss *SetStatment) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Name.String())
	out.WriteString(" = ")

	if ss.Value != nil {
		out.WriteString(ss.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Value Expression
	Token Token
}

func (i *ExpressionStatement) statementNode()       {}
func (i *ExpressionStatement) TokenLiteral() string { return i.Token.Value }
func (es *ExpressionStatement) String() string      { return es.Value.String() }

// Expressions
type Identifier struct {
	Value string
	Token Token
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string       { return i.Value }

type IntegerExpression struct {
	Value int64
	Token Token
}

func (i *IntegerExpression) expressionNode()      {}
func (i *IntegerExpression) TokenLiteral() string { return i.Token.Value }
func (i *IntegerExpression) String() string       { return fmt.Sprintf("%v", i.Token.Value) }
