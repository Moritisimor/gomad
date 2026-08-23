package iolib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) {
	args := []value.Value{}
	for _, a := range os.Args[1:] {
		args = append(args, value.NewString(a))
	}

	env.Bindings["args"] = value.NewList(args)

	env.RegisterNative("print", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		for i, exp := range e {
			evaluated, err := eval.Eval(exp, env)
			if err != nil {
				return helpers.Err("Error in argument %d to print: %s", i+1, err.Error())
			}

			fmt.Print(evaluated)
		}

		return value.Unit{}, nil
	})

	env.RegisterNative("println", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		for i, exp := range e {
			evaluated, err := eval.Eval(exp, env)
			if err != nil {
				return helpers.Err("Error in argument %d to print:\n\t%s", i+1, err.Error())
			}

			fmt.Print(evaluated.String())
		}

		fmt.Print("\n")
		return value.Unit{}, nil
	})

	env.RegisterNative("readln", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		for i, exp := range e {
			evaluated, err := eval.Eval(exp, env)
			if err != nil {
				return helpers.Err("Error in argument %d to readln:\n\t%s", i+1, err.Error())
			}

			fmt.Print(evaluated.String())
		}

		scanner := bufio.NewReader(os.Stdout)
		read, err := scanner.ReadString('\n')
		if err != nil {
			helpers.Err("Error while reading line from stdin:\n\t%s", err.Error())
		}

		return value.String{Val: strings.TrimRight(read, "\n")}, nil
	})
}
