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

func New(code ErrorCode, line, col int, format string, args ...interface{}) object.Value {
	msg := fmt.Sprintf("%s: %s", code.String(), fmt.Sprintf(format, args...))

	return object.NewError(msg, line, col, true)
}

func ExpectType(line, col int, val object.Value, expectedKind object.ValueKind) object.Value {
	if val.Kind() != expectedKind {
		return New(
			TypeError,
			line,
			col,
			"expected %s, got %s",
			expectedKind.String(),
			val.Kind().String(),
		)
	}

	return object.NewNull()
}

func ExpectOneOfTypes(
	line, col int,
	val object.Value,
	expectedTypes ...object.ValueKind,
) object.Value {
	for _, expectedType := range expectedTypes {
		if val.Kind() == expectedType {
			return object.NewNull()
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
		val.Kind().String(),
	)
}

func ExpectNumberOfArguments(line, col, expected int, args []object.Value) object.Value {
	if len(args) != expected {
		return New(ArgumentError, line, col, "expected %d arguments, got %d", expected, len(args))
	}
	return object.NewNull()
}

func ExpectMinNumberOfArguments(line, col, expected int, args []object.Value) object.Value {
	if len(args) < expected {
		return New(ArgumentError, line, col, "expected at least %d arguments, got %d", expected, len(args))
	}
	return object.NewNull()
}

func NewIndexNotSupportedError(line, col int, val object.Value) object.Value {
	return New(IndexNotSupportedError, line, col, "%s", val.Kind().String())
}

func NewIndexOutOfBoundsError(line, col int, val object.Value) object.Value {
	return New(IndexOutOfBoundsError, line, col, "%s", val.Kind().String())
}

func NewNonObjectPropertyAccessError(line, col int, val object.Value) object.Value {
	return New(PropertyAccessOnNonObjectError, line, col, "%s", val.Kind().String())
}

func NewUndefinedVariableError(line, col int, identifier string) object.Value {
	return New(UndefinedVariableError, line, col, "%s", identifier)
}

func NewUnusableAsHashKeyError(line, col int, val object.Value) object.Value {
	return New(KeyError, line, col, "%s", val.Kind().String())
}

func NewInvalidAssignmentTargetError(line, col int, name string) object.Value {
	return New(InvalidAssignmentError, line, col, "%s", name)
}

func NewTypeError(line, col int, format string, args ...any) object.Value {
	return New(TypeError, line, col, format, args...)
}

func NewTypeMismatchError(
	line, col int,
	left object.ValueKind,
	operator string,
	right object.ValueKind,
) object.Value {
	return New(TypeMismatchError, line, col, "%s %s %s", left.String(), operator, right.String())
}

func NewVariableReassignmentError(line, col int, identifier string) object.Value {
	return New(VariableReassignmentError, line, col, "%s", identifier)
}

func NewNotCallableError(line, col int, obj object.Value) object.Value {
	return New(NotCallableError, line, col, "%s", obj.Kind().String())
}

func NewImportError(line, col int, message string) object.Value {
	return New(ImportError, line, col, "%s", message)
}

func NewUnknownOperatorError(
	line, col int,
	left object.Value,
	operator string,
	right object.Value,
) object.Value {
	return New(
		UnknownOperatorError,
		line,
		col,
		"%s %s %s",
		left.Kind().String(),
		operator,
		right.Kind().String(),
	)
}

func NewUnknownPrefixOperatorError(
	line, col int,
	operator string,
	right object.Value,
) object.Value {
	return New(UnknownOperatorError, line, col, "%s%s", operator, right.Kind().String())
}

func NewDivisionByZeroError(line, col int) object.Value {
	return New(DivisionByZeroError, line, col, "")
}

func NewRuntimeError(line, col int, format string, args ...any) object.Value {
	return New(RuntimeError, line, col, format, args...)
}

func NewConstantReassignmentError(line, col int, identifier string) object.Value {
	return New(ConstantReassignmentError, line, col, "%s", identifier)
}

func NewArgumentError(line, col int, format string, args ...any) object.Value {
	return New(ArgumentError, line, col, format, args...)
}

func NewFeatureNotImplementedError(line, col int, feature string) object.Value {
	return New(FeatureNotImplementedError, line, col, "%s", feature)
}
