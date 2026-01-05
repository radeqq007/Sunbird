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
	"math":   math.New(),
	"io":     io.New(),
	"array":  array.New(),
	"string": str.New(),
	"random": random.New(),
	"errors": errors.New(),
	"json":   json.New(),
	"http":   http.New(),
}

func Get(name string) (*object.Hash, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
