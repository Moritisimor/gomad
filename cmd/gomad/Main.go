package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/pkg/expr"
	"github.com/Moritisimor/gomad/pkg/lexer"
	"github.com/Moritisimor/gomad/pkg/parser"
)


func main() {
	tokens, err := lexer.Tokenize("(\"string\" true false unit 10) (+ 1 2)")
	if err != nil {
		fmt.Printf("Error while tokenizing: %s\n", err.Error())
	}

	for _, t := range tokens {
		lexer.PrintToken(t)
	}

	parsedExprs, _, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("Error while parsing: %s\n", err.Error())
	}

	for _, e := range parsedExprs {
		expr.PrintExpr(e)
	}
}
