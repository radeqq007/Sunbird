package errors

import (
	"fmt"
	"sunbird/internal/object"
)

const (
	SyntaxError = iota
	TypeError
	TypeMismatch
	UndefinedVariable
	DivisionByZero
	RuntimeError
	IndexNotSupported
	IndexOutOfBounds
	KeyError
	ImportError
	VariableReassignmentError
	NotCallable
	InvalidAssignment
	ArgumentError
	ParseError
	PropertyAccessOnNonObject
	UnknownOperator
)

type ErrorCode int

func (ec ErrorCode) String() string {
	switch ec {
	case SyntaxError:
		return "SyntaxError"
	case TypeError:
		return "TypeError"
	case TypeMismatch:
		return "TypeMismatch"
	case UndefinedVariable:
		return "UndefinedVariable"
	case DivisionByZero:
		return "DivisionByZero"
	case RuntimeError:
		return "RuntimeError"
	case IndexNotSupported:
		return "IndexNotSupported"
	case IndexOutOfBounds:
		return "IndexOutOfBounds"
	case KeyError:
		return "KeyError"
	case ImportError:
		return "ImportError"
	case VariableReassignmentError:
		return "VariableReassignmentError"
	case NotCallable:
		return "NotCallable"
	case InvalidAssignment:
		return "InvalidAssignment"
	case ArgumentError:
		return "ArgumentError"
	case ParseError:
		return "ParseError"
	case PropertyAccessOnNonObject:
		return "PropertyAccessOnNonObject"
	case UnknownOperator:
		return "UnknownOperator"
	default:
		return "UnknownError"
	}
}

func New(code ErrorCode, line, col int, format string, args ...interface{}) *object.Error {
	return &object.Error{Line: line, Col: col, Message: fmt.Sprintf("%s: %s", code.String(), fmt.Sprintf(format, args...))}
}

func ExpectType(line, col int, obj object.Object, expectedType object.ObjectType) *object.Error {
	if obj.Type() != expectedType {
		return New(TypeError, line, col, "expected %s, got %s", expectedType.String(), obj.Type().String())
	}
	return nil
}

func ExpectOneOfTypes(line, col int, obj object.Object, expectedTypes ...object.ObjectType) *object.Error {
	for _, expectedType := range expectedTypes {
		if obj.Type() == expectedType {
			return nil
		}
	}
	return New(TypeError, line, col, "expected one of %s, got %s", expectedTypes, obj.Type().String())
}

func ExpectNumberOfArguments(line, col int, expected int, args ...object.Object) *object.Error {
	if len(args) != expected {
		return New(ArgumentError, line, col, "expected %d arguments, got %d", expected, len(args))
	}
	return nil
}

func NewIndexNotSupportedError(line, col int, obj object.Object) *object.Error {
	return New(IndexNotSupported, line, col, "index operator not supported: %s", obj.Type().String())
}

func NewIndexOutOfBoundsError(line, col int, obj object.Object) *object.Error {
	return New(IndexOutOfBounds, line, col, "index out of bounds: %s", obj.Type().String())
}

func NewNonObjectPropertyAccessError(line, col int, obj object.Object) *object.Error {
	return New(PropertyAccessOnNonObject, line, col, "property access on non-object: %s", obj.Type().String())
}

func NewUndefinedVariableError(line, col int, identifier string) *object.Error {
	return New(UndefinedVariable, line, col, "undefined variable: %s", identifier)
}

func NewUnusableAsHashKeyError(line, col int, obj object.Object) *object.Error {
	return New(KeyError, line, col, "unusable as hash key: %s", obj.Type().String())
}

func NewInvalidAssignmentTargetError(line, col int, name string) *object.Error {
	return New(InvalidAssignment, line, col, "invalid assignment target: %s", name)
}

func NewTypeError(line, col int, format string, args ...interface{}) *object.Error {
	return New(TypeError, line, col, format, args...)
}

func NewTypeMismatchError(line, col int, left object.ObjectType, operator string, right object.ObjectType) *object.Error {
	return New(TypeMismatch, line, col, "type mismatch: %s %s %s", left.String(), operator, right.String())
}

func NewVariableReassignmentError(line, col int, identifier string) *object.Error {
	return New(VariableReassignmentError, line, col, "variable reassignment: %s", identifier)
}

func NewNotCallableError(line, col int, obj object.Object) *object.Error {
	return New(NotCallable, line, col, "not callable: %s", obj.Type().String())
}

func NewImportError(line, col int, message string) *object.Error {
	return New(ImportError, line, col, message)
}

func NewUnknownOperatorError(line, col int, operator string) *object.Error {
	return New(UnknownOperator, line, col, "unknown operator: %s", operator)
}

func NewDivisionByZeroError(line, col int) *object.Error {
	return New(DivisionByZero, line, col, "division by zero")
}
