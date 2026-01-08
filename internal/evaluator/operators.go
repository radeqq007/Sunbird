package evaluator

import (
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalPrefixExpression(operator string, right object.Object, line, col int) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)

	case "-":
		return evalMinusPrefixOperator(right, line, col)
	default:
		return errors.NewUnknownPrefixOperatorError(line, col, operator, right)
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

func evalMinusPrefixOperator(right object.Object, line, col int) object.Object {
	if right.Type() == object.IntegerObj {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	if right.Type() == object.FloatObj {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	if right.Type() != object.IntegerObj && right.Type() != object.FloatObj {
		return errors.NewUnknownPrefixOperatorError(line, col, "-", right)
	}

	return NULL
}

func evalInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	switch {
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))

	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))

	case operator == "==":
		if left.Type() == object.StringObj && right.Type() == object.StringObj {
			return nativeBoolToBooleanObject(
				left.(*object.String).Value == right.(*object.String).Value,
			)
		}
		return nativeBoolToBooleanObject(left.Inspect() == right.Inspect())

	case operator == "!=":
		if left.Type() == object.StringObj && right.Type() == object.StringObj {
			return nativeBoolToBooleanObject(
				left.(*object.String).Value != right.(*object.String).Value,
			)
		}
		return nativeBoolToBooleanObject(left.Inspect() != right.Inspect())

	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right, line, col)

	case left.Type() == object.StringObj || right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right, line, col)

	case left.Type() == object.FloatObj || right.Type() == object.FloatObj:
		return evalFloatInfixExpression(operator, left, right)

	// TODO: this probably should be a different error
	case left.Type() != right.Type():
		return errors.NewTypeMismatchError(line, col, left.Type(), operator, right.Type())

	default:
		return errors.NewUnknownOperatorError(line, col, left, operator, right)
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
	line, col int,
) object.Object {
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
			return errors.NewDivisionByZeroError(line, col)
		}

		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return errors.NewDivisionByZeroError(line, col)
		}

		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return errors.NewUnknownOperatorError(line, col, left, operator, right)
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
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return NULL
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
	line, col int,
) object.Object {
	if operator != "+" && operator != "==" && operator != "!=" && operator != "&&" &&
		operator != "||" {
		return errors.NewUnknownOperatorError(line, col, left, operator, right)
	}

	var leftVal, rightVal string

	if s, ok := left.(*object.String); ok {
		leftVal = s.Value
	} else {
		leftVal = left.Inspect()
	}

	if s, ok := right.(*object.String); ok {
		rightVal = s.Value
	} else {
		rightVal = right.Inspect()
	}

	return &object.String{Value: leftVal + rightVal}
}
