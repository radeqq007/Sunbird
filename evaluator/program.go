package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value

			case *object.Error:
				return result
		}
	}

	return result
}