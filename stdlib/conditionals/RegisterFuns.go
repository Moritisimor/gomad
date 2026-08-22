package conditionals

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("if", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 3 {
			return helpers.WrongArgs("if", 3, len(e))
		}

		cond := e[0]
		yes := e[1]
		no := e[2]

		evaluatedCond, err := eval.GetBoolean(cond, env)
		if err != nil {
			return helpers.Err("Error in condition of if:\n\t%s", err.Error())
		}

		if evaluatedCond {
			evaluated, err := eval.Eval(yes, env)
			if err != nil {
				return helpers.Err("Error in yes-branch of if:\n\t%s", err.Error())
			}

			return evaluated, nil
		}

		evaluated, err := eval.Eval(no, env)
		if err != nil {
			return helpers.Err("Error in no-branch of if:\n\t%s", err.Error())
		}

		return evaluated, nil
	})

	env.RegisterNative("=", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("=", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of =:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of =:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.String); ok {
			if b, ok := rhs.(value.String); ok {
				return value.NewBool(a.Val == b.Val), nil
			}
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.NewBool(a.Val == b.Val), nil
			}
		}

		if a, ok := lhs.(value.Boolean); ok {
			if b, ok := rhs.(value.Boolean); ok {
				return value.NewBool(a.Val == b.Val), nil
			}
		}

		if a, ok := lhs.(value.List); ok {
			if b, ok := rhs.(value.List); ok {
				if len(a.Val) != len(b.Val) {
					return value.NewBool(false), nil
				}

				for i := range len(a.Val) {
					if a.Val[i] != b.Val[i] {
						return value.NewBool(false), nil
					}
				}

				return value.NewBool(true), nil
			}
		}

		if _, ok := lhs.(value.Unit); ok {
			if _, ok := rhs.(value.Unit); ok {
				return value.NewBool(true), nil
			}
		}

		return value.NewBool(false), nil
	})

	env.RegisterNative(">", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs(">", 2, len(e))
		}

		lhs, err := eval.GetNumber(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of >:\n\t%s", err.Error())
		}

		rhs, err := eval.GetNumber(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of >:\n\t%s", err.Error())
		}

		return value.NewBool(lhs > rhs), nil
	})

	env.RegisterNative(">=", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs(">=", 2, len(e))
		}

		lhs, err := eval.GetNumber(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of >=:\n\t%s", err.Error())
		}

		rhs, err := eval.GetNumber(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of >=:\n\t%s", err.Error())
		}

		return value.NewBool(lhs >= rhs), nil
	})

	env.RegisterNative("<", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("<", 2, len(e))
		}

		lhs, err := eval.GetNumber(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of <:\n\t%s", err.Error())
		}

		rhs, err := eval.GetNumber(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of <:\n\t%s", err.Error())
		}

		return value.NewBool(lhs < rhs), nil
	})

	env.RegisterNative("<=", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("<=", 2, len(e))
		}

		lhs, err := eval.GetNumber(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of <=:\n\t%s", err.Error())
		}

		rhs, err := eval.GetNumber(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of <=:\n\t%s", err.Error())
		}

		return value.NewBool(lhs <= rhs), nil
	})
}
