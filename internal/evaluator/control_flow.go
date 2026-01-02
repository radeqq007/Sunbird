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

		if result != nil {
			switch result.Type() {
			case object.ReturnValueObj:
				return result
			case object.BreakObj:
				return NULL
			case object.ContinueObj:
				// Fallthrough to update
			}
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

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	var result object.Object = NULL

	for {
		condition := Eval(ws.Condition, loopEnv)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ws.Body, loopEnv)
		if isError(result) {
			return result
		}

		if result != nil {
			switch result.Type() {
			case object.ReturnValueObj:
				return result
			case object.BreakObj:
				return NULL
			case object.ContinueObj:
				// Fallthrough to update
			}
		}
	}

	return result
}

func evalTryCatchStatement(tcs *ast.TryCatchStatement, env *object.Environment) object.Object {
	tryResult := Eval(tcs.Try, env)

	result := tryResult

	if isError(tryResult) {
		errObj := tryResult.(*object.Error)
		errObj.Propagating = false

		catchEnv := object.NewEnclosedEnvironment(env)
		catchEnv.Set(tcs.Param.Value, tryResult)

		result = Eval(tcs.Catch, catchEnv)
	}

	if tcs.Finally != nil {
		finallyResult := Eval(tcs.Finally, env)
		if isError(finallyResult) {
			return finallyResult
		}

		// If finally produces a control flow change (Return, Break, Continue),
		// it overrides the result from try/catch.
		if finallyResult != nil {
			rt := finallyResult.Type()
			if rt == object.ReturnValueObj || rt == object.BreakObj || rt == object.ContinueObj {
				return finallyResult
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	blockEnv := object.NewEnclosedEnvironment(env)

	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj || isError(result) || rt == object.BreakObj || rt == object.ContinueObj {
				return result
			}
		}
	}

	return result
}
