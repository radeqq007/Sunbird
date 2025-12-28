package str

import (
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
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

func concat(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	str1 := args[0].(*object.String)
	str2 := args[1].(*object.String)

	return &object.String{Value: str1.Value + str2.Value}
}

func isEmpty(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)

	return &object.Boolean{Value: str.Value == ""}
}

func startsWith(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)
	startStr := args[1].(*object.String)

	return &object.Boolean{Value: strings.HasPrefix(str.Value, startStr.Value)}
}

func endsWith(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)
	endStr := args[1].(*object.String)

	return &object.Boolean{Value: strings.HasSuffix(str.Value, endStr.Value)}
}

func contains(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)
	subStr := args[1].(*object.String)

	return &object.Boolean{Value: strings.Contains(str.Value, subStr.Value)}
}

func toUpper(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)

	return &object.String{Value: strings.ToUpper(str.Value)}
}

func toLower(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)

	return &object.String{Value: strings.ToLower(str.Value)}
}

func trim(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)

	return &object.String{Value: strings.TrimSpace(str.Value)}
}

func split(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)
	sep := args[1].(*object.String)

	strs := strings.Split(str.Value, sep.Value)
	var objects []object.Object = make([]object.Object, len(strs))
	for i, s := range strs {
		objects[i] = &object.String{Value: s}
	}

	return &object.Array{Elements: objects}
}

func repeat(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.IntegerObj)
	if err != nil {
		return err
	}

	str := args[0].(*object.String)
	count := args[1].(*object.Integer)

	return &object.String{Value: strings.Repeat(str.Value, int(count.Value))}
}
