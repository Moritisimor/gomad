package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/stdlib"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	tokens, err := lexer.Tokenize("((lambda (n) (+ \"Hello, \" n)) \"John\")")
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
		Parent:   nil,
	}

	stdlib.RegisterStdlib(&env)
	for _, e := range parsedExprs {
		val, err := eval.Eval(e, &env)
		if val == nil {
			fmt.Println("this is nil, but WHY???")
			return
		}

		if err != nil {
			fmt.Printf("Error while evaluating:\n\t%s\n", err.Error())
			return
		}

		fmt.Println(val.String())
	}
}
