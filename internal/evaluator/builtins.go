package evaluator

import (
	"fmt"
	"sunbird/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to `append` must be an array, got %s", args[0].Type())
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
}
