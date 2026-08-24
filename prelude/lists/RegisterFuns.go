package lists

import (
	"slices"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("list", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		acc := []value.Value{}
		for i, e := range e {
			evaluated, err := eval.Eval(e, env)
			if err != nil {
				return helpers.Err("Error in argument %d to list:\n\t%s", i+1, err.Error())
			}

			acc = append(acc, evaluated)
		}

		return value.NewList(acc), nil
	})

	env.RegisterNative("car", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("car", 1, len(e))
		}

		evaluated, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to car:\n\t%s", err.Error())
		}

		if len(evaluated) == 0 {
			return value.NewUnit(), nil
		}

		return evaluated[0], nil
	})

	env.RegisterNative("cdr", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("cdr", 1, len(e))
		}

		evaluated, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to cdr:\n\t%s", err.Error())
		}

		if len(evaluated) == 0 {
			return value.NewUnit(), nil
		}

		return value.NewList(evaluated[1:]), nil
	})

	env.RegisterNative("cons", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("cons", 2, len(e))
		}

		elem, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to cons:\n\t%s", err.Error())
		}

		l, err := eval.GetList(e[1], env)
		if err != nil {
			return helpers.Err("Error in argument 2 to cons:\n\t%s", err.Error())
		}

		return value.NewList(slices.Insert(l, 0, elem)), nil
	})

	env.RegisterNative("append", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("append", 2, len(e))
		}

		l1, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to append:\n\t%s", err.Error())
		}

		l2, err := eval.GetList(e[1], env)
		if err != nil {
			return helpers.Err("Error in argument 2 to append:\n\t%s", err.Error())
		}

		return value.NewList(slices.Concat(l1, l2)), nil
	})

	env.RegisterNative("push", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("push", 2, len(e))
		}

		elem, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to push:\n\t%s", err.Error())
		}

		l, err := eval.GetList(e[1], env)
		if err != nil {
			return helpers.Err("Error in argument 2 to push:\n\t%s", err.Error())
		}

		return value.NewList(append(l, elem)), nil
	})

	env.RegisterNative("first", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("first", 1, len(e))
		}

		l, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to first:\n\t%s", err.Error())
		}

		if len(l) == 0 {
			return helpers.Err("List holds no such index")
		}

		return l[0], nil
	})

	env.RegisterNative("last", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("last", 1, len(e))
		}

		l, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to last:\n\t%s", err.Error())
		}

		if len(l) == 0 {
			return helpers.Err("List holds no such index")
		}

		return l[len(l)-1], nil
	})

	env.RegisterNative("len", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("len", 1, len(e))
		}

		l, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to len:\n\t%s", err.Error())
		}

		return value.NewNumber(float64(len(l))), nil
	})

	env.RegisterNative("nth", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("nth", 1, len(e))
		}

		evaluated, err := eval.GetList(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to nth:\n\t%s", err.Error())
		}

		index, err := eval.GetNumber(e[1], env)
		if err != nil {
			return helpers.Err("Error in argument 2 to nth:\n\t%s", err.Error())
		}

		if len(e) < int(index) || index < 0 {
			return helpers.Err("List has no such index")
		}

		return evaluated[int(index)], nil
	})
}
