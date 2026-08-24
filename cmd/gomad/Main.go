package main

import (
	"fmt"
	"os"

	"github.com/Moritisimor/EpsilonFetch/pkg/color"
	"github.com/Moritisimor/gomad/internal/modes"
	"github.com/Moritisimor/gomad/interpreter"
)

func main() {
	interp, err := interpreter.New()
	if err != nil {
		fmt.Printf("Error while creating gomad interpreter: %s\n", err.Error())
		return
	}

	if len(os.Args) == 1 {
		modes.Repl(interp, color.SprintMagenta("Gomad λ "))
		return
	}

	_, err = interp.DoFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
