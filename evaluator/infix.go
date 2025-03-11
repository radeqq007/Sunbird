package evaluator

import "sunbird/object"

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)


	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return  newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
			
	case "||":
		// TODO: Update isTruthy() to handle this
		return nativeBoolToBooleanObject(leftVal != 0 || rightVal != 0)
	
	case "&&":
		// TODO: Update isTruthy() to handle this
		return nativeBoolToBooleanObject(leftVal != 0 && rightVal != 0)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	var leftVal, rightVal float64

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
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
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
			
	case "||":
		// TODO: Update isTruthy() to handle this
		return nativeBoolToBooleanObject(leftVal != 0 || rightVal != 0)
	
	case "&&":
		// TODO: Update isTruthy() to handle this
		return nativeBoolToBooleanObject(leftVal != 0 && rightVal != 0)

	default:
		return NULL
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" && operator != "==" && operator != "!=" && operator != "&&" && operator != "||" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.Inspect()
	rightVal := right.Inspect()

	if operator == "==" {
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}

	if operator == "!=" {
		return nativeBoolToBooleanObject(leftVal != rightVal)
	}

	if operator == "&&" {
		return nativeBoolToBooleanObject(leftVal != "" && rightVal != "")
	}

	if operator == "||" {
		return nativeBoolToBooleanObject(leftVal != "" || rightVal != "")
	}
			
	return &object.String{Value: leftVal + rightVal}
}