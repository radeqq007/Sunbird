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

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	// This sucks but uhhh I don't know how to do it better
	if assign, ok := fs.Init.(*ast.AssignExpression); ok {
		loopEnv.Set(assign.Name.String(), NULL)
	}

	if fs.Init != nil {
		initResult := Eval(fs.Init, loopEnv)
		if isError(initResult) {
			return initResult
		}
	}

	var result object.Object = NULL

	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, loopEnv)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fs.Body, loopEnv)
		if isError(result) {
			return result
		}

		// handle return statements
		if result != nil && result.Type() == object.ReturnValueObj {
			return result
		}

		if fs.Update != nil {
			updateResult := Eval(fs.Update, loopEnv)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj || rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}
