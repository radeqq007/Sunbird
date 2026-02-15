package evaluator

import (
	"fmt"
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

	case ast.Statement:
		return evalStatement(node, env)

	case ast.Expression:
		return evalExpression(node, env)
	}
	return NULL
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

func evalStatement(node ast.Statement, env *object.Environment) object.Value {
	switch stmt := node.(type) {
	case *ast.ImportStatement:
		return evalImportStatement(stmt, env)

	case *ast.ExportStatement:
		return evalExportStatement(stmt, env)

	case *ast.ExpressionStatement:
		return Eval(stmt.Expression, env)

	case *ast.BreakStatement:
		return BREAK

	case *ast.ContinueStatement:
		return CONTINUE

	case *ast.BlockStatement:
		return evalBlockStatement(stmt, env)

	case *ast.ForStatement:
		return evalForStatement(stmt, env)

	case *ast.WhileStatement:
		return evalWhileStatement(stmt, env)

	case *ast.TryCatchStatement:
		return evalTryCatchStatement(stmt, env)

	case *ast.ReturnStatement:
		var val object.Value
		if stmt.ReturnValue != nil {
			val = Eval(stmt.ReturnValue, env)
			if isError(val) {
				return val
			}
		} else {
			val = NULL
		}
		return object.NewReturnValue(val)
	}

	return NULL
}

func evalExpression(node ast.Expression, env *object.Environment) object.Value {
	switch exp := node.(type) {
	case *ast.InfixExpression:
		left := Eval(exp.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(exp.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(exp.Operator, left, right, exp.Token.Line, exp.Token.Col)

	case *ast.PrefixExpression:
		right := Eval(exp.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(exp.Operator, right, exp.Token.Line, exp.Token.Col)

	case *ast.IfExpression:
		return evalIfExpression(exp, env)

	case *ast.LetExpression:
		return evalLetExpression(exp, env)

	case *ast.AssignExpression:
		val := Eval(exp.Value, env)
		if isError(val) {
			return val
		}

		return evalAssignment(exp.Name, val, env)

	case *ast.CompoundAssignExpression:
		val := Eval(exp.Value, env)
		if isError(val) {
			return val
		}

		return evalCompoundAssignExpression(exp, val, env)

	case *ast.ConstExpression:
		return evalConstExpression(exp, env)

	case *ast.PropertyExpression:
		return evalPropertyExpression(exp, env)

	case *ast.CallExpression:
		if _, ok := exp.Function.(*ast.PropertyExpression); ok {
			return evalMethodCallExpression(exp, env)
		}

		// Regular function call
		function := Eval(exp.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(exp.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, exp.Token.Line, exp.Token.Col)

	case *ast.IndexExpression:
		left := Eval(exp.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(exp.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index, exp.Token.Line, exp.Token.Col)

	case *ast.RangeExpression:
		return evalRangeExpression(exp, env)

	case *ast.StringLiteral:
		return object.NewString(exp.Value)

	case *ast.IntegerLiteral:
		return object.NewInt(exp.Value)

	case *ast.FloatLiteral:
		return object.NewFloat(exp.Value)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(exp.Value)

	case *ast.NullLiteral:
		return NULL

	case *ast.Identifier:
		return evalIdentifier(exp, env)

	case *ast.FunctionLiteral:
		params := exp.Parameters
		body := exp.Body
		return object.NewFunction(params, exp.ReturnType, body, env)

	case *ast.HashLiteral:
		return evalHashLiteral(exp, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(exp.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return object.NewArray(elements)
	}
	return NULL
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Value {
	result := make([]object.Value, 0, len(exps))

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

func evalLetExpression(exp *ast.LetExpression, env *object.Environment) object.Value {
	if env.Has(exp.Name.String()) {
		return errors.NewVariableReassignmentError(exp.Token.Line, exp.Token.Col, exp.Name.String())
	}

	val := Eval(exp.Value, env)
	if isError(val) {
		return val
	}

	if exp.Type != nil {
		if err := checkType(exp.Type, val, exp.Token.Line, exp.Token.Col); err.IsError() {
			return err
		}
	}

	if val.IsArray() {
		fmt.Printf("DEBUG: Let %s ptr=%p\n", exp.Name.String(), val.AsArray())
	}

	env.SetWithType(exp.Name.String(), val, exp.Type)
	return val
}

func evalConstExpression(exp *ast.ConstExpression, env *object.Environment) object.Value {
	if env.Has(exp.Name.String()) {
		return errors.NewVariableReassignmentError(exp.Token.Line, exp.Token.Col, exp.Name.String())
	}

	val := Eval(exp.Value, env)
	if isError(val) {
		return val
	}

	if exp.Type != nil {
		if err := checkType(exp.Type, val, exp.Token.Line, exp.Token.Col); err.IsError() {
			return err
		}
	}

	env.SetConstWithType(exp.Name.String(), val, exp.Type)
	return val
}
