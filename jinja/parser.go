package jinja

import "fmt"

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	Errors    []Error
}

type Error struct {
	Value    string
	Position int
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		Errors: []Error{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	fmt.Println(p.curToken.Value, p.peekToken.Value)
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
	return false
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
		return p.parseSetStatement()
	}

	return nil
}

func (p *Parser) parseSetStatement() Statement {
	p.nextToken()
	switch p.curToken.Token {
	case SET:
		return p.parseSetStatment()
	}
	return nil
}

func (p *Parser) parseSetStatment() *SetStatment {
	stmt := &SetStatment{Token: p.curToken}

	if !p.peekTokenIs(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.peekTokenIs(ASSIGN) {
		return nil
	}

	p.nextToken()

	return stmt
}
