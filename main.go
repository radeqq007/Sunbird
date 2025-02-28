package main

import (
	"fmt"
	"os"
	"sunbird/repl"
)

func main() {
	fmt.Println("Welcome to the sunbird programming language!")
	fmt.Printf("Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}