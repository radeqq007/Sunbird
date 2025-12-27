package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ImportStatement:
		return evalImportStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, node.Token.Line, node.Token.Col)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Line, node.Token.Col)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarExpression:
		if _, ok := env.Get(node.Name.String()); ok {
			return NewError(node.Token.Line, node.Token.Col, "Identifier '%s' has already been declared.", node.Name.String())
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.String(), val)
		return val

	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return evalAssignment(node.Name, val, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.PropertyExpression:
		return evalPropertyExpression(node, env)

	case *ast.CallExpression:
		// Check if it's an object method call
		if propExpr, ok := node.Function.(*ast.PropertyExpression); ok {
			obj := Eval(propExpr.Object, env)
			if isError(obj) {
				return obj
			}

			hash, ok := obj.(*object.Hash)
			if ok {
				key := &object.String{Value: propExpr.Property.Value}
				method := evalHashIndexExpression(hash, key, node.Token.Line, node.Token.Col)
				if isError(method) {
					return method
				}

				args := evalExpressions(node.Arguments, env)
				if len(args) == 1 && isError(args[0]) {
					return args[0]
				}

				return evalMethodCall(hash, method, args, node.Token.Line, node.Token.Col)
			}
		}

		// Regular function call
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, node.Token.Line, node.Token.Col)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index, node.Token.Line, node.Token.Col)
	}
	return nil
}

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

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		switch obj := obj.(type) {
		case *object.String:
			return obj.Value != ""
		case *object.Integer:
			return obj.Value != 0
		case *object.Float:
			return obj.Value != 0.0
		default:
			return true
		}
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
