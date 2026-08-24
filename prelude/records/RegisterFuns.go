package records

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	env.RegisterNative("record", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		rec := value.NewRecord(map[string]value.Value{})
		if len(e) < 1 {
			return helpers.Err("Record expected at least one argument, got 0")
		}

		for i, field := range e {
			f, ok := field.(expr.List)
			if !ok {
				return helpers.Err("Error in field %d of record: Expected a list", i+1)
			}

			if len(f.Val) != 2 {
				return helpers.Err("Error in field %d of record: Expected a list with length of 2, got: %d", i+1, len(f.Val))
			}

			fieldName, fieldExpr := f.Val[0], f.Val[1]
			if name, ok := fieldName.(expr.Symbol); ok {
				evaluated, err := eval.Eval(fieldExpr, env)
				if err != nil {
					return helpers.Err("Error in field %d of record:\n\t%s", i+1, err.Error())
				}

				rec.Val[name.Val] = evaluated
				continue
			}

			return helpers.Err("Error in field %d of record: Field name was expected to be a symbol (%s)", i+1, fieldName.String())
		}

		return rec, nil
	})

	env.RegisterNative(".", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return helpers.WrongArgs(".", 2, len(e))
		}

		rec, err := eval.GetRecord(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to .:\n\t%s", err.Error())
		}

		fieldName, ok := e[1].(expr.Symbol)
		if !ok {
			return helpers.Err("Error in argument 2 to .: Symbol expected (%s)", e[1].String())
		}

		val, ok := rec[fieldName.Val]
		if !ok {
			return helpers.Err("Record has no such field: '%s'", fieldName.Val)
		}

		return val, nil
	})

	env.RegisterNative("record_mut", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 3 {
			return helpers.WrongArgs("record_mut", 3, len(e))
		}

		rec, err := eval.GetRecord(e[0], env)
		if err != nil {
			return helpers.Err("Error in argument 1 to record_mut:\n\t%s", err.Error())
		}

		fieldName, ok := e[1].(expr.Symbol)
		if !ok {
			return helpers.Err("Error in argument 2 to record_mut: Symbol expected (%s)", e[1].String())
		}

		newExpr, err := eval.Eval(e[2], env)
		if err != nil {
			return helpers.Err("Error in argument 3 to record_mut:\n\t%s", err.Error())
		}

		rec[fieldName.Val] = newExpr
		return value.NewUnit(), nil
	})
}
