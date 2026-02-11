package evaluator

import (
	"path/filepath"
	"strings"
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var moduleCache = NewModuleCache()

func evalImportStatement(stmt *ast.ImportStatement, env *object.Environment) object.Value {
	path := stmt.Path.Value

	module, err := moduleCache.loadModule(path)
	if err != nil {
		return errors.NewImportError(stmt.Token.Line, stmt.Token.Col, err.Error())
	}

	// Determine the name to bind to
	name := path
	if stmt.Alias != nil {
		name = stmt.Alias.Value
	} else {
		// If it's a file path, extract filename without extension
		base := filepath.Base(path)
		ext := filepath.Ext(base)

		if ext != "" {
			name = strings.TrimSuffix(base, ext)
		}
	}

	// Bind module to environment
	env.SetConst(name, module)

	return NULL
}
