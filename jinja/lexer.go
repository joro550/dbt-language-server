package jinja

import (
	"bufio"
	"bytes"
	"io"
)

type TOKEN int

type Token struct {
	Value string
	Token TOKEN
}

var eof = rune(0)

const (
	LEFT_BRACE TOKEN = iota
	RIGHT_BRACE
	PERCENT
	IDENT
	PIPE
	LEFT_BRACKET
	RIGHT_BRACKET

	START_EXPRESSION
	START_STATEMENT
	START_COMMENT

	END_EXPRESSION
	END_STATEMENT
	END_COMMENT

	NUMBER

	ILLEGAL
	WS

	// special characters
	TEXT
	EOF
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) ScanAll() []Token {
	tokens := []Token{}

	token, value := s.scan(false)
	withinJinja := false
	for token != EOF {

		switch token {
		case START_EXPRESSION, START_STATEMENT, START_COMMENT:
			withinJinja = true
		case END_EXPRESSION, END_STATEMENT, END_COMMENT:
			withinJinja = false
		}

		tokens = append(tokens, Token{Value: value, Token: token})
		token, value = s.scan(withinJinja)
	}

	return tokens
}

func (s *Scanner) scan(withinJinja bool) (TOKEN, string) {
	ch := s.read()

	// if we are not within a jinja template, we can only scan text
	if !withinJinja {
		if ch == eof {
			return EOF, ""
		} else if isWhitespace(ch) {
			s.unread()
			return s.scanWhitespace()
		} else {
			s.scanText()
		}
	}

	if ch == eof {
		return EOF, ""
	} else if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	} else if isJinjaIdentifier(ch) {
		s.unread()
		s.scanJinjaTemplate()
	}

	return ILLEGAL, ""
}

func (s *Scanner) scanWhitespace() (tok TOKEN, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok TOKEN, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanText() (tok TOKEN, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return TEXT, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanNumber() (tok TOKEN, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return NUMBER, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanJinjaTemplate() (tok TOKEN, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isJinjaIdentifier(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	stringVal := buf.String()
	switch stringVal {
	case "{{":
		return START_EXPRESSION, stringVal
	case "}}":
		return END_EXPRESSION, stringVal

	case "{%":
		return START_STATEMENT, stringVal
	case "%}":
		return END_STATEMENT, stringVal

	case "{#":
		return START_COMMENT, stringVal
	case "#}":
		return END_COMMENT, stringVal
	}

	s.unread()
	return ILLEGAL, stringVal
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\r' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isJinjaIdentifier(ch rune) bool {
	return ch == '%' || ch == '#' || ch == '{' || ch == '}'
}
