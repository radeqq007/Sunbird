package evaluator

import "sunbird/internal/object"

func applyFunction(fn object.Object, args []object.Object, line, col int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(line, col, "wrong number of arguments: expected %d, got %d", len(fn.Parameters), len(args))
		}

		extendedEnv := extendFunctionEnv(fn, args)

		// set 'this' binding
		if thisVal, ok := extendedEnv.Get("this"); ok {
			extendedEnv.Set("this", thisVal)
		}

		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError(line, col, "not a function: %s", fn.Type().String())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
