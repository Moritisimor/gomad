package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
	"github.com/chzyer/readline"
)

func main() {
	gomad := interpreter.New()
	editor, err := readline.New("Gomad >> ")
	if err != nil {
		fmt.Printf("Error while setting up readline: %s\n", err.Error())
		os.Exit(1)
	}

	gomad.RegisterNative("hello", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		fmt.Println("Hello from Go!")
		return value.NewUnit(), nil
	})

	for {
		input, err := editor.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}

			if err == io.EOF {
				fmt.Println("Bye")
				return
			}

			fmt.Printf("Error while reading with readline: %s\n", err.Error())
			os.Exit(1)
		}

		evaluated, err := gomad.DoString(input)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Evaluates to: %s\n", evaluated.String())
	}
}
