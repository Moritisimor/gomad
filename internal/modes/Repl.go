package modes

import (
	"fmt"
	"io"
	"os"

	"github.com/Moritisimor/EpsilonFetch/pkg/color"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
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
			if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrExit {
				os.Exit(ve.Code)
				return
			}
			color.PrintRedln(reportError(err))
			continue
		}

		color.PrintCyan("Evaluates to: ")
		color.PrintGreenln(evaluated.String())
	}
}

func reportError(err error) string {
	if ve, ok := err.(*value.Error); ok {
		switch ve.Kind {
		case value.ErrParse:
			return "Error while parsing: " + ve.Msg
		case value.ErrTokenize:
			return "Error while tokenizing: " + ve.Msg
		case value.ErrIo:
			return "Error while reading file: " + ve.Msg
		default:
			return "Error while evaluating: " + ve.Msg
		}
	}
	return "Error while evaluating: " + err.Error()
}
