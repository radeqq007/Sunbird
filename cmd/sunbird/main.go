package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"sunbird/internal/repl"
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: sunbird [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Welcome to the sunbird programming language!")
		fmt.Printf("Type in 'exit' to exit.\n")
		repl.Start(os.Stdin, os.Stdout)

		os.Exit(0)
	}

	if len(args) == 1 {
		src, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		defer func() {
			_ = src.Close()
		}()

		content, err := io.ReadAll(src)
		if err != nil {
			fmt.Printf("Error: %s\n", err)

			_ = src.Close() // Close file before exiting
			os.Exit(1)
		}

		l := lexer.New(string(content))
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			fmt.Println("Parser errors:")
			for _, msg := range p.Errors() {
				fmt.Printf("\t%s\n", msg)
			}
			os.Exit(1)
		}

		env := object.NewEnvironment()

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			if errObj, ok := evaluated.(*object.Error); ok {
				fmt.Println(errObj.Inspect())
				os.Exit(1)
			}
		}

		os.Exit(0)
	}
}
