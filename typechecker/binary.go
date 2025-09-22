package typechecker

import (
	"fmt"

	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/types"
)

type typeResolver func(left, right *types.Type) (*types.Type, error)

var binaryResolvers = map[token.TokenType]typeResolver{
	token.PLUS:     resolvePlus,
	token.MINUS:    resolveArithmetic,
	token.MULTIPLY: resolveArithmetic,
	token.SLASH:    resolveArithmetic,
	token.EQ:       resolveComparison,
	token.NEQ:      resolveComparison,
	token.LT:       resolveComparison,
	token.LTE:      resolveComparison,
	token.GT:       resolveComparison,
	token.GTE:      resolveComparison,
	token.AND:      resolveLogical,
	token.OR:       resolveLogical,
}

func resolvePlus(left, right *types.Type) (*types.Type, error) {
	if left == types.StringType || right == types.StringType {
		if left == types.StringType && right == types.StringType {
			return types.StringType, nil
		}
		return nil, fmt.Errorf("Cannot add '%s' and '%s'", left, right)
	}

	return resolveArithmetic(left, right)
}

func resolveArithmetic(left, right *types.Type) (*types.Type, error) {
	if (left == types.IntType || left == types.FloatType) && (right == types.FloatType || right == types.IntType) {
		if left == types.FloatType || right == types.FloatType {
			return types.FloatType, nil
		}

		return types.IntType, nil
	}

	return nil, fmt.Errorf("invalid operands for arithmetic: %s, %s", left, right)
}
func resolveComparison(left, right *types.Type) (*types.Type, error) {

	if left == right {
		return types.BoolType, nil
	}

	if (left == types.IntType || left == types.FloatType) && (right == types.IntType || right == types.FloatType) {
		return types.BoolType, nil
	}
	return nil, fmt.Errorf("invalid operands for comparison: %s, %s", left, right)
}

func resolveLogical(left, right *types.Type) (*types.Type, error) {
	if left == types.BoolType && right == types.BoolType {
		return types.BoolType, nil
	}

	return nil, fmt.Errorf("logical operators require bool, got %s and %s", left, right)
}
