package parser

import (
	"log"

	"github.com/pspiagicw/fener/ast"
	"github.com/pspiagicw/fener/lexer"
	"github.com/pspiagicw/fener/token"
)

const (
	_ = iota
	LOWEST
	BINARY
)

var precedence = map[token.TokenType]int{
	token.PLUS:  BINARY,
	token.MINUS: BINARY,
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:            l,
		prefixParseFnMap: map[token.TokenType]prefixParseFn{},
		infixParseFnMap:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefixFn(token.NUMBER, p.parseNumberExpression)

	p.registerInfixFn(token.PLUS, p.parseBinaryExpression)
	p.registerInfixFn(token.MINUS, p.parseBinaryExpression)

	p.advance()
	return p
}

func (p *Parser) registerPrefixFn(tokentype token.TokenType, fn prefixParseFn) {
	p.prefixParseFnMap[tokentype] = fn
}
func (p *Parser) registerInfixFn(tokentype token.TokenType, fn infixParseFn) {
	p.infixParseFnMap[tokentype] = fn
}
func appendStatement(a *ast.AST, statement ast.Statement) {
	if statement != nil {
		a.Statements = append(a.Statements, statement)
	}
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	lexer            *lexer.Lexer
	current          *token.Token
	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap  map[token.TokenType]infixParseFn
}

func (p *Parser) advance() {
	p.current = p.lexer.Next()
}

func (p *Parser) ParseAST() *ast.AST {
	a := &ast.AST{}

	for !p.lexer.EOF {
		statement := p.parseStatement()

		appendStatement(a, statement)
	}

	return a
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.parseLetStatements()
	default:
		log.Fatalf("Can't start statement with '%s'", p.current.Type)
		return nil
	}
}
func (p *Parser) parseLetStatements() ast.LetStatement {
	p.advance()

	let := ast.LetStatement{}

	p.expect(p.current, token.IDENTIFIER)

	let.Name = p.current.Value

	p.advance()

	p.expect(p.current, token.ASSIGN)

	p.advance()

	let.Value = p.parseExpression(LOWEST)

	return let
}
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFnMap[p.current.Type]

	if prefixFn == nil {
		log.Fatalf("Can't find prefix fn for type '%s'", p.current.Type)
	}

	// Should have a advance function automatically in it.
	left := prefixFn()

	for p.current.Type != token.EOF && precedence < p.currentPrecedence() {

		infixFn := p.infixParseFnMap[p.current.Type]

		if infixFn == nil {
			return left
		}

		left = infixFn(left)

		if left == nil {
			return nil
		}

	}

	return left
}
func (p *Parser) currentPrecedence() int {
	val, ok := precedence[p.current.Type]

	if ok {
		return val
	}

	return LOWEST
}
func (p *Parser) parseNumberExpression() ast.Expression {
	n := ast.NumberExpression{
		Value: p.current.Value,
	}

	p.advance()
	return n
}
func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.current

	b := ast.BinaryExpression{
		Left:     left,
		Operator: operator,
	}

	operatorPrecedence := p.currentPrecedence()

	p.advance()

	b.Right = p.parseExpression(operatorPrecedence)

	return b
}
func (p *Parser) expect(t *token.Token, tokentype token.TokenType) {
	if t.Type != tokentype {
		log.Fatalf("Expected '%s', got '%s'", tokentype, t.Type)
	}
}
