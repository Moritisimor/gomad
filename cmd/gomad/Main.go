package main

import (
	"os"
	"fmt"

	"github.com/Moritisimor/gomad/internal/modes"
	"github.com/Moritisimor/gomad/interpreter"
)

func main() {
	interp := interpreter.New()
	if len(os.Args) == 1 {
		modes.Repl(interp)
		return
	}

	_, err := interp.DoFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
