package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/pkg/lexer"
)


func main() {
	tokens, err := lexer.Tokenize("(\"string\" true false unit 10) (+ 1 2)")
	if err != nil {
		fmt.Printf("Error while tokenizing: %s\n", err.Error())
	}

	for _, t := range tokens {
		lexer.PrintToken(t)
	}
}
