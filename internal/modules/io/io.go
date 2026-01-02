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
		if s, ok := arg.(*object.String); ok {
			fmt.Print(s.Value)
		} else {
			fmt.Print(arg.Inspect())
		}
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
		err := errors.ExpectNumberOfArguments(0, 0, 1, args)
		if err != nil {
			return err
		}
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
		err := errors.ExpectNumberOfArguments(0, 0, 1, args)
		if err != nil {
			return err
		}
		err = errors.ExpectType(0, 0, args[0], object.StringObj)
		if err != nil {
			return err
		}
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
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}
	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	var format string
	if s, ok := args[0].(*object.String); ok {
		format = s.Value
	} else {
		format = args[0].Inspect()
	}


	for _, arg := range args[1:] {
		var val string
		if s, ok := arg.(*object.String); ok {
			val = s.Value
		} else {
			val = arg.Inspect()
		}
		format = strings.Replace(format, "{}", val, 1)
	}

	fmt.Print(format)

	return nil
}

func sprintf(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}
	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	var format string
	if s, ok := args[0].(*object.String); ok {
		format = s.Value
	} else {
		format = args[0].Inspect()
	}

	for _, arg := range args[1:] {
		var val string
		if s, ok := arg.(*object.String); ok {
			val = s.Value
		} else {
			val = arg.Inspect()
		}
		format = strings.Replace(format, "{}", val, 1)
	}

	return &object.String{
		Value: format,
	}
}

func clear(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
	}
	fmt.Print("\033[H")
	fmt.Print("\033[2J")
	return nil
}

func beep(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
	}
	fmt.Print("\a")
	return nil
}

func getArgsArray(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
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
