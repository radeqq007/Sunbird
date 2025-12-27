package math

import (
	"math"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
	AddFunction("abs", abs).
	AddFunction("max", max).
	AddFunction("min", min).
	Build()

func abs(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(0, 0, "wrong number of arguments: expected 1, got %d", len(args))
	}

	if args[0].Type() != object.IntegerObj && args[0].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	return &object.Integer{
		Value: int64(args[0].(*object.Integer).Value),
	}
}

func max(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}

	if args[0].Type() != object.IntegerObj && args[0].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	if args[1].Type() != object.IntegerObj && args[1].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	return &object.Integer{
		Value: int64(math.Max(float64(args[0].(*object.Integer).Value), float64(args[1].(*object.Integer).Value))),
	}
}

func min(args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(0, 0, "wrong number of arguments: expected 2, got %d", len(args))
	}

	if args[0].Type() != object.IntegerObj && args[0].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	if args[1].Type() != object.IntegerObj && args[1].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	return &object.Integer{
		Value: int64(math.Min(float64(args[0].(*object.Integer).Value), float64(args[1].(*object.Integer).Value))),
	}
}
