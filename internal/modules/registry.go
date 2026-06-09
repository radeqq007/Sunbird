package modules

import (
	"github.com/radeqq007/sunbird/internal/modules/array"
	"github.com/radeqq007/sunbird/internal/modules/errors"
	"github.com/radeqq007/sunbird/internal/modules/fs"
	"github.com/radeqq007/sunbird/internal/modules/http"
	"github.com/radeqq007/sunbird/internal/modules/io"
	"github.com/radeqq007/sunbird/internal/modules/json"
	"github.com/radeqq007/sunbird/internal/modules/math"
	"github.com/radeqq007/sunbird/internal/modules/random"
	"github.com/radeqq007/sunbird/internal/modules/str"
	"github.com/radeqq007/sunbird/internal/modules/time"
	"github.com/radeqq007/sunbird/internal/object"
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
