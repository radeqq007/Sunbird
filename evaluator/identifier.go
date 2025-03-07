package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}


	return newError("identifier not found: " + node.Value)
}