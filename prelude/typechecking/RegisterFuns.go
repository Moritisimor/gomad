package typechecking

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("isunit", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isunit", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isunit:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Unit)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isstr", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isstr", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isstr:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.String)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isnum", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isnum", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isnum:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Number)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isbool", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isbool", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isbool:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Boolean)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("islist", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("islist", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to islist:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.List)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isrecord", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isrecord", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isrecord:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Record)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isfun", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isfun", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isfun:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Lambda)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("isnative", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("isnative", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to isnative:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.NativeFunction)
		return value.NewBool(ok), nil
	})

	env.RegisterNative("ismac", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 1 {
			return helpers.WrongArgs("ismac", 1, len(e))
		}

		evaluated, err := eval.Eval(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument to ismac:\n\t%s", err.Error())
		}

		_, ok := evaluated.(value.Macro)
		return value.NewBool(ok), nil
	})
} 
