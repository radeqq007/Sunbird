package str

import (
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("concat", concat).
		AddFunction("is_empty", isEmpty).
		AddFunction("starts_with", startsWith).
		AddFunction("ends_with", endsWith).
		AddFunction("to_upper", toUpper).
		AddFunction("to_lower", toLower).
		AddFunction("trim", trim).
		AddFunction("split", split).
		AddFunction("repeat", repeat).
		AddFunction("contains", contains).
		Build()
}

func concat(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str1 := args[0].AsString()
	str2 := args[1].AsString()

	return object.NewString(str1.Value + str2.Value)
}

func isEmpty(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()

	return object.NewBool(str.Value == "")
}

func startsWith(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()
	startStr := args[1].AsString()

	return object.NewBool(strings.HasPrefix(str.Value, startStr.Value))
}

func endsWith(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()
	endStr := args[1].AsString()

	return object.NewBool(strings.HasSuffix(str.Value, endStr.Value))
}

func contains(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()
	subStr := args[1].AsString()

	return object.NewBool(strings.Contains(str.Value, subStr.Value))
}

func toUpper(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()

	return object.NewString(strings.ToUpper(str.Value))
}

func toLower(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()

	return object.NewString(strings.ToLower(str.Value))
}

func trim(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()

	return object.NewString(strings.TrimSpace(str.Value))
}

func split(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()
	sep := args[1].AsString()

	strs := strings.Split(str.Value, sep.Value)
	objects := make([]object.Value, len(strs))
	for i, s := range strs {
		objects[i] = object.NewString(s)
	}

	return object.NewArray(objects)
}

func repeat(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.IntKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString()
	count := args[1].AsInt()

	return object.NewString(strings.Repeat(str.Value, int(count)))
}
