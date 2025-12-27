package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
	AddFunction("print", print).
	AddFunction("println", println).
	AddFunction("readln", readln).
	AddFunction("read", read).
	AddFunction("printf", printf).
	AddFunction("sprintf", sprintf).
	AddFunction("clear", clear).
	AddFunction("beep", beep).
	AddValue("args", getArgsArray()).
	Build()

var stdin = bufio.NewReader(os.Stdin)

func print(args ...object.Object) object.Object {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.Inspect())
	}
	return nil
}

func println(args ...object.Object) object.Object {
	print(args...)
	fmt.Print("\n")
	return nil
}

func readln(args ...object.Object) object.Object {
	if len(args) > 0 {
		print(args[0])
	}

	input, err := stdin.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil
	}

	return &object.String{
		Value: strings.TrimRight(input, "\r\n"),
	}
}

func read(args ...object.Object) object.Object {
	if len(args) > 0 {
		print(args[0])
	}

	input, err := stdin.ReadString(' ')
	if err != nil && err != io.EOF {
		return nil
	}

	return &object.String{
		Value: strings.TrimSpace(input),
	}
}

func printf(args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.StringObj {
		return object.NewError(0, 0, "argument must be a string, got %s", args[0].Type().String())
	}

	format := args[0].Inspect()

	for _, arg := range args[1:] {
		format = strings.Replace(format, "{}", arg.Inspect(), 1)
	}

	fmt.Print(format)

	return nil
}

func sprintf(args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.StringObj {
		return object.NewError(0, 0, "argument must be a string, got %s", args[0].Type().String())
	}

	format := args[0].Inspect()

	for _, arg := range args[1:] {
		format = strings.Replace(format, "{}", arg.Inspect(), 1)
	}

	return &object.String{
		Value: format,
	}
}

func clear(args ...object.Object) object.Object {
	if len(args) > 0 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=0", len(args))
	}
	fmt.Print("\033[H")
	fmt.Print("\033[2J")
	return nil
}

func beep(args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=0", len(args))
	}
	fmt.Print("\a")
	return nil
}

func getArgsArray(args ...object.Object) object.Object {
	if len(args) > 0 {
		return object.NewError(0, 0, "wrong number of arguments. got=%d, want=0", len(args))
	}

	var elements []object.Object
	if len(os.Args) > 2 {
		osArgs := os.Args[2:]
		for _, arg := range osArgs {
			elements = append(elements, &object.String{Value: arg})
		}
	}

	return &object.Array{
		Elements: elements,
	}
}
