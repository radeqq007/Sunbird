package evaluator

import (
	"math"
	"sunbird/lexer"
	"sunbird/object"
	"sunbird/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  } {
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
  } {
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

func testEval(input string) object.Object {
  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()

  return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
  result, ok := obj.(*object.Integer)
  
  if !ok {
    t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
    return false
  }

  if result.Value != expected {
    t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
    return false
  }

  return true
}

// I have to do it that way cause result.Value == expected doesn't work
const floatTolerance = 1e-9

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
result, ok := obj.(*object.Float)

if !ok {
  t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
  return false
}

if math.Abs(result.Value-expected) > floatTolerance {
  t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
  return false
}

return true
}

func TestEvalBooleanExpression(t *testing.T) {
  tests := []struct {
    input string
    expected bool
  } {
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

 func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
  result, ok := obj.(*object.Boolean)

  if !ok {
    t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
    return false
  }

  if result.Value != expected {
    t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
    return false
  }

  return true
}

func TestBangOperator(t *testing.T) {
  tests := []struct {
    input string
    expected bool
  } {
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
    input string
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

func testNullObject(t *testing.T, obj object.Object) bool {
  if obj != NULL {
    t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
    return false
  }

  return true
}

func TestReturnStatements(t *testing.T) {
  tests := []struct{
    input    string
    expected int64
  } {
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
    input string
    expectedMessage string
  }{
    {
      "5 + true;",
      "type mismatch: INTEGER + BOOLEAN",
    },
    {
      "5 + true; 5;",
      "type mismatch: INTEGER + BOOLEAN",
    },
    {
      "-true",
      "unknown operator: -BOOLEAN",
    },
    {
      "true + false;",
      "unknown operator: BOOLEAN + BOOLEAN",
    },
    {
      "5; true + false; 5",
      "unknown operator: BOOLEAN + BOOLEAN",
    },
    {
      "if 10 > 1 { true + false; }",
      "unknown operator: BOOLEAN + BOOLEAN",
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
    "unknown operator: BOOLEAN + BOOLEAN",
    },
  }
  for _, tt := range tests {
  evaluated := testEval(tt.input)
  errObj, ok := evaluated.(*object.Error)
  if !ok {
  t.Errorf("no error object returned. got=%T(%+v)",
  evaluated, evaluated)
  continue
  }
  if errObj.Message != tt.expectedMessage {
  t.Errorf("wrong error message. expected=%q, got=%q",
  tt.expectedMessage, errObj.Message)
  }
  }
  }