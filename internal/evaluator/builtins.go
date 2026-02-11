package evaluator

import (
	"os"
	"strconv"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var builtins = map[string]object.Value{
	"len": object.NewBuiltin(
		func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			err = errors.ExpectOneOfTypes(0, 0, args[0], object.StringKind, object.ArrayKind)
			if !err.IsNull() {
				return err
			}

			// switch arg := args[0].(type) {
			// case *object.String:
			// 	return &object.Integer{Value: int64(len(arg.Value))}

			// case *object.Array:
			// 	return &object.Integer{Value: int64(len(arg.Elements))}
			// }

			arg := args[0]
			if arg.IsString() {
				return object.NewInt(int64(len(arg.AsString().Value)))
			}

			if arg.IsArray() {
				return object.NewInt(int64(len(arg.AsArray().Elements)))
			}

			return NULL
		},
	),

	"append": object.NewBuiltin(
		func(args ...object.Value) object.Value {
			err := errors.ExpectMinNumberOfArguments(0, 0, 2, args)
			if !err.IsNull() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.ArrayKind)
			if !err.IsNull() {
				return err
			}

			arr := args[0].AsArray()

			newElements := append(arr.Elements, args[1:]...)

			return object.NewArray(newElements)
		},
	),

	"type": object.NewBuiltin(
		func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			return object.NewString(args[0].Kind().String())
		},
	),

	"string": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			if args[0].IsString() {
				val := args[0].AsString().Value
				return object.NewString(val)
			}

			return object.NewString(args[0].Inspect())
		},
	),

	"int": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			switch args[0].Kind() {
			case object.IntKind:
				return args[0]

			case object.FloatKind:
				arg := args[0].AsFloat()
				return object.NewInt(int64(arg))

			case object.StringKind:
				arg := args[0].AsString().Value
				num, err := strconv.Atoi(arg)
				if err != nil {
					return errors.NewTypeError(0, 0, "failed to convert string to int: %s", arg)
				}
				return object.NewInt(int64(num))

			case object.BoolKind:
				arg := args[0].AsBool()
				if arg {
					return object.NewInt(1)
				}
				return object.NewInt(0)

			default:
				return errors.NewTypeError(0, 0, "argument to `int` not supported, got %s", args[0].Kind().String())
			}
		},
	),

	"float": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			switch args[0].Kind() {
			case object.IntKind:
				arg := args[0].AsInt()
				return object.NewFloat(float64(arg))

			case object.FloatKind:
				return args[0]

			case object.StringKind:
				arg := args[0].AsString().Value
				num, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return errors.NewTypeError(0, 0, "failed to convert string to float: %s", arg)
				}
				return object.NewFloat(num)

			case object.BoolKind:
				arg := args[0].AsBool()
				if arg {
					return object.NewFloat(1.0)
				}
				return object.NewFloat(0.0)

			default:
				return errors.NewTypeError(0, 0, "argument to `float` not supported, got %s", args[0].Kind().String())
			}
		},
	),

	"bool": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			switch args[0].Kind() {
			case object.IntKind:
				arg := args[0].AsInt()
				if arg == 0 {
					return object.NewBool(false)
				}
				return object.NewBool(true)

			case object.FloatKind:
				arg := args[0].AsFloat()
				if arg == 0.0 {
					return object.NewBool(false)
				}
				return object.NewBool(true)

			case object.StringKind:
				arg := args[0].AsString().Value
				if arg == "" {
					return object.NewBool(false)
				}
				return object.NewBool(true)

			case object.BoolKind:
				return args[0]

			default:
				return errors.NewTypeError(0, 0, "argument to `bool` not supported, got %s", args[0].Kind().String())
			}
		},
	),

	"exit": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if !err.IsNull() {
				return err
			}

			os.Exit(0)
			return object.NewNull()
		},
	),

	"error": object.NewBuiltin(func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if !err.IsNull() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if !err.IsNull() {
				return err
			}

			var msg string
			if args[0].IsString() {
				msg = args[0].AsString().Value
			} else {
				msg = args[0].Inspect()
			}

			return errors.New(errors.RuntimeError, 0, 0, "%s", msg)
		},
	),
}
