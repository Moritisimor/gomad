package lexer

import (
	"fmt"
	"strings"
)

func parseSymbol(left []rune) (string, []rune, error) {
	acc := strings.Builder{}
	wasEnded := false
	steps := 0

	for _, c := range left {
		if isTerminator(c) {
			wasEnded = true
			break
		}

		acc.WriteRune(c)
		steps++
	}

	if !wasEnded {
		return "", []rune{}, fmt.Errorf("Symbol was never terminated! Got to \"%s\"", acc.String())
	}

	return acc.String(), left[steps:], nil
}
