package exceptions

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("try", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("try", 2, len(e))
		}

		yes, err := eval.Eval(e[0], env)
		if err != nil {
			no, err := eval.Eval(e[1], env)
			if err != nil {
				return helpers.Err("Error in catch-block of try:\n\t%s", err.Error())
			}

			return no, nil
		}

		return yes, nil
	})

	env.RegisterNative("throw", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("throw", 1, len(e))
		}

		throwExpr, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in throw-expression:\n\t%s", err.Error())
		}

		return helpers.Err("%s", throwExpr)
	})
}
