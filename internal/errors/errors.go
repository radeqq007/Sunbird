package errors

import (
	"fmt"
	"strings"
	"sunbird/internal/object"
)

type ErrorCode int

const (
	SyntaxError ErrorCode = iota
	TypeError
	TypeMismatchError
	UndefinedVariableError
	DivisionByZeroError
	ConstantReassignmentError
	RuntimeError
	IndexNotSupportedError
	IndexOutOfBoundsError
	KeyError
	ImportError
	VariableReassignmentError
	NotCallableError
	InvalidAssignmentError
	ArgumentError
	ParseError
	PropertyAccessOnNonObjectError
	UnknownOperatorError
	FeatureNotImplementedError
)

//go:generate stringer -type=ErrorCode

func New(code ErrorCode, line, col int, format string, args ...interface{}) *object.Error {
	return &object.Error{
		Line:        line,
		Col:         col,
		Message:     fmt.Sprintf("%s: %s", code.String(), fmt.Sprintf(format, args...)),
		Propagating: true,
	}
}

func ExpectType(line, col int, obj object.Object, expectedType object.ObjectType) *object.Error {
	if obj.Type() != expectedType {
		return New(
			TypeError,
			line,
			col,
			"expected %s, got %s",
			expectedType.String(),
			obj.Type().String(),
		)
	}
	return nil
}

func ExpectOneOfTypes(
	line, col int,
	obj object.Object,
	expectedTypes ...object.ObjectType,
) *object.Error {
	for _, expectedType := range expectedTypes {
		if obj.Type() == expectedType {
			return nil
		}
	}

	typeNames := make([]string, len(expectedTypes))
	for i, t := range expectedTypes {
		typeNames[i] = t.String()
	}
	return New(
		TypeError,
		line,
		col,
		"expected one of %s, got %s",
		strings.Join(typeNames, ", "),
		obj.Type().String(),
	)
}

func ExpectNumberOfArguments(line, col int, expected int, args []object.Object) *object.Error {
	if len(args) != expected {
		return New(ArgumentError, line, col, "expected %d arguments, got %d", expected, len(args))
	}
	return nil
}

func NewIndexNotSupportedError(line, col int, obj object.Object) *object.Error {
	return New(IndexNotSupportedError, line, col, "%s", obj.Type().String())
}

func NewIndexOutOfBoundsError(line, col int, obj object.Object) *object.Error {
	return New(IndexOutOfBoundsError, line, col, "%s", obj.Type().String())
}

func NewNonObjectPropertyAccessError(line, col int, obj object.Object) *object.Error {
	return New(PropertyAccessOnNonObjectError, line, col, "%s", obj.Type().String())
}

func NewUndefinedVariableError(line, col int, identifier string) *object.Error {
	return New(UndefinedVariableError, line, col, "%s", identifier)
}

func NewUnusableAsHashKeyError(line, col int, obj object.Object) *object.Error {
	return New(KeyError, line, col, "%s", obj.Type().String())
}

func NewInvalidAssignmentTargetError(line, col int, name string) *object.Error {
	return New(InvalidAssignmentError, line, col, "%s", name)
}

func NewTypeError(line, col int, format string, args ...interface{}) *object.Error {
	return New(TypeError, line, col, format, args...)
}

func NewTypeMismatchError(
	line, col int,
	left object.ObjectType,
	operator string,
	right object.ObjectType,
) *object.Error {
	return New(TypeMismatchError, line, col, "%s %s %s", left.String(), operator, right.String())
}

func NewVariableReassignmentError(line, col int, identifier string) *object.Error {
	return New(VariableReassignmentError, line, col, "%s", identifier)
}

func NewNotCallableError(line, col int, obj object.Object) *object.Error {
	return New(NotCallableError, line, col, "%s", obj.Type().String())
}

func NewImportError(line, col int, message string) *object.Error {
	return New(ImportError, line, col, "%s", message)
}

func NewUnknownOperatorError(
	line, col int,
	left object.Object,
	operator string,
	right object.Object,
) *object.Error {
	return New(
		UnknownOperatorError,
		line,
		col,
		"%s %s %s",
		left.Type().String(),
		operator,
		right.Type().String(),
	)
}

func NewUnknownPrefixOperatorError(
	line, col int,
	operator string,
	right object.Object,
) *object.Error {
	return New(UnknownOperatorError, line, col, "%s%s", operator, right.Type().String())
}

func NewDivisionByZeroError(line, col int) *object.Error {
	return New(DivisionByZeroError, line, col, "")
}

func NewRuntimeError(line, col int, format string, args ...interface{}) *object.Error {
	return New(RuntimeError, line, col, format, args...)
}

func NewConstantReassignmentError(line, col int, identifier string) *object.Error {
	return New(ConstantReassignmentError, line, col, "%s", identifier)
}

func NewArgumentError(line, col int, format string, args ...interface{}) *object.Error {
	return New(ArgumentError, line, col, format, args...)
}

func NewFeatureNotImplementedError(line, col int, feature string) *object.Error {
	return New(FeatureNotImplementedError, line, col, "%s", feature)
}
