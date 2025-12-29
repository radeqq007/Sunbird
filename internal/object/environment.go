package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, constants: c}
}

type Environment struct {
	store     map[string]Object
	constants map[string]bool
	outer     *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) SetConst(name string, val Object) Object {
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

func (e *Environment) GetFromCurrentScope(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Has(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Update(name string, val Object) bool {
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

func (e *Environment) GetStore() map[string]Object {
	return e.store
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
