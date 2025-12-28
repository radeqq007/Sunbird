package evaluator

import (
	"sunbird/internal/object"
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
}
