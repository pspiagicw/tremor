package parser

import (
	"fmt"
	"testing"

	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/token"
)

func TestLetStatementError(t *testing.T) {
	input := `let a int 1`

	expected := fmt.Sprintf(FAILED_EXPECT_MESSAGE, token.ASSIGN, token.INTEGER)

	testParserError(t, input, expected)
}

func TestStatementError(t *testing.T) {
	input := `=`

	expected := fmt.Sprintf(FAILED_PREFIX_MESSAGE, token.ASSIGN)

	testParserError(t, input, expected)
}

// func TestLetStatementTypeError(t *testing.T) {
// 	input := `let a b = 1`
//
// 	expected := fmt.Sprintf(FAILED_EXPECT_MESSAGE, token.ASSIGN, token.IDENTIFIER)
//
// 	testParserError(t, input, expected)
//
// }

func testParserError(t *testing.T, input string, message string) {
	l := lexer.NewLexer(input)
	p := NewParser(l)

	_ = p.ParseAST()

	errs := p.Errors()

	if len(errs) == 0 {
		t.Fatalf("Expected some errors, got zero!")
	}

	firsterr := errs[0].Error()

	if firsterr != message {
		t.Fatalf("Error message doesn't match, expected %s, got %s", message, firsterr)
	}
}
