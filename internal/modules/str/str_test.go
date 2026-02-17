package str

import (
	"sunbird/internal/object"
	"testing"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Value
		expected string
		isError  bool
	}{
		{"valid concat", []object.Value{object.NewString("foo"), object.NewString("bar")}, "foobar", false},
		{"empty strings", []object.Value{object.NewString(""), object.NewString("")}, "", false},
		{"wrong arg count", []object.Value{object.NewString("foo")}, "", true},
		{"wrong type", []object.Value{object.NewInt(1), object.NewString("foo")}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := concat(object.NewCallContext(0, 0), tt.args...)
			if tt.isError {
				if !got.IsError() {
					t.Errorf("expected error, got %v", got.Kind())
				}
				return
			}
			if got.AsString().Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got.AsString().Value)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Value
		expected bool
		isError  bool
	}{
		{"is empty", []object.Value{object.NewString("")}, true, false},
		{"is not empty", []object.Value{object.NewString("hi")}, false, false},
		{"wrong type", []object.Value{object.NewInt(0)}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmpty(object.NewCallContext(0, 0), tt.args...)
			if tt.isError {
				if !got.IsError() {
					t.Errorf("expected error")
				}
				return
			}
			if got.AsBool() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got.AsBool())
			}
		})
	}
}

func TestStartsWithAndEndsWith(t *testing.T) {
	t.Run("starts_with", func(t *testing.T) {
		res := startsWith(object.NewCallContext(0, 0), object.NewString("hello"), object.NewString("he"))
		if !res.AsBool() {
			t.Errorf("expected true")
		}
	})

	t.Run("ends_with", func(t *testing.T) {
		res := endsWith(object.NewCallContext(0, 0), object.NewString("hello"), object.NewString("lo"))
		if !res.AsBool() {
			t.Errorf("expected true")
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Value
		expected bool
	}{
		{"contains substring", []object.Value{object.NewString("sunbird"), object.NewString("bird")}, true},
		{"does not contain", []object.Value{object.NewString("sunbird"), object.NewString("dog")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(object.NewCallContext(0, 0), tt.args...)
			if got.AsBool() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got.AsBool())
			}
		})
	}
}

func TestCaseAndTrim(t *testing.T) {
	t.Run("to_upper", func(t *testing.T) {
		res := toUpper(object.NewCallContext(0, 0), object.NewString("abc"))
		if res.AsString().Value != "ABC" {
			t.Errorf("got %s", res.AsString().Value)
		}
	})

	t.Run("to_lower", func(t *testing.T) {
		res := toLower(object.NewCallContext(0, 0), object.NewString("XYZ"))
		if res.AsString().Value != "xyz" {
			t.Errorf("got %s", res.AsString().Value)
		}
	})

	t.Run("trim", func(t *testing.T) {
		res := trim(object.NewCallContext(0, 0), object.NewString("  space  "))
		if res.AsString().Value != "space" {
			t.Errorf("got %q", res.AsString().Value)
		}
	})
}

func TestSplit(t *testing.T) {
	input := object.NewString("a,b,c")
	sep := object.NewString(",")
	res := split(object.NewCallContext(0, 0), input, sep)

	if !res.IsArray() {
		t.Fatalf("expected array, got %v", res.Kind())
	}

	elements := res.AsArray().Elements
	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(elements))
	}

	if elements[0].AsString().Value != "a" || elements[2].AsString().Value != "c" {
		t.Errorf("unexpected elements: %v", elements)
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Value
		expected string
		isError  bool
	}{
		{"repeat 3 times", []object.Value{object.NewString("a"), object.NewInt(3)}, "aaa", false},
		{"repeat 0 times", []object.Value{object.NewString("a"), object.NewInt(0)}, "", false},
		{"invalid count type", []object.Value{object.NewString("a"), object.NewString("3")}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repeat(object.NewCallContext(0, 0), tt.args...)
			if tt.isError {
				if !got.IsError() {
					t.Errorf("expected error for %s", tt.name)
				}
				return
			}
			if got.AsString().Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got.AsString().Value)
			}
		})
	}
}

func TestModuleExport(t *testing.T) {
	module := New()
	if !module.IsModule() {
		t.Fatalf("New() did not return a module, got %v", module.Kind())
	}

	exports := module.AsModule().Exports
	expectedFuncs := []string{
		"concat", "is_empty", "starts_with", "ends_with",
		"to_upper", "to_lower", "trim", "split", "repeat", "contains",
	}

	for _, name := range expectedFuncs {
		if _, ok := exports[name]; !ok {
			t.Errorf("module missing exported function: %s", name)
		}
	}
}
