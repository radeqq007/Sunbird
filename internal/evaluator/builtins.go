package evaluator

import (
	"fmt"
	"os"
	"strconv"
	"sunbird/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return NewError(0, 0, "argument to `len` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"append": {
		Fn: func(args ...object.Object) object.Object {
			arr, ok := args[0].(*object.Array)
			if !ok {
				return NewError(
					0, 0,
					"first argument to `append` must be an array, got %s",
					args[0].Type().String(),
				)
			}

			newElements := append(arr.Elements, args[1:]...)

			return &object.Array{Elements: newElements}
		},
	},

	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()

			return nil
		},
	},

	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}

			return nil
		},
	},

	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &object.String{Value: args[0].Type().String()}
		},
	},

	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},

	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return arg

			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}

			case *object.String:
				num, err := strconv.Atoi(arg.Value)
				if err != nil {
					return NewError(0, 0, "failed to convert string to int: %s", arg.Value)
				}
				return &object.Integer{Value: int64(num)}

			case *object.Boolean:
				if arg.Value {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}

			default:
				return NewError(0, 0, "argument to `int` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}

			case *object.Float:
				return arg

			case *object.String:
				num, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return NewError(0, 0, "failed to convert string to float: %s", arg.Value)
				}
				return &object.Float{Value: num}

			case *object.Boolean:
				if arg.Value {
					return &object.Float{Value: 1.0}
				}
				return &object.Float{Value: 0.0}

			default:
				return NewError(0, 0, "argument to `float` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"bool": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
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
				return NewError(0, 0, "argument to `bool` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"exit": {
		Fn: func(args ...object.Object) object.Object {
			os.Exit(0)
			return nil
		},
	},

	"error": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.StringObj {
				return NewError(0, 0, "argument to `error` must be a string, got %s", args[0].Type().String())
			}

			return NewError(0, 0, "%s", args[0].Inspect())
		},
	},
}
