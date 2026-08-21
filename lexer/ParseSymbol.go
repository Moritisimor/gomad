package lexer

import (
	"fmt"
	"strings"
)

func ParseSymbol(left []byte) (string, []byte, error) {
	acc := strings.Builder{}
	wasEnded := false
	steps := 0

	for _, c := range left {
		if isTerminator(c) {
			wasEnded = true
			break
		}

		acc.WriteByte(c)
		steps++
	}

	if !wasEnded {
		return "", []byte{}, fmt.Errorf("Symbol was never terminated! Got to \"%s\"", acc.String())
	}

	return acc.String(), left[steps:], nil
}
