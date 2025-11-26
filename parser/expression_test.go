package parser

import (
	"testing"
)

func TestAddition(t *testing.T) {
	input := "1 + 2"

	expected := "(1 + 2)"

	testParser(t, input, expected)
}
func TestSubstraction(t *testing.T) {
	input := "1 - 2"

	expected := "(1 - 2)"

	testParser(t, input, expected)
}
func TestMultiplication(t *testing.T) {
	input := "2 * 3"

	expected := "(2 * 3)"

	testParser(t, input, expected)
}

func TestDivision(t *testing.T) {
	input := "6 / 3"

	expected := "(6 / 3)"

	testParser(t, input, expected)
}

func TestModulus(t *testing.T) {
	input := "7 % 3"

	expected := "(7 % 3)"

	testParser(t, input, expected)
}

func TestExponentiation(t *testing.T) {
	input := "2 ^ 3"

	expected := "(2 ^ 3)"

	testParser(t, input, expected)
}

func TestConcatenation(t *testing.T) {
	input := "\"hello\" .. \" world\""

	expected := "(\"hello\" .. \" world\")"

	testParser(t, input, expected)
}

func TestUnaryMinus(t *testing.T) {
	input := "- 5"

	expected := "(- 5)"

	testParser(t, input, expected)
}

func TestNotOperator(t *testing.T) {
	input := "not true"

	expected := "(not true)"

	testParser(t, input, expected)
}

func TestLogicalAnd(t *testing.T) {
	input := "true and false"

	expected := "(true and false)"

	testParser(t, input, expected)
}

func TestLogicalOr(t *testing.T) {
	input := "false or true"

	expected := "(false or true)"

	testParser(t, input, expected)
}

func TestEqualityComparison(t *testing.T) {
	input := "1 == 1"

	expected := "(1 == 1)"

	testParser(t, input, expected)
}

func TestNotEqualComparison(t *testing.T) {
	input := "1 != 2"

	expected := "(1 != 2)"

	testParser(t, input, expected)
}

func TestPrecedenceWithParentheses(t *testing.T) {
	input := "(1 + 2) * 3"

	expected := "((1 + 2) * 3)"

	testParser(t, input, expected)
}

func TestStringSingleQuotes(t *testing.T) {
	input := "'single quoted string'"

	expected := "'single quoted string'"

	testParser(t, input, expected)
}

func TestStringDoubleQuotes(t *testing.T) {
	input := "\"double quoted string\""

	expected := "\"double quoted string\""

	testParser(t, input, expected)
}

func TestLongBracketString(t *testing.T) {
	input := "[[multi\nline\nstring]]"

	expected := "[[multi\nline\nstring]]"

	testParser(t, input, expected)
}
func TestFunctionCall(t *testing.T) {
	input := "foo(1, 2, 3)"

	expected := "foo(1, 2, 3)"

	testParser(t, input, expected)
}

func TestNestedFunctionCall(t *testing.T) {
	input := "outer(inner(1), 2)"

	expected := "outer(inner(1), 2)"

	testParser(t, input, expected)
}

func TestIndexExpression(t *testing.T) {
	input := "arr[1]"

	expected := "arr[1]"

	testParser(t, input, expected)
}

func TestNestedIndexExpression(t *testing.T) {
	input := "arr[foo(1)]"

	expected := "arr[foo(1)]"

	testParser(t, input, expected)
}

func TestTableAccess(t *testing.T) {
	input := "obj.field"

	expected := "obj.field"

	testParser(t, input, expected)
}

func TestChainedTableAccess(t *testing.T) {
	input := "obj.inner.value"

	expected := "obj.inner.value"

	testParser(t, input, expected)
}

func TestPrecedenceMultipleOperators(t *testing.T) {
	input := "1 + 2 * 3 ^ 2 - 4 / 2"

	// ^ has highest precedence, then *, /, then +, -
	expected := "((1 + (2 * (3 ^ 2))) - (4 / 2))"

	testParser(t, input, expected)
}

func TestPrecedenceWithFunctionCall(t *testing.T) {
	input := "foo(1 + 2) * 3"

	expected := "(foo((1 + 2)) * 3)"

	testParser(t, input, expected)
}
func TestPrecedence_AddMul(t *testing.T) {
	input := "1 + 2 * 3"
	expected := "(1 + (2 * 3))"
	testParser(t, input, expected)
}

func TestPrecedence_MulAdd(t *testing.T) {
	input := "2 * 3 + 4"
	expected := "((2 * 3) + 4)"
	testParser(t, input, expected)
}

func TestPrecedence_ExpMul(t *testing.T) {
	input := "2 ^ 3 * 4"
	expected := "((2 ^ 3) * 4)"
	testParser(t, input, expected)
}

func TestPrecedence_MulExp(t *testing.T) {
	input := "2 * 3 ^ 4"
	expected := "(2 * (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrecedence_AddConcat(t *testing.T) {
	input := "1 + 2 .. 3"
	expected := "((1 + 2) .. 3)"
	testParser(t, input, expected)
}

func TestPrecedence_ConcatAdd(t *testing.T) {
	input := "1 .. 2 + 3"
	expected := "(1 .. (2 + 3))"
	testParser(t, input, expected)
}

func TestPrecedence_NotAnd(t *testing.T) {
	input := "not true and false"
	expected := "((not true) and false)"
	testParser(t, input, expected)
}

func TestPrecedence_AndOr(t *testing.T) {
	input := "true or false and true"
	expected := "(true or (false and true))"
	testParser(t, input, expected)
}

func TestPrecedence_ComparisonAdd(t *testing.T) {
	input := "1 + 2 == 3"
	expected := "((1 + 2) == 3)"
	testParser(t, input, expected)
}

func TestPrecedence_AddComparison(t *testing.T) {
	input := "1 == 2 + 3"
	expected := "(1 == (2 + 3))"
	testParser(t, input, expected)
}

func TestPrecedence_IndexMul(t *testing.T) {
	input := "arr[1] * 2"
	expected := "(arr[1] * 2)"
	testParser(t, input, expected)
}

func TestPrecedence_MulIndex(t *testing.T) {
	input := "2 * arr[1]"
	expected := "(2 * arr[1])"
	testParser(t, input, expected)
}

func TestPrecedence_FuncCallAdd(t *testing.T) {
	input := "foo(1) + 2"
	expected := "(foo(1) + 2)"
	testParser(t, input, expected)
}

func TestPrecedence_AddFuncCall(t *testing.T) {
	input := "1 + foo(2)"
	expected := "(1 + foo(2))"
	testParser(t, input, expected)
}

func TestPrecedence_TableAccessMul(t *testing.T) {
	input := "obj.field * 2"
	expected := "(obj.field * 2)"
	testParser(t, input, expected)
}

func TestPrecedence_MulTableAccess(t *testing.T) {
	input := "2 * obj.field"
	expected := "(2 * obj.field)"
	testParser(t, input, expected)
}

func TestPrecedence_FuncCallIndex(t *testing.T) {
	input := "foo(1)[2]"
	expected := "foo(1)[2]"
	testParser(t, input, expected)
}

func TestPrecedence_IndexFuncCall(t *testing.T) {
	input := "arr[foo(1)]"
	expected := "arr[foo(1)]"
	testParser(t, input, expected)
}

func TestPrecedence_TableAccessFuncCall(t *testing.T) {
	input := "obj.method(1)"
	expected := "obj.method(1)"
	testParser(t, input, expected)
}

func TestPrecedence_MixedComplex(t *testing.T) {
	input := "foo(1 + 2)[3] * obj.field - 4 ^ 2"
	expected := "((foo((1 + 2))[3] * obj.field) - (4 ^ 2))"
	testParser(t, input, expected)
}

// 1. Exponentiation precedence
func TestPrec_Exp1(t *testing.T) {
	input := "2 ^ 3 ^ 4"
	expected := "((2 ^ 3) ^ 4)"
	testParser(t, input, expected)
}

func TestPrec_Exp2(t *testing.T) {
	input := "(2 ^ 3) ^ 4"
	expected := "((2 ^ 3) ^ 4)"
	testParser(t, input, expected)
}

// 2. Unary minus with exponent
// TODO: Check the precedence once!
func TestPrec_UnaryMinusExp(t *testing.T) {
	input := "-2 ^ 3"
	expected := "((- 2) ^ 3)"
	testParser(t, input, expected)
}

func TestPrec_ExpUnaryMinus(t *testing.T) {
	input := "(-2) ^ 3"
	expected := "((- 2) ^ 3)"
	testParser(t, input, expected)
}

// 3. Mul / Div / Mod with exponent
func TestPrec_MulExp(t *testing.T) {
	input := "2 * 3 ^ 4"
	expected := "(2 * (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrec_DivExp(t *testing.T) {
	input := "8 / 2 ^ 3"
	expected := "(8 / (2 ^ 3))"
	testParser(t, input, expected)
}

func TestPrec_ModExp(t *testing.T) {
	input := "10 % 3 ^ 2"
	expected := "(10 % (3 ^ 2))"
	testParser(t, input, expected)
}

// 4. Mul / Div / Mod with addition / subtraction
func TestPrec_MulAdd(t *testing.T) {
	input := "2 * 3 + 4"
	expected := "((2 * 3) + 4)"
	testParser(t, input, expected)
}

func TestPrec_AddMul(t *testing.T) {
	input := "1 + 2 * 3"
	expected := "(1 + (2 * 3))"
	testParser(t, input, expected)
}

func TestPrec_DivSub(t *testing.T) {
	input := "8 / 4 - 2"
	expected := "((8 / 4) - 2)"
	testParser(t, input, expected)
}

func TestPrec_SubDiv(t *testing.T) {
	input := "8 - 4 / 2"
	expected := "(8 - (4 / 2))"
	testParser(t, input, expected)
}

// 5. Concatenation precedence
func TestPrec_Concat1(t *testing.T) {
	input := "1 .. 2 .. 3"
	expected := "((1 .. 2) .. 3)"
	testParser(t, input, expected)
}

func TestPrec_AddConcat(t *testing.T) {
	input := "1 + 2 .. 3"
	expected := "((1 + 2) .. 3)"
	testParser(t, input, expected)
}

func TestPrec_ConcatAdd(t *testing.T) {
	input := "1 .. 2 + 3"
	expected := "(1 .. (2 + 3))"
	testParser(t, input, expected)
}

func TestPrec_MulConcat(t *testing.T) {
	input := "2 * 3 .. 4"
	expected := "((2 * 3) .. 4)"
	testParser(t, input, expected)
}

func TestPrec_ConcatMul(t *testing.T) {
	input := "2 .. 3 * 4"
	expected := "(2 .. (3 * 4))"
	testParser(t, input, expected)
}

// 6. Comparison precedence
func TestPrec_CompAdd(t *testing.T) {
	input := "1 + 2 == 3"
	expected := "((1 + 2) == 3)"
	testParser(t, input, expected)
}

func TestPrec_AddComp(t *testing.T) {
	input := "1 == 2 + 3"
	expected := "(1 == (2 + 3))"
	testParser(t, input, expected)
}

func TestPrec_MulComp(t *testing.T) {
	input := "2 * 3 > 4"
	expected := "((2 * 3) > 4)"
	testParser(t, input, expected)
}

func TestPrec_CompMul(t *testing.T) {
	input := "2 > 3 * 4"
	expected := "(2 > (3 * 4))"
	testParser(t, input, expected)
}

// 7. Logical AND / OR precedence
func TestPrec_AndOr1(t *testing.T) {
	input := "1 and 2 or 3"
	expected := "((1 and 2) or 3)"
	testParser(t, input, expected)
}

func TestPrec_OrAnd1(t *testing.T) {
	input := "1 or 2 and 3"
	expected := "(1 or (2 and 3))"
	testParser(t, input, expected)
}

func TestPrec_CompAnd(t *testing.T) {
	input := "1 == 2 and 3"
	expected := "((1 == 2) and 3)"
	testParser(t, input, expected)
}

func TestPrec_AndComp(t *testing.T) {
	input := "1 and 2 == 3"
	expected := "(1 and (2 == 3))"
	testParser(t, input, expected)
}

// 8. Not operator precedence
func TestPrec_NotAdd(t *testing.T) {
	input := "not 1 + 2"
	expected := "((not 1) + 2)"
	testParser(t, input, expected)
}

func TestPrec_AddNot(t *testing.T) {
	input := "1 + not 2"
	expected := "(1 + (not 2))"
	testParser(t, input, expected)
}

func TestPrec_NotExp(t *testing.T) {
	input := "not 2 ^ 3"
	expected := "((not 2) ^ 3)"
	testParser(t, input, expected)
}

func TestPrec_ExpNot(t *testing.T) {
	input := "2 ^ not 3"
	expected := "(2 ^ (not 3))"
	testParser(t, input, expected)
}

// 9. Mixed multi-level
func TestPrec_Mixed1(t *testing.T) {
	input := "1 + 2 * 3 ^ 4"
	expected := "(1 + (2 * (3 ^ 4)))"
	testParser(t, input, expected)
}

func TestPrec_Mixed2(t *testing.T) {
	input := "(1 + 2) * 3 ^ 4"
	expected := "((1 + 2) * (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrec_Mixed3(t *testing.T) {
	input := "1 + 2 ^ 3 * 4"
	expected := "(1 + ((2 ^ 3) * 4))"
	testParser(t, input, expected)
}

func TestPrec_Mixed4(t *testing.T) {
	input := "1 .. 2 * 3 ^ 4"
	expected := "(1 .. (2 * (3 ^ 4)))"
	testParser(t, input, expected)
}

func TestPrec_Mixed5(t *testing.T) {
	input := "(1 .. 2) * 3 ^ 4"
	expected := "((1 .. 2) * (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrec_Mixed6(t *testing.T) {
	input := "1 + 2 .. 3 ^ 4"
	expected := "((1 + 2) .. (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrec_Mixed7(t *testing.T) {
	input := "(1 + 2) .. 3 ^ 4"
	expected := "((1 + 2) .. (3 ^ 4))"
	testParser(t, input, expected)
}

func TestPrec_Mixed8(t *testing.T) {
	input := "1 and 2 + 3"
	expected := "(1 and (2 + 3))"
	testParser(t, input, expected)
}

func TestPrec_Mixed9(t *testing.T) {
	input := "1 + 2 and 3"
	expected := "((1 + 2) and 3)"
	testParser(t, input, expected)
}

func TestPrec_Mixed10(t *testing.T) {
	input := "not 1 and 2"
	expected := "((not 1) and 2)"
	testParser(t, input, expected)
}

// 1. Function call with arithmetic
func TestPrec_FuncCallMul(t *testing.T) {
	input := "foo() * 2"
	expected := "(foo() * 2)"
	testParser(t, input, expected)
}

func TestPrec_MulFuncCall(t *testing.T) {
	input := "2 * foo()"
	expected := "(2 * foo())"
	testParser(t, input, expected)
}

func TestPrec_FuncCallExp(t *testing.T) {
	input := "foo() ^ 2"
	expected := "(foo() ^ 2)"
	testParser(t, input, expected)
}

func TestPrec_ExpFuncCall(t *testing.T) {
	input := "2 ^ foo()"
	expected := "(2 ^ foo())"
	testParser(t, input, expected)
}

// 2. Indexing with arithmetic
func TestPrec_IndexAdd(t *testing.T) {
	input := "arr[1] + 2"
	expected := "(arr[1] + 2)"
	testParser(t, input, expected)
}

func TestPrec_AddIndex(t *testing.T) {
	input := "2 + arr[1]"
	expected := "(2 + arr[1])"
	testParser(t, input, expected)
}

func TestPrec_IndexExp(t *testing.T) {
	input := "arr[1] ^ 2"
	expected := "(arr[1] ^ 2)"
	testParser(t, input, expected)
}

func TestPrec_ExpIndex(t *testing.T) {
	input := "2 ^ arr[1]"
	expected := "(2 ^ arr[1])"
	testParser(t, input, expected)
}

// 3. Table access with arithmetic
func TestPrec_TableAccessMul(t *testing.T) {
	input := "obj.field * 2"
	expected := "(obj.field * 2)"
	testParser(t, input, expected)
}

func TestPrec_MulTableAccess(t *testing.T) {
	input := "2 * obj.field"
	expected := "(2 * obj.field)"
	testParser(t, input, expected)
}

func TestPrec_TableAccessExp(t *testing.T) {
	input := "obj.field ^ 2"
	expected := "(obj.field ^ 2)"
	testParser(t, input, expected)
}

func TestPrec_ExpTableAccess(t *testing.T) {
	input := "2 ^ obj.field"
	expected := "(2 ^ obj.field)"
	testParser(t, input, expected)
}

// 4. Mixed indexing + function call
func TestPrec_IndexFuncCall(t *testing.T) {
	input := "arr[foo()] + 2"
	expected := "(arr[foo()] + 2)"
	testParser(t, input, expected)
}

func TestPrec_FuncCallIndex(t *testing.T) {
	input := "foo()[1] + 2"
	expected := "(foo()[1] + 2)"
	testParser(t, input, expected)
}

func TestPrec_IndexIndex(t *testing.T) {
	input := "arr1[arr2[1]] * 3"
	expected := "(arr1[arr2[1]] * 3)"
	testParser(t, input, expected)
}

// 5. Mixed table access + indexing
func TestPrec_TableAccessIndex(t *testing.T) {
	input := "obj.field[1] + 2"
	expected := "(obj.field[1] + 2)"
	testParser(t, input, expected)
}

func TestPrec_IndexTableAccess(t *testing.T) {
	input := "arr[1].field * 2"
	expected := "(arr[1].field * 2)"
	testParser(t, input, expected)
}

// 6. Function call + table access
func TestPrec_FuncCallTableAccess(t *testing.T) {
	input := "foo().bar + 2"
	expected := "(foo().bar + 2)"
	testParser(t, input, expected)
}

func TestPrec_TableAccessFuncCall(t *testing.T) {
	input := "obj.method(2) * 3"
	expected := "(obj.method(2) * 3)"
	testParser(t, input, expected)
}

// 7. Mixed precedence with calls, indexing, and table access
func TestPrec_Mixed11(t *testing.T) {
	input := "foo(1 + 2)[3] * obj.field"
	expected := "(foo((1 + 2))[3] * obj.field)"
	testParser(t, input, expected)
}

func TestPrec_Mixed12(t *testing.T) {
	input := "arr[foo(1 ^ 2)] - obj.field"
	expected := "(arr[foo((1 ^ 2))] - obj.field)"
	testParser(t, input, expected)
}

func TestPrec_Mixed13(t *testing.T) {
	input := "foo()[1] ^ obj.bar"
	expected := "(foo()[1] ^ obj.bar)"
	testParser(t, input, expected)
}

func TestPrec_Mixed14(t *testing.T) {
	input := "obj.method(1)[2] + foo().bar"
	expected := "(obj.method(1)[2] + foo().bar)"
	testParser(t, input, expected)
}

func TestPrec_Mixed15(t *testing.T) {
	input := "arr[1].field ^ foo(2 * 3)"
	expected := "(arr[1].field ^ foo((2 * 3)))"
	testParser(t, input, expected)
}

func TestPrec_Mixed16(t *testing.T) {
	input := "foo(1)[2].bar * arr[3]"
	expected := "(foo(1)[2].bar * arr[3])"
	testParser(t, input, expected)
}

func TestPrecedenceLteSimple(t *testing.T) {
	input := "2 <= 3"
	expected := "(2 <= 3)"
	testParser(t, input, expected)
}

func TestPrecedenceLteWithAdditionMultiplication(t *testing.T) {
	input := "(1 + 2) <= (3 * 4)"
	expected := "((1 + 2) <= (3 * 4))"
	testParser(t, input, expected)
}

func TestPrecedenceGteFunctionCalls(t *testing.T) {
	input := "foo() >= bar()"
	expected := "(foo() >= bar())"
	testParser(t, input, expected)
}

func TestPrecedenceLteIndexExpressions(t *testing.T) {
	input := "x[1] <= y[2]"
	expected := "(x[1] <= y[2])"
	testParser(t, input, expected)
}

func TestPrecedenceGteTableAccess(t *testing.T) {
	input := "tbl.key >= 42"
	expected := "(tbl.key >= 42)"
	testParser(t, input, expected)
}

func TestPrecedenceGteAdditionSubtraction(t *testing.T) {
	input := "(a + b) >= (c - d)"
	expected := "((a + b) >= (c - d))"
	testParser(t, input, expected)
}

func TestPrecedenceLteFunctionCallAndMultiplication(t *testing.T) {
	input := "some(1, 2) <= (3 * 4)"
	expected := "(some(1, 2) <= (3 * 4))"
	testParser(t, input, expected)
}

func TestPrecedenceGteConcatenation(t *testing.T) {
	input := "(x .. y) >= z"
	expected := "((x .. y) >= z)"
	testParser(t, input, expected)
}

func TestPrecedenceLteAndGteCombinedAnd(t *testing.T) {
	input := "1 <= 2 and 3 >= 4"
	expected := "((1 <= 2) and (3 >= 4))"
	testParser(t, input, expected)
}

func TestPrecedenceLteAndGteCombinedOr(t *testing.T) {
	input := "foo() <= bar() or baz >= qux"
	expected := "((foo() <= bar()) or (baz >= qux))"
	testParser(t, input, expected)
}

func TestLambdaExpression(t *testing.T) {
	input := `return fn() void then print("Something") end`

	testParser(t, input, input)
}

func TestLambdaExpressionReturn(t *testing.T) {
	input := `return fn(a int, b int) int then return (a + b) end`

	testParser(t, input, input)
}
func TestLambdaExpressionWithArgs(t *testing.T) {
	input := `return fn(a int, b int) int then return (a + b) end`

	testParser(t, input, input)
}

func TestPrefixExpressionNegation(t *testing.T) {
	input := `return (- 5)`

	testParser(t, input, input)
}

func TestPrefixExpressionNot(t *testing.T) {
	input := `return (not true)`

	testParser(t, input, input)
}

// TOOD: Complete this sometime.
//	func TestPrec_MegaAll(t *testing.T) {
//		input := `foo(1 + 2)[bar(3 ^ 2)].baz * arr[4 + 5] ^ obj.method(6 / 2)
//		          and x or y .. z < 100 == otherFunc()[idx].field`
//
//		expected := `((((foo((1 + 2)))[(bar((3 ^ 2)))]).baz) * ((arr[(4 + 5)]) ^ ((obj.method)((6 / 2)))))
//		              and ((x) or (((y) .. (z)) < (100)) == (((otherFunc())[idx]).field)))`
//
//		testParser(t, input, expected)
//	}
// func TestPrec_MegaAll(t *testing.T) {
// 	input := `foo(1 + 2)[bar(3 ^ 2)].baz * arr[4 + 5] ^ obj.method(6 / 2)
// 	          and x or y .. z < 100 == otherFunc()[idx].field`
//
// 	expected := `(foo(1 + 2)[bar((3 ^ 2))].baz * (arr[4 + 5] ^ obj.method((6 / 2))))
// 	              and (x or (((y .. z) < 100) == otherFunc()[idx].field))`
//
// 	testParser(t, input, expected)
// }
