package modes

import (
	"fmt"
	"io"
	"os"

	"github.com/Moritisimor/EpsilonFetch/pkg/color"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/chzyer/readline"
)

func Repl(interp *interpreter.Interpreter, prompt string) {
	color.PrintGreenln("Welcome to the Gomad REPL!")
	editor, err := readline.New(prompt)
	if err != nil {
		color.PrintRedln(fmt.Sprintf("Error while setting up readline: %s\n", err.Error()))
		os.Exit(1)
	}

	for {
		input, err := editor.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}

			if err == io.EOF {
				color.PrintBlueln("Bye")	
				return
			}

			color.PrintRedln(fmt.Sprintf("Error while reading with readline: %s\n", err.Error()))
			os.Exit(1)
		}

		evaluated, err := interp.DoString(input)
		if err != nil {
			color.PrintRedln(err.Error())
			continue
		}

		color.PrintCyan("Evaluates to: ")
		color.PrintGreenln(evaluated.String())
	}
}
