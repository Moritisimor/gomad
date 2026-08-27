package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/interpreter"
)

func main() {
	interp := interpreter.NewNoStdlib()
	_, err := interp.DoString("(if true (println \"hi\") (println \"bye\"))")
	interp.PrintGlobals()
	if err != nil {
		fmt.Printf("Error while evaluating: %s\n", err.Error())
	}
}
