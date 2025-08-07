package token

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

const (
	TOKEN_EOF = "TOKEN_EOF"
)
