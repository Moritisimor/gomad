package eval

import (
	"reflect"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

func ifImpl(args []expr.Expr, env *value.Env) (value.Value, error) {
	if len(args) != 3 {
		return nil, arityErr("if", 3, len(args))
	}
	cond, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}
	switch c := cond.(type) {
	case value.Boolean:
		if c.Val {
			return Eval(args[1], env)
		}
		return Eval(args[2], env)
	default:
		return nil, value.EvalErrf(
			"Condition of if-construct does not evaluate to a bool: %s", cond)
	}
}

func doImpl(args []expr.Expr, env *value.Env) (value.Value, error) {
	return EvalSeq(args, env)
}

func switchImpl(args []expr.Expr, env *value.Env) (value.Value, error) {
	if len(args) == 0 {
		return nil, arityErr("switch", 2, 0)
	}
	scrutinee, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}
	for _, caseExpr := range args[1:] {
		items, ok := caseExpr.(expr.List)
		if !ok || len(items.Val) != 2 {
			return nil, value.EvalErr("Malformed switch-arm syntax")
		}
		matcher, onMatch := items.Val[0], items.Val[1]
		if m, ok := matcher.(expr.Symbol); ok && m.Val == "_" {
			return Eval(onMatch, env)
		}
		matcherVal, err := Eval(matcher, env)
		if err != nil {
			return nil, err
		}
		if value.Equal(matcherVal, scrutinee) {
			return Eval(onMatch, env)
		}
	}
	return value.Unit{}, nil
}

func scopedImpl(args []expr.Expr, env *value.Env) (value.Value, error) {
	if len(args) != 2 {
		return nil, arityErr("scoped", 2, len(args))
	}
	bindingPairs, ok := args[0].(expr.List)
	if !ok {
		return nil, value.EvalErr(
			"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')")
	}
	thisEnv := value.NewEnv(env)
	for _, pair := range bindingPairs.Val {
		items, ok := pair.(expr.List)
		if !ok || len(items.Val) != 2 {
			return nil, value.EvalErr(
				"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')")
		}
		nameSym, ok := items.Val[0].(expr.Symbol)
		if !ok {
			return nil, value.EvalErr(
				"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')")
		}
		val, err := Eval(items.Val[1], env)
		if err != nil {
			return nil, err
		}
		if err := thisEnv.Set(nameSym.Val, val); err != nil {
			return nil, err
		}
	}
	return Eval(args[1], thisEnv)
}

func tryImpl(args []expr.Expr, env *value.Env) (value.Value, error) {
	if len(args) != 2 {
		return nil, arityErr("try", 2, len(args))
	}
	v, err := Eval(args[0], env)
	if err != nil {
		if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrEval {
			return Eval(args[1], env)
		}
		return nil, err
	}
	return v, nil
}

var coreFormHandles []CoreBinding

var coreFormPtrs map[uintptr]string

var coreFormNames map[string]int

func init() {
	coreFormHandles = []CoreBinding{
		{"if", ifImpl},
		{"do", doImpl},
		{"switch", switchImpl},
		{"scoped", scopedImpl},
		{"try", tryImpl},
	}
	coreFormPtrs = make(map[uintptr]string, len(coreFormHandles))
	coreFormNames = make(map[string]int, len(coreFormHandles))
	for _, c := range coreFormHandles {
		coreFormPtrs[reflect.ValueOf(c.Fn).Pointer()] = c.Name
		coreFormNames[c.Name] = 1
	}
}

func isCoreFormName(name string) bool {
	_, ok := coreFormNames[name]
	return ok
}

func coreFormName(fn value.NativeFunc) string {
	return coreFormPtrs[reflect.ValueOf(fn).Pointer()]
}

type CoreBinding struct {
	Name string
	Fn   value.NativeFunc
}

func CoreFormBindings() []CoreBinding {
	return coreFormHandles
}
