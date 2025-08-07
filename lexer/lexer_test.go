package lexer

import "testing"
import "github.com/pspiagicw/fener/token"

func TestSimple(t *testing.T) {
}

func TestSymbol(t *testing.T) {
	input := "+ - * /"

	expectedTokens := []token.Token{
		{Type: token.PLUS, Value: "+"},
		{Type: token.MINUS, Value: "-"},
		{Type: token.ASTERISK, Value: "*"},
		{Type: token.SLASH, Value: "/"},
	}

	testToken(t, input, expectedTokens)
}

func TestComparisonOperators(t *testing.T) {
	input := "== ~= <= >="
	expected := []token.Token{
		{Type: token.EQ, Value: "=="},
		{Type: token.NEQ, Value: "~="},
		{Type: token.LTE, Value: "<="},
		{Type: token.GTE, Value: ">="},
	}
	testToken(t, input, expected)

}

func TestSingleCharComparisonAndAssignment(t *testing.T) {
	input := "< > ="
	expected := []token.Token{
		{Type: token.LT, Value: "<"},
		{Type: token.GT, Value: ">"},
		{Type: token.ASSIGN, Value: "="},
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
		{Type: token.LBRACKET, Value: "["},
		{Type: token.RBRACKET, Value: "]"},
	}
	testToken(t, input, expected)
}

func TestDots(t *testing.T) {
	input := ".. ..."
	expected := []token.Token{
		{Type: token.CONCAT, Value: ".."},
		{Type: token.ELLIPSIS, Value: "..."},
	}
	testToken(t, input, expected)
}

func TestKeywords(t *testing.T) {
	input := "if else return fn end let"
	expected := []token.Token{
		{Type: token.IF, Value: "if"},
		{Type: token.ELSE, Value: "else"},
		{Type: token.RETURN, Value: "return"},
		{Type: token.FN, Value: "fn"},
		{Type: token.END, Value: "end"},
	}
	testToken(t, input, expected)
}

func TestLiterals(t *testing.T) {
	input := "nil true false"
	expected := []token.Token{
		{Type: token.NIL, Value: "nil"},
		{Type: token.TRUE, Value: "true"},
		{Type: token.FALSE, Value: "false"},
	}
	testToken(t, input, expected)
}

func TestIdentifiers(t *testing.T) {
	input := "foo bar _baz"
	expected := []token.Token{
		{Type: token.IDENTIFIER, Value: "foo"},
		{Type: token.IDENTIFIER, Value: "bar"},
		{Type: token.IDENTIFIER, Value: "_baz"},
	}
	testToken(t, input, expected)
}
func TestNumbers(t *testing.T) {
	input := "123 3.14 0x1A"
	expected := []token.Token{
		{Type: token.NUMBER, Value: "123"},
		{Type: token.NUMBER, Value: "3.14"},
		{Type: token.NUMBER, Value: "0x1A"},
	}
	testToken(t, input, expected)
}
func TestStrings(t *testing.T) {
	input := `"hello" 'world' [[long string]]`
	expected := []token.Token{
		{Type: token.STRING, Value: "hello"},
		{Type: token.STRING, Value: "world"},
		{Type: token.STRING, Value: "long string"},
	}
	testToken(t, input, expected)
}

func TestSingleLineCommentOnly(t *testing.T) {
	input := "-- this is a comment"
	expected := []token.Token{} // No tokens returned for comment-only input
	testToken(t, input, expected)
}

func TestMultilineCommentOnly(t *testing.T) {
	input := "--[[ this is a \n multiline comment ]]"
	expected := []token.Token{}
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
		{Type: token.LOCAL, Value: "local"},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.NUMBER, Value: "10"},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.PLUS, Value: "+"},
		{Type: token.NUMBER, Value: "1"},
	}
	testToken(t, input, expected)
}

func testToken(t *testing.T, input string, expectedTokens []token.Token) {
	lexer := newLexer(input)

	actual := lexer.Next()

	for i := 0; i < len(expectedTokens); i++ {
		expected := expectedTokens[i]

		if expected.Type != actual.Type {
			t.Fatalf("Token Type not matching: got %s expected %s", actual.Type, expected.Type)
		}
		if expected.Value != actual.Value {
			t.Fatalf("Token Value not matching: got %s expected %s", actual.Value, expected.Value)
		}
	}
}
