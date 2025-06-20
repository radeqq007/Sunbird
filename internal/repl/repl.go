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

	if f, err := os.Open(history_fn); err == nil {
		_, _ = line.ReadHistory(f)

		f.Close()
	}

	keywords := []string{"func", "var", "true", "false", "if", "else", "return", "null"}

	line.SetCompleter(func(line string) (c []string) {
		for _, keyword := range keywords {
			if strings.HasPrefix(keyword, line) {
				c = append(c, keyword)
			}
		}
		return
	})

	env := object.NewEnvironment()

	for {
		input, err := line.Prompt(PROMPT)
		if err != nil && err == liner.ErrPromptAborted {
			f, _ := os.Create(history_fn)
			line.WriteHistory(f)
			f.Close()
		}

		if input == "exit" {
			f, _ := os.Create(history_fn)

			_, _ = line.WriteHistory(f)
			f.Close()
			return
		}

		line.AppendHistory(input)

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
