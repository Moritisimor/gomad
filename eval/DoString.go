package eval

import (
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/value"
)

func DoString(sourceCode string, env *value.Env) (value.Value, error) {
	tokens, err := lexer.Tokenize(sourceCode)
	if err != nil {
		return nil, &value.Error{Kind: value.ErrTokenize, Msg: err.Error()}
	}

	ast, err := parser.ParseProgram(tokens)
	if err != nil {
		return nil, &value.Error{Kind: value.ErrParse, Msg: err.Error()}
	}

	return EvalSeq(ast, env)
}

func DoStringSeq(sourceCode string, env *value.Env) (value.Value, error) {
	return DoString(sourceCode, env)
}
