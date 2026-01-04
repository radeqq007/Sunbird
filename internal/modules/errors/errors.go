package errors

import (
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
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
	Build()

func wrapError(args []object.Object, internalFunc func(int, int, string) *object.Error) object.Object {
	if err := errors.ExpectNumberOfArguments(0, 0, 1, args); err != nil {
		return err
	}
	if err := errors.ExpectType(0, 0, args[0], object.StringObj); err != nil {
		return err
	}
	msg := args[0].(*object.String).Value
	return internalFunc(0, 0, msg)
}

func NewTypeError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewTypeError(l, c, m)
	})
}

func NewRuntimeError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewRuntimeError(l, c, m)
	})
}

func NewImportError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewImportError(l, c, m)
	})
}

func NewDivisionByZeroError(args ...object.Object) object.Object {
	// Division by zero doesn't need a custom message
	if err := errors.ExpectNumberOfArguments(0, 0, 0, args); err != nil {
		return err
	}

	return errors.NewDivisionByZeroError(0, 0)
}

func NewConstantReassignmentError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewConstantReassignmentError(l, c, m)
	})
}

func NewIndexNotSupportedError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.IndexNotSupportedError, l, c, m)
	})
}

func NewIndexOutOfBoundsError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.IndexOutOfBoundsError, l, c, m)
	})
}

func NewKeyError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.KeyError, l, c, m)
	})
}

func NewVariableReassignmentError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewVariableReassignmentError(l, c, m)
	})
}

func NewNotCallableError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.NotCallableError, l, c, m)
	})
}

func NewInvalidAssignmentError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.NewInvalidAssignmentTargetError(l, c, m)
	})
}

func NewArgumentError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.ArgumentError, l, c, m)
	})
}

func NewPropertyAccessError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.PropertyAccessOnNonObjectError, l, c, m)
	})
}

func NewUnknownOperatorError(args ...object.Object) object.Object {
	return wrapError(args, func(l, c int, m string) *object.Error {
		return errors.New(errors.UnknownOperatorError, l, c, m)
	})
}
