package array

import (
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
	AddFunction("push", push).
	AddFunction("pop", pop).
	AddFunction("shift", shift).
	AddFunction("unshift", unshift).
	AddFunction("reverse", reverse).
	AddFunction("indexOf", indexOf).
	AddFunction("slice", slice).
	AddFunction("clear", clear).
	AddFunction("join", join).
	AddFunction("concat", concat).
	AddFunction("contains", contains).
	Build()

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	array.Elements = append(array.Elements, args[1])
	return nil
}

func pop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	if len(array.Elements) == 0 {
		return object.NewError(0, 0, "array is empty")
	}

	lastElement := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]

	return lastElement
}

func shift(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	if len(array.Elements) == 0 {
		return object.NewError(0, 0, "array is empty")
	}

	firstElement := array.Elements[0]
	array.Elements = array.Elements[1:]

	return firstElement
}

func unshift(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	array.Elements = append([]object.Object{args[1]}, array.Elements...)
	return nil
}

func reverse(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)

	reversed := make([]object.Object, len(array.Elements))
	for i, v := range array.Elements {
		reversed[len(array.Elements)-1-i] = v
	}

	array.Elements = reversed
	return nil
}

func join(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	if args[1].Type() != object.StringObj {
		return object.NewError(0, 0, "argument must be a string, got %s", args[1].Type().String())
	}

	array := args[0].(*object.Array)
	separator := args[1].(*object.String)

	result := ""
	for i, v := range array.Elements {
		if i > 0 {
			result += separator.Value
		}
		result += v.Inspect()
	}

	return &object.String{Value: result}
}

func slice(args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want 2 or 3", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "first argument must be an array, got %s", args[0].Type().String())
	}

	if args[1].Type() != object.IntegerObj {
		return object.NewError(0, 0, "second argument must be an integer, got %s", args[1].Type().String())
	}

	if len(args) == 3 && args[2].Type() != object.IntegerObj {
		return object.NewError(0, 0, "third argument must be an integer, got %s", args[2].Type().String())
	}

	array := args[0].(*object.Array)
	start := args[1].(*object.Integer).Value
	var end int64 = int64(len(array.Elements))

	if len(args) == 3 {
		end = args[2].(*object.Integer).Value
	}

	if start < 0 || end < 0 || start > end {
		return object.NewError(0, 0, "invalid slice arguments")
	}

	result := make([]object.Object, end-start)
	copy(result, array.Elements[start:end])

	return &object.Array{Elements: result}
}

func indexOf(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "first argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	value := args[1]

	for i, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return &object.Integer{Value: int64(i)}
		}
	}

	return &object.Integer{Value: -1}
}

func contains(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "first argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	value := args[1]

	for _, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return &object.Boolean{Value: true}
		}
	}

	return &object.Boolean{Value: false}
}

func concat(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "first argument must be an array, got %s", args[0].Type().String())
	}

	if args[1].Type() != object.ArrayObj {
		return object.NewError(0, 0, "second argument must be an array, got %s", args[1].Type().String())
	}

	arr1 := args[0].(*object.Array)
	arr2 := args[1].(*object.Array)

	result := make([]object.Object, len(arr1.Elements)+len(arr2.Elements))
	copy(result, arr1.Elements)
	copy(result[len(arr1.Elements):], arr2.Elements)

	return &object.Array{Elements: result}
}

func clear(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewError(0, 0, "argument must be an array, got %s", args[0].Type().String())
	}

	array := args[0].(*object.Array)
	array.Elements = []object.Object{} // should this be set to nil?
	return nil
}
