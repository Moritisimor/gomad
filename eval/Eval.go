package eval

import (
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func Eval(e expr.Expression, env *value.Env) (value.Value, error) {
	switch exp := e.(type) {
	case expr.Number:
		return value.Number{Val: exp.Val}, nil

	case expr.String:
		return value.String{Val: exp.Val}, nil

	case expr.Boolean:
		return value.Boolean{Val: exp.Val}, nil

	case expr.Unit:
		return value.Unit{}, nil

	case expr.Symbol:
		return env.GetBinding(exp.Val)

	case expr.List:
		if len(exp.Val) == 0 {
			return value.List{Val: []value.Value{}}, nil
		}

		invocationArgs := exp.Val[1:]
		funExpr, err := Eval(exp.Val[0], env)
		if err != nil {
			return value.Unit{}, err
		}

		switch fun := funExpr.(type) {
		case value.NativeFunction:
			return fun.Callback(invocationArgs, env)

		case value.Lambda:
			if len(fun.Params) != len(invocationArgs) {
				return helpers.Err(
					"Lambda invoked with wrong amount of args. Expected: %d, Got: %d",
					len(fun.Params), len(invocationArgs),
				)
			}

			thisEnv := value.Env{
				Bindings: map[string]value.Value{},
				Parent:   fun.Captured,
			}

			for i := range len(invocationArgs) {
				evaluated, err := Eval(invocationArgs[i], env)
				if err != nil {
					return helpers.Err("Error while evaluating argument %d of lambda-invocation:\n\t%s", i+1, err.Error())
				}

				thisEnv.Bindings[fun.Params[i]] = evaluated
			}

			evaluated, err := Eval(fun.Body, &thisEnv)
			if err != nil {
				return helpers.Err("Error in body of lambda:\n\t%s", err.Error())
			}

			return evaluated, nil

		case value.Macro:
			

		default:
			return helpers.Err("Attempt to invoke non-invocable value: '%s'", e.String())
		}
	}

	return helpers.Err("Unknown expression, this should not have happened")
}
