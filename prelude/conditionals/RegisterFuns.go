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

	env.RegisterNative("switch", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) < 2 {
			return helpers.Err("Switch expected at least 2 arguments, but it got %d", len(e))
		}

		scrutinee := e[0]
		cases := e[1:]

		evaluatedScrutinee, err := eval.Eval(scrutinee, env)
		if err != nil {
			return helpers.Err("Error in scrutinee of match-expression:\n\t%s", err.Error())
		}

		for i, c := range cases {
			if l, ok := c.(expr.List); ok {
				if len(l.Val) != 2 {
					return helpers.Err("Switch-arm %d has received bad syntax (2 elements in list expected)", i+1)
				}

				if s, ok := l.Val[0].(expr.Symbol); ok {
					if s.Val == "_" {
						evaluated, err := eval.Eval(l.Val[1], env)
						if err != nil {
							return helpers.Err("Error in value part of switch-arm %d:\n\t%s", i+1, err.Error())
						}

						return evaluated, nil
					}
				}

				matcher, err := eval.Eval(l.Val[0], env)
				if err != nil {
					return helpers.Err("Error in case part of switch-arm %d:\n\t%s", i+1, err.Error())
				}

				if matcher == evaluatedScrutinee {
					evaluated, err := eval.Eval(l.Val[1], env)
					if err != nil {
						return helpers.Err("Error in value part of switch-arm %d:\n\t%s", i+1, err.Error())
					}

					return evaluated, nil
				}

				continue
			}

			return helpers.Err("Switch-arm %d has received bad syntax (list expected)", i+1)
		}

		return value.NewUnit(), nil
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

	env.RegisterNative("not", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("not", 1, len(e))
		}

		cond, err := eval.GetBoolean(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to not:\n\t%s", err.Error())
		}

		return value.NewBool(!cond), nil
	})

	env.RegisterNative("and", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("and", 2, len(e))
		}

		lhs, err := eval.GetBoolean(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of and:\n\t%s", err.Error())
		}

		if !lhs {
			return value.NewBool(false), nil
		}

		rhs, err := eval.GetBoolean(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of and:\n\t%s", err.Error())
		}

		return value.NewBool(rhs), nil
	})

	env.RegisterNative("or", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("or", 2, len(e))
		}

		lhs, err := eval.GetBoolean(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of or:\n\t%s", err.Error())
		}

		if lhs {
			return value.NewBool(true), nil
		}

		rhs, err := eval.GetBoolean(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of or:\n\t%s", err.Error())
		}

		return value.NewBool(rhs), nil
	})
}
