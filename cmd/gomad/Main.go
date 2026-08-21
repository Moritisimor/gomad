package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
)

func main() {
	tokens, err := lexer.Tokenize("((+ 1 2) (* \"hello \" 4) (/ 10 50)) (hi)")
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
	}
}
