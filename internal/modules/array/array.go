package array

import (
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("push", push).
		AddFunction("pop", pop).
		AddFunction("shift", shift).
		AddFunction("unshift", unshift).
		AddFunction("reverse", reverse).
		AddFunction("index_of", indexOf).
		AddFunction("slice", slice).
		AddFunction("clear", clearArray).
		AddFunction("join", join).
		AddFunction("concat", concat).
		AddFunction("contains", contains).
		Build()
}

func push(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	array.Elements = append(array.Elements, args[1])
	return object.NewNull()
}

func pop(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	if len(array.Elements) == 0 {
		return errors.NewRuntimeError(ctx.Line, ctx.Col, "array is empty")
	}

	lastElement := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]

	return lastElement
}

func shift(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	if len(array.Elements) == 0 {
		return errors.NewRuntimeError(ctx.Line, ctx.Col, "array is empty")
	}

	firstElement := array.Elements[0]
	array.Elements = array.Elements[1:]

	return firstElement
}

func unshift(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	array.Elements = append([]object.Value{args[1]}, array.Elements...)

	return object.NewNull()
}

func reverse(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()

	reversed := make([]object.Value, len(array.Elements))
	for i, v := range array.Elements {
		reversed[len(array.Elements)-1-i] = v
	}

	array.Elements = reversed
	return object.NewNull()
}

func join(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	separator := args[1].AsString()

	var b strings.Builder
	for i, v := range array.Elements {
		if i > 0 {
			b.WriteString(separator.Value)
		}
		b.WriteString(v.Inspect())
	}

	return object.NewString(b.String())
}

func slice(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		err = errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 3, args)
		if err.IsError() {
			return err
		}
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.IntKind)
	if err.IsError() {
		return err
	}

	if len(args) == 3 {
		err = errors.ExpectType(ctx.Line, ctx.Col, args[2], object.IntKind)
		if err.IsError() {
			return err
		}
	}

	array := args[0].AsArray()
	start := args[1].AsInt()
	end := int64(len(array.Elements))

	if len(args) == 3 {
		end = args[2].AsInt()
	}

	if start < 0 {
		return errors.NewIndexOutOfBoundsError(ctx.Line, ctx.Col, args[0])
	}

	if end < 0 {
		return errors.NewIndexOutOfBoundsError(ctx.Line, ctx.Col, args[0])
	}

	if start > end {
		return errors.NewRuntimeError(ctx.Line, ctx.Col, "start index is greater than end index")
	}

	result := make([]object.Value, end-start)
	copy(result, array.Elements[start:end])

	return object.NewArray(result)
}

func indexOf(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	value := args[1]

	for i, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return object.NewInt(int64(i))
		}
	}

	return object.NewInt(-1)
}

func contains(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	value := args[1]

	for _, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return object.NewBool(true)
		}
	}

	return object.NewBool(false)
}

func concat(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.ArrayKind)
	if err.IsError() {
		return err
	}

	arr1 := args[0].AsArray()
	arr2 := args[1].AsArray()

	result := make([]object.Value, len(arr1.Elements)+len(arr2.Elements))
	copy(result, arr1.Elements)
	copy(result[len(arr1.Elements):], arr2.Elements)

	return object.NewArray(result)
}

func clearArray(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	array := args[0].AsArray()
	array.Elements = []object.Value{} // should this be set to nil?
	return object.NewNull()
}
