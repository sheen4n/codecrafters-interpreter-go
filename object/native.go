package object

type NativeFunction struct {
	Fn func(args ...Object) Object
}

func (n *NativeFunction) Type() ObjectType { return NATIVE_FUNCTION_OBJ }
func (n *NativeFunction) Inspect() string  { return "<native fn>" }
