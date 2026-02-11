package errors

import (
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() *object.Hash {
	return modbuilder.NewModuleBuilder().
		AddFunction("type_error", NewTypeError).
		AddFunction("runtime_error", NewRuntimeError).
		AddFunction("import_error", NewImportError).
		AddFunction("division_by_zero_error", NewDivisionByZeroError).
		AddFunction("constant_reassignment_error", NewConstantReassignmentError).
		AddFunction("index_not_supported_error", NewIndexNotSupportedError).
		AddFunction("index_out_of_bounds_error", NewIndexOutOfBoundsError).
		AddFunction("key_error", NewKeyError).
		AddFunction("variable_reassignment_error", NewVariableReassignmentError).
		AddFunction("not_callable_error", NewNotCallableError).
		AddFunction("invalid_assignment_error", NewInvalidAssignmentError).
		AddFunction("argument_error", NewArgumentError).
		AddFunction("property_access_error", NewPropertyAccessError).
		AddFunction("feature_not_implemented_error", NewFeatureNotImplementedError).
		Build()
}

func wrapError(
	args []object.Value,
	internalFunc func(int, int, string) object.Value,
) object.Value {
	if err := errors.ExpectNumberOfArguments(0, 0, 1, args); !err.IsNull() {
		return err
	}
	if err := errors.ExpectType(0, 0, args[0], object.StringKind); !err.IsNull() {
		return err
	}
	msg := args[0].AsString().Value
	return internalFunc(0, 0, msg)
}

func NewTypeError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.NewTypeError(l, c, "%s", m)
	})
}

func NewRuntimeError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.NewRuntimeError(l, c, "%s", m)
	})
}

func NewImportError(args ...object.Value) object.Value {
	return wrapError(args, errors.NewImportError)
}

func NewDivisionByZeroError(args ...object.Value) object.Value {
	// Division by zero doesn't need a custom message
	if err := errors.ExpectNumberOfArguments(0, 0, 0, args); !err.IsNull() {
		return err
	}

	return errors.NewDivisionByZeroError(0, 0)
}

func NewConstantReassignmentError(args ...object.Value) object.Value {
	return wrapError(args, errors.NewConstantReassignmentError)
}

func NewIndexNotSupportedError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.IndexNotSupportedError, l, c, "%s", m)
	})
}

func NewIndexOutOfBoundsError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.IndexOutOfBoundsError, l, c, "%s", m)
	})
}

func NewKeyError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.KeyError, l, c, "%s", m)
	})
}

func NewVariableReassignmentError(args ...object.Value) object.Value {
	return wrapError(args, errors.NewVariableReassignmentError)
}

func NewNotCallableError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.NotCallableError, l, c, "%s", m)
	})
}

func NewInvalidAssignmentError(args ...object.Value) object.Value {
	return wrapError(args, errors.NewInvalidAssignmentTargetError)
}

func NewArgumentError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.ArgumentError, l, c, "%s", m)
	})
}

func NewPropertyAccessError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.PropertyAccessOnNonObjectError, l, c, "%s", m)
	})
}

func NewUnknownOperatorError(args ...object.Value) object.Value {
	return wrapError(args, func(l, c int, m string) object.Value {
		return errors.New(errors.UnknownOperatorError, l, c, "%s", m)
	})
}

func NewFeatureNotImplementedError(args ...object.Value) object.Value {
	return wrapError(args, errors.NewFeatureNotImplementedError)
}
