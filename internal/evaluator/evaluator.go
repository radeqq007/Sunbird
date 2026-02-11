package evaluator

import (
	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var (
	NULL     = object.NewNull()
	TRUE     = object.NewBool(true)
	FALSE    = object.NewBool(false)
	BREAK    = object.NewBreak()
	CONTINUE = object.NewContinue()
)

func Eval(node ast.Node, env *object.Environment) object.Value {
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
		return object.NewString(node.Value)

	case *ast.IntegerLiteral:
		return object.NewInt(node.Value)

	case *ast.FloatLiteral:
		return object.NewFloat(node.Value)

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
		var val object.Value
		if node.ReturnValue != nil {
			val = Eval(node.ReturnValue, env)
			if isError(val) {
				return val
			}
		} else {
			val = NULL
		}
		return object.NewReturnValue(val)

	case *ast.LetExpression:
		if env.Has(node.Name.String()) {
			return errors.NewVariableReassignmentError(node.Token.Line, node.Token.Col, node.Name.String())
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Type != nil {
			if err := checkType(node.Type, val, node.Token.Line, node.Token.Col); err.IsError() {
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

	case *ast.CompoundAssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return evalCompoundAssignExpression(node, val, env)

	case *ast.ConstExpression:
		if env.Has(node.Name.String()) {
			return errors.NewVariableReassignmentError(node.Token.Line, node.Token.Col, node.Name.String())
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Type != nil {
			if err := checkType(node.Type, val, node.Token.Line, node.Token.Col); err.IsError() {
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
		return object.NewFunction(params, node.ReturnType, body, env)

	case *ast.PropertyExpression:
		return evalPropertyExpression(node, env)

	case *ast.CallExpression:
		// Check if it's an object method call
		propExpr, ok := node.Function.(*ast.PropertyExpression)
		if ok {

			obj := Eval(propExpr.Object, env)
			if isError(obj) {
				return obj
			}

			if obj.IsHash() {
				key := object.NewString(propExpr.Property.Value)
				method := evalHashIndexExpression(obj, key, node.Token.Line, node.Token.Col)
				if isError(method) {
					return method
				}

				args := evalExpressions(node.Arguments, env)
				if len(args) == 1 && isError(args[0]) {
					return args[0]
				}

				return evalMethodCall(obj, method, args, node.Token.Line, node.Token.Col)
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
		return object.NewArray(elements)

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
	return object.NewNull()
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Value {
	var result object.Value

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result.Kind() {
		case object.ReturnValueKind:
			return result.AsReturnValue().Value

		case object.ErrorKind:
			if result.AsError().Propagating {
				return result
			}
		}
	}

	return result
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Value {
	var result []object.Value

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Value{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func isTruthy(obj object.Value) bool {
	switch obj.Kind() {
	case object.NullKind:
		return false
	case object.BoolKind:
		return obj.AsBool()
	case object.StringKind:
		return obj.AsString().Value != ""
	case object.IntKind:
		return obj.AsInt() != 0
	case object.FloatKind:
		return obj.AsFloat() != 0.0
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) object.Value {
	if input {
		return TRUE
	}

	return FALSE
}

func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Value {
	start := Eval(node.Start, env)
	if isError(start) {
		return start
	}

	end := Eval(node.End, env)
	if isError(end) {
		return end
	}

	err := errors.ExpectType(node.Token.Line, node.Token.Col, start, object.IntKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(node.Token.Line, node.Token.Col, end, object.IntKind)
	if err.IsError() {
		return err
	}

	startVal := start.AsInt()

	endVal := end.AsInt()
	stepVal := int64(1)

	if node.Step != nil {
		step := Eval(node.Step, env)
		if isError(step) {
			return step
		}

		err = errors.ExpectType(node.Token.Line, node.Token.Col, step, object.IntKind)
		if err.IsError() {
			return err
		}

		stepVal = step.AsInt()
	}

	return object.NewRange(startVal, endVal, stepVal)
}
