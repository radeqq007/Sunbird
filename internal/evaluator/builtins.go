package evaluator

import (
	"fmt"
	"sunbird/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError(0, 0, "argument to `len` not supported, got %s", args[0].Type().String())
			}
		},
	},

	"append": {
		Fn: func(args ...object.Object) object.Object {
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError(
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
				return newError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &object.String{Value: args[0].Type().String()}
		},
	},

	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},
}
