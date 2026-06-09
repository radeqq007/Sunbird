package errors

import (
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
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
	ctx object.CallContext,
	args []object.Value,
	internalFunc func(int, int, string) object.Value,
) object.Value {
	if err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args); err.IsError() {
		return err
	}
	if err := errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind); err.IsError() {
		return err
	}
	msg := args[0].AsString().Value
	return internalFunc(ctx.Line, ctx.Col, msg)
}

func NewTypeError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.NewTypeError(l, c, "%s", m)
	})
}

func NewRuntimeError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.NewRuntimeError(l, c, "%s", m)
	})
}

func NewImportError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, errors.NewImportError)
}

func NewDivisionByZeroError(ctx object.CallContext, args ...object.Value) object.Value {
	// Division by zero doesn't need a custom message
	if err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 0, args); err.IsError() {
		return err
	}

	return errors.NewDivisionByZeroError(ctx.Line, ctx.Col)
}

func NewConstantReassignmentError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, errors.NewConstantReassignmentError)
}

func NewIndexNotSupportedError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.IndexNotSupportedError, l, c, "%s", m)
	})
}

func NewIndexOutOfBoundsError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.IndexOutOfBoundsError, l, c, "%s", m)
	})
}

func NewKeyError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.KeyError, l, c, "%s", m)
	})
}

func NewVariableReassignmentError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, errors.NewVariableReassignmentError)
}

func NewNotCallableError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.NotCallableError, l, c, "%s", m)
	})
}

func NewInvalidAssignmentError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, errors.NewInvalidAssignmentTargetError)
}

func NewArgumentError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.ArgumentError, l, c, "%s", m)
	})
}

func NewPropertyAccessError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.PropertyAccessOnNonObjectError, l, c, "%s", m)
	})
}

func NewUnknownOperatorError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, func(l, c int, m string) object.Value {
		return errors.New(errors.UnknownOperatorError, l, c, "%s", m)
	})
}

func NewFeatureNotImplementedError(ctx object.CallContext, args ...object.Value) object.Value {
	return wrapError(ctx, args, errors.NewFeatureNotImplementedError)
}
