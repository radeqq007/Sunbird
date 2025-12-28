package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var moduleCache = NewModuleCache()

func evalImportStatement(stmt *ast.ImportStatement, env *object.Environment) object.Object {
	path := stmt.Path.Value

	module, err := moduleCache.loadModule(path, env)
	if err != nil {
		return errors.NewImportError(stmt.Token.Line, stmt.Token.Col, err.Error())
	}

	// Determine the name to bind to
	name := path
	if stmt.Alias != nil {
		name = stmt.Alias.Value
	}

	// Bind module to environment
	env.Set(name, module)

	return NULL
}
