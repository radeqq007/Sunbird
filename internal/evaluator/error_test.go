package evaluator_test

import (
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"testing"
)

func TestErrorLineNumbers(t *testing.T) {
	tests := []struct {
		input        string
		expectedLine int
		expectedCol  int
		expectedMsg  string
	}{
		{
			"5 + true;",
			1, 3,
			"TypeMismatchError: INTEGER + BOOLEAN",
		},
		{
			"foobar",
			1, 1,
			"UndefinedVariableError: foobar",
		},
		{
			"-true",
			1, 1,
			"UnknownOperatorError: -BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			1, 20,
			"UnknownOperatorError: BOOLEAN + BOOLEAN",
		},
		{
			`
var a = 5;
a + true;
`,
			3, 3,
			"TypeMismatchError: INTEGER + BOOLEAN",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for input: %q. got=%T(%+v)",
				tt.input, evaluated, evaluated)
			continue
		}

		if errObj.Line != tt.expectedLine {
			t.Errorf("wrong line number for input: %q. expected=%d, got=%d",
				tt.input, tt.expectedLine, errObj.Line)
		}

		if errObj.Col != tt.expectedCol {
			t.Errorf("wrong column number for input: %q. expected=%d, got=%d",
				tt.input, tt.expectedCol, errObj.Col)
		}

		if errObj.Message != tt.expectedMsg {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMsg, errObj.Message)
		}
	}
}
