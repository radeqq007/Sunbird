package modbuilder

import "sunbird/internal/object"

type ModuleBuilder struct {
	pairs map[string]object.Value
}

func NewModuleBuilder() *ModuleBuilder {
	return &ModuleBuilder{
		pairs: make(map[string]object.Value),
	}
}

func (mb *ModuleBuilder) AddFunction(name string, fn object.BuiltinFunction) *ModuleBuilder {
	mb.pairs[name] = object.NewBuiltin(fn)
	return mb
}

func (mb *ModuleBuilder) AddValue(name string, value object.Value) *ModuleBuilder {
	mb.pairs[name] = value
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
	// No need for module name, the registry handles that
	return object.NewModule("", mb.pairs)
}

type HashBuilder struct {
	pairs map[object.HashKey]object.HashPair
}

func NewHashBuilder() *HashBuilder {
	return &HashBuilder{
		pairs: make(map[object.HashKey]object.HashPair),
	}
}

func (hb *HashBuilder) AddFunction(name string, fn object.BuiltinFunction) *HashBuilder {
	key := object.NewString(name)
	hashKey := key.HashKey()
	hb.pairs[hashKey] = object.NewHashPair(key, object.NewBuiltin(fn))
	return hb
}

func (hb *HashBuilder) AddValue(name string, value object.Value) *HashBuilder {
	key := object.NewString(name)
	hashKey := key.HashKey()
	hb.pairs[hashKey] = object.NewHashPair(key, value)
	return hb
}

func (hb *HashBuilder) AddInteger(name string, value int64) *HashBuilder {
	return hb.AddValue(name, object.NewInt(value))
}

func (hb *HashBuilder) AddFloat(name string, value float64) *HashBuilder {
	return hb.AddValue(name, object.NewFloat(value))
}

func (hb *HashBuilder) AddString(name string, value string) *HashBuilder {
	return hb.AddValue(name, object.NewString(value))
}

func (hb *HashBuilder) AddBoolean(name string, value bool) *HashBuilder {
	return hb.AddValue(name, object.NewBool(value))
}

func (hb *HashBuilder) Build() object.Value {
	return object.NewHash(hb.pairs)
}
