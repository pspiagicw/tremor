package token

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

const (
	EOF     = "EOF"
	INVALID = "INVALID"

	PLUS     = "PLUS"
	MINUS    = "MINUS"
	MULTIPLY = "MULTIPLY"
	SLASH    = "SLASH"
	BANG     = "BANG"

	EQ  = "EQ"
	NEQ = "NEQ"
	LTE = "LTE"
	GTE = "GTE"

	LT     = "LT"
	GT     = "GT"
	ASSIGN = "ASSIGN"

	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	LBRACE  = "LBRACE"
	RBRACE  = "RBRACE"
	LSQUARE = "LBRACKET"
	RSQUARE = "RBRACKET"

	CONCAT   = "CONCAT"
	ELLIPSIS = "ELLIPSIS"

	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	FN     = "FN"
	END    = "END"
	LET    = "LET"

	NIL   = "NIL"
	TRUE  = "TRUE"
	FALSE = "FALSE"

	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"
)
