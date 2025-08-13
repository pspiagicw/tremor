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
	MODULUS  = "MODULUS"
	EXPONENT = "EXPONENT"
	COMMA    = "COMMA"

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
	LSQUARE = "LSQUARE"
	RSQUARE = "RSQUARE"

	CONCAT   = "CONCAT"
	ELLIPSIS = "ELLIPSIS"
	DOT      = "DOT"

	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	FN     = "FN"
	END    = "END"
	LET    = "LET"

	NIL   = "NIL"
	TRUE  = "TRUE"
	FALSE = "FALSE"
	NOT   = "NOT"
	AND   = "AND"
	OR    = "OR"

	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"

	STRING_DOUBLE    = "STRING_DOUBLE"
	STRING_SINGLE    = "STRING_SINGLE"
	STRING_MULTILINE = "STRING_MULTILINE"
)
