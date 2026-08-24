package stringlib

import (
	"strconv"
	"strings"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("splitws", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("splitws", 1, len(e))
		}

		fields, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to splitws:\n\t%s", err.Error())
		}

		return value.NewList(func(strList []string) []value.Value {
			acc := []value.Value{}
			for _, s := range strList {
				acc = append(acc, value.NewString(s))
			}

			return acc
		}(strings.Fields(fields))), nil
	})

	env.RegisterNative("trim", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("trim", 1, len(e))
		}

		str, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to trim:\n\t%s", err.Error())
		}

		return value.NewString(strings.TrimSpace(str)), nil
	})

	env.RegisterNative("lower", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("lower", 1, len(e))
		}

		str, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to lower:\n\t%s", err.Error())
		}

		return value.NewString(strings.ToLower(str)), nil
	})

	env.RegisterNative("upper", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("upper", 1, len(e))
		}

		str, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to upper:\n\t%s", err.Error())
		}

		return value.NewString(strings.ToUpper(str)), nil
	})

	env.RegisterNative("sprint", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		acc := strings.Builder{}
		for i, exp := range e {
			evaluated, err := eval.Eval(exp, env)
			if err != nil {
				return helpers.Err("Error in argument %d of call to sprint:\n\t%s", i+1, err.Error())
			}

			acc.WriteString(evaluated.String())
		}

		return value.NewString(acc.String()), nil
	})

	env.RegisterNative("string_to_num", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("string_to_num", 1, len(e))
		}

		str, err := eval.GetString(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to string_to_num:\n\t%s", err.Error())
		}

		parsedNum, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return helpers.Err("Could not parse \"%s\" to a number", str)
		}

		return value.NewNumber(parsedNum), nil
	})
}
