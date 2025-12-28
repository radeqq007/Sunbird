package math

import (
	"math"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
	AddFunction("abs", abs).
	AddFunction("max", max).
	AddFunction("min", min).
	AddFunction("pow", pow).
	AddFunction("sqrt", sqrt).
	AddFunction("floor", floor).
	AddFunction("ceil", ceil).
	AddFunction("round", round).
	AddFunction("sign", sign).
	AddFunction("clamp", clamp).
	AddFunction("sin", sin).
	AddFunction("cos", cos).
	AddFunction("tan", tan).
	AddFloat("PI", math.Pi).
	AddFloat("E", math.E).
	Build()

func abs(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Abs(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Abs(getFloat64(args[0]))),
	}
}

func max(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj || args[1].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Max(getFloat64(args[0]), getFloat64(args[1])),
		}
	}

	return &object.Integer{
		Value: int64(math.Max(getFloat64(args[0]), getFloat64(args[1]))),
	}
}

func min(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() != object.IntegerObj && args[0].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	if args[1].Type() != object.IntegerObj && args[1].Type() != object.FloatObj {
		return object.NewError(0, 0, "argument must be an integer or float")
	}

	if args[0].Type() == object.FloatObj || args[1].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Min(getFloat64(args[0]), getFloat64(args[1])),
		}
	}

	return &object.Integer{
		Value: int64(math.Min(getFloat64(args[0]), getFloat64(args[1]))),
	}
}

func pow(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj || args[1].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Pow(getFloat64(args[0]), getFloat64(args[1])),
		}
	}

	return &object.Integer{
		Value: int64(math.Pow(getFloat64(args[0]), getFloat64(args[1]))),
	}
}

func sqrt(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Sqrt(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Sqrt(getFloat64(args[0]))),
	}
}

func floor(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Floor(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Floor(getFloat64(args[0]))),
	}
}

func ceil(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Ceil(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Ceil(getFloat64(args[0]))),
	}
}

func round(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Round(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Round(getFloat64(args[0]))),
	}
}

func clamp(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 3, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[2], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj || args[1].Type() == object.FloatObj || args[2].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Max(getFloat64(args[0]), math.Min(getFloat64(args[1]), getFloat64(args[2]))),
		}
	}

	return &object.Integer{
		Value: int64(math.Max(getFloat64(args[0]), math.Min(getFloat64(args[1]), getFloat64(args[2])))),
	}
}

func sign(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		if getFloat64(args[0]) > 0 {
			return &object.Float{Value: 1}
		} else if getFloat64(args[0]) < 0 {
			return &object.Float{Value: -1}
		} else {
			return &object.Float{Value: 0}
		}
	}

	if getFloat64(args[0]) > 0 {
		return &object.Integer{Value: 1}
	} else if getFloat64(args[0]) < 0 {
		return &object.Integer{Value: -1}
	} else {
		return &object.Integer{Value: 0}
	}

}

func sin(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Sin(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Sin(getFloat64(args[0]))),
	}
}

func cos(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Cos(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Cos(getFloat64(args[0]))),
	}
}

func tan(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	if args[0].Type() == object.FloatObj {
		return &object.Float{
			Value: math.Tan(getFloat64(args[0])),
		}
	}

	return &object.Integer{
		Value: int64(math.Tan(getFloat64(args[0]))),
	}
}

func getFloat64(obj object.Object) float64 {
	switch val := obj.(type) {
	case *object.Integer:
		return float64(val.Value)
	case *object.Float:
		return val.Value
	default:
		return 0
	}
}
