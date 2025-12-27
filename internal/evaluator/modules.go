// internal/evaluator/modules.go
package evaluator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sunbird/internal/lexer"
	"sunbird/internal/modules"
	"sunbird/internal/object"
	"sunbird/internal/parser"
)

type ModuleCache struct {
	modules map[string]*object.Hash
}

func NewModuleCache() *ModuleCache {
	return &ModuleCache{
		modules: make(map[string]*object.Hash),
	}
}

func (mc *ModuleCache) loadModule(path string, env *object.Environment) (*object.Hash, error) {
	if module, ok := mc.modules[path]; ok {
		return module, nil
	}

	// Check if it's a built-in module
	if builtinModule, ok := modules.BuiltinModules[path]; ok {
		mc.modules[path] = builtinModule
		return builtinModule, nil
	}

	// Load from file
	return mc.loadFileModule(path, env)
}

func (mc *ModuleCache) loadFileModule(path string, env *object.Environment) (*object.Hash, error) {
	fullPath := path
	if !filepath.IsAbs(path) {
		// Look relative to current directory
		if _, err := os.Stat(path); err == nil {
			fullPath = path
		} else {
			// Try with .sb extension
			withExt := path + ".sb"
			if _, err := os.Stat(withExt); err == nil {
				fullPath = withExt
			} else {
				return nil, fmt.Errorf("module not found: %s", path)
			}
		}
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors in module: %v", p.Errors())
	}

	moduleEnv := object.NewEnvironment()

	Eval(program, moduleEnv)

	pairs := make(map[object.HashKey]object.HashPair)
	for name, value := range moduleEnv.GetStore() {
		key := &object.String{Value: name}
		hashKey := key.HashKey()
		pairs[hashKey] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}

	module := &object.Hash{Pairs: pairs}
	mc.modules[path] = module

	return module, nil
}
