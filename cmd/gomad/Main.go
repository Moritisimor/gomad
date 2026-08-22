package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	gomad := interpreter.New()
	gomad.RegisterNative("hello", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		fmt.Println("Hello from Go!")
		return value.NewUnit(), nil
	})

	evaluated, err := gomad.DoString("(hello)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Evaluates to: %s\n", evaluated.String())
}
