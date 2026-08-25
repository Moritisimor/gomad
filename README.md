# gomad

<p align="left">
  <img src="./assets/gomad.svg" alt="gomad logo" width="360">
</p>

                                              Embeddable Nomad Lisp interpreter written in Go.

## What is this project about?
At its core, this is simply an implementation of [nomad](https://github.com/Moritisimor/nomad-lisp).

However, this implementation is more minimal and specifically made for embedding into go applications.

The difference to nomad is that gomad's API is much more suited for embedding and most of nomad's original I/O capabilities have been stripped from the prelude.

Also, be sure to check out [RobertFlexx's](https://github.com/robertflexx) nomad implementation written in rust, [romad](https://github.com/robertflexx/romad).

And if you want a BEAM port of nomad, be sure to check out **RobertFlexx's** [bomad](https://github.com/RobertFlexx/bomad)

## Getting started
Import the repository:
```go
import "github.com/Moritisimor/gomad"
```

Then update your dependencies:
```bash
go mod tidy
```

And now you can use the `gomad` package within your project!

## Examples
### Instantiating the interpreter, registering natives and calling them through gomad
```go
package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	interp := interpreter.New() // Create new interpreter
	interp.RegisterNative("hello", func(
        e []expr.Expression, 
        env *value.Env,
    ) (value.Value, error) {
		fmt.Println("Hello from Go!")
		return value.NewUnit(), nil // Unit is similar to null 
		// nil because no error occurred
	})

	evaluated, err := interp.DoString("(hello)") // S-Expressions
	if err != nil {
		fmt.Printf("Error while evaluating: %s\n", err.Error())
	} else {
		fmt.Printf("Evaluates to: %s\n", evaluated.String())
	}
}
```

### A simple logging function
```go
package main

import (
	"fmt"
	"log"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
	"github.com/Moritisimor/gomad/eval"
)

func main() {
	interp := interpreter.New() // Create new interpreter
	interp.RegisterNative("log", func(
        e []expr.Expression, 
        env *value.Env,
    ) (value.Value, error) {
		if len(e) != 1 { // Check argument length
			return value.NewUnit(), fmt.Errorf("Error in call to log: Expected one argument, got %d", len(e))
		}

		// We expect a string as the argument, anything else is an error.
		logString, err := eval.GetString(e[0], env)
		if err != nil {
			return value.NewUnit(), fmt.Errorf("Error in call to log: %s", err.Error())
		}

		log.Println(logString)
		return value.NewUnit(), nil
	})

	interp.DoString("(log \"Gomad interpreter running!\")")
}
```
