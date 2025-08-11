package parser

import (
	"log"

	"github.com/pspiagicw/fener/ast"
	"github.com/pspiagicw/fener/lexer"
	"github.com/pspiagicw/fener/token"
)

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.advance()
	return p
}

func appendStatement(a *ast.AST, statement ast.Statement) {
	if statement != nil {
		a.Statements = append(a.Statements, statement)
	}
}

type Parser struct {
	lexer   *lexer.Lexer
	current *token.Token
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

	let.Value = p.parseExpression()

	return let
}
func (p *Parser) parseExpression() ast.Expression {
	switch p.current.Type {
	case token.NUMBER:
		return p.parseNumberExpression()
	default:
		return nil
	}
}
func (p *Parser) parseNumberExpression() ast.NumberExpression {
	n := ast.NumberExpression{
		Value: p.current.Value,
	}

	p.advance()
	return n
}
func (p *Parser) expect(t *token.Token, tokentype token.TokenType) {
	if t.Type != tokentype {
		log.Fatalf("Expected '%s', got '%s'", tokentype, t.Type)
	}
}
