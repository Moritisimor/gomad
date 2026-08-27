# gomad

<img src="./assets/gomad.svg" alt="gomad logo" width="290">

**An embeddable [Nomad Lisp](https://github.com/Moritisimor/nomad-lisp) interpreter for Go.**

## What is this project about?

At its core, gomad is simply an implementation of Nomad Lisp: a small,
dynamically typed Lisp with functions, macros, lexical scopes, records, and a
compact standard library.

The main reason this version exists is embedding. If you want to let users
write a little configuration language, add scripting to a Go program, or just
experiment with a Lisp without leaving the Go ecosystem, gomad is meant to be
easy to drop in. An interpreter owns its environment, Go functions can be
registered as native Nomad functions, and every failure comes back as an
ordinary Go `error`.

Gomad can also be used as a standalone interpreter. Its command-line interface
includes a REPL, expression evaluation, and script files, so the same code can
be tried interactively before it is embedded in an application.

This project follows the behavior of the original
[nomad-lisp](https://github.com/Moritisimor/nomad-lisp) closely while fixing a
handful of rough edges around parsing, Unicode, error handling, and deep
recursion. It is also kept in step with
[romad](https://github.com/robertflexx/romad), the Rust implementation. If you
are interested in a BEAM version, take a look at
[bomad](https://github.com/RobertFlexx/bomad).

## Getting started

Clone the repository and build the command:

```bash
go build ./cmd/gomad
```

You can then start the REPL or evaluate some Nomad code directly:

```bash
./gomad
./gomad -e '(+ (* 10 5) (- 1000 250))'
./gomad path/to/my_script.nomad
```

The complete command-line interface is:

```text
gomad                       Start the REPL
gomad --repl                Start the REPL explicitly
gomad -e 'EXPRESSION'       Evaluate an expression and print its value
gomad FILE.nomad            Run a script
gomad --help                Show the help screen
```

Arguments passed to a script are available inside Nomad through the `args`
list.

## Embedding gomad

Add gomad to your project:

```bash
go get github.com/Moritisimor/gomad
```

Creating an interpreter and evaluating source code only takes a few lines:

```go
package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/interpreter"
)

func main() {
	interp := interpreter.New()

	result, err := interp.DoString("(+ 20 22)")
	if err != nil {
		panic(err)
	}

	fmt.Println(result) // 42
}
```

A new interpreter includes the native functions and the Nomad standard
library. It keeps its own global environment, so creating several independent
interpreters in one process is fine.

### Registering a Go function

Native functions receive the original, unevaluated Nomad expressions and the
calling environment. This makes ordinary functions possible, but also leaves
room for special forms that decide which arguments to evaluate.

```go
package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	interp := interpreter.New()

	err := interp.RegisterNative("square", func(
		args []expr.Expr,
		env *value.Env,
	) (value.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("square expects one argument")
		}

		n, err := eval.GetNumber(args[0], env)
		if err != nil {
			return nil, err
		}
		return value.NewNumber(n * n), nil
	})
	if err != nil {
		panic(err)
	}

	result, err := interp.DoString("(square 9)")
	if err != nil {
		panic(err)
	}
	fmt.Println(result) // 81
}
```

The embedding API also provides `DoFile`, `EvalExpr`, `GetGlobal`, `Set`,
`Unset`, and `Parse`, along with direct access to the interpreter environment.
Use `NewNoStdlib` when only the native runtime is wanted. Older embedding names
such as `expr.Expression`, `value.NewUnit`, and `GetBinding` remain available
for existing programs.

## Language and runtime

Gomad supports the full core expected from Nomad:

- numbers, strings, booleans, units, lists, and mutable records;
- lexical scopes, closures, named functions, mutation, and currying;
- macros whose arguments are expanded in the caller's scope;
- `if`, `switch`, `do`, `scoped`, `try`, and catchable evaluation errors;
- list helpers such as `map`, `filter`, `foldl`, `range`, and `list_init`;
- string, console, filesystem, process, and environment functions;
- decimal and OCaml-style hexadecimal number literals;
- Unicode-aware character handling and CRLF-compatible source files.

Lists are persistent cons lists rather than ordinary Go slices. `cons`, `car`,
and `cdr` are constant-time operations, and different lists can safely share
the same tail. This matters for Lisp code that builds lists recursively and
avoids repeatedly copying values.

Tail calls are evaluated through a trampoline. Lambda calls, macro expansion,
and the tail positions of `if`, `do`, `switch`, `scoped`, and `try` do not keep
growing the Go stack. A properly tail-recursive Nomad function can therefore
run for hundreds of thousands of iterations without a stack overflow.

Evaluation stops at the first error. Calls to `exit` and `bye` are returned to
an embedding host as controlled exit errors rather than terminating the whole
Go process. The command-line program is the layer that turns those signals
into actual process exit codes.

## How does it compare to romad?

Gomad and romad understand the same language and provide nearly the same
runtime features. Scripts using the standard library, records, macros, tail
calls, filesystem functions, or command-line arguments should behave the same
under both interpreters. They also intentionally preserve a few odd diagnostic
messages from the original Nomad implementation because some programs rely on
them.

The important difference is on the host side. Gomad exposes Go interfaces,
pointers, callbacks, and `error` values. Romad exposes Rust enums, `Rc` values,
and its `NomadError` type. Gomad is generally more convenient inside a Go
application, while romad's tagged values and Rust runtime use noticeably less
memory.

For some rough context, these measurements were taken on an Intel i7-9700 with
optimized binaries. They are snapshots rather than promises, but they show the
general trade-off fairly well:

| Workload | gomad | romad |
|---|---:|---:|
| 1,000,000-step tail-recursive loop | ~0.49 s | ~0.52 s |
| Mapping over 10,000 values | ~0.03 s | ~0.03 s |
| 100 fresh command invocations | ~0.13 s | ~0.09 s |
| Peak memory during the tail loop | ~10.8 MB | ~2.8 MB |
| Stripped command size | ~2.4 MB | ~1.1 MB |

In other words, evaluator speed is very close. Romad currently starts a little
faster and wins clearly on memory and binary size. Gomad can reuse its parsed
standard-library syntax tree when an application creates multiple interpreters
inside one process, which makes that embedding case cheaper after the first
interpreter.

## Testing

Run the full test suite and static checks with:

```bash
go test ./...
go vet ./...
```

There are also benchmarks for interpreter startup, tail recursion, and list
processing:

```bash
go test ./interpreter -run '^$' -bench . -benchmem
```

The conformance tests cover the tokenizer, numeric syntax, closures, macros,
records, lists, Unicode behavior, exact errors, exit handling, and deep tail
calls. Several of romad's example programs are also useful as side-by-side
smoke tests.

Gomad is available under the MIT license.
