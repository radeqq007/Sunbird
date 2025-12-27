package evaluator

import (
	"sunbird/internal/object"
)

func NewError(line, col int, format string, a ...interface{}) *object.Error {
	return object.NewError(line, col, format, a...)
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
}
