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

func (mb *ModuleBuilder) Build() *object.Hash {
	return &object.Hash{Pairs: mb.pairs}
}

func CreateModule(pairs map[object.HashKey]object.HashPair) *object.Hash {
	return &object.Hash{Pairs: pairs}
}

// NewHashBuilder is just a NewModuleBuilder() wrapper, but it makes creating hash objects more descriptive
func NewHashBuilder() *ModuleBuilder {
	return NewModuleBuilder()
}
