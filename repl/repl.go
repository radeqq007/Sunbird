package repl

import (
	"io"
	"os"
	"path/filepath"
	"sunbird/evaluator"
	"sunbird/lexer"
	"sunbird/object"
	"sunbird/parser"

	"github.com/peterh/liner"
)

const PROMPT = "$ "

var history_fn = filepath.Join(os.TempDir(), ".sunbird-history")

func Start(in io.Reader, out io.Writer) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	env := object.NewEnvironment()

	for {
		if input, err := line.Prompt(PROMPT); err == nil {

			if input == "exit" {

				if f, err := os.Create(history_fn); err == nil {
					line.WriteHistory(f)
					f.Close()
				}
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

		} else if err == liner.ErrPromptAborted {
			if f, err := os.Create(history_fn); err == nil {
				line.WriteHistory(f)
				f.Close()
			}
			return
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
