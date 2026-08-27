package eval

import (
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

func GetString(e expr.Expr, env *value.Env) (string, error) {
	v, err := Eval(e, env)
	if err != nil {
		return "", err
	}
	if s, ok := v.(value.String); ok {
		return s.Val, nil
	}
	return "", typeErr("string", e, v)
}

func GetNumber(e expr.Expr, env *value.Env) (float64, error) {
	v, err := Eval(e, env)
	if err != nil {
		return 0, err
	}
	if n, ok := v.(value.Number); ok {
		return n.Val, nil
	}
	return 0, typeErr("number", e, v)
}

func GetBool(e expr.Expr, env *value.Env) (bool, error) {
	v, err := Eval(e, env)
	if err != nil {
		return false, err
	}
	if b, ok := v.(value.Boolean); ok {
		return b.Val, nil
	}
	return false, typeErr("bool", e, v)
}

func GetList(e expr.Expr, env *value.Env) (value.List, error) {
	v, err := Eval(e, env)
	if err != nil {
		return value.NewNil(), err
	}
	if l, ok := v.(value.List); ok {
		return l, nil
	}
	return value.NewNil(), typeErr("list", e, v)
}

func GetRecord(e expr.Expr, env *value.Env) (*value.Record, error) {
	v, err := Eval(e, env)
	if err != nil {
		return nil, err
	}
	if r, ok := v.(*value.Record); ok {
		return r, nil
	}
	return nil, typeErr("record", e, v)
}

func typeErr(expected string, e expr.Expr, v value.Value) *value.Error {
	return value.EvalErrf(
		"This expression was expected to evaluate to a %s, but it didn't: %s (%s)",
		expected, e, v,
	)
}
