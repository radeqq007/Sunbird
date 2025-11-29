package evaluator

import (
	"fmt"
	"sunbird/internal/object"
)

func newError(line, col int, format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...), Line: line, Col: col}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
}
