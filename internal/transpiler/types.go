package transpiler

import "sunbird/internal/ast"

func transpileType(t ast.TypeAnnotation) string {
	if t == nil {
		return ""
	}

	switch t := t.(type) {
	case *ast.SimpleType:
		return simpleTypeToTS(t.Name)
	case *ast.ArrayType:
		if t.ElementType != nil {
			return transpileType(t.ElementType) + "[]"
		}
		return "unkown[]"

	case *ast.HashType:
		return "Record<string | number, unknown>"

	case *ast.FunctionType:
		return "(...args: unknown[]) => unknown"

	case *ast.OptionalType:
		return transpileType(t.BaseType) + " | null"
	}

	return "unknown"
}

func simpleTypeToTS(t string) string {
	switch t {
	case "Int", "Float":
		return "number"
	case "String":
		return "string"
	case "Bool":
		return "boolean"
	case "Void":
		return "void"
	case "Array":
		return "unknown[]"
	case "Hash":
		return "Record<string | number, unknown>"
	case "Func":
		return "(...args: unknown[]) => unknown"
	default:
		return "unknown"
	}
}
