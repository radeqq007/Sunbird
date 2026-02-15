package evaluator

import (
	"fmt"
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Value {
	pairs := make(map[object.HashKey]object.HashPair) // ! take a look at this

	for _, pair := range node.Pairs {
		key := Eval(pair.Key, env)
		if isError(key) {
			return key
		}

		value := Eval(pair.Value, env)
		if isError(value) {
			return value
		}

		hashed := key.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return object.NewHash(pairs)
}

func evalIndexExpression(left, index object.Value, line, col int) object.Value {
	switch {
	case left.IsArray() && index.IsInt():
		return evalArrayIndexExpression(left, index, line, col)

	case left.IsHash():
		return evalHashIndexExpression(left, index, line, col)

	case left.IsString() && index.IsInt():
		return evalStringIndexExpression(left, index, line, col)

	default:
		return errors.NewIndexNotSupportedError(line, col, left)
	}
}

func evalArrayIndexExpression(left, index object.Value, line, col int) object.Value {
	if !left.IsArray() {
		return errors.NewIndexNotSupportedError(line, col, left)
	}

	if !index.IsInt() {
		return errors.NewIndexNotSupportedError(line, col, index)
	}

	array := left.AsArray()
	idx := index.AsInt()

	maxIdx := int64(len(array.Elements) - 1)

	if idx > maxIdx {
		return errors.NewIndexOutOfBoundsError(line, col, left)
	}

	if idx < 0 {
		return array.Elements[maxIdx+1+idx]
	}

	return array.Elements[idx]
}

func evalHashIndexExpression(left, index object.Value, line, col int) object.Value {
	if !left.IsHash() {
		return errors.NewIndexNotSupportedError(line, col, left)
	}

	hash := left.AsHash()

	pair, ok := hash.Pairs[index.HashKey()]
	if ok {
		return pair.Value
	}

	if hash.Proto != nil {
		return evalHashIndexExpression(left, index, line, col)
	}

	return NULL
}

func evalStringIndexExpression(left, index object.Value, line, col int) object.Value {
	if !left.IsString() {
		return errors.NewIndexNotSupportedError(line, col, left)
	}

	// idxObj, _ := index.(*object.Integer)
	// TODO: check if it's an int
	idx := index.AsInt()
	str := left.AsString().Value

	maxIdx := int64(len(str) - 1)

	if idx > maxIdx {
		return errors.NewIndexOutOfBoundsError(line, col, left)
	}

	if idx < 0 {
		return object.NewString(string(str[maxIdx+1+idx]))
	}

	return object.NewString(string(str[idx]))
}

func evalPropertyExpression(pe *ast.PropertyExpression, env *object.Environment) object.Value {
	obj := Eval(pe.Object, env)
	if isError(obj) {
		return obj
	}

	if obj.IsModule() {
		module := obj.AsModule()
		propertyName := pe.Property.Value

		// Look up the property in the modules exports
		if val, ok := module.Exports[propertyName]; ok {
			return val
		}

		// Property not found in module
		return errors.NewUndefinedVariableError(
			pe.Token.Line,
			pe.Token.Col,
			fmt.Sprintf("%s.%s", module.Name, propertyName),
		)
	}

	if !obj.IsHash() {
		return errors.NewNonObjectPropertyAccessError(pe.Token.Line, pe.Token.Col, obj)
	}

	key := object.NewString(pe.Property.Value)
	return evalHashIndexExpression(obj, key, pe.Token.Line, pe.Token.Col)
}

func evalMethodCall(
	obj object.Value,
	method object.Value,
	args []object.Value,
	line, col int,
) object.Value {
	if !method.IsFunction() {
		// Not a function, just call normally
		return applyFunction(method, args, line, col)
	}

	fn := method.AsFunction()

	newEnv := object.NewEnclosedEnvironment(fn.Env)
	newEnv.Set("this", obj)

	boundFn := object.NewFunction(fn.Parameters, fn.ReturnType, fn.Body, fn.Env)

	return applyFunction(boundFn, args, line, col)
}

func evalMethodCallExpression(exp *ast.CallExpression, env *object.Environment) object.Value {
	propExp, _ := exp.Function.(*ast.PropertyExpression)
	obj := Eval(propExp.Object, env)
	if isError(obj) {
		return obj
	}

	if obj.IsHash() {
		key := object.NewString(propExp.Property.Value)
		method := evalHashIndexExpression(obj, key, exp.Token.Line, exp.Token.Col)
		if isError(method) {
			return method
		}

		args := evalExpressions(exp.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return evalMethodCall(obj, method, args, exp.Token.Line, exp.Token.Col)
	}

	return NULL
}
