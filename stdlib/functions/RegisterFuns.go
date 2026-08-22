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
}
