package evaluator

import (
	"sunbird/ast"
	"sunbird/object"
)

var (
	NULL  = &object.Null{} 
	TRUE 	= &object.Boolean{Value: true}
  FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}	

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		
		return evalInfixExpression(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		
		case *ast.Identifier:
			return evalIdentifier(node, env)

		case *ast.FunctionLiteral:
			params := node.Parameters
			body := node.Body
			return &object.Function{Parameters: params, Body: body, Env: env}
	}
	return nil
}
