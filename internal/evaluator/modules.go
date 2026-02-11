package evaluator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sunbird/internal/lexer"
	"sunbird/internal/modules"
	"sunbird/internal/object"
	"sunbird/internal/parser"
	"sunbird/internal/pkg"
	"sync"
)

type ModuleCache struct {
	modules map[string]object.Value
	mu      sync.RWMutex
}

func NewModuleCache() *ModuleCache {
	return &ModuleCache{
		modules: make(map[string]object.Value),
	}
}

func (mc *ModuleCache) loadModule(path string) (object.Value, error) {
	mc.mu.RLock()
	module, ok := mc.modules[path]
	mc.mu.RUnlock()
	if ok {
		return module, nil
	}

	// Check if it's a built-in module
	if builtinModule, ok := modules.BuiltinModules[path]; ok {
		mc.mu.Lock()
		mc.modules[path] = builtinModule
		mc.mu.Unlock()
		return builtinModule, nil
	}

	// Check in .sb_modules directory
	modulesDir := ".sb_modules"
	if _, err := os.Stat(modulesDir); err == nil {
		if module, err := mc.tryLoadFromModulesDir(path); err == nil {
			mc.mu.Lock()
			mc.modules[path] = module
			mc.mu.Unlock()
			return module, nil
		}
	}

	// Load from file
	return mc.loadFileModule(path)
}

func (mc *ModuleCache) loadFileModule(path string) (object.Value, error) {
	mainFileDir := ""
	if len(os.Args) > 1 {
		mainFileDir = filepath.Dir(os.Args[1]) // TODO: don't use os.Args
	}

	fullPath := filepath.Join(mainFileDir, path)
	if !filepath.IsAbs(path) {
		// Look relative to current directory
		if _, err := os.Stat(fullPath); err != nil {
			// Try with .sb extension
			withExt := fullPath + ".sb"
			if _, err = os.Stat(withExt); err == nil {
				fullPath = withExt
			} else {
				return object.NewNull(), fmt.Errorf("module not found: %s", path)
			}
		}
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return object.NewNull(), err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return object.NewNull(), err
	}

	file.Close()

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return object.NewNull(), fmt.Errorf("parse errors in module: %v", p.Errors())
	}

	moduleEnv := object.NewEnvironment()

	Eval(program, moduleEnv)

	pairs := make(map[object.HashKey]object.HashPair)
	for name, value := range moduleEnv.GetStore() {
		key := object.NewString(name)
		hashKey := key.HashKey()

		pairs[hashKey] = object.NewHashPair(key, value)

	}

	module := object.NewHash(pairs)
	mc.mu.Lock()
	mc.modules[path] = module
	mc.mu.Unlock()

	return module, nil
}

func (mc *ModuleCache) tryLoadFromModulesDir(path string) (object.Value, error) {
	modulesDir := ".sb_modules"
	packagePath := filepath.Join(modulesDir, path)

	if _, err := os.Stat(packagePath); err != nil {
		return object.NewNull(), err
	}

	moduleConf, err := pkg.LoadConfig(filepath.Join(packagePath, "sunbird.toml"))
	if err != nil {
		return object.NewNull(), errors.New("failed to load module config: " + err.Error())
	}

	if moduleConf.Package.Main == "" {
		return object.NewNull(), errors.New("module config missing main file")
	}

	return mc.loadFileModule(filepath.Join(packagePath, moduleConf.Package.Main))
}
