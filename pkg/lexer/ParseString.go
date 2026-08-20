package lexer

import (
	"strings"
	"fmt"
)

func parseString(left []byte) (string, []byte, error) {
	acc := strings.Builder{}
	wasFinished := false
	steps := 0

	for _, c := range left {
		steps++
		if c == '"' {
			wasFinished = true
			break
		}

		if c == '\n' {
			acc.WriteByte('\n')
			continue
		}

		if c == '\t' {
			acc.WriteByte('\t')
			continue
		}

		if c == '\b' {
			acc.WriteByte('\b')
		}

		acc.WriteByte(c)
	}

	if !wasFinished {
		return "", []byte{}, fmt.Errorf("String literal was never terminated! Got to: \"%s\"", acc.String())
	}

	return acc.String(), left[steps:], nil
}
