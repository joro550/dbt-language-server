package jinja

type TokenType int

type Token struct {
	Value string
	Token TokenType
}

var eof = rune(0)

var keywords = map[string]TokenType{
	"set":      SET,
	"for":      FOR,
	"in":       IN,
	"macro":    MACRO,
	"endmacro": END_MACRO,
	"if":       IF,
}

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	PERCENT
	IDENT
	INT
	PIPE
	LEFT_BRACKET
	RIGHT_BRACKET
	ASSIGN
	PLUS
	MINUS
	SLASH
	ASTERIKS
	COMMA
	BANG
	COLLECTION
	SEMI_COLON
	START_COLLECTION
	END_COLLECTION
	QUOTE

	LT
	GT
	EQ
	NOT_EQ

	START_EXPRESSION
	START_STATEMENT
	START_COMMENT

	END_EXPRESSION
	END_STATEMENT
	END_COMMENT

	SET
	FOR
	IN
	MACRO
	END_MACRO
	IF

	ILLEGAL
	WS

	// special characters
	TEXT
	EOF
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewJinjaLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
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
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {

	case '{':
		nextChar := l.peekChar()

		switch nextChar {
		case '{':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: START_EXPRESSION, Value: literal}

		case '%':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: START_STATEMENT, Value: literal}
		case '#':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: START_COMMENT, Value: literal}
		default:
			tok = newToken(LEFT_BRACE, l.ch)
		}

	case '%':
		nextChar := l.peekChar()
		if nextChar == '}' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: END_STATEMENT, Value: literal}

		} else {
			tok = newToken(ASSIGN, l.ch)
		}

	case '#':
		nextChar := l.peekChar()
		if nextChar == '}' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: END_COMMENT, Value: literal}

		} else {
			tok = newToken(ASSIGN, l.ch)
		}

	case '}':
		nextChar := l.peekChar()
		if nextChar == '}' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: END_EXPRESSION, Value: literal}

		} else {
			tok = newToken(RIGHT_BRACE, l.ch)
		}

	case '=':
		nextChar := l.peekChar()

		if nextChar == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: EQ, Value: literal}

		} else {
			tok = newToken(ASSIGN, l.ch)
		}
	case '(':
		tok = newToken(LEFT_BRACKET, l.ch)
	case ')':
		tok = newToken(RIGHT_BRACKET, l.ch)
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '!':
		nextChar := l.peekChar()
		if nextChar == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: NOT_EQ, Value: literal}
		} else {
			tok = newToken(BANG, l.ch)
		}

	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case ';':
		tok = newToken(SEMI_COLON, l.ch)
	case '[':
		tok = newToken(START_COLLECTION, l.ch)
	case ']':
		tok = newToken(END_COLLECTION, l.ch)
	case '"':
		tok = newToken(QUOTE, l.ch)
	case 0:
		tok.Token = EOF
		tok.Value = ""

	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Token = LookupIdent(tok.Value)
			return tok
		} else if isDigit(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Token = INT
			return tok
		}

	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch == '_'
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Token: tokenType, Value: string(ch)}
}
