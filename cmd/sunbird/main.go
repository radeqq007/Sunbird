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
	"github.com/radeqq007/sunbird/internal/repl"
	"github.com/radeqq007/sunbird/internal/transpiler"
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

dependencies = []
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
	build [file] [--target TARGET]   Transpile to TypeScript
                                   Targets: node (default), deno, bun, web
  init                Initialize a new Sunbird project
  install, i          Install dependencies from sunbird.toml
  get <package>       Download and install a specific package
  update              Update all dependencies
  tidy                Remove unused dependencies
  help, -h, --help    Show this help message
  version, -v         Show version information

Examples:
  sunbird init
  sunbird build main.sb
  sunbird build main.sb --target deno
  sunbird build main.sb -t bun
  sunbird add github.com/user/package@v1.0.0
  sunbird install

For more information, visit: https://github.com/radeqq007/sunbird
`

	fmt.Println(help)
}

func printVersion() {

}


func handleBuild() {
	target := transpiler.DefaultTarget
	var extraArgs []string

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--target" || arg == "-t":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: %s requires a value\n", arg)
				os.Exit(1)
			}
			i++
			t, err := transpiler.ParseTarget(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			target = t
		case len(arg) > 9 && arg[:9] == "--target=":
			t, err := transpiler.ParseTarget(arg[9:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			target = t
		default:
			extraArgs = append(extraArgs, arg)
		}
	}

	// Temporarily patch os.Args so resolveFilePath works (it reads os.Args[2]).
	// We rebuild the slice so the file path, if supplied, lands at index 2.
	if len(extraArgs) > 0 {
		os.Args = append(os.Args[:2], extraArgs...)
	}

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

	if err := transpiler.WriteRuntime(outDir, target); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing runtime: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ %s → %s  [target: %s]\n", filePath, outFile, target)
}
