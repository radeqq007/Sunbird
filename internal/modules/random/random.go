package random

import (
	"math/rand/v2"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
	"time"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("int", randInt).
		AddFunction("float", randFloat).
		AddFunction("bool", randBool).
		AddFunction("shuffle", shuffle).
		AddFunction("choice", choice).
		AddFunction("seed", newSeed).
		Build()
}

var seed uint64 = uint64(time.Now().UnixNano())

var r = rand.New(rand.NewPCG(seed, seed))

func newSeed(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.IntKind)
	if err.IsError() {
		return err
	}

	seed = uint64(args[0].AsInt())
	r = rand.New(rand.NewPCG(seed, seed))

	return object.NewNull()
}

func randInt(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.IntKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.IntKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.IntKind)
	if err.IsError() {
		return err
	}

	minVal := args[0].AsInt()
	maxVal := args[1].AsInt()

	return object.NewInt(r.Int64N(maxVal-minVal) + minVal)
}

func randFloat(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.FloatKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[1], object.FloatKind)
	if err.IsError() {
		return err
	}

	minVal := args[0].AsFloat()
	maxVal := args[1].AsFloat()

	return object.NewFloat(r.Float64()*(maxVal-minVal) + minVal)
}

func randBool(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 0, args)
	if err.IsError() {
		return err
	}

	return object.NewBool(r.Int64N(2) == 1)
}

func choice(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	arr := args[0].AsArray()
	return arr.Elements[r.Int64N(int64(len(arr.Elements)))]
}

func shuffle(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(ctx.Line, ctx.Col, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(ctx.Line, ctx.Col, args[0], object.ArrayKind)
	if err.IsError() {
		return err
	}

	arr := args[0].AsArray()
	shuffled := make([]object.Value, len(arr.Elements))
	copy(shuffled, arr.Elements)

	for i := range shuffled {
		j := r.Int64N(int64(i + 1))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return object.NewArray(shuffled)
}
