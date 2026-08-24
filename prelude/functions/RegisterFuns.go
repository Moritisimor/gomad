package functions

import (
	"github.com/Moritisimor/gomad/eval"
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
			return helpers.Err("Function name was expected to be a symbol, got: %s", e[0].String())
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
			Params:   params,
			Body:     e[2],
			Captured: env,
		}); err != nil {
			return helpers.Err("Error while binding function '%s':\n\t%s", funName, err.Error())
		}

		return value.NewUnit(), nil
	})

	env.RegisterNative("do", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		var lastExpr value.Value
		lastExpr = value.NewUnit()
		for i, exp := range e {
			evaluated, err := eval.Eval(exp, env)
			if err != nil {
				return helpers.Err("Error in expression %d of do:\n\t%s", i, err.Error())
			}

			lastExpr = evaluated
		}

		return lastExpr, nil
	})

	env.RegisterNative("scoped", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("scoped", 2, len(e))
		}

		bindings, ok := e[0].(expr.List)
		if !ok {
			return helpers.Err("Error in argument 1 to scoped: expected list (%s)", e[0].String())
		}

		acc := map[string]value.Value{}
		for i, b := range bindings.Val {
			l, ok := b.(expr.List)
			if !ok {
				return helpers.Err("Error in element %d of bindings-list in scoped: expected list(%s)", i+1, b.String())
			}

			if len(l.Val) != 2 {
				return helpers.Err(
					"Error in element %d of bindings-list in scoped: list was expected to have length of 2, got: %d",
					i+1, len(e),
				)
			}

			s, ok := l.Val[0].(expr.Symbol)
			if !ok {
				return helpers.Err(
					"Error in element %d of bindings-list in scoped: binding-name was expected to be a symbol (%s)",
					i+1, l.Val[0].String(),
				)
			}

			evaluated, err := eval.Eval(l.Val[1], env)
			if err != nil {
				return helpers.Err("Error in element %d of bindings-list in scoped:\n\t%s", i+1, err.Error())
			}

			acc[s.Val] = evaluated
		}

		thisScope := value.Env{
			Bindings: acc,
			Parent: env,
		}

		evaluated, err := eval.Eval(e[1], &thisScope)
		if err != nil {
			return helpers.Err("Error while evaluating expression of scoped:\n\t%s", err.Error())
		}

		return evaluated, nil
	})
}
