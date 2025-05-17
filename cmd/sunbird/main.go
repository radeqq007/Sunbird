package main

import (
	"flag"
	"fmt"
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
		defer src.Close()

		content := make([]byte, 100)
		_, err = src.Read(content)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
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
			fmt.Println(evaluated.Inspect())
		}

		os.Exit(0)
	}
}
