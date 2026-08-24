package eval

import (
	"github.com/Moritisimor/gomad/internal/helpers"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/value"
)

func DoString(sourceCode string, env *value.Env) (value.Value, error) {
	tokens, err := lexer.Tokenize(sourceCode)
	if err != nil {
		return helpers.Err("Error while tokenizing: %s", err.Error())
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return helpers.Err("Error while parsing: %s", err.Error())
	}

	var lastExpr value.Value
	lastExpr = value.NewUnit()
	for _, node := range ast {
		evaluated, err := Eval(node, env)
		if err != nil {
			return helpers.Err("Uncaught Error:\n\t%s", err.Error())
		}

		lastExpr = evaluated
	}

	return lastExpr, nil
}
