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
				Elements: []object.Object{
					&object.Integer{Value: 1},
					&object.Integer{Value: 2},
					&object.Integer{Value: 3},
				},
			},
		},
		{
			input: "import 'array'; let a = [1, 2]; array.push(a, 'abc'); a",
			want: object.Array{
				Elements: []object.Object{
					&object.Integer{Value: 1},
					&object.Integer{Value: 2},
					&object.String{Value: "abc"},
				},
			},
		},
		{
			input: "import 'array'; let a = []; array.push(a, true); a",
			want:  object.Array{Elements: []object.Object{&object.Boolean{Value: true}}},
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
	want := object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}

	testArrayObject(t, testEval(input), want)
}

func TestReverse(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.reverse(a); a"
	want := object.Array{Elements: []object.Object{
		&object.Integer{Value: 3},
		&object.Integer{Value: 2},
		&object.Integer{Value: 1},
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
		result, ok := val.(*object.Integer)
		if !ok {
			t.Fatalf("expected Integer, got=%T", val)
		}
		if result.Value != tt.want {
			t.Errorf("indexOf wrong. want=%d, got=%d", tt.want, result.Value)
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
		array, ok := val.(*object.Array)
		if !ok {
			t.Fatalf("expected Array, got=%T", val)
		}
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

	str, ok := val.(*object.String)
	if !ok {
		t.Fatalf("expected String, got=%T", val)
	}
	if str.Value != "1-2-3" {
		t.Errorf("join result wrong. want='1-2-3', got='%s'", str.Value)
	}
}

func TestConcat(t *testing.T) {
	input := "import 'array'; let a = [1]; let b = [2]; array.concat(a, b)"
	want := object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
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
		boolean, ok := val.(*object.Boolean)
		if !ok {
			t.Fatalf("expected Boolean, got=%T", val)
		}
		if boolean.Value != tt.want {
			t.Errorf("contains wrong. want=%t, got=%t", tt.want, boolean.Value)
		}
	}
}

func TestClear(t *testing.T) {
	input := "import 'array'; let a = [1, 2, 3]; array.clear(a); a"
	val := testEval(input)

	array, ok := val.(*object.Array)
	if !ok || len(array.Elements) != 0 {
		t.Errorf("clear failed. array not empty, got=%+v", val)
	}
}

func testArrayObject(t *testing.T, obj object.Object, expected object.Array) {
	array, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
	}

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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}
