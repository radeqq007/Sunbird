package modbuilder

import "sunbird/internal/object"

type ModuleBuilder struct {
	pairs map[object.HashKey]object.HashPair
}

func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		pairs: make(map[object.HashKey]object.HashPair),
	}
}

func (mb *ModuleBuilder) AddFunction(name string, fn object.BuiltinFunction) *ModuleBuilder {
	key := object.NewString(name)
	hashKey := key.HashKey()
	mb.pairs[hashKey] = object.NewHashPair(key, object.NewBuiltin(fn))
	return mb
}

func (mb *ModuleBuilder) AddValue(name string, value object.Value) *ModuleBuilder {
	key := object.NewString(name)
	hashKey := key.HashKey()
	mb.pairs[hashKey] = object.NewHashPair(key, value)
	return mb
}

func (mb *ModuleBuilder) AddInteger(name string, value int64) *ModuleBuilder {
	return mb.AddValue(name, object.NewInt(value))
}

func (mb *ModuleBuilder) AddFloat(name string, value float64) *ModuleBuilder {
	return mb.AddValue(name, object.NewFloat(value))
}

func (mb *ModuleBuilder) AddString(name string, value string) *ModuleBuilder {
	return mb.AddValue(name, object.NewString(value))
}

func (mb *ModuleBuilder) AddBoolean(name string, value bool) *ModuleBuilder {
	return mb.AddValue(name, object.NewBool(value))
}

func (mb *ModuleBuilder) Build() object.Value {
	return object.NewHash(mb.pairs)
}

func CreateModule(pairs map[object.HashKey]object.HashPair) object.Value {
	return object.NewHash(pairs)
}

// NewHashBuilder is just a NewModuleBuilder() wrapper, but it makes creating hash objects more descriptive
func NewHashBuilder() *ModuleBuilder {
	return NewModuleBuilder()
}
