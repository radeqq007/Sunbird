package evaluator

import "sunbird/internal/object"

func evalIndexExpression(left, index object.Object, line, col int) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index, line, col)

	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index, line, col)

	default:
		return newError(line, col, "index operator not supported: %s", left.Type().String())
	}
}

func evalArrayIndexExpression(left, index object.Object, line, col int) object.Object {
	array, ok := left.(*object.Array)
	if !ok {
		return newError(line, col, "index operator not supported: %s", left.Type())
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

func evalHashIndexExpression(left, index object.Object, line, col int) object.Object {
	hash, ok := left.(*object.Hash)
	if !ok {
		return newError(line, col, "index operator not supported: %s", left.Type())
	}

	key, ok := index.(object.Hashable)
	if !ok {
		return newError(line, col, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if ok {
		return pair.Value
	}

	if hash.Proto != nil {
		return evalHashIndexExpression(hash.Proto, index, line, col)
	}

	return NULL
}
