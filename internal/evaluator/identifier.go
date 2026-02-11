package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Value {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return errors.NewUndefinedVariableError(node.Token.Line, node.Token.Col, node.Value)
}
