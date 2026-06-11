package evaluator

import (
	"github.com/radeqq007/sunbird/internal/errors"
	"github.com/radeqq007/sunbird/internal/object"
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

		return unwrapReturnValue(evaluated)

	case object.BuiltinKind:
		return fn.AsBuiltin().Fn(object.NewCallContext(line, col), args...)

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
		env.Set(param.Value, args[i])
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
