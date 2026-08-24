package interpreter

import (
	"os"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/prelude"
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
	tokens, err := lexer.Tokenize(sourceCode)
	if err != nil {
		return helpers.Err("Error while tokenizing: %s", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return helpers.Err("Error while parsing: %s", err)
	}

	var lastExpr value.Value
	lastExpr = value.Unit{}
	for _, exp := range ast {
		evaluated, err := eval.Eval(exp, i.Env)
		if err != nil {
			return helpers.Err("Uncaught error:\n\t%s", err)
		}

		lastExpr = evaluated
	}

	return lastExpr, nil
}

func (i *Interpreter) DoFile(path string) (value.Value, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return helpers.Err("Error while reading file: %s", err)
	}

	sourceCode := string(content)
	tokens, err := lexer.Tokenize(sourceCode)
	if err != nil {
		return helpers.Err("Error while tokenizing: %s", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return helpers.Err("Error while parsing: %s", err)
	}

	var lastExpr value.Value
	lastExpr = value.Unit{}
	for _, exp := range ast {
		evaluated, err := eval.Eval(exp, i.Env)
		if err != nil {
			return helpers.Err("Uncaught error: %s\n", err)
		}

		lastExpr = evaluated
	}

	return lastExpr, nil
}
