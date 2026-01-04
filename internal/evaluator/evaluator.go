package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ImportStatement:
		return evalImportStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BreakStatement:
		return BREAK

	case *ast.ContinueStatement:
		return CONTINUE

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

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.TryCatchStatement:
		return evalTryCatchStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetExpression:
		if env.Has(node.Name.String()) {
			return errors.NewVariableReassignmentError(node.Token.Line, node.Token.Col, node.Name.String())
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Type != nil {
			if err := checkType(node.Type, val, node.Token.Line, node.Token.Col); err != nil {
				return err
			}
		}

		env.SetWithType(node.Name.String(), val, node.Type)
		return val

	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return evalAssignment(node.Name, val, env)

	case *ast.ConstExpression:
		if env.Has(node.Name.String()) {
			return errors.NewVariableReassignmentError(node.Token.Line, node.Token.Col, node.Name.String())
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Type != nil {
			if err := checkType(node.Type, val, node.Token.Line, node.Token.Col); err != nil {
				return err
			}
		}

		env.SetConstWithType(node.Name.String(), val, node.Type)
		return val

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, ReturnType: node.ReturnType, Body: body, Env: env}

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

	case *ast.RangeExpression:
		return evalRangeExpression(node, env)
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
			if result.Propagating {
				return result
			}
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

func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	start := Eval(node.Start, env)
	if isError(start) {
		return start
	}

	end := Eval(node.End, env)
	if isError(end) {
		return end
	}

	if start.Type() != object.IntegerObj || end.Type() != object.IntegerObj {
		return errors.NewTypeError(node.Token.Line, node.Token.Col, "range values must be integers (got %s and %s)", start.Type(), end.Type())
	}

	startVal := start.(*object.Integer).Value
	endVal := end.(*object.Integer).Value
	stepVal := int64(1)

	if node.Step != nil {
		step := Eval(node.Step, env)
		if isError(step) {
			return step
		}

		if step.Type() != object.IntegerObj {
			return errors.NewTypeError(node.Token.Line, node.Token.Col, "range step must be integer (got %s)", step.Type())
		}
		stepVal = step.(*object.Integer).Value
	}

	return &object.Range{Start: startVal, End: endVal, Step: stepVal}
}
