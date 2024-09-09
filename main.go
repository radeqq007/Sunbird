package main

import (
	"fmt"
	"os"
	"vex-programming-language/repl"
)

func main() {
	fmt.Println("Welcome to the Vex programming language!")
	fmt.Printf("Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}