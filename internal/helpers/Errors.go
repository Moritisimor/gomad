package helpers

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

func Err(format string, args ...any) (value.Unit, error) {
	return value.Unit{}, fmt.Errorf(format, args...)
}

func BadType[T any](expectedType string, e expr.Expression) (T, error) {
	var t T
	return t, fmt.Errorf("This expression was expected to evaluate to a %s, but it didn't: %s", expectedType, expr.SprintExpr(e))
}

func WrongArgs(name string, expected, actual int) (value.Unit, error) {
	return value.Unit{}, fmt.Errorf(
		"Function '%s' was given the wrong amount of args. Expected: %d, got: %d", name, expected, actual,
	)
}
