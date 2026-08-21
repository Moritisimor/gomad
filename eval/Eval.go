package eval

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

func Eval(e expr.Expression, env *value.Env) (value.Value, error) {
	switch exp := e.(type) {
	case expr.Number:
		return value.Number{ Val: exp.Val }, nil

	case expr.String:
		return value.String{ Val: exp.Val }, nil

	case expr.Boolean:
		return value.Boolean{ Val: exp.Val }, nil

	case expr.Unit:
		return value.Unit{}, nil

	case expr.Symbol:
		return env.GetBinding(exp.Val)

	case expr.List:
		if len(exp.Val) == 0 {
			return value.List{ Val: []value.Value{} }, nil
		}

		funExpr, err := Eval(exp.Val[0], env)
		if err != nil {
			return value.Unit{}, err
		}

		switch fun := funExpr.(type) {
		case value.NativeFunction:
			return fun.Callback(exp.Val[1:], env)

		case value.Lambda:
			panic("Lambdas are not yet implemented!")

		case value.Macro:
			panic("Macros are not yet implemented!")

		default:
			return value.Unit{}, fmt.Errorf("Attempt to invoke non-invocable value: '%s'", fun.String())
		}
	}

	return value.Unit{}, fmt.Errorf("Unknown expression, this should not have happened")
}
