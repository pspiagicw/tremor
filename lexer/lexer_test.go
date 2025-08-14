package lexer

import "testing"
import "github.com/pspiagicw/fener/token"

func TestSimple(t *testing.T) {
}

func TestSymbol(t *testing.T) {
	input := "+ - * / ! % ^ , ."
	expectedTokens := []token.Token{
		{Type: token.PLUS, Value: "+"},
		{Type: token.MINUS, Value: "-"},
		{Type: token.MULTIPLY, Value: "*"},
		{Type: token.SLASH, Value: "/"},
		{Type: token.BANG, Value: "!"},
		{Type: token.MODULUS, Value: "%"},
		{Type: token.EXPONENT, Value: "^"},
		{Type: token.COMMA, Value: ","},
		{Type: token.DOT, Value: "."},
		{Type: token.EOF, Value: ""},
	}

	testToken(t, input, expectedTokens)
}

func TestComparisonOperators(t *testing.T) {
	input := "== != <= >="
	expected := []token.Token{
		{Type: token.EQ, Value: "=="},
		{Type: token.NEQ, Value: "!="},
		{Type: token.LTE, Value: "<="},
		{Type: token.GTE, Value: ">="},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)

}

func TestSingleCharComparisonAndAssignment(t *testing.T) {
	input := "< > ="
	expected := []token.Token{
		{Type: token.LT, Value: "<"},
		{Type: token.GT, Value: ">"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestDelimiters(t *testing.T) {
	input := "( ) { } [ ]"
	expected := []token.Token{
		{Type: token.LPAREN, Value: "("},
		{Type: token.RPAREN, Value: ")"},
		{Type: token.LBRACE, Value: "{"},
		{Type: token.RBRACE, Value: "}"},
		{Type: token.LSQUARE, Value: "["},
		{Type: token.RSQUARE, Value: "]"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestDots(t *testing.T) {
	input := ".. ..."
	expected := []token.Token{
		{Type: token.CONCAT, Value: ".."},
		{Type: token.ELLIPSIS, Value: "..."},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestIdentifiers(t *testing.T) {
	input := "foo bar _baz"
	expected := []token.Token{
		{Type: token.IDENTIFIER, Value: "foo"},
		{Type: token.IDENTIFIER, Value: "bar"},
		{Type: token.IDENTIFIER, Value: "_baz"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestNumbers(t *testing.T) {
	input := "123 3.14"
	expected := []token.Token{
		{Type: token.NUMBER, Value: "123"},
		{Type: token.NUMBER, Value: "3.14"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestKeywords(t *testing.T) {
	input := "if else return fn end let not and or then"
	expected := []token.Token{
		{Type: token.IF, Value: "if"},
		{Type: token.ELSE, Value: "else"},
		{Type: token.RETURN, Value: "return"},
		{Type: token.FN, Value: "fn"},
		{Type: token.END, Value: "end"},
		{Type: token.LET, Value: "let"},
		{Type: token.NOT, Value: "not"},
		{Type: token.AND, Value: "and"},
		{Type: token.OR, Value: "or"},
		{Type: token.THEN, Value: "then"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestLiterals(t *testing.T) {
	input := "nil true false"
	expected := []token.Token{
		{Type: token.NIL, Value: "nil"},
		{Type: token.TRUE, Value: "true"},
		{Type: token.FALSE, Value: "false"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestStrings(t *testing.T) {
	input := `"hello" 'world' [[long string]]`
	expected := []token.Token{
		{Type: token.STRING_DOUBLE, Value: "hello"},
		{Type: token.STRING_SINGLE, Value: "world"},
		{Type: token.STRING_MULTILINE, Value: "long string"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestSingleLineCommentOnly(t *testing.T) {
	input := "-- this is a comment"
	expected := []token.Token{
		{Type: token.EOF, Value: ""},
	} // No tokens returned for comment-only input
	testToken(t, input, expected)
}

func TestMultilineCommentOnly(t *testing.T) {
	input := "--[[ this is a \n multiline comment ]]"
	expected := []token.Token{
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestMultilineCommentBetweenCode(t *testing.T) {
	input := `
	let y = 10
	--[[
		this comment should be ignored
	]]
	y = y + 1
	`
	expected := []token.Token{
		{Type: token.LET, Value: "let"},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.NUMBER, Value: "10"},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.PLUS, Value: "+"},
		{Type: token.NUMBER, Value: "1"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestTypes(t *testing.T) {
	input := `int float string bool`

	expected := []token.Token{
		{Type: token.TYPE, Value: "int"},
		{Type: token.TYPE, Value: "float"},
		{Type: token.TYPE, Value: "string"},
		{Type: token.TYPE, Value: "bool"},
	}

	testToken(t, input, expected)
}

func testToken(t *testing.T, input string, expectedTokens []token.Token) {
	lexer := NewLexer(input)

	for i := range expectedTokens {
		actual := lexer.Next()
		expected := expectedTokens[i]

		if expected.Type != actual.Type {
			t.Fatalf("Token Type not matching: got %s expected %s", actual.Type, expected.Type)
		}
		if expected.Value != actual.Value {
			t.Fatalf("Token Value not matching: got %s expected %s", actual.Value, expected.Value)
		}
	}
}
