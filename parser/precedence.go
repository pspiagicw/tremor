package parser

import "github.com/pspiagicw/fener/token"

const (
	_ = iota
	LOWEST
	OR
	AND
	COMPARISON
	CONCAT
	BINARY
	MULTIPLY
	EXPONENT
	UNARY
	CALL
	INDEX
	FIELD
)

var precedence = map[token.TokenType]int{
	token.PLUS:     BINARY,
	token.MINUS:    BINARY,
	token.MULTIPLY: MULTIPLY,
	token.SLASH:    MULTIPLY,
	token.MODULUS:  MULTIPLY,
	token.EXPONENT: EXPONENT,
	token.CONCAT:   CONCAT,
	token.NOT:      UNARY,
	token.AND:      AND,
	token.OR:       OR,
	token.EQ:       COMPARISON,
	token.NEQ:      COMPARISON,
	token.LPAREN:   CALL,
	token.LSQUARE:  INDEX,
	token.DOT:      FIELD,
	token.GT:       COMPARISON,
	token.LT:       COMPARISON,
	token.LTE:      COMPARISON,
	token.GTE:      COMPARISON,
}
