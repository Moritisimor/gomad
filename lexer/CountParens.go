package lexer

// CountParens counts the parentheses, returning the amount of left and the amount of right parentheses.
// It's used by the tokenizer to detect unbalanced parentheses.
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
