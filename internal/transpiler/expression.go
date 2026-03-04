package transpiler

import (
	"fmt"
	"strings"
	"sunbird/internal/ast"
)

func (t *Transpiler) transpileExpression(node ast.Expression) (string, error) {
	switch exp := node.(type) {
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", exp.Value), nil
	
	case *ast.FloatLiteral:
		return fmt.Sprintf("%f", exp.Value), nil
	
	case *ast.StringLiteral:
		return "\"" + exp.Value + "\"", nil
	
	case *ast.Boolean:
		if exp.Value {
			return "true", nil
		}
		return "false", nil

	case *ast.NullLiteral:
		return "null", nil

	case *ast.PrefixExpression:
		right, err := t.transpileExpression(exp.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s%s)", exp.Operator, right), nil

	case *ast.InfixExpression:
		return t.transpileInfix(exp)
	
	case *ast.FunctionLiteral:
		return t.transpileFunction(exp)

	case *ast.ConstExpression:
		return t.transpileConst(exp)

	case *ast.LetExpression:
		return t.transpileLet(exp)

	case *ast.Identifier:
		return exp.Value, nil

	case *ast.CallExpression:
		return t.transpileCall(exp)

	case *ast.PropertyExpression:
		obj, err := t.transpileExpression(exp.Object)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.%s", obj, exp.Property.Value), nil

	case *ast.AssignExpression:
		return t.transpileAssign(exp)

	case *ast.CompoundAssignExpression:
		return t.transpileCompoundAssign(exp)

	case *ast.RangeExpression:
		return t.transpileRangeExpression(exp)
	}
	
	return "", fmt.Errorf("Unknown expression type: %T", node)
}

func (t *Transpiler) transpileFunction(exp *ast.FunctionLiteral) (string, error) {
	params := []string{}
	for _, p := range exp.Parameters {
		param := p.Value
		if p.Type != nil {
			param += ": " + transpileType(p.Type)
		}
		params = append(params, param)
	}

	retType := ""
	if exp.ReturnType != nil {
		retType = ": " + transpileType(exp.ReturnType)
	}

	body, err := t.transpileBlock(exp.Body)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("(%s)%s => %s", strings.Join(params, ", "), retType, body), nil
}

func (t *Transpiler) transpileConst(exp *ast.ConstExpression) (string, error) {
	name, err := t.transpileExpression(exp.Name)
	if err != nil {
		return "", err
	}

	val, err := t.transpileExpression(exp.Value)
	if err != nil {
		return "", err
	}
	
	tsType := ""
	if exp.Type != nil {
		tsType = ": " + transpileType(exp.Type)
	}

	return fmt.Sprintf("const %s%s = %s", name, tsType, val), nil
}

func (t *Transpiler) transpileLet(exp *ast.LetExpression) (string, error) {
	name, err := t.transpileExpression(exp.Name)
	if err != nil {
		return "", err
	}

	val, err := t.transpileExpression(exp.Value)
	if err != nil {
		return "", err
	}
	
	tsType := ""
	if exp.Type != nil {
		tsType = ": " + transpileType(exp.Type)
	}

	return fmt.Sprintf("let %s%s = %s", name, tsType, val), nil
}

func (t *Transpiler) transpileAssign(exp *ast.AssignExpression) (string, error) {
	name, err := t.transpileExpression(exp.Name)
	if err != nil {
		return "", err 
	}

	val, err := t.transpileExpression(exp.Value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = %s", name, val), nil
}

func (t *Transpiler) transpileCompoundAssign(exp *ast.CompoundAssignExpression) (string, error) {
	name, err := t.transpileExpression(exp.Name)
	if err != nil {
		return "", err 
	}

	val, err := t.transpileExpression(exp.Value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s %s", name, exp.Operator, val), nil
}

func (t *Transpiler) transpileInfix(exp *ast.InfixExpression) (string, error) {
	left, err := t.transpileExpression(exp.Left)
	if err != nil {
		return "", err
	}

	right, err := t.transpileExpression(exp.Right)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s %s", left, exp.Operator, right), nil
}

func (t *Transpiler) transpileCall(exp *ast.CallExpression) (string, error) {
	fn, err := t.transpileExpression(exp.Function)
	if err != nil {
		return "", err
	}

	args := []string{}

	for _, arg := range exp.Arguments {
		a, err := t.transpileExpression(arg)
		if err != nil {
			return "", err
		}
		args = append(args, a)
	}

	return fmt.Sprintf("%s(%s)", fn, strings.Join(args, ", ")), nil
}

func (t *Transpiler) transpileRangeExpression(exp *ast.RangeExpression) (string, error) {
	t.imports["$range"] = "$range" // ensure that the $range helper is imported

	start, err := t.transpileExpression(exp.Start)
	if err != nil {
		return "", err
	}

	end, err := t.transpileExpression(exp.End)
	if err != nil {
		return "", err
	}

	if exp.Step != nil {
		step, err := t.transpileExpression(exp.Step)
		if err != nil {
			return "", err
		}
		return "$range(" + start + ", " + end + ", " + step + ")", nil
	}

	return "$range(" + start + ", " + end + ")", nil
}
