package math_test

import (
	m "math"
	"sunbird/internal/modules/math"
	"sunbird/internal/object"
	"testing"
)

func TestMathFunctions(t *testing.T) {
	mod := math.New()

	tests := []struct {
		name     string
		fn       string
		args     []object.Value
		expected object.Value
	}{
		// abs
		{"abs int", "abs", []object.Value{object.NewInt(-5)}, object.NewInt(5)},
		{"abs float", "abs", []object.Value{object.NewFloat(-3.2)}, object.NewFloat(3.2)},

		// max
		{"max ints", "max", []object.Value{object.NewInt(3), object.NewInt(7)}, object.NewInt(7)},
		{"max float and int", "max", []object.Value{object.NewFloat(2.5), object.NewInt(3)}, object.NewFloat(3)},

		// min
		{"min ints", "min", []object.Value{object.NewInt(3), object.NewInt(7)}, object.NewInt(3)},
		{"min float and int", "min", []object.Value{object.NewFloat(2.5), object.NewInt(3)}, object.NewFloat(2.5)},

		// pow
		{"pow ints", "pow", []object.Value{object.NewInt(2), object.NewInt(3)}, object.NewInt(8)},
		{"pow float", "pow", []object.Value{object.NewFloat(2), object.NewFloat(3.5)}, object.NewFloat(m.Pow(2, 3.5))},

		// sqrt
		{"sqrt int", "sqrt", []object.Value{object.NewInt(16)}, object.NewInt(4)},
		{"sqrt float", "sqrt", []object.Value{object.NewFloat(2)}, object.NewFloat(m.Sqrt(2))},

		// floor
		{"floor int", "floor", []object.Value{object.NewInt(5)}, object.NewInt(5)},
		{"floor float", "floor", []object.Value{object.NewFloat(3.7)}, object.NewFloat(3)},

		// ceil
		{"ceil int", "ceil", []object.Value{object.NewInt(5)}, object.NewInt(5)},
		{"ceil float", "ceil", []object.Value{object.NewFloat(3.3)}, object.NewFloat(4)},

		// round
		{"round int", "round", []object.Value{object.NewInt(5)}, object.NewInt(5)},
		{"round float", "round", []object.Value{object.NewFloat(3.6)}, object.NewFloat(4)},

		// sign
		{"sign positive int", "sign", []object.Value{object.NewInt(5)}, object.NewInt(1)},
		{"sign negative float", "sign", []object.Value{object.NewFloat(-2.3)}, object.NewFloat(-1)},
		{"sign zero", "sign", []object.Value{object.NewInt(0)}, object.NewInt(0)},

		// clamp
		{"clamp int within", "clamp", []object.Value{object.NewInt(5), object.NewInt(1), object.NewInt(10)}, object.NewInt(5)},
		{"clamp int below", "clamp", []object.Value{object.NewInt(-1), object.NewInt(0), object.NewInt(5)}, object.NewInt(0)},
		{"clamp float above", "clamp", []object.Value{object.NewFloat(7.5), object.NewFloat(0), object.NewFloat(5)}, object.NewFloat(5)},

		// trigonometry
		{"sin float", "sin", []object.Value{object.NewFloat(m.Pi / 2)}, object.NewFloat(1)},
		{"cos float", "cos", []object.Value{object.NewFloat(0)}, object.NewFloat(1)},
		{"tan float", "tan", []object.Value{object.NewFloat(0)}, object.NewFloat(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := mod.AsModule().Exports[tt.fn].AsBuiltin().Fn
			result := fn(object.NewCallContext(0, 0,), tt.args...)

			switch {
			case result.IsInt() && tt.expected.IsInt():
				if result.AsInt() != tt.expected.AsInt() {
					t.Errorf("%s: expected %d, got %d", tt.fn, tt.expected.AsInt(), result.AsInt())
				}
			case result.IsFloat() && tt.expected.IsFloat():
				if m.Abs(result.AsFloat()-tt.expected.AsFloat()) > 1e-9 {
					t.Errorf("%s: expected %f, got %f", tt.fn, tt.expected.AsFloat(), result.AsFloat())
				}
			default:
				t.Errorf("%s: unexpected type %s", tt.fn, result.Kind())
			}
		})
	}
}
