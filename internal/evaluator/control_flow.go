package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Value {
	ifEnv := object.NewEnclosedEnvironment(env)
	condition := Eval(ie.Condition, ifEnv)
	if isError(condition) {
		return condition
	}

	switch {
	case isTruthy(condition):
		return Eval(ie.Consequence, ifEnv)
	case ie.Alternative != nil:
		return Eval(ie.Alternative, ifEnv)
	default:
		return NULL
	}
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Value {
	loopEnv := object.NewEnclosedEnvironment(env)

	iterable := Eval(fs.Iterable, loopEnv)
	if isError(iterable) {
		return iterable
	}

	switch iterable.Kind() {
	case object.RangeKind:
		return evalRangeLoop(fs, iterable.AsRange(), loopEnv)
	case object.ArrayKind:
		return evalArrayLoop(fs, iterable.AsArray(), loopEnv)
	case object.StringKind:
		return evalStringLoop(fs, iterable.AsString(), loopEnv)

	default:
		return errors.NewTypeError(fs.Token.Line, fs.Token.Col, "cannot iterate over %s", iterable.Kind().String())
	}
}

func evalRangeLoop(fs *ast.ForStatement, iterable *object.Range, env *object.Environment) object.Value {
	step := iterable.Step
	if step == 0 {
		step = 1
	}

	if step > 0 {
		for i := iterable.Start; i < iterable.End; i += step {
			env.Set(fs.Variable.Value, object.NewInt(i))

			result := Eval(fs.Body, env)
			if isError(result) {
				return result
			}

			if !result.IsNull() {
				switch result.Kind() {
				case object.ReturnValueKind:
					return result
				case object.BreakKind:
					return NULL
				case object.ContinueKind:
					continue
				}
			}
		}
	} else {
		for i := iterable.Start; i > iterable.End; i += step {
			env.Set(fs.Variable.Value, object.NewInt(i))

			result := Eval(fs.Body, env)
			if isError(result) {
				return result
			}

			if !result.IsNull() {
				switch result.Kind() {
				case object.ReturnValueKind:
					return result
				case object.BreakKind:
					return NULL
				case object.ContinueKind:
					continue
				}
			}
		}
	}
	return NULL
}

func evalArrayLoop(fs *ast.ForStatement, iterable *object.Array, env *object.Environment) object.Value {
	for _, element := range iterable.Elements {
		env.Set(fs.Variable.Value, element)

		result := Eval(fs.Body, env)
		if isError(result) {
			return result
		}

		if !result.IsNull() {
			switch result.Kind() {
			case object.ReturnValueKind:
				return result
			case object.BreakKind:
				return NULL
			case object.ContinueKind:
				continue
			}
		}
	}

	return NULL
}

func evalStringLoop(fs *ast.ForStatement, iterable *object.String, env *object.Environment) object.Value {
	for _, element := range iterable.Value {
		env.Set(fs.Variable.Value, object.NewString(string(element)))

		result := Eval(fs.Body, env)
		if isError(result) {
			return result
		}

		if !result.IsNull() {
			switch result.Kind() {
			case object.ReturnValueKind:
				return result
			case object.BreakKind:
				return NULL
			case object.ContinueKind:
				continue
			}
		}
	}

	return NULL
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Value {
	result := NULL

	for {
		loopEnv := object.NewEnclosedEnvironment(env)
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

		if !result.IsNull() {
			switch result.Kind() {
			case object.ReturnValueKind:
				return result
			case object.BreakKind:
				return NULL
			case object.ContinueKind:
				// Fallthrough to update
			}
		}
	}

	return result
}

func evalTryCatchStatement(tcs *ast.TryCatchStatement, env *object.Environment) object.Value {
	tryResult := Eval(tcs.Try, env)

	result := tryResult

	if isError(tryResult) {
		err := tryResult.AsError()
		caughtError := object.NewError(err.Message, err.Line, err.Col, false)

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
		if !finallyResult.IsNull() {
			kind := finallyResult.Kind()
			if kind == object.ReturnValueKind || kind == object.BreakKind || kind == object.ContinueKind {
				return finallyResult
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Value {
	blockEnv := object.NewEnclosedEnvironment(env)

	var result object.Value

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if !result.IsNull() {
			kind := result.Kind()
			if kind == object.ReturnValueKind || isError(result) || kind == object.BreakKind || kind == object.ContinueKind {
				return result
			}
		}
	}

	return result
}
