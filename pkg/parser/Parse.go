package parser

import (
	"fmt"

	"github.com/Moritisimor/gomad/pkg/expr"
	"github.com/Moritisimor/gomad/pkg/lexer"
)

func Parse(tokens []lexer.Token) ([]expr.Expression, int, error) {
	acc := []expr.Expression{}
	left := tokens
	steps := 0

	for len(left) != 0 {
		steps++

		t := left[0]
		lexer.PrintToken(t)

		switch tok := t.(type) {
		case lexer.NUMLIT:
			acc = append(acc, expr.Number{ Val: tok.Val })
			left = left[1:]

		case lexer.STRINGLIT:
			acc = append(acc, expr.String{ Val: tok.Val })
			left = left[1:]

		case lexer.BOOLLIT:
			acc = append(acc, expr.Boolean{ Val: tok.Val })
			left = left[1:]

		case lexer.UNITLIT:
			acc = append(acc, expr.Boolean{})
			left = left[1:]

		case lexer.SYMBOL:
			acc = append(acc, expr.Symbol{ Val: tok.Val })
			left = left[1:]

		case lexer.LPAREN:
			exprs, steps, err := Parse(left[1:])
			if err != nil {
				return acc, 0, err
			}

			acc = append(acc, expr.List{ Val: exprs })
			left = left[steps+1:]

		case lexer.RPAREN:
			left = left[1:]
			return acc, steps, nil
		}
	}

	return acc, steps, fmt.Errorf("List was never closed")
}
