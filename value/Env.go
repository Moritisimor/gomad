package value

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
)

type Env struct {
	Bindings map[string]Value
	Parent   *Env
}

func (e *Env) GetBinding(name string) (Value, error) {
	val, ok := e.Bindings[name]
	if !ok {
		if e.Parent == nil {
			return Unit{}, fmt.Errorf("No such binding: '%s'", name)
		}

		return e.Parent.GetBinding(name)
	}

	return val, nil
}

func (e *Env) SetBinding(name string, val Value) error {
	if _, ok := e.Bindings[name]; !ok {
		e.Bindings[name] = val
		return nil
	}

	return fmt.Errorf("Cannot set binding '%s': Already exists within this scope", name)
}

func (e *Env) MutateBinding(name string, val Value) error {
	_, ok := e.Bindings[name]
	if !ok {
		if e.Parent == nil {
			return fmt.Errorf("No such binding: '%s'", name)
		}

		return e.Parent.MutateBinding(name, val)
	}

	e.Bindings[name] = val
	return nil
}

func (e *Env) RegisterNative(name string, val func(e []expr.Expression, env *Env) (Value, error)) {
	e.SetBinding(name, NativeFunction{val})
}
