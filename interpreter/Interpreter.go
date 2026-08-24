package interpreter

import (
	"os"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/prelude"
	"github.com/Moritisimor/gomad/preludesrc"
	"github.com/Moritisimor/gomad/value"
)

type Interpreter struct {
	Env *value.Env
}

func New() *Interpreter {
	env := &value.Env{
		Parent:   nil,
		Bindings: map[string]value.Value{},
	}

	prelude.RegisterCommonPrelude(env)
	preludesrc.RegisterPreludeSrc(env)

	return &Interpreter{Env: env}
}

func NewNoSrcPrelude() *Interpreter {
	env := &value.Env{
		Parent: nil,
		Bindings: map[string]value.Value{},
	}

	prelude.RegisterCommonPrelude(env)
	return &Interpreter{Env: env}
}

func NewNoStdlib() *Interpreter {
	return &Interpreter{
		Env: &value.Env{
			Parent:   nil,
			Bindings: map[string]value.Value{},
		},
	}
}

func (i *Interpreter) Set(name string, val value.Value) {
	i.Env.Bindings[name] = val
}

func (i *Interpreter) Unset(name string) {
	delete(i.Env.Bindings, name)
}

func (i *Interpreter) RegisterNative(
	name string,
	fun func(e []expr.Expression, env *value.Env) (value.Value, error),
) {
	i.Env.RegisterNative(name, fun)
}

func (i *Interpreter) DoString(sourceCode string) (value.Value, error) {
	return eval.DoString(sourceCode, i.Env)
}

func (i *Interpreter) DoFile(path string) (value.Value, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return helpers.Err("Error while reading file: %s", err)
	}

	sourceCode := string(content)
	return eval.DoString(sourceCode, i.Env)
}
