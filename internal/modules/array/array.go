package array

import (
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() *object.Hash {
	return modbuilder.NewModuleBuilder().
		AddFunction("push", push).
		AddFunction("pop", pop).
		AddFunction("shift", shift).
		AddFunction("unshift", unshift).
		AddFunction("reverse", reverse).
		AddFunction("indexOf", indexOf).
		AddFunction("slice", slice).
		AddFunction("clear", clearArray).
		AddFunction("join", join).
		AddFunction("concat", concat).
		AddFunction("contains", contains).
		Build()
}

func push(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	array.Elements = append(array.Elements, args[1])
	return &object.Null{}
}

func pop(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	if len(array.Elements) == 0 {
		return errors.NewRuntimeError(0, 0, "array is empty")
	}

	lastElement := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]

	return lastElement
}

func shift(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	if len(array.Elements) == 0 {
		return errors.NewRuntimeError(0, 0, "array is empty")
	}

	firstElement := array.Elements[0]
	array.Elements = array.Elements[1:]

	return firstElement
}

func unshift(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	array.Elements = append([]object.Object{args[1]}, array.Elements...)
	return &object.Null{}
}

func reverse(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)

	reversed := make([]object.Object, len(array.Elements))
	for i, v := range array.Elements {
		reversed[len(array.Elements)-1-i] = v
	}

	array.Elements = reversed
	return &object.Null{}
}

func join(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
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
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		err = errors.ExpectNumberOfArguments(0, 0, 3, args)
		if err != nil {
			return err
		}
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.IntegerObj)
	if err != nil {
		return err
	}

	if len(args) == 3 {
		err = errors.ExpectType(0, 0, args[2], object.IntegerObj)
		if err != nil {
			return err
		}
	}

	array, _ := args[0].(*object.Array)
	start := args[1].(*object.Integer).Value
	end := int64(len(array.Elements))

	if len(args) == 3 {
		end = args[2].(*object.Integer).Value
	}

	if start < 0 {
		return errors.NewIndexOutOfBoundsError(0, 0, array)
	}

	if end < 0 {
		return errors.NewIndexOutOfBoundsError(0, 0, array)
	}

	if start > end {
		return errors.NewRuntimeError(0, 0, "start index is greater than end index")
	}

	result := make([]object.Object, end-start)
	copy(result, array.Elements[start:end])

	return &object.Array{Elements: result}
}

func indexOf(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	value := args[1]

	for i, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return &object.Integer{Value: int64(i)}
		}
	}

	return &object.Integer{Value: -1}
}

func contains(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	value := args[1]

	for _, v := range array.Elements {
		if v.Inspect() == value.Inspect() {
			return &object.Boolean{Value: true}
		}
	}

	return &object.Boolean{Value: false}
}

func concat(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.ArrayObj)
	if err != nil {
		return err
	}

	arr1, _ := args[0].(*object.Array)
	arr2 := args[1].(*object.Array)

	result := make([]object.Object, len(arr1.Elements)+len(arr2.Elements))
	copy(result, arr1.Elements)
	copy(result[len(arr1.Elements):], arr2.Elements)

	return &object.Array{Elements: result}
}

func clearArray(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	array, _ := args[0].(*object.Array)
	array.Elements = []object.Object{} // should this be set to nil?
	return &object.Null{}
}
