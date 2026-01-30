package evaluator

import (
	"sunbird/internal/object"
)

func isError(obj object.Object) bool {
	if obj != nil && obj.Type() == object.ErrorObj {
		errObj := obj.(*object.Error)
		return errObj.Propagating
	}

	return false
}
