package evaluator

import "sunbird/object"

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
  
  default:
	 return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)

	if idx > max {
		return NULL
	}

	if idx < 0 {
		return array.Elements[max + 1 + idx]
	}

	return array.Elements[idx]
}