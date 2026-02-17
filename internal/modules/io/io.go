package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("print", printObject).
		AddFunction("println", printlnObject).
		AddFunction("readln", readln).
		AddFunction("read", read).
		AddFunction("printf", printf).
		AddFunction("printfn", printfn).
		AddFunction("sprintf", sprintf).
		AddFunction("clear", clearScreen).
		AddFunction("beep", beep).
		AddValue("args", getArgsArray()).
		Build()
}

var stdin = bufio.NewReader(os.Stdin)

func printObject(ctx object.CallContext, args ...object.Value) object.Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		if arg.IsString() {
			fmt.Print(arg.AsString().Value)
		} else {
			fmt.Print(arg.Inspect())
		}
	}
	return object.NewNull()
}

func printlnObject(ctx object.CallContext, args ...object.Value) object.Value {
	printObject(ctx, args...)
	fmt.Print("\n")
	return object.NewNull()
}

func readln(ctx object.CallContext, args ...object.Value) object.Value {
	if len(args) > 0 {
		err := errors.ExpectNumberOfArguments(0, 0, 1, args)
		if err.IsError() {
			return err
		}
		printObject(ctx, args[0])
	}

	input, err := stdin.ReadString('\n')
	if err != nil && err != io.EOF {
		return object.NewNull()
	}

	return object.NewString(strings.TrimRight(input, "\r\n"))
}

func read(ctx object.CallContext, args ...object.Value) object.Value {
	if len(args) > 0 {
		err := errors.ExpectNumberOfArguments(0, 0, 1, args)
		if err.IsError() {
			return err
		}

		err = errors.ExpectType(0, 0, args[0], object.StringKind)
		if err.IsError() {
			return err
		}
		printObject(ctx, args[0])
	}

	input, err := stdin.ReadString(' ')
	if err != nil && err != io.EOF {
		return object.NewNull()
	}

	return object.NewString(strings.TrimSpace(input))
}

func printf(ctx object.CallContext, args ...object.Value) object.Value {
	if len(args) < 1 {
		return errors.NewArgumentError(0, 0, "expected minimum 1 argument, got %v", len(args))
	}

	err := errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	var format string
	if args[0].IsString() {
		format = args[0].AsString().Value
	} else {
		format = args[0].Inspect()
	}

	for _, arg := range args[1:] {
		var val string
		if arg.IsString() {
			val = arg.AsString().Value
		} else {
			val = arg.Inspect()
		}
		format = strings.Replace(format, "{}", val, 1)
	}

	fmt.Print(format)

	return object.NewNull()
}

func printfn(ctx object.CallContext, args ...object.Value) object.Value {
	if len(args) < 1 {
		return errors.NewArgumentError(0, 0, "expected minimum 1 argument, got %v", len(args))
	}

	err := errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	var format string
	if args[0].IsString() {
		format = args[0].AsString().Value
	} else {
		format = args[0].Inspect()
	}

	for _, arg := range args[1:] {
		var val string
		if arg.IsString() {
			val = arg.AsString().Value
		} else {
			val = arg.Inspect()
		}
		format = strings.Replace(format, "{}", val, 1)
	}

	fmt.Println(format)

	return object.NewNull()
}

func sprintf(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}
	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	var format string
	if args[0].IsString() {
		format = args[0].AsString().Value
	} else {
		format = args[0].Inspect()
	}

	for _, arg := range args[1:] {
		var val string
		if arg.IsString() {
			val = arg.AsString().Value
		} else {
			val = arg.Inspect()
		}
		format = strings.Replace(format, "{}", val, 1)
	}

	return object.NewString(format)
}

func clearScreen(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}
	fmt.Print("\033[H")
	fmt.Print("\033[2J")
	return object.NewNull()
}

func beep(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}
	fmt.Print("\a")
	return object.NewNull()
}

func getArgsArray(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	var elements []object.Value
	userScriptIndex := 2
	if len(os.Args) > userScriptIndex {
		osArgs := os.Args[userScriptIndex:]
		for _, arg := range osArgs {
			elements = append(elements, object.NewString(arg))
		}
	}

	return object.NewArray(elements)
}
