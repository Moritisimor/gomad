package helpers

import (
	"fmt"

	"github.com/Moritisimor/gomad/value"
)

func Err(format string, args... any) (value.Unit, error) {
	return value.Unit{}, fmt.Errorf(format, args...)
}

func WrongArgs(name string, expected, actual int) (value.Unit, error) {
	return value.Unit{}, fmt.Errorf(
		"Function '%s' was given the wrong amount of args. Expected: %d, got: %d", name, expected, actual,
	)
}
