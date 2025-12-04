package lexer

import (
	"strings"

	"github.com/pspiagicw/tremor/token"
)

type Lexer struct {
	input   string
	curPos  int
	readPos int
	length  int
	current string
	EOF     bool
}

func (l *Lexer) peek() string {
	if l.readPos < l.length {
		return string(l.input[l.readPos])
	}
	return ""
}

func (l *Lexer) advance() {
	if l.readPos == l.length {
		l.EOF = true
		l.current = ""
	} else {
		l.curPos = l.readPos
		l.readPos += 1
		l.current = string(l.input[l.curPos])
	}
}

func isWhitespace(input string) bool {
	if input == " " || input == "\n" || input == "\t" {
		return true
	}
	return false
}

func (l *Lexer) whitespace() {
	for !l.EOF && isWhitespace(l.current) {
		l.advance()
	}
}

func newToken(tokentype token.TokenType, value string) *token.Token {
	return &token.Token{Type: tokentype, Value: value}
}
func isAlpha(input string) bool {
	if len(input) == 0 {
		return false
	}
	letter := input[0]
	if ('a' <= letter && letter <= 'z') || ('A' <= letter && letter <= 'Z') || (letter == '_') {
		return true
	}
	return false
}
func isDigit(input string) bool {
	if len(input) == 0 {
		return false
	}
	letter := input[0]
	if '0' <= letter && letter <= '9' {
		return true
	}
	return false
}

func (l *Lexer) identifier() string {
	start := l.curPos
	for !l.EOF && (isAlpha(l.peek()) || isDigit(l.peek())) {
		l.advance()
	}
	end := l.curPos
	return l.input[start : end+1]
}
func (l *Lexer) number() string {
	start := l.curPos
	for !l.EOF && isDigit(l.peek()) || l.peek() == "." {
		l.advance()
	}
	end := l.curPos
	return l.input[start : end+1]
}
func (l *Lexer) longString() string {
	l.advance()
	start := l.curPos
	for !l.EOF && !(l.current == "]" && l.peek() == "]") {
		l.advance()
	}
	end := l.curPos
	l.advance() // Skip over the first ], the second ] will be automatically skipped over
	return l.input[start:end]
}
func (l *Lexer) string(endValue string) string {
	l.advance()
	start := l.curPos
	for !l.EOF && l.current != endValue {
		l.advance()
	}
	end := l.curPos
	return l.input[start:end]
}
func (l *Lexer) comment() {
	l.advance()
	l.advance() // Skip the 2 dashes
	multiline := false

	if l.current == "[" && l.peek() == "[" {
		multiline = true
	}

	for !l.EOF && multiline && !(l.current == "]" && l.peek() == "]") {
		l.advance()
	}

	if !l.EOF && multiline {
		l.advance()
	}

	for !l.EOF && !multiline && (l.current != "\n") {
		l.advance()
	}

}
func predictNumber(input string) token.TokenType {
	// TOOD: Complete lexing of integer and floats.

	dotcount := strings.Count(input, ".")

	if dotcount == 0 {
		return token.INTEGER
	}
	if dotcount == 1 {
		return token.FLOAT
	}
	return token.INVALID
}

func predictType(input string) token.TokenType {
	switch input {
	case "if":
		return token.IF
	case "else":
		return token.ELSE
	case "return":
		return token.RETURN
	case "fn":
		return token.FN
	case "end":
		return token.END
	case "let":
		return token.LET
	case "nil":
		return token.NIL
	case "true":
		return token.TRUE
	case "false":
		return token.FALSE
	case "not":
		return token.NOT
	case "and":
		return token.AND
	case "or":
		return token.OR
	case "then":
		return token.THEN
	case "int":
		fallthrough
	case "void":
		fallthrough
	case "float":
		fallthrough
	case "string":
		fallthrough
	case "bool":
		return token.TYPE
	default:
		return token.IDENTIFIER

	}

}

func (l *Lexer) Next() *token.Token {
	l.advance()
	l.whitespace()
	if l.EOF {
		return newToken(token.EOF, "")
	}

	switch l.current {
	case "+":
		return newToken(token.PLUS, l.current)
	case "-":
		if l.peek() == "-" {
			l.comment()
			return l.Next()
		}
		return newToken(token.MINUS, l.current)
	case "%":
		return newToken(token.MODULUS, l.current)
	case ",":
		return newToken(token.COMMA, l.current)
	case "^":
		return newToken(token.EXPONENT, l.current)
	case "*":
		return newToken(token.MULTIPLY, l.current)
	case "/":
		return newToken(token.SLASH, l.current)
	case ":":
		return newToken(token.COLON, l.current)
	case "(":
		return newToken(token.LPAREN, l.current)
	case ")":
		return newToken(token.RPAREN, l.current)
	case "{":
		return newToken(token.LBRACE, l.current)
	case "}":
		return newToken(token.RBRACE, l.current)
	case "[":
		if l.peek() == "[" {
			l.advance()
			value := l.longString()
			return newToken(token.STRING_MULTILINE, value)
		}
		return newToken(token.LSQUARE, l.current)
	case "]":
		return newToken(token.RSQUARE, l.current)
	case ".":
		if l.peek() == "." {
			l.advance()
			if l.peek() == "." {
				l.advance()
				return newToken(token.ELLIPSIS, "...")
			}
			return newToken(token.CONCAT, "..")
		}
		return newToken(token.DOT, l.current)
	case "=":
		if l.peek() == "=" {
			l.advance()
			return newToken(token.EQ, "==")
		}
		return newToken(token.ASSIGN, l.current)
	case "!":
		if l.peek() == "=" {
			l.advance()
			return newToken(token.NEQ, "!=")
		}
		return newToken(token.BANG, l.current)
	case "<":
		if l.peek() == "=" {
			l.advance()
			return newToken(token.LTE, "<=")
		}
		return newToken(token.LT, l.current)
	case ">":
		if l.peek() == "=" {
			l.advance()
			return newToken(token.GTE, ">=")
		}
		return newToken(token.GT, l.current)
	case "'":
		value := l.string(l.current)
		return newToken(token.STRING_SINGLE, value)
	case "\"":
		value := l.string(l.current)
		return newToken(token.STRING_DOUBLE, value)
	default:
		if isAlpha(l.current) {
			value := l.identifier()
			tokentype := predictType(value)
			return newToken(tokentype, value)
		} else if isDigit(l.current) {
			value := l.number()
			tokentype := predictNumber(value)
			return newToken(tokentype, value)
		}
		return newToken(token.INVALID, l.current)
	}
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:   input,
		curPos:  -1,
		readPos: 0,
		length:  len(input),
	}
}
