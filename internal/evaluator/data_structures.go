package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for _, pair := range node.Pairs {
		key := Eval(pair.Key, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return NewError(node.Token.Line, node.Token.Col, "unusable as hash key: %s", key.Type())
		}

		value := Eval(pair.Value, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object, line, col int) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index, line, col)

	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index, line, col)

	case left.Type() == object.StringObj:
		return evalStringIndexExpression(left, index, line, col)

	default:
		return NewError(line, col, "index operator not supported: %s", left.Type().String())
	}
}

func evalArrayIndexExpression(left, index object.Object, line, col int) object.Object {
	array, ok := left.(*object.Array)
	if !ok {
		return NewError(line, col, "index operator not supported: %s", left.Type())
	}

	idx := index.(*object.Integer).Value
	maxIdx := int64(len(array.Elements) - 1)

	if idx > maxIdx {
		return NewError(line, col, "index out of range: %d", idx)
	}

	if idx < 0 {
		return array.Elements[maxIdx+1+idx]
	}

	return array.Elements[idx]
}

func evalHashIndexExpression(left, index object.Object, line, col int) object.Object {
	hash, ok := left.(*object.Hash)
	if !ok {
		return NewError(line, col, "index operator not supported: %s", left.Type())
	}

	key, ok := index.(object.Hashable)
	if !ok {
		return NewError(line, col, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if ok {
		return pair.Value
	}

	if hash.Proto != nil {
		return evalHashIndexExpression(hash.Proto, index, line, col)
	}

	return NULL
}

func evalStringIndexExpression(left, index object.Object, line, col int) object.Object {
	str, ok := left.(*object.String)
	if !ok {
		return NewError(line, col, "index operator not supported: %s", left.Type())
	}

	idx := index.(*object.Integer).Value
	maxIdx := int64(len(str.Value) - 1)

	if idx > maxIdx {
		return NewError(line, col, "index out of range: %d", idx)
	}

	if idx < 0 {
		return &object.String{Value: string(str.Value[maxIdx+1+idx])}
	}

	return &object.String{Value: string(str.Value[idx])}
}

func evalPropertyExpression(pe *ast.PropertyExpression, env *object.Environment) object.Object {
	obj := Eval(pe.Object, env)
	if isError(obj) {
		return obj
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return NewError(pe.Token.Line, pe.Token.Col,
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
		return NewError(stmt.Token.Line, stmt.Token.Col,
			"identifier not found: %s", stmt.Object.String())
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return NewError(stmt.Token.Line, stmt.Token.Col,
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
