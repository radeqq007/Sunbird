package evaluator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/radeqq007/sunbird/internal/lexer"
	"github.com/radeqq007/sunbird/internal/modules"
	"github.com/radeqq007/sunbird/internal/object"
	"github.com/radeqq007/sunbird/internal/parser"
	"github.com/radeqq007/sunbird/internal/pkg"
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
	if len(os.Args) > 2 {
		mainFileDir = filepath.Dir(os.Args[2]) // TODO: don't use os.Args
	}

	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = filepath.Clean(path)
	} else {
		joined, err := safeJoin(mainFileDir, path)
		if err != nil {
			return object.NewNull(), fmt.Errorf("module not found: %s", path)
		}
		fullPath = joined

		if _, err := os.Stat(fullPath); err != nil {
			withExt := fullPath + ".sb"
			if _, err2 := os.Stat(withExt); err2 == nil {
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

	err = file.Close()
	if err != nil {
		return object.NewNull(), err
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return object.NewNull(), fmt.Errorf("parse errors in module: %v", p.Errors())
	}

	moduleEnv := object.NewEnvironment()

	result := Eval(program, moduleEnv)
	if isError(result) {
		return result, nil
	}

	moduleName := filepath.Base(path)
	if ext := filepath.Ext(moduleName); ext != "" {
		moduleName = moduleName[:len(moduleName)-len(ext)]
	}

	exports := moduleEnv.GetExports()

	module := object.NewModule(moduleName, exports)
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

// safeJoin joins base and rel and verifies the result stays within base.
func safeJoin(base, rel string) (string, error) {
    absBase, err := filepath.Abs(base)
    if err != nil {
        return "", fmt.Errorf("cannot resolve base directory: %w", err)
    }
    full := filepath.Clean(filepath.Join(absBase, rel))
    if !strings.HasPrefix(full, absBase+string(os.PathSeparator)) && full != absBase {
        return "", fmt.Errorf("path %q escapes base directory", rel)
    }
    return full, nil
}

