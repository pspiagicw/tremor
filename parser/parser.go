package parser

import (
	"fmt"

	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/token"
)

type ParserError error

type Parser struct {
	lexer            *lexer.Lexer
	current          *token.Token
	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap  map[token.TokenType]infixParseFn
	errors           []ParserError
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:            l,
		prefixParseFnMap: map[token.TokenType]prefixParseFn{},
		infixParseFnMap:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefixFn(token.NUMBER, p.parseNumberExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.NOT, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBooleanExpression)
	p.registerPrefixFn(token.FALSE, p.parseBooleanExpression)
	p.registerPrefixFn(token.LPAREN, p.parseParenthesisExpression)
	p.registerPrefixFn(token.STRING_DOUBLE, p.parseStringExpression)
	p.registerPrefixFn(token.STRING_SINGLE, p.parseStringExpression)
	p.registerPrefixFn(token.STRING_MULTILINE, p.parseStringExpression)
	p.registerPrefixFn(token.IDENTIFIER, p.parseIdentifierExpression)
	p.registerPrefixFn(token.FN, p.parseLambdaExpression)

	p.registerInfixFn(token.PLUS, p.parseBinaryExpression)
	p.registerInfixFn(token.MINUS, p.parseBinaryExpression)
	p.registerInfixFn(token.MULTIPLY, p.parseBinaryExpression)
	p.registerInfixFn(token.SLASH, p.parseBinaryExpression)
	p.registerInfixFn(token.MODULUS, p.parseBinaryExpression)
	p.registerInfixFn(token.EXPONENT, p.parseBinaryExpression)
	p.registerInfixFn(token.CONCAT, p.parseBinaryExpression)
	p.registerInfixFn(token.AND, p.parseBinaryExpression)
	p.registerInfixFn(token.OR, p.parseBinaryExpression)
	p.registerInfixFn(token.EQ, p.parseBinaryExpression)
	p.registerInfixFn(token.NEQ, p.parseBinaryExpression)
	p.registerInfixFn(token.LPAREN, p.parseFunctionCallExpression)
	p.registerInfixFn(token.LSQUARE, p.parseIndexExpression)
	p.registerInfixFn(token.DOT, p.parseFieldExpression)
	p.registerInfixFn(token.GT, p.parseBinaryExpression)
	p.registerInfixFn(token.LT, p.parseBinaryExpression)
	p.registerInfixFn(token.LTE, p.parseBinaryExpression)
	p.registerInfixFn(token.GTE, p.parseBinaryExpression)

	p.advance()
	return p
}

func (p *Parser) registerPrefixFn(tokentype token.TokenType, fn prefixParseFn) {
	p.prefixParseFnMap[tokentype] = fn
}
func (p *Parser) registerInfixFn(tokentype token.TokenType, fn infixParseFn) {
	p.infixParseFnMap[tokentype] = fn
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

func (p *Parser) advance() {
	p.current = p.lexer.Next()
}

func (p *Parser) ParseAST() *ast.AST {
	a := &ast.AST{}

	for !p.lexer.EOF {
		statement := p.parseStatement()

		if statement != nil {
			a.Statements = append(a.Statements, statement)
		} else {
			p.advance()
		}
	}

	return a
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}
func (p *Parser) registerError(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	p.errors = append(p.errors, err)
}
