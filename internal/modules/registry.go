package modules

import (
	"sunbird/internal/modules/io"
	"sunbird/internal/modules/math"
	"sunbird/internal/object"
)

var BuiltinModules = map[string]*object.Hash{
	"math": math.Module,
	"io":   io.Module,
}

func Get(name string) (*object.Hash, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
