package repl

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"

	"github.com/peterh/liner"
)

const PROMPT = "$ "

var history_fn = filepath.Join(os.TempDir(), ".sunbird-history")

func Start(in io.Reader, out io.Writer) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	loadHistory(line)

	setupCompleter(line)

	env := object.NewEnvironment()

	replLoop(line, env, out)

}

func loadHistory(line *liner.State) {
	if f, err := os.Open(history_fn); err == nil {
		_, _ = line.ReadHistory(f)

		_ = f.Close()
	}
}

func saveHistory(line *liner.State) {
	f, _ := os.Create(history_fn)
	line.WriteHistory(f)
	_ = f.Close()
}

func setupCompleter(line *liner.State) {
	keywords := []string{"func", "var", "true", "false", "if", "else", "return", "null"}

	line.SetCompleter(func(input string) []string {
		var completions []string
		for _, keyword := range keywords {
			if strings.HasPrefix(keyword, input) {
				completions = append(completions, keyword)
			}
		}
		return completions
	})
}

func replLoop(line *liner.State, env *object.Environment, out io.Writer) {
	for {
		input, err := line.Prompt(PROMPT)
		if err != nil && err == liner.ErrPromptAborted {
			saveHistory(line)
		}

		if input == "exit" {
			f, _ := os.Create(history_fn)

			_, _ = line.WriteHistory(f)
			_ = f.Close()
			return
		}

		line.AppendHistory(input)

		evalInput(input, env, out)
	}
}

func evalInput(input string, env *object.Environment, out io.Writer) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}

}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
