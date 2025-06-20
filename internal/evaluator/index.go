package evaluator

import "sunbird/internal/object"

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index)

	default:
		return newError("index operator not supported: %s", left.Type().String())
	}
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	array, ok := left.(*object.Array)
	if !ok {
		return newError("index operator not supported: %s", left.Type())
	}

	idx := index.(*object.Integer).Value
	maxIdx := int64(len(array.Elements) - 1)

	if idx > maxIdx {
		return NULL
	}

	if idx < 0 {
		return array.Elements[maxIdx+1+idx]
	}

	return array.Elements[idx]
}
