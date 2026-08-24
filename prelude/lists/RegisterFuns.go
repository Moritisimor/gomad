package lists

import (
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
