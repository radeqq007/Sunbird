package transpiler

import (
	"fmt"
	"strings"
	"sunbird/internal/ast"
)

func (t *Transpiler) transpileStatement(node ast.Statement) (string, error) {
	switch stmt := node.(type) {
	case *ast.ExpressionStatement:
		expr, err := t.transpileExpression(stmt.Expression)
		if err != nil {
			return "", err
		}
		return t.indentStr() + expr + ";", nil

	case *ast.ReturnStatement:
		if stmt.ReturnValue == nil {
			return t.indentStr() + "return;", nil
		}
		val, err := t.transpileExpression(stmt.ReturnValue)
		if err != nil {
			return "", err
		}
		return t.indentStr() + "return " + val + ";", nil
	
	case *ast.BlockStatement:
		return t.transpileBlock(stmt)

	case *ast.ImportStatement:
		return t.transpileImportStatement(stmt)
	
	case *ast.WhileStatement:
		return t.transpileWhileStatement(stmt)
	}

	return "", fmt.Errorf("unknown statement type: %T", node)
}

func (t *Transpiler) transpileBlock(node *ast.BlockStatement) (string, error) {
	var out strings.Builder
	out.WriteString("{\n")
	t.pushIndent()
	
	for _, stmt := range node.Statements {
		s, err := t.transpileStatement(stmt)
		if err != nil {
			return "", err
		}

		out.WriteString(s)
		out.WriteByte('\n')
	}
	out.WriteString("}")

	t.popIndent()

	return out.String(), nil
}

func (t *Transpiler) transpileWhileStatement(stmt *ast.WhileStatement) (string, error) {
	cond, err := t.transpileExpression(stmt.Condition)
	if err != nil {
		return "", err
	}

	body, err := t.transpileStatement(stmt.Body)
	if err != nil {
		return "", err
	}

	return "while (" + cond + ") " + body, nil
}

func (t *Transpiler) transpileImportStatement(stmt *ast.ImportStatement) (string, error) {
	path := stmt.Path.Value
	if isBuiltinModule(path) {
		alias := path
		if stmt.Alias != nil {
			alias = stmt.Alias.Value
		}
		t.imports[alias] = path
		return "", nil
	}

	alias := stmt.Path.Value
	if stmt.Alias != nil {
		alias = stmt.Alias.Value
	}

	return fmt.Sprintf("import * as %s from \"./%s.js\";", alias, path), nil
}

func isBuiltinModule(name string) bool {
	switch name {
	case "io", "math", "array", "string", "random", "json", "http", "fs", "time", "errors":
		return true
	}
	return false
}
