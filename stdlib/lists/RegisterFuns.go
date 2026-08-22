package lists

import (
	"fmt"

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
				return value.NewUnit(), fmt.Errorf("Error in argument %d to list:\n\t%s", i, err.Error())
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
			return value.NewUnit(), fmt.Errorf("Error in argument to car:\n\t%s", err.Error())
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
			return value.NewUnit(), fmt.Errorf("Error in argument to cdr:\n\t%s", err.Error())
		}

		if len(evaluated) == 0 {
			return value.NewUnit(), nil
		}

		return value.NewList(evaluated[1:]), nil
	})
}
