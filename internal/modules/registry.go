package modules

import (
	"sunbird/internal/modules/array"
	"sunbird/internal/modules/io"
	"sunbird/internal/modules/math"
	"sunbird/internal/modules/str"
	"sunbird/internal/object"
)

var BuiltinModules = map[string]*object.Hash{
	"math":   math.Module,
	"io":     io.Module,
	"array":  array.Module,
	"string": str.Module,
}

func Get(name string) (*object.Hash, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
