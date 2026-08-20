package lexer

import (
	"strings"
	"strconv"
	"fmt"
)

func parseNumber(left []byte) (float64, []byte, error) {
	acc := strings.Builder{}
	wasFinished := false
	steps := 0

	for _, c := range left {
		if isTerminator(c) {
			wasFinished = true
			break
		}

		acc.WriteByte(c)
		steps++
	}

	num, err := strconv.ParseFloat(acc.String(), 64)
	if err != nil {
		return 0, []byte{}, fmt.Errorf("Error while parsing number literal: %s", err.Error())
	}

	if !wasFinished {
		return 0, []byte{}, fmt.Errorf("Number literal was never terminated! Got to: \"%s\"", acc.String())
	}

	return num, left[steps:], nil
}
