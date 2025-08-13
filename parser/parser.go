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
func (p *Parser) parseStringExpression() ast.Expression {
	s := ast.StringExpression{
		Value: p.current.Value,
	}

	switch p.current.Type {
	case token.STRING_DOUBLE:
		s.Type = ast.DOUBLE_QUOTED
	case token.STRING_SINGLE:
		s.Type = ast.SINGLE_QUOTED
	default:
		s.Type = ast.MULTILINE
	}

	p.advance()
	return s
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	i := ast.IdentifierExpression{
		Value: p.current,
	}

	p.advance()
	return i
}
func (p *Parser) parseBooleanExpression() ast.Expression {
	b := ast.BooleanExpression{
		Value: p.current,
	}
	p.advance()

	return b
}
func (p *Parser) parseParenthesisExpression() ast.Expression {
	exp := ast.ParenthesisExpression{}

	p.advance()

	exp.Inside = p.parseExpression(LOWEST)

	// Skip over the ending round brackets.
	p.advance()

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := ast.PrefixExpression{
		Operator: p.current,
	}

	p.advance()
	// Hard coded precedence value!
	exp.Right = p.parseExpression(UNARY)

	return exp
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
func (p *Parser) parseFunctionCallExpression(left ast.Expression) ast.Expression {
	f := ast.FunctionCallExpression{
		Caller: left,
	}

	p.advance()

	f.Arguments = []ast.Expression{}

	for p.current.Type != token.RPAREN {
		f.Arguments = append(f.Arguments, p.parseExpression(LOWEST))

		if p.current.Type == token.COMMA {
			p.advance()
		}
	}

	p.advance()

	return f
}
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	i := ast.IndexExpression{
		Caller: left,
	}

	p.advance()

	i.Index = p.parseExpression(LOWEST)

	p.advance()

	return i
}
func (p *Parser) parseFieldExpression(left ast.Expression) ast.Expression {
	f := ast.FieldExpression{
		Caller: left,
	}

	p.advance()

	// Hard coded dot precedence
	f.Field = p.parseExpression(FIELD)

	return f
}
func (p *Parser) expect(t *token.Token, tokentype token.TokenType) {
	if t.Type != tokentype {
		log.Fatalf("Expected '%s', got '%s'", tokentype, t.Type)
	}
}
