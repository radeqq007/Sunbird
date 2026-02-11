package array_test

import (
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"testing"
)

func TestPush(t *testing.T) {
	tests := []struct {
		input string
		want  object.Array
	}{
		{
			input: "import 'array'; let a = [1, 2]; array.push(a, 3); a",
			want: object.Array{
				Elements: []object.Value{
					object.NewInt(1),
					object.NewInt(2),
					object.NewInt(3),
				},
			},
		},
		{
			input: "import 'array'; let a = [1, 2]; array.push(a, 'abc'); a",
			want: object.Array{
				Elements: []object.Value{
					object.NewInt(1),
					object.NewInt(2),
					object.NewString("abc"),
				},
			},
		},
		{
			input: "import 'array'; let a = []; array.push(a, true); a",
			want:  object.Array{Elements: []object.Value{object.NewBool(true)}},
		},
	}

	for _, tt := range tests {
		testArrayObject(t, testEval(tt.input), tt.want)
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		input     string
		wantValue string
		wantArray []string
	}{
		{
			input:     "import 'array'; let a = [1, 2, 3]; array.pop(a)",
			wantValue: "3",
			wantArray: []string{"1", "2"},
		},
		{
			input:     "import 'array'; let a = ['hi']; array.pop(a)",
			wantValue: "\"hi\"",
			wantArray: []string{},
		},
	}

	for _, tt := range tests {
		// Test the return value of pop
		val := testEval(tt.input)
		if val.Inspect() != tt.wantValue {
			t.Errorf("pop returned wrong value. want=%s, got=%s", tt.wantValue, val.Inspect())
		}
	}
}

func TestShift(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.shift(a)"
	val := testEval(input)

	if val.Inspect() != "1" {
		t.Errorf("shift returned wrong value. want=1, got=%s", val.Inspect())
	}
}

func TestUnshift(t *testing.T) {
	input := "import 'array'; let a = [2, 3]; array.unshift(a, 1); a"
	want := object.Array{Elements: []object.Value{
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
	}}

	testArrayObject(t, testEval(input), want)
}

func TestReverse(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.reverse(a); a"
	want := object.Array{Elements: []object.Value{
		object.NewInt(3),
		object.NewInt(2),
		object.NewInt(1),
	}}

	testArrayObject(t, testEval(input), want)
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"import 'array'; let a = ['a', 'b', 'c']; array.index_of(a, 'b')", 1},
		{"import 'array'; let a = ['a', 'b', 'c']; array.index_of(a, 'z')", -1},
	}

	for _, tt := range tests {
		val := testEval(tt.input)
		if !val.IsInt() {
			t.Fatalf("expected Integer, got=%T", val)
		}

		value := val.AsInt()
		if value != tt.want {
			t.Errorf("indexOf wrong. want=%d, got=%d", tt.want, value)
		}
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"import 'array'; let a = [0, 1, 2, 3, 4]; array.slice(a, 1, 3)", []string{"1", "2"}},
		{"import 'array'; let a = [0, 1, 2]; array.slice(a, 1)", []string{"1", "2"}},
	}

	for _, tt := range tests {
		val := testEval(tt.input)
		if !val.IsArray() {
			t.Fatalf("expected Array, got=%T", val)
		}

		array := val.AsArray()
		
		if len(array.Elements) != len(tt.want) {
			t.Fatalf("slice length wrong. want=%d, got=%d", len(tt.want), len(array.Elements))
		}
		
		for i, wantStr := range tt.want {
			if array.Elements[i].Inspect() != wantStr {
				t.Errorf(
					"element %d wrong. want=%s, got=%s",
					i,
					wantStr,
					array.Elements[i].Inspect(),
				)
			}
		}
	}
}

func TestJoin(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.join(a, '-')"
	val := testEval(input)

	if !val.IsString() {
		t.Fatalf("expected String, got=%T", val)
	}

	str := val.AsString()
	
	if str.Value != "1-2-3" {
		t.Errorf("join result wrong. want='1-2-3', got='%s'", str.Value)
	}
}

func TestConcat(t *testing.T) {
	input := "import 'array'; let a = [1]; let b = [2]; array.concat(a, b)"
	want := object.Array{Elements: []object.Value{
		object.NewInt(1),
		object.NewInt(2),
	}}

	testArrayObject(t, testEval(input), want)
}

func TestContains(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"import 'array'; let a = [1, 2, 3]; array.contains(a, 2)", true},
		{"import 'array'; let a = [1, 2, 3]; array.contains(a, 5)", false},
	}

	for _, tt := range tests {
		val := testEval(tt.input)
		if !val.IsBool() {
			t.Fatalf("expected Boolean, got=%T", val)
		}

		boolean := val.AsBool()
		
		if boolean != tt.want {
			t.Errorf("contains wrong. want=%t, got=%t", tt.want, boolean)
		}
	}
}

func TestClear(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.clear(a); a"
	val := testEval(input)

	if !val.IsArray() || len(val.AsArray().Elements) != 0 {
		t.Errorf("clear failed. array not empty, got=%+v", val)
	}
}

func testArrayObject(t *testing.T, obj object.Value, expected object.Array) {
	if !obj.IsArray() {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
	}

	array := obj.AsArray()

	if len(array.Elements) != len(expected.Elements) {
		t.Errorf(
			"array has wrong number of elements. want=%d, got=%d",
			len(expected.Elements),
			len(array.Elements),
		)
	}

	for i, a := range array.Elements {
		if a.Inspect() != expected.Elements[i].Inspect() {
			t.Errorf(
				"element %d is not equal. want=%+v, got=%+v",
				i,
				expected.Elements[i].Inspect(),
				a.Inspect(),
			)
		}
	}
}

func testEval(input string) object.Value {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}
