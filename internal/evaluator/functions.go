package evaluator

import (
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func applyFunction(fn object.Object, args []object.Object, line, col int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		err := errors.ExpectNumberOfArguments(line, col, len(fn.Parameters), args)
		if err != nil {
			return err
		}

		extendedEnv, err := extendFunctionEnv(fn, args)
		if err != nil {
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
			if typeErr := checkType(fn.ReturnType, result, line, col); typeErr != nil {
				return typeErr
			}
		}

		return result

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return errors.NewNotCallableError(line, col, fn)
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		if param.Type != nil {
			if err := checkType(param.Type, args[i], param.Token.Line, param.Token.Col); err != nil {
				return nil, err
			}
		}
		env.SetWithType(param.Value, args[i], param.Type)
	}

	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func init() {
	object.ApplyFunction = func(fn object.Object, args []object.Object) object.Object {
		return applyFunction(fn, args, 0, 0)
	}
}
