package jinja

type TokenType int

type Token struct {
	Value string
	Token TokenType
}

var eof = rune(0)

var keywords = map[string]TokenType{
	"set":       SET,
	"endset":    END_SET,
	"for":       FOR,
	"endfor":    END_FOR,
	"in":        IN,
	"macro":     MACRO,
	"endmacro":  END_MACRO,
	"if":        IF,
	"elif":      ELIF,
	"is":        IS,
	"block":     BLOCK,
	"endblock":  END_BLOCK,
	"extends":   EXTENDS,
	"scoped":    SCOPED,
	"call":      CALL,
	"endcall":   END_CALL,
	"filter":    FILTER,
	"endfilter": END_FILTER,
	"not":       NOT,
}

const (
	// general language things
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	PERCENT
	IDENT
	INT
	PIPE
	TILDA
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
	SINGLE_QUOTE
	QUOTE
	DOT
	LT
	GT
	EQ
	NOT_EQ

	// jinja opertations
	START_EXPRESSION
	START_STATEMENT
	START_COMMENT

	END_EXPRESSION
	END_STATEMENT
	END_COMMENT

	// keywords
	SET
	END_SET
	FOR
	END_FOR
	IN
	BLOCK
	END_BLOCK
	EXTENDS
	IS
	MACRO
	END_MACRO
	IF
	ELIF
	SCOPED
	CALL
	END_CALL
	FILTER
	END_FILTER
	NOT

	// other stuff
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
	withinJinja  bool
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

	if l.ch == 0 {
		return Token{Token: EOF, Value: ""}
	} else if l.withinJinja {
		return l.nextJinjaToken()
	}

	position := l.position
	isNextCharacterJinja := func(ch byte) bool {
		if ch != '{' {
			return false
		}

		peekChar := l.peekChar()
		return peekChar == '{' || peekChar == '%' || peekChar == '#'
	}

	// are we starting out with a jinja block?
	if isNextCharacterJinja(l.ch) {
		l.withinJinja = true
		return l.nextJinjaToken()
	}

	for !isNextCharacterJinja(l.ch) && l.ch != 0 {
		l.readChar()
	}

	tok.Value = l.input[position:l.position]
	tok.Token = TEXT
	return tok
}

func (l *Lexer) nextJinjaToken() Token {
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
			l.withinJinja = false

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
			l.withinJinja = false

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
			l.withinJinja = false

			literal := string(ch) + string(l.ch)
			tok = Token{Token: END_EXPRESSION, Value: literal}

		} else {
			tok = newToken(RIGHT_BRACE, l.ch)
		}

	case '=':
		nextChar := l.peekChar()

		switch nextChar {
		case '=':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: EQ, Value: literal}
		case '!':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Token: NOT_EQ, Value: literal}
		default:

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
	case '\'':
		tok = newToken(SINGLE_QUOTE, l.ch)
	case '.':
		tok = newToken(DOT, l.ch)
	case '~':
		tok = newToken(TILDA, l.ch)
	case '|':
		tok = newToken(PIPE, l.ch)
	case 0:
		tok.Token = EOF
		tok.Value = ""

	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Token = LookupIdent(tok.Value)
			return tok
		} else if isDigit(l.ch) {
			tok.Value = l.readNumber()
			tok.Token = INT
			return tok
		}
		tok.Token = ILLEGAL
		tok.Value = ""
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
	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Token: tokenType, Value: string(ch)}
}
