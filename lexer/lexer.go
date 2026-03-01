package lexer

import (
	"strings"

	"github.com/pspiagicw/tremor/token"
)

type Lexer struct {
	input   string
	file    string
	curPos  int
	readPos int
	length  int
	current string
	line    int
	column  int
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
		if !l.EOF {
			if l.current == "\n" {
				l.line += 1
				l.column = 1
			} else {
				l.column += 1
			}
		}
		l.EOF = true
		l.current = ""
	} else {
		if l.current == "\n" {
			l.line += 1
			l.column = 1
		} else {
			l.column += 1
		}
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

func newToken(tokentype token.TokenType, value string, offset int, line int, column int) *token.Token {
	return &token.Token{
		Type:   tokentype,
		Value:  value,
		Offset: offset,
		Line:   line,
		Column: column,
	}
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
	case "class":
		return token.CLASS
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

	startOffset := l.curPos
	startLine := l.line
	startColumn := l.column

	emit := func(tokentype token.TokenType, value string) *token.Token {
		return newToken(tokentype, value, startOffset, startLine, startColumn)
	}

	if l.EOF {
		return newToken(token.EOF, "", l.readPos, l.line, l.column)
	}

	switch l.current {
	case "+":
		return emit(token.PLUS, l.current)
	case "-":
		if l.peek() == "-" {
			l.comment()
			return l.Next()
		}
		return emit(token.MINUS, l.current)
	case "%":
		return emit(token.MODULUS, l.current)
	case ",":
		return emit(token.COMMA, l.current)
	case "^":
		return emit(token.EXPONENT, l.current)
	case "*":
		return emit(token.MULTIPLY, l.current)
	case "/":
		return emit(token.SLASH, l.current)
	case ":":
		return emit(token.COLON, l.current)
	case "(":
		return emit(token.LPAREN, l.current)
	case ")":
		return emit(token.RPAREN, l.current)
	case "{":
		return emit(token.LBRACE, l.current)
	case "}":
		return emit(token.RBRACE, l.current)
	case "[":
		if l.peek() == "[" {
			l.advance()
			value := l.longString()
			return emit(token.STRING_MULTILINE, value)
		}
		return emit(token.LSQUARE, l.current)
	case "]":
		return emit(token.RSQUARE, l.current)
	case ".":
		if l.peek() == "." {
			l.advance()
			if l.peek() == "." {
				l.advance()
				return emit(token.ELLIPSIS, "...")
			}
			return emit(token.CONCAT, "..")
		}
		return emit(token.DOT, l.current)
	case "=":
		if l.peek() == "=" {
			l.advance()
			return emit(token.EQ, "==")
		}
		return emit(token.ASSIGN, l.current)
	case "!":
		if l.peek() == "=" {
			l.advance()
			return emit(token.NEQ, "!=")
		}
		return emit(token.BANG, l.current)
	case "<":
		if l.peek() == "=" {
			l.advance()
			return emit(token.LTE, "<=")
		}
		return emit(token.LT, l.current)
	case ">":
		if l.peek() == "=" {
			l.advance()
			return emit(token.GTE, ">=")
		}
		return emit(token.GT, l.current)
	case "'":
		value := l.string(l.current)
		return emit(token.STRING_SINGLE, value)
	case "\"":
		value := l.string(l.current)
		return emit(token.STRING_DOUBLE, value)
	default:
		if isAlpha(l.current) {
			value := l.identifier()
			tokentype := predictType(value)
			return emit(tokentype, value)
		} else if isDigit(l.current) {
			value := l.number()
			tokentype := predictNumber(value)
			return emit(tokentype, value)
		}
		return emit(token.INVALID, l.current)
	}
}

func NewLexer(input string) *Lexer {
	return NewLexerWithFile(input, "<input>")
}

func NewLexerWithFile(input string, file string) *Lexer {
	if file == "" {
		file = "<input>"
	}
	return &Lexer{
		input:   input,
		file:    file,
		curPos:  -1,
		readPos: 0,
		length:  len(input),
		line:    1,
		column:  0,
	}
}

func (l *Lexer) Source() string {
	return l.input
}

func (l *Lexer) FileName() string {
	return l.file
}
