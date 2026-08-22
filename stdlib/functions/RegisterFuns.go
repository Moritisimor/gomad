package functions

import (
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("lambda", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("lambda", 2, len(e))
		}

		if params, ok := e[0].(expr.List); ok {
			paramNames := []string{}
			for i, p := range params.Val {
				if s, ok := p.(expr.Symbol); ok {
					paramNames = append(paramNames, s.Val)
					continue
				}

				return helpers.Err("Non-symbol in parameter list of lambda (argument %d)", i+1)
			}

			return value.Lambda{
				Captured: env,
				Body:     e[1],
				Params:   paramNames,
			}, nil
		}

		return helpers.Err("Expected parameter list after lambda")
	})

	env.RegisterNative("letfun", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 3 {
			helpers.WrongArgs("letfun", 3, len(e))
		}

		var funName string
		if s, ok := e[0].(expr.Symbol); ok {
			funName = s.Val
		} else {
			return helpers.Err("Function name was expected to be a symbol, got: %s", expr.SprintExpr(e[0]))
		}

		params := []string{}
		if p, ok := e[1].(expr.List); ok {
			for i, p := range p.Val {
				if s, ok := p.(expr.Symbol); ok {
					params = append(params, s.Val)
					continue
				}

				return helpers.Err("Non-symbol in parameter list of letfun (argument %d)", i+1)
			}
		}

		if err := env.SetBinding(funName, value.Lambda{
			Params: params,
			Body: e[2],
			Captured: env,
		}); err != nil {
			return helpers.Err("Error while binding function '%s':\n\t%s", funName, err.Error())
		}

		return value.NewUnit(), nil
	})
}
