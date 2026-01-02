package evaluator

import (
	"strings"

	"sunbird/internal/ast"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

func checkType(expected ast.TypeAnnotation, actual object.Object, line, col int) *object.Error {
	if expected == nil {
		return nil
	}

	switch expectedType := expected.(type) {
	case *ast.SimpleType:
		expectedObjType := typeAnnotationToObjectType(expectedType.Name)
		if expectedObjType != actual.Type() {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Type().String())
		}

	case *ast.ArrayType:
		array, ok := actual.(*object.Array)
		if !ok {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Type().String())
		}

		if expectedType.ElementType != nil {
			for _, element := range array.Elements {
				if err := checkType(expectedType.ElementType, element, line, col); err != nil {
					return err
				}
			}
		}

	case *ast.OptionalType:
		if actual.Type() == object.NullObj {
			return nil
		}
		return checkType(expectedType.BaseType, actual, line, col)

	case *ast.HashType:
		if actual.Type() != object.HashObj {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Type().String())
		}

	case *ast.FunctionType:
		if actual.Type() != object.FunctionObj && actual.Type() != object.BuiltinObj {
			return errors.NewTypeError(line, col,
				"expected %s, got %s", expected.String(), actual.Type().String())
		}
	}

	return nil
}

func typeAnnotationToObjectType(typeName string) object.ObjectType {
	switch strings.ToLower(typeName) {
	case "int":
		return object.IntegerObj
	case "float":
		return object.FloatObj
	case "string":
		return object.StringObj
	case "bool":
		return object.BooleanObj
	case "array":
		return object.ArrayObj
	case "hash":
		return object.HashObj
	default:
		return object.NullObj
	}
}
