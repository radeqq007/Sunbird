package modules

import (
	"sunbird/internal/modules/array"
	"sunbird/internal/modules/errors"
	"sunbird/internal/modules/fs"
	"sunbird/internal/modules/http"
	"sunbird/internal/modules/io"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/math"
	"sunbird/internal/modules/random"
	"sunbird/internal/modules/str"
	"sunbird/internal/modules/time"
	"sunbird/internal/object"
)

func init() {
	registerModule("math", math.New())
	registerModule("io", io.New())
	registerModule("array", array.New())
	registerModule("string", str.New())
	registerModule("random", random.New())
	registerModule("errors", errors.New())
	registerModule("json", json.New())
	registerModule("http", http.New())
	registerModule("fs", fs.New())
	registerModule("time", time.New())
}

var BuiltinModules = make(map[string]object.Value)

func registerModule(name string, module object.Value) {
	if module.IsModule() {
		mod := module.AsModule()
		mod.Name = name
		BuiltinModules[name] = module
	}
}

func Get(name string) (object.Value, bool) {
	module, ok := BuiltinModules[name]
	return module, ok
}
