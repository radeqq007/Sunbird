package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

func evalAssignment(name ast.Expression, val object.Object, env *object.Environment) object.Object {
	switch node := name.(type) {
	case *ast.Identifier:
		return evalIdentifierAssignment(node, val, env)

	case *ast.PropertyExpression:
		return evalPropertyExpressionAssignment(node, val, env)

	case *ast.IndexExpression:
		return evalIndexExpressionAssignment(node, val, env)

	default:
		return newError(0, 0, "invalid assignment target: %s", name.String())
	}
}

func evalIdentifierAssignment(node *ast.Identifier, val object.Object, env *object.Environment) object.Object {
	if _, ok := env.Get(node.Value); !ok {
		return newError(node.Token.Line, node.Token.Col, "Identifier '%s' has not been declared.", node.Value)
	}

	env.Set(node.Value, val)
	return val
}

func evalPropertyExpressionAssignment(node *ast.PropertyExpression, val object.Object, env *object.Environment) object.Object {
	obj := Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return newError(node.Token.Line, node.Token.Col, "only hash objects have properties, got %s", obj.Type())
	}

	key := &object.String{Value: node.Property.Value}
	hashKey := key.HashKey()

	pair := object.HashPair{Key: key, Value: val}
	hash.Pairs[hashKey] = pair

	return val
}

func evalIndexExpressionAssignment(node *ast.IndexExpression, val object.Object, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch obj := left.(type) {
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if !ok {
			return newError(node.Token.Line, node.Token.Col, "index must be an integer, got %s", index.Type())
		}

		if idx.Value < 0 || idx.Value >= int64(len(obj.Elements)) {
			return newError(node.Token.Line, node.Token.Col, "index out of bounds")
		}

		obj.Elements[idx.Value] = val
		return val

	case *object.Hash:
		key, ok := index.(object.Hashable)
		if !ok {
			return newError(node.Token.Line, node.Token.Col, "unusable as hash key: %s", index.Type())
		}

		pair := object.HashPair{Key: index, Value: val}
		obj.Pairs[key.HashKey()] = pair
		return val

	default:
		return newError(node.Token.Line, node.Token.Col, "index operator not supported: %s", left.Type())
	}
}
