package evaluator

import (
	"github.com/radeqq007/sunbird/internal/ast"
	"github.com/radeqq007/sunbird/internal/errors"
	"github.com/radeqq007/sunbird/internal/object"
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
