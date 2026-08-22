package variables

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("let", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("let", 2, len(e))
		}

		var varName string
		if s, ok := e[0].(expr.Symbol); ok {
			varName = s.Val
		} else {
			return helpers.Err("Name of let-binding was expected to be a symbol (%s)", expr.SprintExpr(e[0]))
		}

		evaluated, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error while evaluating value of let-binding:\n\t%s", err.Error())
		}

		if err := env.SetBinding(varName, evaluated); err != nil {
			return helpers.Err("Error while binding variable '%s':\n\t%s", varName, err.Error())
		}

		return value.NewUnit(), nil
	})

	env.RegisterNative("mut", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("mut", 2, len(e))
		}

		var varName string
		if s, ok := e[0].(expr.Symbol); ok {
			varName = s.Val
		} else {
			return helpers.Err("First argument to mut was expected to be a symbol (%s)", expr.SprintExpr(e[0]))
		}

		evaluated, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error while evaluating new value of mutation of variable '%s':\n\t%s", varName, err.Error())
		}

		if err := env.MutateBinding(varName, evaluated); err != nil {
			return helpers.Err("Error while mutating binding '%s':\n\t%s", varName, err.Error())
		}
		
		return value.NewUnit(), nil
	})
}
