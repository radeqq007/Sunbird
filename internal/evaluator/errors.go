package evaluator

import (
	"sunbird/internal/object"
)

func isError(obj object.Object) bool {
	if obj != nil && obj.Type() == object.ErrorObj {
		return obj.(*object.Error).Propagating
	}

	return false
}
