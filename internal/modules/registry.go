package modules

import (
	"sunbird/internal/modules/math"
	"sunbird/internal/object"
)

var BuiltinModules = map[string]*object.Hash{
	"math": math.Module,
}

func Get(name string) (*object.Hash, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
