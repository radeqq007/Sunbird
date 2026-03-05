package evaluator

import (
	"github.com/radeqq007/sunbird/internal/object"
)

func isError(obj object.Value) bool {
	isErr := obj.IsError()
	isPropagating := false
	if isErr {
		isPropagating = obj.AsError().Propagating
	}

	return isErr && isPropagating
}
