package object

type Environment struct {
	store     map[string]Value
	constants map[string]bool
	exports   map[string]bool
	outer     *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Value)
	c := make(map[string]bool)
	e := make(map[string]bool)
	return &Environment{store: s, constants: c, exports: e}
}

func (e *Environment) Get(name string) (Value, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) SetConst(name string, val Value) Value {
	e.store[name] = val
	e.constants[name] = true
	return val
}

func (e *Environment) IsConst(name string) bool {
	if isConst, ok := e.constants[name]; ok {
		return isConst
	}
	if e.outer != nil {
		return e.outer.IsConst(name)
	}
	return false
}

func (e *Environment) GetFromCurrentScope(name string) (Value, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Has(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}

func (e *Environment) Update(name string, val Value) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}

	if e.outer != nil {
		return e.outer.Update(name, val)
	}

	// Variable not found anywhere
	return false
}

func (e *Environment) GetStore() map[string]Value {
	return e.store
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) MarkAsExported(name string) {
	e.exports[name] = true
}

func (e *Environment) IsExported(name string) bool {
	return e.exports[name]
}

func (e *Environment) GetExports() map[string]Value {
	exports := make(map[string]Value)
	for name, value := range e.store {
		if e.exports[name] {
			exports[name] = value
		}
	}
	return exports
}
