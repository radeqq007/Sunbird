package repl_test

import (
	"bytes"
	"strings"
	"testing"

	"sunbird/internal/object"
	"sunbird/internal/repl"
)

func TestEvalInput(t *testing.T) {
	env := object.NewEnvironment()
	tests := []struct {
		input    string
		expected string
	}{
		{"let x = 5;", "5"},
		{"x = 10; x;", "10"},
		{"true;", "true"},
		{"false;", "false"},
		{"1 + 2;", "3"},
	}

	for _, tt := range tests {
		out := &bytes.Buffer{}
		err := repl.EvalInput(tt.input, env, out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := strings.TrimSpace(out.String())
		if got != tt.expected {
			t.Errorf("input %q: expected %q, got %q", tt.input, tt.expected, got)
		}
	}
}
