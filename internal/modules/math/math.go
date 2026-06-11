package math

import (
	"github.com/radeqq007/sunbird/internal/errors"
	"github.com/radeqq007/sunbird/internal/modules/modbuilder"
	"github.com/radeqq007/sunbird/internal/object"
	"math"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("abs", abs).
		AddFunction("max", maxValue).
		AddFunction("min", minValue).
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

// evalBinaryNumeric validates two numeric args, calls fn(a, b), and returns
// Float if either arg is float, Int otherwise.
func evalBinaryNumeric(
	ctx object.CallContext,
	args []object.Value,
	fn func(a, b float64) float64,
) object.Value {
	if err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args); err.IsError() {
		return err
	}
	if err := errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind); err.IsError() {
		return err
	}
	if err := errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[1], object.IntKind, object.FloatKind); err.IsError() {
		return err
	}
 
	result := fn(getFloat64(args[0]), getFloat64(args[1]))
 
	if args[0].IsFloat() || args[1].IsFloat() {
		return object.NewFloat(result)
	}
	return object.NewInt(int64(result))
}


func abs(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].Kind() == object.FloatKind {
		return object.NewFloat(math.Abs(args[0].AsFloat()))
	}

	return object.NewInt(int64(math.Abs(float64(args[0].AsInt()))))
}

func maxValue(ctx object.CallContext, args ...object.Value) object.Value {
	return evalBinaryNumeric(ctx, args, math.Max)
}

func minValue(ctx object.CallContext, args ...object.Value) object.Value {
	return evalBinaryNumeric(ctx, args, math.Min)
}

func pow(ctx object.CallContext, args ...object.Value) object.Value {
	return evalBinaryNumeric(ctx, args, math.Pow)
}

func sqrt(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Sqrt(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Sqrt(getFloat64(args[0]))))
}

func floor(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Floor(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Floor(getFloat64(args[0]))))
}

func ceil(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Ceil(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Ceil(getFloat64(args[0]))))
}

func round(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Round(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Round(getFloat64(args[0]))))
}

func clamp(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 3, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[1], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[2], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	val := getFloat64(args[0])
	minVal := getFloat64(args[1])
	maxVal := getFloat64(args[2])

	res := math.Max(minVal, math.Min(val, maxVal))

	if args[0].IsFloat() || args[1].IsFloat() || args[2].IsFloat() {
		return object.NewFloat(res)
	}

	return object.NewInt(int64(res))
}

func sign(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
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

func sin(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Sin(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Sin(getFloat64(args[0]))))
}

func cos(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
	if err.IsError() {
		return err
	}

	if args[0].IsFloat() {
		return object.NewFloat(math.Cos(getFloat64(args[0])))
	}

	return object.NewInt(int64(math.Cos(getFloat64(args[0]))))
}

func tan(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(ctx.Line, ctx.Col, args[0], object.IntKind, object.FloatKind)
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
