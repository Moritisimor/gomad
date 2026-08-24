package arithmetics

import (
	"fmt"
	"math"
	"strings"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("+", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("+", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of +:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of +:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.Number{Val: a.Val + b.Val}, nil
			}
		}

		if a, ok := lhs.(value.String); ok {
			if b, ok := rhs.(value.String); ok {
				return value.String{Val: a.Val + b.Val}, nil
			}
		}

		return value.Unit{}, fmt.Errorf("Cannot add these values: (%s) and (%s)", lhs.String(), rhs.String())
	})

	env.RegisterNative("-", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("-", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return value.Unit{}, fmt.Errorf("Error in LHS of -:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of -:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.Number{Val: a.Val - b.Val}, nil
			}
		}

		return helpers.Err("Cannot subtract these expressions: (%s) and (%s)", lhs.String(), rhs.String())
	})

	env.RegisterNative("*", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("*", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of *:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of *:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.Number{Val: a.Val * b.Val}, nil
			}

			if b, ok := rhs.(value.String); ok {
				acc := strings.Builder{}
				for range int64(a.Val) {
					acc.WriteString(b.Val)
				}

				return value.String{Val: acc.String()}, nil
			}
		}

		if a, ok := lhs.(value.String); ok {
			if b, ok := rhs.(value.Number); ok {
				acc := strings.Builder{}
				for range int64(b.Val) {
					acc.WriteString(a.Val)
				}

				return value.String{Val: acc.String()}, nil
			}
		}

		return helpers.Err("Cannot multiply these expressions: (%s) and (%s)", lhs.String(), rhs.String())
	})

	env.RegisterNative("/", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			helpers.WrongArgs("/", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of /:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in RHS of /:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				if b.Val == 0 {
					return helpers.Err("Division by zero")
				}

				return value.Number{Val: a.Val / b.Val}, nil
			}
		}

		return helpers.Err("Cannot divide these expressions: (%s) and (%s)", lhs.String(), rhs.String())
	})

	env.RegisterNative("mod", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs("mod", 2, len(e))
		}

		lhs, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in LHS of mod:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return helpers.Err("Error in LHS of mod:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.Number{Val: math.Mod(a.Val, b.Val)}, nil
			}
		}

		return helpers.Err("Cannot apply modulo on these expressions: (%s) and (%s)", lhs.String(), rhs.String())
	})
}
