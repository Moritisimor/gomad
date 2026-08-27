package interpreter

import (
	"os"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/value"
)

type Interpreter struct {
	Env *value.Env
}

func New(args ...string) *Interpreter {
	env := value.RootEnv()

	eval.RegisterNatives(env)
	eval.RegisterOSNatives(env)

	argVals := make([]value.Value, len(args))
	for i, a := range args {
		argVals[i] = value.String{Val: a}
	}
	_ = env.Set("args", value.FromVec(argVals))

	if err := eval.LoadStdlib(env); err != nil {
		panic("stdlib failed to load: " + err.Error())
	}

	return &Interpreter{Env: env}
}

func NewNoStdlib(args ...string) *Interpreter {
	env := value.RootEnv()
	eval.RegisterNatives(env)
	eval.RegisterOSNatives(env)
	argVals := make([]value.Value, len(args))
	for i, a := range args {
		argVals[i] = value.String{Val: a}
	}
	_ = env.Set("args", value.FromVec(argVals))
	return &Interpreter{Env: env}
}

func NewNoSrcPrelude(args ...string) *Interpreter { return NewNoStdlib(args...) }

func (i *Interpreter) RegisterNative(name string, fn value.NativeFunc) error {
	return i.Env.Set(name, fn)
}

func (i *Interpreter) DoString(source string) (value.Value, error) {
	return eval.DoString(source, i.Env)
}

func (i *Interpreter) DoFile(path string) (value.Value, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, &value.Error{Kind: value.ErrIo, Msg: path + ": " + err.Error()}
	}
	return eval.DoString(string(content), i.Env)
}

func (i *Interpreter) EvalExpr(e expr.Expr) (value.Value, error) {
	return eval.Eval(e, i.Env)
}

func Parse(source string) ([]expr.Expr, error) {
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		return nil, &value.Error{Kind: value.ErrTokenize, Msg: err.Error()}
	}
	return parser.ParseProgram(tokens)
}

func (i *Interpreter) GetGlobal(name string) (value.Value, error) {
	return i.Env.Get(name)
}

func (i *Interpreter) GetBinding(name string) (value.Value, error) { return i.GetGlobal(name) }
func (i *Interpreter) Set(name string, val value.Value)            { i.Env.Define(name, val) }
func (i *Interpreter) Unset(name string)                           { delete(i.Env.Bindings, name) }
