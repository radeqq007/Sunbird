package random

import (
	"math/rand"
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

var seed = time.Now().UnixNano()

var r = rand.New(rand.NewSource(seed))

func newSeed(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
	if err != nil {
		return err
	}

	seed = args[0].(*object.Integer).Value
	r = rand.New(rand.NewSource(seed))

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

	min := args[0].(*object.Integer).Value
	max := args[1].(*object.Integer).Value

	return &object.Integer{Value: r.Int63n(max-min) + min}
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

	min := args[0].(*object.Float).Value
	max := args[1].(*object.Float).Value

	return &object.Float{Value: r.Float64()*(max-min) + min}
}

func randBool(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
	}

	return &object.Boolean{Value: r.Intn(2) == 1}
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

	return args[0].(*object.Array).Elements[r.Intn(len(args[0].(*object.Array).Elements))]
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
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return &object.Array{Elements: shuffled}
}
