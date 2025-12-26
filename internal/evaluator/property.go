package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

func evalPropertyExpression(pe *ast.PropertyExpression, env *object.Environment) object.Object {
	obj := Eval(pe.Object, env)
	if isError(obj) {
		return obj
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return newError(pe.Token.Line, pe.Token.Col,
			"property access on non-object: %s", obj.Type())
	}

	// Convert property name to string key
	key := &object.String{Value: pe.Property.Value}
	return evalHashIndexExpression(hash, key, pe.Token.Line, pe.Token.Col)
}

func evalPropertyAssignment(stmt *ast.PropertyAssignStatement, env *object.Environment) object.Object {
	obj := Eval(stmt.Object, env)
	if isError(obj) {
		return obj
	}

	if obj == nil {
		return newError(stmt.Token.Line, stmt.Token.Col,
			"identifier not found: %s", stmt.Object.String())
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return newError(stmt.Token.Line, stmt.Token.Col,
			"property assignment on non-object: %s", obj.Type())
	}

	value := Eval(stmt.Value, env)
	if isError(value) {
		return value
	}

	key := &object.String{Value: stmt.Property.Value}
	hashKey := key.HashKey()

	if hash.Pairs == nil {
		hash.Pairs = make(map[object.HashKey]object.HashPair)
	}

	hash.Pairs[hashKey] = object.HashPair{Key: key, Value: value}

	return value
}

func evalMethodCall(obj *object.Hash, method object.Object, args []object.Object, line, col int) object.Object {
	fn, ok := method.(*object.Function)
	if !ok {
		// Not a function, just call normally
		return applyFunction(method, args, line, col)
	}

	newEnv := object.NewEnclosedEnvironment(fn.Env)
	newEnv.Set("this", obj)

	boundFn := &object.Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		Env:        newEnv,
	}

	return applyFunction(boundFn, args, line, col)
}
