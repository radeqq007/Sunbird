package math

import (
	"math"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
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
		AddFloat("pi", math.Pi).
		AddFloat("e", math.E).
		Build()
}

func abs(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].Kind() == object.FloatKind {
		return object.NewFloat(math.Abs(args[0].AsFloat()))
	}

	return object.NewInt(int64(math.Abs(float64(args[0].AsInt()))))
}

func max(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() || args[1].IsFloat() {
		return object.NewFloat(math.Max(getFloat64(args[0]), getFloat64(args[1])))
	}

	return object.NewInt(int64(math.Max(getFloat64(args[0]), getFloat64(args[1]))))
}

func min(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() || args[1].IsFloat() {
		return object.NewFloat(math.Min(getFloat64(args[0]), getFloat64(args[1])))
	}

	return object.NewInt(int64(math.Min(getFloat64(args[0]), getFloat64(args[1]))))
}

func pow(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() || args[1].IsFloat() {
		return object.NewFloat(math.Pow(getFloat64(args[0]), getFloat64(args[1])))
	}

	return object.NewInt(int64(math.Pow(getFloat64(args[0]), getFloat64(args[1]))))
}

func sqrt(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Sqrt(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Sqrt(getFloat64(args[0]))))
}

func floor(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Floor(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Floor(getFloat64(args[0]))))
}

func ceil(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Ceil(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Ceil(getFloat64(args[0]))))
}

func round(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Round(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Round(getFloat64(args[0]))))
}

func clamp(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 3, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[1], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[2], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() || args[1].IsFloat() || args[2].IsFloat() {
		if args[2].IsFloat() {
			return object.NewFloat(math.Max(
				getFloat64(args[0]),
				math.Min(getFloat64(args[1]), getFloat64(args[2])),
			))
		}
	}

	return object.NewInt(int64(math.Max(
		getFloat64(args[0]),
		math.Min(getFloat64(args[1]), getFloat64(args[2])),
	)))
}

func sign(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		floatVar := getFloat64(args[0])
		switch {
		case floatVar > 0:
			return object.NewFloat(1)
		case floatVar < 0:
			return object.NewFloat(-1)
		default:
			return object.NewFloat(0)
		}
	}

	floatVar := getFloat64(args[0])
	switch {
	case floatVar > 0:
		return object.NewInt(1)
	case floatVar < 0:
		return object.NewInt(-1)
	default:
		return object.NewInt(0)
	}
}

func sin(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Sin(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Sin(getFloat64(args[0]))))
}

func cos(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Cos(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Cos(getFloat64(args[0]))))
}

func tan(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Tan(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Tan(getFloat64(args[0]))))
}

func getFloat64(obj object.Value) float64 {
	switch obj.Kind() {
	case object.IntKind:
		return float64(obj.AsInt())
	case object.FloatKind:
		return obj.AsFloat()
	default:
		return 0
	}
}
