package evaluator

import (
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func applyFunction(fn object.Value, args []object.Value, line, col int) object.Value {
	switch fn.Kind() {
	case object.FunctionKind:
		fn := fn.AsFunction()
		err := errors.ExpectNumberOfArguments(line, col, len(fn.Parameters), args)
		if err.IsError() {
			return err
		}

		extendedEnv, err := extendFunctionEnv(fn, args)
		if err.IsError() {
			return err
		}

		// set 'this' binding
		if thisVal, ok := extendedEnv.Get("this"); ok {
			extendedEnv.Set("this", thisVal)
		}

		evaluated := Eval(fn.Body, extendedEnv)

		if isError(evaluated) {
			return evaluated
		}

		result := unwrapReturnValue(evaluated)

		// Type check return value
		if fn.ReturnType != nil {
			if typeErr := checkType(fn.ReturnType, result, line, col); typeErr.IsError() {
				return typeErr
			}
		}

		return result

	case object.BuiltinKind:
		return fn.AsBuiltin().Fn(args...)

	default:
		return errors.NewNotCallableError(line, col, fn)
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Value,
) (*object.Environment, object.Value) {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		if param.Type != nil {
			if err := checkType(param.Type, args[i], param.Token.Line, param.Token.Col); err.IsError() {
				return nil, err
			}
		}
		env.SetWithType(param.Value, args[i], param.Type)
	}

	return env, NULL
}

func unwrapReturnValue(obj object.Value) object.Value {
	if obj.Kind() == object.ReturnValueKind {
		return obj.AsReturnValue().Value
	}

	return obj
}

func init() {
	object.ApplyFunction = func(fn object.Value, args []object.Value) object.Value {
		return applyFunction(fn, args, 0, 0)
	}
}
