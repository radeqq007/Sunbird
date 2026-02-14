package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalAssignment(name ast.Expression, val object.Value, env *object.Environment) object.Value {
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

func evalIdentifierAssignment(
	node *ast.Identifier,
	val object.Value,
	env *object.Environment,
) object.Value {
	if env.IsConst(node.Value) {
		return errors.NewConstantReassignmentError(node.Token.Line, node.Token.Col, node.Value)
	}

	if t, ok := env.GetType(node.Value); ok && t != nil {
		if err := checkType(t, val, node.Token.Line, node.Token.Col); err.IsError() {
			return err
		}
	}

	if env.Update(node.Value, val) {
		return val
	}

	return errors.NewUndefinedVariableError(node.Token.Line, node.Token.Col, node.Value)
}

func evalPropertyExpressionAssignment(
	node *ast.PropertyExpression,
	val object.Value,
	env *object.Environment,
) object.Value {
	obj := Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	if !obj.IsHash() {
		return errors.NewNonObjectPropertyAccessError(node.Token.Line, node.Token.Col, obj)
	}

	hash := obj.AsHash()

	key := object.NewString(node.Property.Value)
	hashKey := key.HashKey()

	pair := object.HashPair{Key: key, Value: val}
	hash.Pairs[hashKey] = pair

	return val
}

func evalIndexExpressionAssignment(
	node *ast.IndexExpression,
	val object.Value,
	env *object.Environment,
) object.Value {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	kind := left.Kind()
	switch kind {
	case object.ArrayKind:
		obj := left.AsArray()

		if !index.IsInt() {
			return errors.NewIndexNotSupportedError(node.Token.Line, node.Token.Col, left)
		}

		idx := index.AsInt()

		if idx < 0 || idx >= int64(len(obj.Elements)) {
			return errors.NewIndexOutOfBoundsError(node.Token.Line, node.Token.Col, left)
		}

		obj.Elements[idx] = val
		return val

	case object.HashKind:
		obj := left.AsHash()

		switch index.Kind() {
		case object.IntKind, object.StringKind:
			// Hashable
		default:
			return errors.NewUnusableAsHashKeyError(
				node.Token.Line,
				node.Token.Col,
				index,
			)
		}

		hashKey := index.HashKey()

		pair := object.NewHashPair(index, val)

		obj.Pairs[hashKey] = pair
		return val

	default:
		return errors.NewIndexNotSupportedError(node.Token.Line, node.Token.Col, left)
	}
}

func evalCompoundAssignExpression(
	node *ast.CompoundAssignExpression,
	val object.Value,
	env *object.Environment,
) object.Value {
	currentVal := Eval(node.Name, env)
	if isError(currentVal) {
		return currentVal
	}

	newVal := evalInfixExpression(node.Operator, currentVal, val, node.Token.Line, node.Token.Col)
	if isError(newVal) {
		return newVal
	}

	return evalAssignment(node.Name, newVal, env)
}
