package evaluator

import (
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalPrefixExpression(operator string, right object.Value, line, col int) object.Value {
	switch operator {
	case "!":
		return evalBangOperator(right)

	case "-":
		return evalMinusPrefixOperator(right, line, col)
	default:
		return errors.NewUnknownPrefixOperatorError(line, col, operator, right)
	}
}

func evalBangOperator(right object.Value) object.Value {
	switch right {
	case TRUE:
		return FALSE

	case FALSE, NULL:
		return TRUE

	default:
		return FALSE
	}
}

func evalMinusPrefixOperator(right object.Value, line, col int) object.Value {
	if right.IsInt() {
		value := right.AsInt()
		return object.NewInt(-value)
	}

	if right.IsFloat() {
		value := right.AsFloat()
		return object.NewFloat(-value)
	}

	return errors.NewUnknownPrefixOperatorError(line, col, "-", right)
}

func evalInfixExpression(operator string, left, right object.Value, line, col int) object.Value {
	switch {
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))

	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))

	case operator == "==":
		if left.IsString() && right.IsString() {
			return nativeBoolToBooleanObject(
				left.AsString().Value == right.AsString().Value,
			)
		}
		return nativeBoolToBooleanObject(left.Inspect() == right.Inspect())

	case operator == "!=":
		if left.IsString() && right.IsString() {
			return nativeBoolToBooleanObject(
				left.AsString().Value != right.AsString().Value,
			)
		}
		return nativeBoolToBooleanObject(left.Inspect() != right.Inspect())

	case left.IsInt() && right.IsInt():
		return evalIntegerInfixExpression(operator, left, right, line, col)

	case left.IsString() || right.IsString():
		return evalStringInfixExpression(operator, left, right, line, col)

	case left.IsFloat() || right.IsFloat():
		return evalFloatInfixExpression(operator, left, right)

	// TODO: this probably should be a different error
	case left.Kind() != right.Kind():
		return errors.NewTypeMismatchError(line, col, left.Kind(), operator, right.Kind())

	default:
		return errors.NewUnknownOperatorError(line, col, left, operator, right)
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Value,
	line, col int,
) object.Value {
	leftVal := left.AsInt()
	rightVal := right.AsInt()

	switch operator {
	case "+":
		return object.NewInt(leftVal + rightVal)
	case "-":
		return object.NewInt(leftVal - rightVal)
	case "*":
		return object.NewInt(leftVal * rightVal)
	case "/":
		if rightVal == 0 {
			return errors.NewDivisionByZeroError(line, col)
		}

		return object.NewInt(leftVal / rightVal)
	case "%":
		if rightVal == 0 {
			return errors.NewDivisionByZeroError(line, col)
		}

		return object.NewInt(leftVal % rightVal)
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

func evalFloatInfixExpression(operator string, left, right object.Value) object.Value {
	var leftVal, rightVal float64

	if left.IsInt() {
		leftVal = float64(left.AsInt())
	} else {
		leftVal = left.AsFloat()
	}

	if right.IsInt() {
		rightVal = float64(right.AsInt())
	} else {
		rightVal = right.AsFloat()
	}

	switch operator {
	case "+":
		return object.NewFloat(leftVal + rightVal)
	case "-":
		return object.NewFloat(leftVal - rightVal)
	case "*":
		return object.NewFloat(leftVal * rightVal)
	case "/":
		return object.NewFloat(leftVal / rightVal)
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
	left, right object.Value,
	line, col int,
) object.Value {
	if operator != "+" && operator != "==" && operator != "!=" && operator != "&&" &&
		operator != "||" {
		return errors.NewUnknownOperatorError(line, col, left, operator, right)
	}

	var leftVal, rightVal string

	if left.IsString() {
		leftVal = left.AsString().Value
	} else {
		leftVal = left.Inspect()
	}

	if right.IsString() {
		rightVal = right.AsString().Value
	} else {
		rightVal = right.Inspect()
	}

	return object.NewString(leftVal + rightVal)
}
