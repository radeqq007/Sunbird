package evaluator

import (
	"os"
	"strconv"
	"sunbird/internal/errors"
	"sunbird/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectOneOfTypes(0, 0, args[0], object.StringObj, object.ArrayObj)
			if err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			}

			return nil
		},
	},

	"append": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 2, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.ArrayObj)
			if err != nil {
				return err
			}

			arr := args[0].(*object.Array)

			newElements := append(arr.Elements, args[1:]...)

			return &object.Array{Elements: newElements}
		},
	},

	"type": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			return &object.String{Value: args[0].Type().String()}
		},
	},

	"string": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},

	"int": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return arg

			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}

			case *object.String:
				num, err := strconv.Atoi(arg.Value)
				if err != nil {
					return errors.NewTypeError(0, 0, "failed to convert string to int: %s", arg.Value)
				}
				return &object.Integer{Value: int64(num)}

			case *object.Boolean:
				if arg.Value {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}

			default:
				return errors.NewTypeError(0, 0, "argument to `int` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"float": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}

			case *object.Float:
				return arg

			case *object.String:
				num, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return errors.NewTypeError(0, 0, "failed to convert string to float: %s", arg.Value)
				}
				return &object.Float{Value: num}

			case *object.Boolean:
				if arg.Value {
					return &object.Float{Value: 1.0}
				}
				return &object.Float{Value: 0.0}

			default:
				return errors.NewTypeError(0, 0, "argument to `float` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"bool": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value == 0 {
					return &object.Boolean{Value: false}
				}
				return &object.Boolean{Value: true}

			case *object.Float:
				if arg.Value == 0.0 {
					return &object.Boolean{Value: false}
				}
				return &object.Boolean{Value: true}

			case *object.String:
				if arg.Value == "" {
					return &object.Boolean{Value: false}
				}
				return &object.Boolean{Value: true}

			case *object.Boolean:
				return arg

			default:
				return errors.NewTypeError(0, 0, "argument to `bool` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"exit": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			os.Exit(0)
			return nil
		},
	},

	"error": {
		Fn: func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			return errors.New(errors.RuntimeError, 0, 0, "%s", args[0].Inspect())
		},
	},
}
