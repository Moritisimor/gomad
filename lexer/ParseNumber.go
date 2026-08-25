package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

func parseNumber(left []rune) (float64, []rune, error) {
	acc := strings.Builder{}
	wasFinished := false
	steps := 0

	for _, c := range left {
		if isTerminator(c) {
			wasFinished = true
			break
		}

		acc.WriteRune(c)
		steps++
	}

	num, err := strconv.ParseFloat(acc.String(), 64)
	if err != nil {
		return 0, []rune{}, fmt.Errorf("Error while parsing number literal: %s", err.Error())
	}

	if !wasFinished {
		return 0, []rune{}, fmt.Errorf("Number literal was never terminated! Got to: \"%s\"", acc.String())
	}

	return num, left[steps:], nil
}
