package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
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
		return errors.NewInvalidAssignmentTargetError(0, 0, name.String())
	}
}

func evalIdentifierAssignment(node *ast.Identifier, val object.Object, env *object.Environment) object.Object {
	if env.Update(node.Value, val) {
		return val
	}

	return errors.NewUndefinedVariableError(node.Token.Line, node.Token.Col, node.Value)
}

func evalPropertyExpressionAssignment(node *ast.PropertyExpression, val object.Object, env *object.Environment) object.Object {
	obj := Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	hash, ok := obj.(*object.Hash)
	if !ok {
		return errors.NewNonObjectPropertyAccessError(node.Token.Line, node.Token.Col, obj)
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
			return errors.NewIndexNotSupportedError(node.Token.Line, node.Token.Col, obj)
		}

		if idx.Value < 0 || idx.Value >= int64(len(obj.Elements)) {
			return errors.NewIndexOutOfBoundsError(node.Token.Line, node.Token.Col, obj)
		}

		obj.Elements[idx.Value] = val
		return val

	case *object.Hash:
		key, ok := index.(object.Hashable)
		if !ok {
			return errors.NewUnusableAsHashKeyError(node.Token.Line, node.Token.Col, index)
		}

		pair := object.HashPair{Key: index, Value: val}
		obj.Pairs[key.HashKey()] = pair
		return val

	default:
		return errors.NewIndexNotSupportedError(node.Token.Line, node.Token.Col, left)
	}
}
