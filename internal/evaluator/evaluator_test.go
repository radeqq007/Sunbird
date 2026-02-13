package evaluator_test

import (
	"math"
	"sunbird/internal/evaluator"
	"sunbird/internal/lexer"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.2", 5.2},
		{"10.52", 10.52},
		{"-10.24", -10.24},
		{"-5.0", -5.0},
		{"5.5 + 5.5 + 5.5 + 5.5 - 12", 10.0},
		{"2.2 * 2.2 * 2", 9.68},
		{"5.0 * 2 + 10.2", 20.2},
		{"5 + 2.5 * 10", 30.0},
		{"3.2 * 3.5 - 2", 9.2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Value {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Value, expected int64) {
	if !obj.IsInt() {
		t.Errorf("object is not Integer. got=%T", obj.Kind().String())
	}

	val := obj.AsInt()
	if obj.AsInt() != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", val, expected)
	}
}

// I have to do it that way cause result.Value == expected doesn't work
const floatTolerance = 1e-9

func testFloatObject(t *testing.T, obj object.Value, expected float64) {
	if !obj.IsFloat() {
		t.Errorf("object is not Float. got=%T", obj.Kind().String())
	}

	val := obj.AsFloat()

	if math.Abs(val-expected) > floatTolerance {
		t.Errorf("object has wrong value. got=%f, want=%f", val, expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Value, expected bool) {
	if !obj.IsBool() {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
	}

	val := obj.AsBool()

	if val != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", val, expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Value) {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 2 * 5; 9", 10},
		{`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }
  return 1;
}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"TypeMismatchError: Integer + Boolean",
		},
		{
			"5 + true; 5;",
			"TypeMismatchError: Integer + Boolean",
		},
		{
			"-true",
			"UnknownOperatorError: -Boolean",
		},
		{
			"true + false;",
			"UnknownOperatorError: Boolean + Boolean",
		},
		{
			"5; true + false; 5",
			"UnknownOperatorError: Boolean + Boolean",
		},
		{
			"if 10 > 1 { true + false; }",
			"UnknownOperatorError: Boolean + Boolean",
		},
		{
			`
      132
      if 10 > 1 {
        if 10 > 1 {
          return true + false;
        }
      return 1;
    }
    `,
			"UnknownOperatorError: Boolean + Boolean",
		},
		{
			"foobar",
			"UndefinedVariableError: foobar",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !evaluated.IsError() {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		err := evaluated.AsError()
		if err.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, err.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)

	if !evaluated.IsFunction() {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	fn := evaluated.AsFunction()

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { x * 2; }; double(5);", 10},
		{"let add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = func(x) {
  func(y) { x + y };
};
let addTwo = newAdder(2);
addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)

	if !evaluated.IsString() {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	str := evaluated.AsString().Value

	if str != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)

	if !evaluated.IsString() {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	str := evaluated.AsString().Value

	if str != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("sunbird")`, 7},
		{`len("hello world")`, 11},
		{`len(1)`, "TypeError: expected one of String, Array, got Integer"},
		{`len("one", "two")`, "ArgumentError: expected 1 arguments, got 2"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`append([], 1)`, []int64{1}},
		{`append(1, 1)`, "TypeError: expected Array, got Integer"},
		{`append([1, 2, 3], 4)`, []int64{1, 2, 3, 4}},
		{`append([1, 2, 3], 4, 5, 6)`, []int64{1, 2, 3, 4, 5, 6}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case []int64:
			if !evaluated.IsArray() {
				t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			arr := evaluated.AsArray()
			if len(arr.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(arr.Elements))
				continue
			}
			for i, expectedElem := range expected {
				testIntegerObject(t, arr.Elements[i], expectedElem)
			}
		case string:
			if !evaluated.IsError() {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			err := evaluated.AsError()
			if err.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, err.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)

	if !evaluated.IsArray() {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	arr := evaluated.AsArray()

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(arr.Elements))
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][-1]", 3},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

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
			"TypeMismatchError: Integer + Boolean",
		},
		{
			"foobar",
			1, 1,
			"UndefinedVariableError: foobar",
		},
		{
			"-true",
			1, 1,
			"UnknownOperatorError: -Boolean",
		},
		{
			"if (10 > 1) { true + false; }",
			1, 20,
			"UnknownOperatorError: Boolean + Boolean",
		},
		{
			`
let a = 5;
a + true;
`,
			3, 3,
			"TypeMismatchError: Integer + Boolean",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		if !evaluated.IsError() {
			t.Errorf("no error object returned for input: %q. got=%T(%+v)",
				tt.input, evaluated, evaluated)
			continue
		}

		err := evaluated.AsError()
		if err.Line != tt.expectedLine {
			t.Errorf("wrong line number for input: %q. expected=%d, got=%d",
				tt.input, tt.expectedLine, err.Line)
		}

		if err.Col != tt.expectedCol {
			t.Errorf("wrong column number for input: %q. expected=%d, got=%d",
				tt.input, tt.expectedCol, err.Col)
		}

		if err.Message != tt.expectedMsg {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMsg, err.Message)
		}
	}
}
func TestFunctionTypeChecking(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let identity = func(x: Int): Int { x; }; identity(5);",
			"", // Should pass
		},
		{
			"let identity = func(x: String): String { x; }; identity(5);",
			"TypeError: expected String, got Integer",
		},
		{
			"let add = func(a: Int, b: Int): Int { a + b; }; add(1, 2);",
			"", // Should pass
		},
		{
			"let add = func(a: Int, b: Int) { a + b; }; add('1', '2');",
			"TypeError: expected Int, got String",
		},
		{
			"let identity = func(x: Float) { x; }; identity(\"hello\");",
			"TypeError: expected Float, got String",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if tt.expectedMessage == "" {
			if isError(evaluated) {
				t.Errorf("expected no error, got=%q", evaluated.AsError().Message)
			}
		} else {
			if !evaluated.IsError() {
				t.Errorf("no error object returned. got=%T(%+v)",
					evaluated, evaluated)
				continue
			}
			err := evaluated.AsError()
			if err.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expectedMessage, err.Message)
			}
		}
	}
}

func isError(obj object.Value) bool {
	return obj.IsError()
}

func BenchmarkIntegerArithmetic(b *testing.B) {
	input := `
		let x = 10
		let y = 20
		x + y * 2 - 5
	`
	benchmarkEval(b, input)
}

func BenchmarkFibonacci(b *testing.B) {
	input := `
		let fib = func(n) {
			if n < 2 {
				return n
			}
			return fib(n-1) + fib(n-2)
		}
		fib(15)
	`
	benchmarkEval(b, input)
}

func BenchmarkArrayOperations(b *testing.B) {
	input := `
		let arr = [1, 2, 3, 4, 5]
		let sum = 0
		for i in arr {
			sum = sum + i
		}
		sum
	`
	benchmarkEval(b, input)
}

func BenchmarkHashAccess(b *testing.B) {
	input := `
		let hash = {"a": 1, "b": 2, "c": 3}
		hash["a"] + hash["b"] + hash["c"]
	`
	benchmarkEval(b, input)
}

func BenchmarkStringConcatenation(b *testing.B) {
	input := `
		let result = ""
		for i in 0..10 {
			result = result + "x"
		}
		result
	`
	benchmarkEval(b, input)
}

func BenchmarkFunctionCalls(b *testing.B) {
	input := `
		let add = func(a, b) { a + b }
		let mul = func(a, b) { a * b }
		
		let result = 0
		for i in 0..20 {
			result = add(mul(i, 2), 1)
		}
		result
	`
	benchmarkEval(b, input)
}

func BenchmarkNestedLoops(b *testing.B) {
	input := `
		let sum = 0
		for i in 0..10 {
			for j in 0..10 {
				sum = sum + 1
			}
		}
		sum
	`
	benchmarkEval(b, input)
}

// Helper function to run benchmarks
func benchmarkEval(b *testing.B, input string) {
	b.ReportAllocs() // Report memory allocations

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		evaluator.Eval(program, env)
	}
}
