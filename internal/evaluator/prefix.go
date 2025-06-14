package evaluator

import "sunbird/internal/object"

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)

	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() == object.IntegerObj {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	if right.Type() == object.FloatObj {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	if right.Type() != object.IntegerObj && right.Type() != object.FloatObj {
		return newError("unknown operator: -%s", right.Type())
	}

	return NULL
}
