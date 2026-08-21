package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	tokens, err := lexer.Tokenize("(+ (+ \"hello \" \"world\") (+ 5 10))")
	if err != nil {
		fmt.Printf("Error while tokenizing: %s\n", err.Error())
		return
	}

	for _, t := range tokens {
		lexer.PrintToken(t)
	}

	parsedExprs, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("Error while parsing: %s\n", err.Error())
		return
	}

	fmt.Printf("Parsed %d expressions\n", len(parsedExprs))
	for _, e := range parsedExprs {
		expr.PrintExpr(e)
		fmt.Print("\n")
	}

	env := value.Env{
		Bindings: map[string]value.Value{},
		Parent: nil,
	}

	env.RegisterNative("+", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		if len(e) != 2 {
			return value.Unit{}, fmt.Errorf("Native function + is given the wrong amount of arguments (expected 2)")
		}

		lhs, err := eval.Eval(e[0], env) 
		if err != nil {
			return value.Unit{}, fmt.Errorf("Error in RHS of +:\n\t%s", err.Error())
		}

		rhs, err := eval.Eval(e[1], env)
		if err != nil {
			return value.Unit{}, fmt.Errorf("Error in LHS of +:\n\t%s", err.Error())
		}

		if a, ok := lhs.(value.Number); ok {
			if b, ok := rhs.(value.Number); ok {
				return value.Number{ Val: a.Val + b.Val }, nil
			}
		}

		return value.Unit{}, fmt.Errorf("Cannot add these values: %s and %s", lhs.String(), rhs.String())
	})

	for _, e := range parsedExprs {
		val, err := eval.Eval(e, &env)
		if err != nil {
			fmt.Printf("Error while evaluating: %s\n", err.Error())
			return
		}

		fmt.Println(val.String())
	}
}
