package evaluator

import (
	"vex/ast"
	"vex/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	}
	
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	
	for _, statement := range stmts {
		result = Eval(statement)
	}
	
	return result
}