package jinja

type Node interface {
	TokenLiteral() string
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

// Statement
type SetStatment struct {
	Value Expression
	Token Token
	Name  *Identifier
}

func (ss *SetStatment) statementNode() {}

func (ss *SetStatment) TokenLiteral() string {
	return ss.Token.Value
}

// Statement
type Identifier struct {
	Value string
	Token Token
}

func (i *Identifier) statementNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Value
}
