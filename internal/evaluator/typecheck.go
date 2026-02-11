package evaluator

import (
	"strings"

	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func checkType(expected ast.TypeAnnotation, actual object.Value, line, col int) object.Value {
	if expected == nil {
		return object.NewNull()
	}

	switch expectedType := expected.(type) {
	case *ast.SimpleType:
		expectedObjType := typeAnnotationToObjectKind(expectedType.Name)
		if expectedObjType != actual.Kind() {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Kind().String())
		}

	case *ast.ArrayType:
		if !actual.IsArray() {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Kind().String())
		}

		arr := actual.AsArray()
		if expectedType.ElementType != nil {
			for _, element := range arr.Elements {
				if err := checkType(expectedType.ElementType, element, line, col); err.IsError() {
					return err
				}
			}
		}

	case *ast.OptionalType:
		if actual.IsNull() {
			return object.NewNull()
		}
		return checkType(expectedType.BaseType, actual, line, col)

	case *ast.HashType:
		if !actual.IsHash() {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Kind().String())
		}

	case *ast.FunctionType:
		if !actual.IsFunction() && !actual.IsBuiltin() {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Kind().String())
		}
	}

	return object.NewNull()
}

func typeAnnotationToObjectKind(typeName string) object.ValueKind {
	switch strings.ToLower(typeName) {
	case "int":
		return object.IntKind
	case "float":
		return object.FloatKind
	case "string":
		return object.StringKind
	case "bool":
		return object.BoolKind
	case "array":
		return object.ArrayKind
	case "hash":
		return object.HashKind
	case "range":
		return object.RangeKind
	default:
		return object.NullKind
	}
}
