package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value

			case *object.Error:
				return result
		}
	}

	return result
}