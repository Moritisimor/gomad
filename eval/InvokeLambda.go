package eval

import (
	"fmt"

	"github.com/Moritisimor/gomad/value"
)

func InvokeLambda(lambda value.Lambda, args... value.Value) (value.Value, error) {
	expected, actual := len(lambda.Params), len(args)
	if expected != actual {
		return value.NewUnit(), fmt.Errorf(
			"Lambda invoked with wrong amount of arguments. Expected: %d, got: %d", 
			expected, actual,
		)
	}

	localEnv := lambda.Captured
	for idx, arg := range args {
		localEnv.SetBinding(lambda.Params[idx], arg)
	}

	return Eval(lambda.Body, localEnv)
}
