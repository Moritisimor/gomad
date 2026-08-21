package parser

import (
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
)

// Parser is a struct that keeps track of the parser's current state, mainly the current token index.
// It's not really meant to be instantiated on its own, just let the Parse function do it for you.
type Parser struct {
	Tokens []lexer.Token
	Cursor int
}

func (p *Parser) parse() ([]expr.Expression, error) {
	acc := []expr.Expression{}

	for p.Cursor < len(p.Tokens) {
		t := p.Tokens[p.Cursor]

		switch tok := t.(type) {
		case lexer.NUMLIT:
			acc = append(acc, expr.Number{ Val: tok.Val })
			p.Cursor++

		case lexer.STRINGLIT:
			acc = append(acc, expr.String{ Val: tok.Val })
			p.Cursor++

		case lexer.BOOLLIT:
			acc = append(acc, expr.Boolean{ Val: tok.Val })
			p.Cursor++

		case lexer.UNITLIT:
			acc = append(acc, expr.Unit{})
			p.Cursor++

		case lexer.SYMBOL:
			acc = append(acc, expr.Symbol{ Val: tok.Val })
			p.Cursor++

		case lexer.LPAREN:
			p.Cursor++
			exprs, err := p.parse()
			if err != nil {
				return acc, err
			}

			acc = append(acc, expr.List{ Val: exprs })

		case lexer.RPAREN:
			p.Cursor++
			return acc, nil
		}
	}

	return acc, nil
}

// Parse is the user-facing parsing function.
// It takes a slice of tokens and transforms them into a slice of expressions.
func Parse(tokens []lexer.Token) ([]expr.Expression, error) {
	parser := Parser{
		Tokens: tokens,
		Cursor: 0,
	}

	return parser.parse()
}
