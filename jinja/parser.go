package jinja

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS    // ==
	LESSGREAT // > or <
	SUM       // +
	PRODUCT   // *
	PREFIX    // -x or +x
	FUNCTION  // myFunction(x)
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []Error

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFn   map[TokenType]infixParseFn
}

type Error struct {
	Value    string
	Position int
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []Error{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(token TokenType, fn prefixParseFn) {
	p.prefixParseFns[token] = fn
}

func (p *Parser) registerInfix(token TokenType, fn infixParseFn) {
	p.infixParseFn[token] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.curToken.Token == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Token == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %v, got %v instead", t, p.peekToken.Token)
	p.errors = append(p.errors, Error{Value: msg, Position: 0})
}

func (p *Parser) GetErrors() []Error {
	return p.errors
}

func (p *Parser) Parse() *File {
	file := &File{}
	file.Statements = []Statement{}

	for p.curToken.Token != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			file.Statements = append(file.Statements, stmt)
		}
		p.nextToken()
	}

	return file
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Token {
	case START_STATEMENT:
		p.nextToken()
		switch p.curToken.Token {
		case SET:
			return p.parseSetStatement()
		}
	}

	return nil
}

func (p *Parser) parseSetStatement() Statement {
	stmt := &SetStatment{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	return stmt
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerExpression{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %v as integer", p.curToken.Value)
		p.errors = append(p.errors, Error{Value: msg, Position: p.l.position})
	}

	lit.Value = value
	return lit
}
