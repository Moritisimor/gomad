package eval

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func GetString(e expr.Expression, env *value.Env) (string, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return "", err
	}

	if s, ok := evaluated.(value.String); ok {
		return s.Val, nil
	}

	return helpers.BadType[string]("string", e)
}

func GetNumber(e expr.Expression, env *value.Env) (float64, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return 0, err
	}

	if f, ok := evaluated.(value.Number); ok {
		return f.Val, nil
	}

	return helpers.BadType[float64]("number", e)
}

func GetBoolean(e expr.Expression, env *value.Env) (bool, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return false, err
	}

	if b, ok := evaluated.(value.Boolean); ok {
		return b.Val, nil
	}

	return helpers.BadType[bool]("boolean", e)
}

func GetList(e expr.Expression, env *value.Env) ([]value.Value, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return []value.Value{}, err
	}

	if b, ok := evaluated.(value.List); ok {
		return b.Val, nil
	}

	return helpers.BadType[[]value.Value]("list", e)
}

func GetRecord(e expr.Expression, env *value.Env) (map[string]value.Value, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return map[string]value.Value{}, err
	}

	if r, ok := evaluated.(value.Record); ok {
		return r.Val, nil
	}

	return helpers.BadType[map[string]value.Value]("record", e)
}

func GetLambda(e expr.Expression, env *value.Env) (value.Lambda, error) {
	evaluated, err := Eval(e, env)
	if err != nil {
		return value.Lambda{}, nil
	}

	if b, ok := evaluated.(value.Lambda); ok {
		return b, nil
	}

	return helpers.BadType[value.Lambda]("lambda", e)
}

func GetNative(
	e expr.Expression,
	env *value.Env,
) (func(e []expr.Expression, env *value.Env) (value.Value, error), error) {
	fun := func(e []expr.Expression, env *value.Env) (value.Value, error) {
		return value.Unit{}, fmt.Errorf("You're not supposed to invoke me!")
	}

	evaluated, err := Eval(e, env)
	if err != nil {
		return fun, err
	}

	if n, ok := evaluated.(value.NativeFunction); ok {
		return n.Callback, nil
	}

	return fun, err
}
