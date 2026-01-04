package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
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

	iterable := Eval(fs.Iterable, loopEnv)
	if isError(iterable) {
		return iterable
	}

	var result object.Object = NULL
	switch iter := iterable.(type) {
	case *object.Range:
		step := iter.Step
		if step == 0 {
			step = 1
		}

		if step > 0 {
			for i := iter.Start; i < iter.End; i += step {
				loopEnv.Set(fs.Variable.Value, &object.Integer{Value: i})

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
						continue
					}
				}
			}
		} else {
			for i := iter.Start; i > iter.End; i += step {
				loopEnv.Set(fs.Variable.Value, &object.Integer{Value: i})

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
						continue
					}
				}
			}
		}
	case *object.Array:
		for _, element := range iter.Elements {
			loopEnv.Set(fs.Variable.Value, element)

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
					continue
				}
			}
		}
	case *object.String:
		for _, element := range iter.Value {
			loopEnv.Set(fs.Variable.Value, &object.String{Value: string(element)})

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
					continue
				}
			}
		}
	default:
		return errors.NewTypeError(fs.Token.Line, fs.Token.Col, "cannot iterate over %s", iterable.Type().String())
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

		caughtError := &object.Error{
			Message:     errObj.Message,
			Line:        errObj.Line,
			Col:         errObj.Col,
			Propagating: false,
		}

		catchEnv := object.NewEnclosedEnvironment(env)
		catchEnv.Set(tcs.Param.Value, caughtError)

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
