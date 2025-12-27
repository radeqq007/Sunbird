// internal/evaluator/modules.go
package evaluator

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sunbird/internal/lexer"
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
	if builtinModule, ok := builtinModules[path]; ok {
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

type ModuleBuilder struct {
	pairs map[object.HashKey]object.HashPair
}

func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		pairs: make(map[object.HashKey]object.HashPair),
	}
}

func (mb *ModuleBuilder) AddFunction(name string, fn object.BuiltinFunction) *ModuleBuilder {
	key := &object.String{Value: name}
	hashKey := key.HashKey()
	mb.pairs[hashKey] = object.HashPair{
		Key:   key,
		Value: &object.Builtin{Fn: fn},
	}
	return mb
}

func (mb *ModuleBuilder) AddValue(name string, value object.Object) *ModuleBuilder {
	key := &object.String{Value: name}
	hashKey := key.HashKey()
	mb.pairs[hashKey] = object.HashPair{
		Key:   key,
		Value: value,
	}
	return mb
}

func (mb *ModuleBuilder) AddInteger(name string, value int64) *ModuleBuilder {
	return mb.AddValue(name, &object.Integer{Value: value})
}

func (mb *ModuleBuilder) AddFloat(name string, value float64) *ModuleBuilder {
	return mb.AddValue(name, &object.Float{Value: value})
}

func (mb *ModuleBuilder) AddString(name string, value string) *ModuleBuilder {
	return mb.AddValue(name, &object.String{Value: value})
}

func (mb *ModuleBuilder) AddBoolean(name string, value bool) *ModuleBuilder {
	return mb.AddValue(name, &object.Boolean{Value: value})
}

// Build the module hash
func (mb *ModuleBuilder) Build() *object.Hash {
	return &object.Hash{Pairs: mb.pairs}
}

func createModule(pairs map[object.HashKey]object.HashPair) *object.Hash {
	return &object.Hash{Pairs: pairs}
}

var builtinModules = map[string]*object.Hash{
	"math": NewModuleBuilder().
		AddFunction("abs", mathAbs).
		Build(),
}

func mathAbs(args ...object.Object) object.Object {
	if len(args) != 1 {
		// return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.IntegerObj {
		// return newError("argument to `abs` must be INTEGER, got %s", args[0].Type())
	}

	return &object.Integer{Value: int64(math.Abs(float64(args[0].(*object.Integer).Value)))}
}
