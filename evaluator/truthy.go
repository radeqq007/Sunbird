package evaluator

import "sunbird/object"

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		// TODO: uhhh nesting switch statements? Can't be good
		switch obj := obj.(type) {
		case *object.String:
				return obj.Value != ""
		case *object.Integer:
				return obj.Value != 0
		case *object.Float:
				return obj.Value != 0.0
		default:
				return true
		}
	}
}