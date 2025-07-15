package evaluator

import "sunbird/internal/object"

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))

	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))

	case operator == "|>":
		return evalPipeExpression(left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.StringObj || right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right)

	case left.Type() == object.FloatObj || right.Type() == object.FloatObj:
		return evalFloatInfixExpression(operator, left, right)

	// TODO: this probably should be a different error
	case left.Type() != right.Type():
		return newError(
			"type mismatch: %s %s %s",
			left.Type().String(),
			operator,
			right.Type().String(),
		)

	default:
		return newError(
			"unknown operator: %s %s %s",
			left.Type().String(),
			operator,
			right.Type().String(),
		)
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}

		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("division by zero")
		}

		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	var leftVal, rightVal float64

	if left.Type() == object.IntegerObj {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.IntegerObj {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	default:
		return NULL
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" && operator != "==" && operator != "!=" && operator != "&&" &&
		operator != "||" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.Inspect()
	rightVal := right.Inspect()

	return &object.String{Value: leftVal + rightVal}
}

func evalPipeExpression(left, right object.Object) object.Object {
	switch fn := right.(type) {
	case *object.Function:
		return applyFunction(fn, []object.Object{left})

	case *object.Builtin:
		return fn.Fn(left)
	}

	return newError("right side of pipe operator is not a function: %s", right.Type())
}
