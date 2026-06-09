package repl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/modules"
	"sunbird/internal/object"
	"sunbird/internal/parser"

	"github.com/c-bata/go-prompt"
)

const PROMPT = "$ "

var keywords = []prompt.Suggest{
	{Text: "func", Description: "Define a function"},
	{Text: "let", Description: "Variable declaration"},
	{Text: "const", Description: "Variable declaration"},
	{Text: "true", Description: "Boolean true"},
	{Text: "false", Description: "Boolean false"},
	{Text: "if", Description: "Conditional statement"},
	{Text: "else", Description: "Conditional alternative"},
	{Text: "return", Description: "Return value"},
	{Text: "import", Description: "Import a module"},
	{Text: "as", Description: "Alias for imported module"},
	{Text: "for", Description: "Loop statement"},
	{Text: "while", Description: "While loop statement"},
	{Text: "null", Description: "Null value"},
	{Text: "break", Description: "Break out of a loop"},
	{Text: "continue", Description: "Continue to the next iteration of a loop"},
	{Text: "Int", Description: "Integer type"},
	{Text: "Float", Description: "Float type"},
	{Text: "String", Description: "String type"},
	{Text: "Bool", Description: "Boolean type"},
	{Text: "Void", Description: "Void type"},
	{Text: "Array", Description: "Array type"},
	{Text: "Hash", Description: "Hash type"},
	{Text: "Func", Description: "Function type"},
	{Text: "try", Description: "Try block for exception handling"},
	{Text: "catch", Description: "Catch block for exception handling"},
	{Text: "finally", Description: "Finally block for exception handling"},
	{Text: "in", Description: "Iteration keyword"},
	{Text: "exit", Description: "Exit the REPL"},
}

func Start(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()

	executor := func(input string) {
		input = strings.TrimSpace(input)
		if input == "" {
			return
		}

		if input == "exit" {
			os.Exit(0)
		}

		if err := EvalInput(input, env, out); err != nil {
			// REPL probably can’t recover meaningfully, but log it
			fmt.Fprintln(os.Stderr, "write error:", err)
		}
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		word := d.GetWordBeforeCursor()
		textBefore := d.TextBeforeCursor()

		// Avoid suggesting keywords after a dot (.)
		if strings.HasSuffix(textBefore, ".") {
			return []prompt.Suggest{}
		}

		// Suggest defined modules after 'import' keyword
		// TODO: suggest module names also after 'import "' etc.
		if strings.HasSuffix(textBefore, "import ") {
			var moduleSuggestions []prompt.Suggest
			for mod := range modules.BuiltinModules {
				moduleSuggestions = append(moduleSuggestions, prompt.Suggest{Text: "\"" + mod + "\"", Description: "Imported module"})
			}
			return prompt.FilterHasPrefix(moduleSuggestions, word, true)
		}

		return prompt.FilterHasPrefix(keywords, d.GetWordBeforeCursor(), true)
	}

	// highlight := func(input string) []prompt. {
	// return sunbirdLexer(input)
	// }

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("$ "),
		prompt.OptionTitle("sunbird-repl"),
	)

	p.Run()
}

func EvalInput(input string, env *object.Environment, out io.Writer) error {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return printParserErrors(out, p.Errors())
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated.IsNull() {
		return nil
	}

	if _, err := io.WriteString(out, evaluated.Inspect()); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "\n"); err != nil {
		return err
	}

	return nil
}

func printParserErrors(out io.Writer, errors []string) error {
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			return err
		}
	}
	return nil
}
