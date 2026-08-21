package lexer

func CountParens(tokens []Token) (int, int) {
	right := 0
	left := 0
	for _, t := range tokens {
		if _, ok := t.(RPAREN); ok {
			right++
		}

		if _, ok := t.(LPAREN); ok {
			left++
		}
	}

	return left, right
}
