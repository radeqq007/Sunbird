package modules

import (
	"sunbird/internal/modules/array"
	"sunbird/internal/modules/errors"
	"sunbird/internal/modules/http"
	"sunbird/internal/modules/io"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/math"
	"sunbird/internal/modules/random"
	"sunbird/internal/modules/str"
	"sunbird/internal/object"
)

var BuiltinModules = map[string]*object.Hash{
	"math":   math.Module,
	"io":     io.Module,
	"array":  array.Module,
	"string": str.Module,
	"random": random.Module,
	"errors": errors.Module,
	"json":   json.Module,
	"http":   http.Module,
}

func Get(name string) (*object.Hash, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
