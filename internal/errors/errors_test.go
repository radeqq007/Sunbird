package errors_test

import (
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/object"
	"testing"
)

func assertError(t *testing.T, val object.Value, wantCode string) *object.Error {
	t.Helper()
	if !val.IsError() {
		t.Fatalf("expected an error value, got kind=%s", val.Kind())
	}
	err := val.AsError()
	if !strings.HasPrefix(err.Message, wantCode) {
		t.Errorf("wrong error code prefix: expected=%q, got=%q", wantCode, err.Message)
	}
	return err
}

func assertNotError(t *testing.T, val object.Value) {
	t.Helper()
	if val.IsError() {
		t.Fatalf("expected no error, got: %s", val.AsError().Message)
	}
}

func TestNew_MessageContainsCodeAndText(t *testing.T) {
	val := errors.New(errors.RuntimeError, 3, 7, "something went wrong")
	err := assertError(t, val, "RuntimeError")

	if !strings.Contains(err.Message, "something went wrong") {
		t.Errorf("message missing detail text, got=%q", err.Message)
	}
}

func TestNew_LineAndColAreStored(t *testing.T) {
	val := errors.New(errors.TypeError, 10, 5, "bad type")
	err := val.AsError()

	if err.Line != 10 {
		t.Errorf("wrong line: expected=10, got=%d", err.Line)
	}
	if err.Col != 5 {
		t.Errorf("wrong col: expected=5, got=%d", err.Col)
	}
}

func TestNew_PropagatingIsTrue(t *testing.T) {
	val := errors.New(errors.RuntimeError, 0, 0, "oops")
	if !val.AsError().Propagating {
		t.Error("expected Propagating=true for newly created error")
	}
}

func TestNew_FormattedMessage(t *testing.T) {
	val := errors.New(errors.ArgumentError, 0, 0, "expected %d args, got %d", 2, 5)
	err := val.AsError()
	if !strings.Contains(err.Message, "expected 2 args, got 5") {
		t.Errorf("formatted message wrong, got=%q", err.Message)
	}
}

func TestExpectType_Match(t *testing.T) {
	val := errors.ExpectType(0, 0, object.NewInt(1), object.IntKind)
	assertNotError(t, val)
}

func TestExpectType_Mismatch(t *testing.T) {
	val := errors.ExpectType(1, 2, object.NewString("hi"), object.IntKind)
	err := assertError(t, val, "TypeError")

	if err.Line != 1 || err.Col != 2 {
		t.Errorf("wrong position: line=%d col=%d", err.Line, err.Col)
	}
	if !strings.Contains(err.Message, "Integer") {
		t.Errorf("message should mention expected type, got=%q", err.Message)
	}
	if !strings.Contains(err.Message, "String") {
		t.Errorf("message should mention actual type, got=%q", err.Message)
	}
}

func TestExpectType_AllKinds(t *testing.T) {
	cases := []struct {
		val      object.Value
		kind     object.ValueKind
		wantPass bool
	}{
		{object.NewInt(1), object.IntKind, true},
		{object.NewFloat(1.0), object.FloatKind, true},
		{object.NewString("x"), object.StringKind, true},
		{object.NewBool(true), object.BoolKind, true},
		{object.NewNull(), object.NullKind, true},
		{object.NewInt(1), object.FloatKind, false},
		{object.NewFloat(1.0), object.IntKind, false},
		{object.NewString("x"), object.BoolKind, false},
	}

	for _, c := range cases {
		result := errors.ExpectType(0, 0, c.val, c.kind)
		if c.wantPass {
			assertNotError(t, result)
		} else {
			assertError(t, result, "TypeError")
		}
	}
}

func TestExpectOneOfTypes_FirstMatch(t *testing.T) {
	val := errors.ExpectOneOfTypes(0, 0, object.NewInt(1), object.IntKind, object.FloatKind)
	assertNotError(t, val)
}

func TestExpectOneOfTypes_SecondMatch(t *testing.T) {
	val := errors.ExpectOneOfTypes(0, 0, object.NewFloat(2.0), object.IntKind, object.FloatKind)
	assertNotError(t, val)
}

func TestExpectOneOfTypes_NoMatch(t *testing.T) {
	val := errors.ExpectOneOfTypes(0, 0, object.NewString("x"), object.IntKind, object.FloatKind)
	err := assertError(t, val, "TypeError")

	if !strings.Contains(err.Message, "String") {
		t.Errorf("message should mention actual type, got=%q", err.Message)
	}
}

func TestExpectOneOfTypes_MessageListsExpected(t *testing.T) {
	val := errors.ExpectOneOfTypes(0, 0, object.NewBool(true), object.IntKind, object.FloatKind, object.StringKind)
	err := assertError(t, val, "TypeError")

	for _, expected := range []string{"Integer", "Float", "String"} {
		if !strings.Contains(err.Message, expected) {
			t.Errorf("message should mention %q, got=%q", expected, err.Message)
		}
	}
}

func TestExpectNumberOfArguments_Exact(t *testing.T) {
	args := []object.Value{object.NewInt(1), object.NewInt(2)}
	val := errors.ExpectNumberOfArguments(0, 0, 2, args)
	assertNotError(t, val)
}

func TestExpectNumberOfArguments_TooFew(t *testing.T) {
	args := []object.Value{object.NewInt(1)}
	val := errors.ExpectNumberOfArguments(0, 0, 2, args)
	err := assertError(t, val, "ArgumentError")

	if !strings.Contains(err.Message, "2") {
		t.Errorf("message should mention expected count, got=%q", err.Message)
	}
}

func TestExpectNumberOfArguments_TooMany(t *testing.T) {
	args := []object.Value{object.NewInt(1), object.NewInt(2), object.NewInt(3)}
	val := errors.ExpectNumberOfArguments(0, 0, 1, args)
	assertError(t, val, "ArgumentError")
}

func TestExpectNumberOfArguments_Zero(t *testing.T) {
	val := errors.ExpectNumberOfArguments(0, 0, 0, []object.Value{})
	assertNotError(t, val)
}

func TestExpectNumberOfArguments_PositionStored(t *testing.T) {
	args := []object.Value{}
	val := errors.ExpectNumberOfArguments(4, 8, 1, args)
	err := val.AsError()
	if err.Line != 4 || err.Col != 8 {
		t.Errorf("wrong position: line=%d col=%d", err.Line, err.Col)
	}
}

func TestExpectMinNumberOfArguments_Exact(t *testing.T) {
	args := []object.Value{object.NewInt(1)}
	val := errors.ExpectMinNumberOfArguments(0, 0, 1, args)
	assertNotError(t, val)
}

func TestExpectMinNumberOfArguments_More(t *testing.T) {
	args := []object.Value{object.NewInt(1), object.NewInt(2), object.NewInt(3)}
	val := errors.ExpectMinNumberOfArguments(0, 0, 2, args)
	assertNotError(t, val)
}

func TestExpectMinNumberOfArguments_TooFew(t *testing.T) {
	args := []object.Value{object.NewInt(1)}
	val := errors.ExpectMinNumberOfArguments(0, 0, 3, args)
	err := assertError(t, val, "ArgumentError")

	if !strings.Contains(err.Message, "3") {
		t.Errorf("message should mention minimum count, got=%q", err.Message)
	}
}

func TestExpectMinNumberOfArguments_Zero(t *testing.T) {
	val := errors.ExpectMinNumberOfArguments(0, 0, 0, []object.Value{})
	assertNotError(t, val)
}

func TestNewTypeError(t *testing.T) {
	val := errors.NewTypeError(1, 2, "expected %s, got %s", "Int", "String")
	err := assertError(t, val, "TypeError")
	if err.Line != 1 || err.Col != 2 {
		t.Errorf("wrong position: line=%d col=%d", err.Line, err.Col)
	}
	if !strings.Contains(err.Message, "expected Int, got String") {
		t.Errorf("wrong message: %q", err.Message)
	}
}

func TestNewTypeMismatchError(t *testing.T) {
	val := errors.NewTypeMismatchError(3, 4, object.IntKind, "+", object.StringKind)
	err := assertError(t, val, "TypeMismatchError")
	if !strings.Contains(err.Message, "Integer") || !strings.Contains(err.Message, "String") || !strings.Contains(err.Message, "+") {
		t.Errorf("message missing components, got=%q", err.Message)
	}
}

func TestNewRuntimeError(t *testing.T) {
	val := errors.NewRuntimeError(5, 6, "division by zero %s", "again")
	err := assertError(t, val, "RuntimeError")
	if !strings.Contains(err.Message, "division by zero again") {
		t.Errorf("wrong message: %q", err.Message)
	}
	if err.Line != 5 || err.Col != 6 {
		t.Errorf("wrong position: line=%d col=%d", err.Line, err.Col)
	}
}

func TestNewDivisionByZeroError(t *testing.T) {
	val := errors.NewDivisionByZeroError(7, 8)
	err := assertError(t, val, "DivisionByZeroError")
	if err.Line != 7 || err.Col != 8 {
		t.Errorf("wrong position: line=%d col=%d", err.Line, err.Col)
	}
}

func TestNewUndefinedVariableError(t *testing.T) {
	val := errors.NewUndefinedVariableError(1, 1, "foobar")
	err := assertError(t, val, "UndefinedVariableError")
	if !strings.Contains(err.Message, "foobar") {
		t.Errorf("message should mention variable name, got=%q", err.Message)
	}
}

func TestNewConstantReassignmentError(t *testing.T) {
	val := errors.NewConstantReassignmentError(2, 3, "myConst")
	err := assertError(t, val, "ConstantReassignmentError")
	if !strings.Contains(err.Message, "myConst") {
		t.Errorf("message should mention constant name, got=%q", err.Message)
	}
}

func TestNewVariableReassignmentError(t *testing.T) {
	val := errors.NewVariableReassignmentError(0, 0, "x")
	err := assertError(t, val, "VariableReassignmentError")
	if !strings.Contains(err.Message, "x") {
		t.Errorf("message should mention variable name, got=%q", err.Message)
	}
}

func TestNewIndexNotSupportedError(t *testing.T) {
	val := errors.NewIndexNotSupportedError(0, 0, object.NewBool(true))
	err := assertError(t, val, "IndexNotSupportedError")
	if !strings.Contains(err.Message, "Boolean") {
		t.Errorf("message should mention type, got=%q", err.Message)
	}
}

func TestNewIndexOutOfBoundsError(t *testing.T) {
	val := errors.NewIndexOutOfBoundsError(0, 0, object.NewArray([]object.Value{}))
	assertError(t, val, "IndexOutOfBoundsError")
}

func TestNewNonObjectPropertyAccessError(t *testing.T) {
	val := errors.NewNonObjectPropertyAccessError(0, 0, object.NewInt(42))
	err := assertError(t, val, "PropertyAccessOnNonObjectError")
	if !strings.Contains(err.Message, "Integer") {
		t.Errorf("message should mention type, got=%q", err.Message)
	}
}

func TestNewUnusableAsHashKeyError(t *testing.T) {
	val := errors.NewUnusableAsHashKeyError(0, 0, object.NewBool(false))
	assertError(t, val, "KeyError")
}

func TestNewInvalidAssignmentTargetError(t *testing.T) {
	val := errors.NewInvalidAssignmentTargetError(0, 0, "5 + 3")
	err := assertError(t, val, "InvalidAssignmentError")
	if !strings.Contains(err.Message, "5 + 3") {
		t.Errorf("message should mention target expression, got=%q", err.Message)
	}
}

func TestNewNotCallableError(t *testing.T) {
	val := errors.NewNotCallableError(0, 0, object.NewInt(99))
	err := assertError(t, val, "NotCallableError")
	if !strings.Contains(err.Message, "Integer") {
		t.Errorf("message should mention type, got=%q", err.Message)
	}
}

func TestNewImportError(t *testing.T) {
	val := errors.NewImportError(0, 0, "module not found: xyz")
	err := assertError(t, val, "ImportError")
	if !strings.Contains(err.Message, "module not found: xyz") {
		t.Errorf("message should contain detail, got=%q", err.Message)
	}
}

func TestNewUnknownOperatorError(t *testing.T) {
	left := object.NewInt(1)
	right := object.NewString("x")
	val := errors.NewUnknownOperatorError(0, 0, left, "**", right)
	err := assertError(t, val, "UnknownOperatorError")

	for _, part := range []string{"Integer", "**", "String"} {
		if !strings.Contains(err.Message, part) {
			t.Errorf("message missing %q, got=%q", part, err.Message)
		}
	}
}

func TestNewUnknownPrefixOperatorError(t *testing.T) {
	val := errors.NewUnknownPrefixOperatorError(0, 0, "-", object.NewBool(true))
	err := assertError(t, val, "UnknownOperatorError")
	if !strings.Contains(err.Message, "-") || !strings.Contains(err.Message, "Boolean") {
		t.Errorf("message missing components, got=%q", err.Message)
	}
}

func TestNewArgumentError(t *testing.T) {
	val := errors.NewArgumentError(1, 1, "expected at least %d, got %d", 2, 0)
	err := assertError(t, val, "ArgumentError")
	if !strings.Contains(err.Message, "expected at least 2, got 0") {
		t.Errorf("wrong message: %q", err.Message)
	}
}

func TestNewFeatureNotImplementedError(t *testing.T) {
	val := errors.NewFeatureNotImplementedError(0, 0, "pattern matching")
	err := assertError(t, val, "FeatureNotImplementedError")
	if !strings.Contains(err.Message, "pattern matching") {
		t.Errorf("message should mention feature name, got=%q", err.Message)
	}
}
