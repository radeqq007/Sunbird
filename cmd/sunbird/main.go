package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/radeqq007/sunbird/internal/lexer"
	"github.com/radeqq007/sunbird/internal/parser"
	"github.com/radeqq007/sunbird/internal/pkg"
	"github.com/radeqq007/sunbird/internal/transpiler"
)

func main() {
	if len(os.Args) < 2 {
		// TODO: do something
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit()
	case "build":
		handleBuild()
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		printVersion()
	}

	os.Exit(0)
}

func resolveFilePath() (string, error) {
	// Check CLI Arguments
	if len(os.Args) >= 3 {
		return os.Args[2], nil
	}

	// Check Config File
	if config, err := pkg.LoadConfig("sunbird.toml"); err == nil && config.Package.Main != "" {
		return config.Package.Main, nil
	}

	// Fallback to searching
	return findMainFile()
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
`

	err := os.WriteFile("sunbird.toml", []byte(template), 0o644)
	if err != nil {
		fmt.Printf("Error creating sunbird.toml: %s\n", err)
		os.Exit(1)
	}

	err = os.Mkdir("src", 0o755)
	if err != nil {
		fmt.Printf("Error creating src directory: %s\n", err)
		os.Exit(1)
	}

	mainTemplate := `import "io"

io.println("Hello, sunbird!")
`

	err = os.WriteFile("./src/main.sb", []byte(mainTemplate), 0o644)
	if err != nil {
		fmt.Printf("Error creating main.sb: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Created sunbird.toml")
	fmt.Println("✓ Created main.sb")
	fmt.Println("✓ Created src/ directory")
	fmt.Println("\nProject initialized successfully!")
}

func printHelp() {
	help := `Sunbird - A dynamically-typed programming language

Usage:
  sunbird [command] [arguments]

Commands:
  build [file]        Transpile to TypeScript
  init                Initialize a new Sunbird project
  help, -h, --help    Show this help message
  version, -v         Show version information

Examples:
  sunbird init
  sunbird build main.sb

For more information, visit: https://github.com/radeqq007/sunbird
`

	fmt.Println(help)
}

func printVersion() {

}

func handleBuild() {
	filePath, err := resolveFilePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No file specified")
		os.Exit(1)
	}

	src, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "\t%s\n", msg)
		}
		os.Exit(1)
	}

	t := transpiler.New()
	output, err := t.Transpile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Transpile error: %s\n", err)
		os.Exit(1)
	}

	outDir := filepath.Dir(filePath)
	outFile := strings.TrimSuffix(filepath.Base(filePath), ".sb") + ".ts"

	if err := os.WriteFile(filepath.Join(outDir, outFile), []byte(output), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %s\n", err)
		os.Exit(1)
	}

	if err := transpiler.WriteRuntime(outDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing runtime: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ %s → %s\n", filePath, outFile)
}
