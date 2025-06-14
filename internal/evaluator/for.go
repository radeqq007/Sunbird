package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	// This sucks but uhhh I don't know how to do it better
	if assign, ok := fs.Init.(*ast.AssignStatement); ok {
		loopEnv.Set(assign.Name.Value, NULL)
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
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
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
