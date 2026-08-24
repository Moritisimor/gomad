# gomad
Embeddable Nomad Lisp interpreter written in Go.

## What is this project about?
At its core, this is simply an implementation of [nomad](https://github.com/Moritisimor/nomad-lisp).

However, this implementation is more minimal and specifically made for embedding into go applications.

The difference to nomad is that gomad's API is much more suited for embedding and most of nomad's original I/O capabilities have been stripped from the prelude.

## Examples
```go
package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	interp := interpreter.New() // Instantiate new interpreter
	interp.RegisterNative("hello", func(e []expr.Expression, env *value.Env) (value.Value, error) {
		fmt.Println("Hello from Go!")
		return value.NewUnit(), nil // Unit is similar to null from other languages
		// nil because no error occurred
	})

	evaluated, err := interp.DoString("(hello)") // S-Expression syntax
	if err != nil {
		fmt.Printf("Error while evaluating: %s\n", err.Error())
	} else {
		fmt.Printf("Evaluates to: %s\n", evaluated.String())
	}
}
```
