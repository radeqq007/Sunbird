package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"sunbird/internal/pkg"
	"sunbird/internal/repl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Welcome to the sunbird programming language!")
		fmt.Printf("Type in 'exit' to exit.\n")
		repl.Start(os.Stdin, os.Stdout)

		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit()
	case "install", "i":
		handleInstall()
	case "add":
		handleAdd()
	case "update":
		handleUpdate()
	case "tidy":
		handleTidy()
	case "run":
		handleRun()
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		printVersion()
	}

	os.Exit(0)
}

func handleRun() {
	var filePath string

	if len(os.Args) >= 3 {
		filePath = os.Args[2]
	} else {
		config, err := pkg.LoadConfig("sunbird.toml")
		if err == nil && config.Package.Main != "" {
			filePath = config.Package.Main
		} else {
			filePath, err = findMainFile()
			if err != nil {
				fmt.Println("Error: No file specified and no main file found")
				fmt.Println("Usage: sunbird run [file.sb]")
				os.Exit(1)
			}
		}
	}

	runFile(filePath)
}

func findMainFile() (string, error) {
	candidates := []string{
		"main.sb",
		"src/main.sb",
		"index.sb",
		"src/index.sb",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("no main file found")
}

func runFile(path string) {
	src, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	
	content, err := io.ReadAll(src)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	
	_ = src.Close()
	
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

	if !evaluated.IsNull() {
		if evaluated.IsError() {
			fmt.Println(evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func handleInit() {
	fmt.Println("Initializing new project...")

	if _, err := os.Stat("sunbird.toml"); err == nil {
		fmt.Println("Error: sunbird.toml already exists")
		os.Exit(1)
	}

	template := `[package]
name = "my_project"
version = "0.1.0"
description = "A Sunbird project"
authors = ["Your Name <you@example.com>"]
main = "./src/main.sb"

dependencies = []
`

	err := os.WriteFile("sunbird.toml", []byte(template), 0644)
	if err != nil {
		fmt.Printf("Error creating sunbird.toml: %s\n", err)
		os.Exit(1)
	}

	err = os.Mkdir("src", 0755)
	if err != nil {
		fmt.Printf("Error creating src directory: %s\n", err)
		os.Exit(1)
	}

	mainTemplate := `import "io"

io.println("Hello, sunbird!")
`

	err = os.WriteFile("./src/main.sb", []byte(mainTemplate), 0644)
	if err != nil {
		fmt.Printf("Error creating main.sb: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Created sunbird.toml")
	fmt.Println("✓ Created main.sb")
	fmt.Println("✓ Created src/ directory")
	fmt.Println("\nProject initialized successfully!")
}

func handleInstall() {

}

func handleAdd() {
	pkgManager := pkg.NewPackageManager()

	if len(os.Args) < 3 {
		fmt.Println("Provide a package url")
		os.Exit(1)
	}

	url := os.Args[2]
	err := pkgManager.Add(url)
	if err != nil {
		fmt.Printf("Error adding dependency: %s\n", err)
	}
}

func handleUpdate() {

}

func handleTidy() {

}

func printHelp() {
	help := `Sunbird - A dynamically-typed programming language

Usage:
  sunbird [command] [arguments]

Commands:
  init                Initialize a new Sunbird project
  install, i          Install dependencies from sunbird.toml
  get <package>       Download and install a specific package
  update              Update all dependencies
  tidy                Remove unused dependencies
  run <file>          Run a Sunbird file with package resolution
  help, -h, --help    Show this help message
  version, -v         Show version information

Running files:
  sunbird <file.sb>   Run a Sunbird file directly (without package resolution)
  sunbird             Start interactive REPL

Examples:
  sunbird init
  sunbird get github.com/user/package@v1.0.0
  sunbird install
  sunbird run main.sb
  sunbird main.sb

For more information, visit: https://github.com/radeqq007/sunbird
`

	fmt.Println(help)
}

func printVersion() {

}
