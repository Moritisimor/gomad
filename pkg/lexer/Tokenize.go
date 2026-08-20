package lexer

import "bytes"

func isNumeric(i byte) bool {
	return i == '0' || i == '1' || i == '2' ||
		i == '3' || i == '4' || i == '5' ||
		i == '6' || i == '7' || i == '8' ||
		i == '9' || i == '-' // for negative numbers
}

func isWhiteSpace(i byte) bool {
	return i == ' ' || i == '\n' || i == '\t'	
}

func isTerminator(i byte) bool {
	return isWhiteSpace(i) || i == ')' || i == '('
}

func skipToNewline(left []byte) []byte {
	steps := 0
	for _, c := range left {
		if c == '\n' {
			return left[steps:]
		}

		steps++
	}

	return []byte{} // nothing after newline, nothing to return
}

func Tokenize(source_code string) ([]Token, error) {
	acc := []Token{}
	left := []byte(source_code)

	for len(left) != 0 {
		if isWhiteSpace(left[0]) {
			left = left[1:]
			continue
		}

		if left[0] == '#' {
			left = skipToNewline(left)
			continue
		}

		if bytes.HasPrefix(left, []byte{ 't', 'r', 'u', 'e' }) {
			acc = append(acc, BOOLLIT{ true })
			left = left[4:]
			continue
		}

		if bytes.HasPrefix(left, []byte{ 'f', 'a', 'l', 's', 'e' }) {
			acc = append(acc, BOOLLIT{ false })
			left = left[5:]
			continue
		}

		if bytes.HasPrefix(left, []byte{ 'u', 'n', 'i', 't' }) {
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

			acc = append(acc, STRINGLIT{ parsedString })
			left = after
			continue
		}

		if isNumeric(left[0]) {
			parsedNumber, after, err := parseNumber(left)
			if err != nil {
				return acc, err
			}

			acc = append(acc, NUMLIT{ parsedNumber })
			left = after
			continue
		}

		symbol, after, err := ParseSymbol(left)
		if err != nil {
			return acc, err
		}

		acc = append(acc, SYMBOL{ symbol })
		left = after
	}

	return acc, nil
}
