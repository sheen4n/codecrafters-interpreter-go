package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Define is used to define and set a value in the current environment.
func (e *Environment) Define(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

// Assign is used to assign a value to an already defined variable
// in the current or outer environment.
func (e *Environment) Assign(name string, obj Object) Object {
	if _, ok := e.store[name]; ok {
		e.store[name] = obj
		return obj
	}
	if e.outer != nil {
		return e.outer.Assign(name, obj)
	}
	return obj
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
