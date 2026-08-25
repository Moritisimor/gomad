package lexer

import (
	"fmt"
	"strings"
)

func parseString(left []rune) (string, []rune, error) {
	acc := strings.Builder{}
	wasFinished := false
	skipNext := false
	steps := 0

	for i, c := range left {
		steps++
		if skipNext {
			skipNext = false
			continue
		}

		if c == '"' {
			wasFinished = true
			break
		}

		if c == '\\' && len(left) >= i+1 && left[i+1] == '"' {
			acc.WriteByte('"')
			skipNext = true
			continue
		}

		if c == '\\' && len(left) >= i+1 && left[i+1] == 'n' {
			acc.WriteByte('\n')
			skipNext = true
			continue
		}

		if c == '\\' && len(left) >= i+1 && left[i+1] == 't' {
			acc.WriteByte('\t')
			skipNext = true
			continue
		}

		if c == '\\' && len(left) >= i+1 && left[i+1] == 'b' {
			acc.WriteByte('\b')
			skipNext = true
			continue
		}

		acc.WriteRune(c)
	}

	if !wasFinished {
		return "", []rune{}, fmt.Errorf("String literal was never terminated! Got to: \"%s\"", acc.String())
	}

	return acc.String(), left[steps:], nil
}
