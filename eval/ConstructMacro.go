package eval

import "github.com/Moritisimor/gomad/expr"

func ConstructMacro(
	acc, left []expr.Expression, 
	kvPairs map[string]expr.Expression,
) expr.Expression {
	if len(left) == 0 {
		return expr.List{Val: acc}
	}

	if s, ok := left[0].(expr.Symbol); ok {
		val, ok := kvPairs[s.Val]
		if !ok {
			return ConstructMacro(append(acc, s), left[1:], kvPairs) 
		}

		return ConstructMacro(append(acc, val), left[1:], kvPairs)
	}

	if l, ok := left[0].(expr.List); ok {
		return ConstructMacro(
			append(acc, ConstructMacro([]expr.Expression{}, l.Val, kvPairs)),
			left[1:],
			kvPairs,
		)
	}

	return ConstructMacro(append(acc, left[0]), left[1:], kvPairs)
}
