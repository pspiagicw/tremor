package lexer

import "testing"
import "github.com/pspiagicw/tremor/token"

func TestSimple(t *testing.T) {
}

func TestSymbol(t *testing.T) {
	input := "+ - * / ! % ^ , . :"
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
		{Type: token.COLON, Value: ":"},
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
		{Type: token.INTEGER, Value: "123"},
		{Type: token.FLOAT, Value: "3.14"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestKeywords(t *testing.T) {
	input := "if else return fn end let not and or then class"
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
		{Type: token.CLASS, Value: "class"},
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
		{Type: token.INTEGER, Value: "10"},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.IDENTIFIER, Value: "y"},
		{Type: token.PLUS, Value: "+"},
		{Type: token.INTEGER, Value: "1"},
		{Type: token.EOF, Value: ""},
	}
	testToken(t, input, expected)
}

func TestArray(t *testing.T) {
	input := `[1,2,3]`

	expected := []token.Token{
		{Type: token.LSQUARE, Value: "["},
		{Type: token.INTEGER, Value: "1"},
		{Type: token.COMMA, Value: ","},
		{Type: token.INTEGER, Value: "2"},
		{Type: token.COMMA, Value: ","},
		{Type: token.INTEGER, Value: "3"},
		{Type: token.RSQUARE, Value: "]"},
		{Type: token.EOF, Value: ""},
	}

	testToken(t, input, expected)
}

func TestTypes(t *testing.T) {
	input := `int float string bool void`

	expected := []token.Token{
		{Type: token.TYPE, Value: "int"},
		{Type: token.TYPE, Value: "float"},
		{Type: token.TYPE, Value: "string"},
		{Type: token.TYPE, Value: "bool"},
		{Type: token.TYPE, Value: "void"},
	}

	testToken(t, input, expected)
}

func TestTokenLocations(t *testing.T) {
	input := "a\n==\n\"hi\""
	lexer := NewLexer(input)

	tests := []struct {
		ttype  token.TokenType
		value  string
		offset int
		line   int
		column int
	}{
		{ttype: token.IDENTIFIER, value: "a", offset: 0, line: 1, column: 1},
		{ttype: token.EQ, value: "==", offset: 2, line: 2, column: 1},
		{ttype: token.STRING_DOUBLE, value: "hi", offset: 5, line: 3, column: 1},
		{ttype: token.EOF, value: "", offset: 9, line: 3, column: 5},
	}

	for _, tt := range tests {
		tok := lexer.Next()
		if tok.Type != tt.ttype {
			t.Fatalf("Token type mismatch: got %s expected %s", tok.Type, tt.ttype)
		}
		if tok.Value != tt.value {
			t.Fatalf("Token value mismatch: got %s expected %s", tok.Value, tt.value)
		}
		if tok.Offset != tt.offset {
			t.Fatalf("Token offset mismatch: got %d expected %d", tok.Offset, tt.offset)
		}
		if tok.Line != tt.line {
			t.Fatalf("Token line mismatch: got %d expected %d", tok.Line, tt.line)
		}
		if tok.Column != tt.column {
			t.Fatalf("Token column mismatch: got %d expected %d", tok.Column, tt.column)
		}
	}
}

func TestTokenLocationsComplexDeclaration(t *testing.T) {
	input := "let total int = value_1 + 42"
	lexer := NewLexer(input)

	tests := []struct {
		ttype  token.TokenType
		value  string
		offset int
		line   int
		column int
	}{
		{ttype: token.LET, value: "let", offset: 0, line: 1, column: 1},
		{ttype: token.IDENTIFIER, value: "total", offset: 4, line: 1, column: 5},
		{ttype: token.TYPE, value: "int", offset: 10, line: 1, column: 11},
		{ttype: token.ASSIGN, value: "=", offset: 14, line: 1, column: 15},
		{ttype: token.IDENTIFIER, value: "value_1", offset: 16, line: 1, column: 17},
		{ttype: token.PLUS, value: "+", offset: 24, line: 1, column: 25},
		{ttype: token.INTEGER, value: "42", offset: 26, line: 1, column: 27},
		{ttype: token.EOF, value: "", offset: 28, line: 1, column: 29},
	}

	for _, tt := range tests {
		tok := lexer.Next()
		if tok.Type != tt.ttype {
			t.Fatalf("Token type mismatch: got %s expected %s", tok.Type, tt.ttype)
		}
		if tok.Value != tt.value {
			t.Fatalf("Token value mismatch: got %s expected %s", tok.Value, tt.value)
		}
		if tok.Offset != tt.offset {
			t.Fatalf("Token offset mismatch: got %d expected %d", tok.Offset, tt.offset)
		}
		if tok.Line != tt.line {
			t.Fatalf("Token line mismatch: got %d expected %d", tok.Line, tt.line)
		}
		if tok.Column != tt.column {
			t.Fatalf("Token column mismatch: got %d expected %d", tok.Column, tt.column)
		}
	}
}

func TestTokenLocationsFunctionImplementation(t *testing.T) {
	input := "fn add(a int, b int) int then\n    return a + b\nend"
	lexer := NewLexer(input)

	tests := []struct {
		ttype  token.TokenType
		value  string
		offset int
		line   int
		column int
	}{
		{ttype: token.FN, value: "fn", offset: 0, line: 1, column: 1},
		{ttype: token.IDENTIFIER, value: "add", offset: 3, line: 1, column: 4},
		{ttype: token.LPAREN, value: "(", offset: 6, line: 1, column: 7},
		{ttype: token.IDENTIFIER, value: "a", offset: 7, line: 1, column: 8},
		{ttype: token.TYPE, value: "int", offset: 9, line: 1, column: 10},
		{ttype: token.COMMA, value: ",", offset: 12, line: 1, column: 13},
		{ttype: token.IDENTIFIER, value: "b", offset: 14, line: 1, column: 15},
		{ttype: token.TYPE, value: "int", offset: 16, line: 1, column: 17},
		{ttype: token.RPAREN, value: ")", offset: 19, line: 1, column: 20},
		{ttype: token.TYPE, value: "int", offset: 21, line: 1, column: 22},
		{ttype: token.THEN, value: "then", offset: 25, line: 1, column: 26},
		{ttype: token.RETURN, value: "return", offset: 34, line: 2, column: 5},
		{ttype: token.IDENTIFIER, value: "a", offset: 41, line: 2, column: 12},
		{ttype: token.PLUS, value: "+", offset: 43, line: 2, column: 14},
		{ttype: token.IDENTIFIER, value: "b", offset: 45, line: 2, column: 16},
		{ttype: token.END, value: "end", offset: 47, line: 3, column: 1},
		{ttype: token.EOF, value: "", offset: 50, line: 3, column: 4},
	}

	for _, tt := range tests {
		tok := lexer.Next()
		if tok.Type != tt.ttype {
			t.Fatalf("Token type mismatch: got %s expected %s", tok.Type, tt.ttype)
		}
		if tok.Value != tt.value {
			t.Fatalf("Token value mismatch: got %s expected %s", tok.Value, tt.value)
		}
		if tok.Offset != tt.offset {
			t.Fatalf("Token offset mismatch: got %d expected %d", tok.Offset, tt.offset)
		}
		if tok.Line != tt.line {
			t.Fatalf("Token line mismatch: got %d expected %d", tok.Line, tt.line)
		}
		if tok.Column != tt.column {
			t.Fatalf("Token column mismatch: got %d expected %d", tok.Column, tt.column)
		}
	}
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
