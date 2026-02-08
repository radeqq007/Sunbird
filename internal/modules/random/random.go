package random

import (
	"math/rand/v2"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
	"time"
)

func New() *object.Hash {
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

func newSeed(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
	if err != nil {
		return err
	}

	seed = uint64(args[0].(*object.Integer).Value)
	r = rand.New(rand.NewPCG(seed, seed))

	return &object.Null{}
}

func randInt(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.IntegerObj)
	if err != nil {
		return err
	}

	minVal := args[0].(*object.Integer).Value
	maxVal := args[1].(*object.Integer).Value

	return &object.Integer{Value: r.Int64N(maxVal-minVal) + minVal}
}

func randFloat(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.FloatObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.FloatObj)
	if err != nil {
		return err
	}

	minVal := args[0].(*object.Float).Value
	maxVal := args[1].(*object.Float).Value

	return &object.Float{Value: r.Float64()*(maxVal-minVal) + minVal}
}

func randBool(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
	}

	return &object.Boolean{Value: r.Int64N(2) == 1}
}

func choice(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	return args[0].(*object.Array).Elements[r.Int64N(int64(len(args[0].(*object.Array).Elements)))]
}

func shuffle(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
	if err != nil {
		return err
	}

	arr, _ := args[0].(*object.Array)
	shuffled := make([]object.Object, len(arr.Elements))
	copy(shuffled, arr.Elements)

	for i := range shuffled {
		j := r.Int64N(int64(i + 1))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return &object.Array{Elements: shuffled}
}
