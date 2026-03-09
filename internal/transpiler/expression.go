package transpiler

import (
	"fmt"
	"strings"

	"github.com/radeqq007/sunbird/internal/ast"
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

	case *ast.ArrayLiteral:
		return t.transpileArray(exp)

	case *ast.IfExpression:
		return t.transpileIfExpression(exp)

	case *ast.IndexExpression:
		return t.transpileIndexEpression(exp)

	case *ast.HashLiteral:
		return t.transpileHashLiteral(exp)
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
	if ident, ok := exp.Function.(*ast.Identifier); ok {
		if name, ok := isBuiltin(ident.Value); ok {
			t.imports[name] = name
		}
	}

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
	t.imports["__range"] = "__range" // ensure that the __range helper is imported

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
		return "__range(" + start + ", " + end + ", " + step + ")", nil
	}

	return "__range(" + start + ", " + end + ")", nil
}

func (t *Transpiler) transpileArray(exp *ast.ArrayLiteral) (string, error) {
	elements := make([]string, len(exp.Elements))
	for i, el := range exp.Elements {
		s, err := t.transpileExpression(el)
		if err != nil {
			return "", err
		}

		elements[i] = s
	}

	return "[" + strings.Join(elements, ", ") + "]", nil
}

// transpileIfExpression generates an if-IIFE in order to keep ifs as expressions
func (t *Transpiler) transpileIfExpression(exp *ast.IfExpression) (string, error) {
	cond, err := t.transpileExpression(exp.Condition)
	if err != nil {
		return "", err
	}

	consequence, err := t.transpileIfBlock(exp.Consequence)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("(() => { if (%s) %s", cond, consequence)

	if exp.Alternative != nil {
		alternative, err := t.transpileIfBlock(exp.Alternative)
		if err != nil {
			return "", err
		}
		result += " else " + alternative
	}

	result += " })()"
	return result, nil
}

func (t *Transpiler) transpileIfBlock(block *ast.BlockStatement) (string, error) {
	if len(block.Statements) == 0 {
		return "{}", nil
	}

	var parts []string
	for i, stmt := range block.Statements {
		// For the last statement, if it's a bare expression, turn it into a return.
		if i == len(block.Statements)-1 {
			if es, ok := stmt.(*ast.ExpressionStatement); ok {
				val, err := t.transpileExpression(es.Expression)
				if err != nil {
					return "", err
				}
				parts = append(parts, "return "+val+";")
				continue
			}
		}
		s, err := t.transpileStatement(stmt)
		if err != nil {
			return "", err
		}
		parts = append(parts, strings.TrimSpace(s))
	}

	return "{ " + strings.Join(parts, " ") + " }", nil
}

func (t *Transpiler) transpileIndexEpression(exp *ast.IndexExpression) (string, error) {
	t.imports["__idx"] = "__idx"

	left, err := t.transpileExpression(exp.Left)
	if err != nil {
		return "", err
	}

	idx, err := t.transpileExpression(exp.Index)
	if err != nil {
		return "", err
	}

	return "__idx(" + left + ", " + idx + ")", nil
}

func (t *Transpiler) transpileHashLiteral(exp *ast.HashLiteral) (string, error) {
	if exp.Pairs == nil {
		return "{}", nil
	}

	pairs := make([]string, 0, len(exp.Pairs))
	for _, pair := range exp.Pairs {
		key, err := t.transpileExpression(pair.Key)
		if err != nil {
			return "", err
		}

		val, err := t.transpileExpression(pair.Value)
		if err != nil {
			return "", err
		}

		// String keys don't need brackets, everything else does
		if _, ok := pair.Key.(*ast.StringLiteral); ok {
			pairs = append(pairs, fmt.Sprintf("%s: %s", key, val))
		} else {
			pairs = append(pairs, fmt.Sprintf("[%s]: %s", key, val))
		}
	}

	return "{ " + strings.Join(pairs, ", ") + " }", nil
}

func isBuiltin(name string) (string, bool) {
	builtins := map[string]string{
		"len": "__len", "append": "__append",
		"type": "type", "string": "__type",
		"int": "__int", "float": "__float",
		"bool": "__bool", "exit": "__exit",
		"error": "__error",
	}

	n, ok := builtins[name]
	return n, ok
}
