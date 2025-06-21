package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	switch {
	case isTruthy(condition):
		return Eval(ie.Consequence, env)
	case ie.Alternative != nil:
		return Eval(ie.Alternative, env)
	default:
		return NULL
	}
}
