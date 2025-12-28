package str

import (
	"strings"
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
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str1, ok1 := args[0].(*object.String)
	str2, ok2 := args[1].(*object.String)
	if !ok1 || !ok2 {
		return object.NewError(0, 0, "wrong arguments: expected 2 strings, got %T and %T", args[0].Type(), args[1].Type())
	}
	
	return &object.String{Value: str1.Value + str2.Value}
}

func isEmpty(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments: expected 1, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}
	
	return &object.Boolean{Value: str.Value == ""}
}

func startsWith(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}

	startStr, ok := args[1].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[1].Type())
	}
	
	return &object.Boolean{Value: strings.HasPrefix(str.Value, startStr.Value)}
}

func endsWith(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}

	endStr, ok := args[1].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[1].Type())
	}
	
	return &object.Boolean{Value: strings.HasSuffix(str.Value, endStr.Value)}
}

func contains(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}

	subStr, ok := args[1].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[1].Type())
	}
	
	return &object.Boolean{Value: strings.Contains(str.Value, subStr.Value)}
}

func toUpper(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments: expected 1, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}
	
	return &object.String{Value: strings.ToUpper(str.Value)}
}

func toLower(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments: expected 1, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}
	
	return &object.String{Value: strings.ToLower(str.Value)}
}

func trim(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments: expected 1, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}
	
	return &object.String{Value: strings.TrimSpace(str.Value)}
}

func split(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}

	sep, ok := args[1].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[1].Type())
	}

	strs := strings.Split(str.Value, sep.Value)
	var objects []object.Object = make([]object.Object, len(strs))
	for i, s := range strs {
		objects[i] = &object.String{Value: s}
	}

	return &object.Array{Elements: objects}
}

func repeat(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}
	
	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected string, got %T", args[0].Type())
	}

	count, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewError(0, 0, "wrong argument: expected integer, got %T", args[1].Type())
	}

	return &object.String{Value: strings.Repeat(str.Value, int(count.Value))}
}

