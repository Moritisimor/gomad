package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Moritisimor/gomad/internal/modes"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

const HELP = `gomad - A fast, tail-call-optimized Lisp interpreter.

Usage:
  gomad                      Run the interactive REPL.
  gomad -e|--eval EXPR       Evaluate EXPR and print the result.
  gomad -r|--repl            Run the interactive REPL.
  gomad FILE.nomad          Run the given script file.
  gomad -h|--help            Show this help.
`

func main() {
	debug.SetMaxStack(1 * 1024 * 1024 * 1024)

	argv := os.Args[1:]

	switch {
	case len(argv) == 1 && (argv[0] == "-h" || argv[0] == "--help"):
		fmt.Print(HELP)

	case len(argv) == 2 && (argv[0] == "-e" || argv[0] == "--eval"):
		interp := interpreter.New(os.Args[1:]...)
		v, err := interp.DoString(argv[1])
		exitOrPrint(v, err)

	case len(argv) == 0 || (len(argv) == 1 && (argv[0] == "-r" || argv[0] == "--repl")):
		interp := interpreter.New(os.Args[1:]...)
		modes.Repl(interp, "Gomad λ ")

	default:
		interp := interpreter.New(os.Args[1:]...)
		v, err := interp.DoFile(argv[0])
		_ = v
		if err != nil {
			if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrExit {
				os.Exit(ve.Code)
			}
			fmt.Println(reportError(err))
			os.Exit(1)
		}
	}
}

func exitOrPrint(v value.Value, err error) {
	if err != nil {
		if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrExit {
			os.Exit(ve.Code)
		}
		fmt.Println(reportError(err))
		os.Exit(1)
	}
	fmt.Println(v)
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
