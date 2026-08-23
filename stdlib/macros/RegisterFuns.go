package macros

import (
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("letmac", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) < 3 {
			return helpers.WrongArgs("letmac", 3, len(e))
		}

		macNameExpr, ok := e[0].(expr.Symbol)
		if !ok {
			return helpers.Err("Macro name was expected to be a symbol, got: %s", e[0].String())
		}

		macName := macNameExpr.Val
		macParamsExpr, ok := e[1].(expr.List)
		if !ok {
			return helpers.Err("Macro params were expected to be a list, got: %s", e[1].String())
		}

		macParams := []string{}
		for i, param := range macParamsExpr.Val {
			if s, ok := param.(expr.Symbol); ok {
				macParams = append(macParams, s.Val)
				continue
			}

			return helpers.Err("Non-symbol in parameter list of letmac '%s' (argument %d)", macName, i+1)
		}
	
		env.SetBinding(macName, value.Macro{
			Params: macParams,
			Expressions: e[2:],
		})

		return value.NewUnit(), nil
	})
}
