package lexer

import (
	"bytes"
	"fmt"
)

func isNumeric(i rune) bool {
	return i == '0' || i == '1' || i == '2' ||
		i == '3' || i == '4' || i == '5' ||
		i == '6' || i == '7' || i == '8' ||
		i == '9'
}

func isWhiteSpace(i rune) bool {
	return i == ' ' || i == '\n' || i == '\t'
}

func isTerminator(i rune) bool {
	return isWhiteSpace(i) || i == ')' || i == '('
}

func skipToNewline(left []rune) []rune {
	steps := 0
	for _, c := range left {
		if c == '\n' {
			return left[steps:]
		}

		steps++
	}

	return []rune{} // nothing after newline, nothing to return
}

func Tokenize(sourceCode string) ([]Token, error) {
	acc := []Token{}
	left := []rune(sourceCode)

	for len(left) != 0 {
		if isWhiteSpace(left[0]) {
			left = left[1:]
			continue
		}

		if left[0] == '#' {
			left = skipToNewline(left)
			continue
		}

		if len(left) > 5 && bytes.HasPrefix([]byte(string(left)), []byte{'t', 'r', 'u', 'e'}) && isTerminator(left[4]) {
			acc = append(acc, BOOLLIT{true})
			left = left[4:]
			continue
		}

		if len(left) > 6 && bytes.HasPrefix([]byte(string(left)), []byte{'f', 'a', 'l', 's', 'e'}) && isTerminator(left[5]) {
			acc = append(acc, BOOLLIT{false})
			left = left[5:]
			continue
		}

		if len(left) > 5 && bytes.HasPrefix([]byte(string(left)), []byte{'u', 'n', 'i', 't'}) && isTerminator(left[4]) {
			acc = append(acc, UNITLIT{})
			left = left[4:]
			continue
		}

		if left[0] == '(' {
			acc = append(acc, LPAREN{})
			left = left[1:]
			continue
		}

		if left[0] == ')' {
			acc = append(acc, RPAREN{})
			left = left[1:]
			continue
		}

		if left[0] == '"' {
			parsedString, after, err := parseString(left[1:])
			if err != nil {
				return acc, err
			}

			acc = append(acc, STRINGLIT{parsedString})
			left = after
			continue
		}

		if isNumeric(left[0]) ||
			(left[0] == '-' && len(left) > 1 && isNumeric(left[1])) {

			parsedNumber, after, err := parseNumber(left)
			if err != nil {
				return acc, err
			}

			acc = append(acc, NUMLIT{parsedNumber})
			left = after
			continue
		}

		symbol, after, err := parseSymbol(left)
		if err != nil {
			return acc, err
		}

		acc = append(acc, SYMBOL{symbol})
		left = after
	}

	l, r := CountParens(acc)
	if l > r {
		return acc, fmt.Errorf("Unbalanced parentheses: One or more unclosed left parentheses. left: %d, right: %d", l, r)
	}

	if r > l {
		return acc, fmt.Errorf("Unbalanced parentheses: One or more superfluous right parentheses. left: %d, right: %d", l, r)
	}

	return acc, nil
}
