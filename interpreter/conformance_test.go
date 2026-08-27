package interpreter_test

import (
	"math"
	"testing"

	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func eval(t *testing.T, source string) value.Value {
	t.Helper()
	v, err := interpreter.New().DoString(source)
	if err != nil {
		t.Fatalf("%q failed: %v", source, err)
	}
	return v
}

func TestLanguageConformance(t *testing.T) {
	tests := []struct{ name, source, want string }{
		{"arithmetic", "(+ (* 10 5) (- 1000 250))", "800"},
		{"negative", "-5", "-5"},
		{"float", "(mod 7.5 2)", "1.50"},
		{"hex", "(- 0xA.8p-2 0)", "2.62"},
		{"leading_hex_fraction", "0x.8", "0.50"},
		{"underscores", "(+ 1_000 2)", "1002"},
		{"keyword_boundary", "(let truest 55) truest", "55"},
		{"closures", "(let add (lambda (x) (lambda (y) (+ x y)))) (let add10 (add 10)) (add10 20)", "30"},
		{"recursion", "(letfun fact (n) (switch n (0 1) (_ (* n (fact (dec n)))))) (fact 10)", "3628800"},
		{"macro", "(letmac m (a) list a a a) (m (+ 1 1))", "(2 2 2)"},
		{"scoped", "(scoped ((x 10) (y 20)) (+ x y))", "30"},
		{"try", "(try (/ 1 0) 42)", "42"},
		{"unit_equality", "(= unit unit)", "true"},
		{"records", "(let p (record (x 1))) (record_mut p x 99) (. p x)", "99"},
		{"record_type", "(typeof (record (a 1)))", "record"},
		{"unicode_chars", "(chars \"hé\")", "(h é)"},
		{"ascii_lower", "(lower \"ÉÀAZ\")", "ÉÀaz"},
		{"map", "(map (lambda (x) (* x x)) (list 1 2 3))", "(1 4 9)"},
		{"filter", "(filter (lambda (x) (< x 3)) (list 1 2 3 4))", "(1 2)"},
		{"range", "(range 1 3 (list 0 1 2 3 4))", "(1 2 3)"},
		{"large_integer_format", "(to_string 1e18)", "1000000000000000000"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := eval(t, tc.source).String(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestExactErrors(t *testing.T) {
	tests := []struct{ source, want string }{
		{"(+ 1 2", "Unbalanced parantheses: one or more unclosed left parantheses"},
		{"+ 1 2)", "Unbalanced parantheses: one or more superfluous right parantheses"},
		{"\"abc", "String literal was never ended (Got \"abc)"},
		{"nope", "No such variable: nope"},
		{"(let x 1) (let x 2)", "Cannot bind x: Already exists in this scope"},
		{"(mut nope 1)", "Cannot mutate non-existant binding: nope"},
		{"((lambda (x x) x) 1 2)", "Cannot bind x: Already exists in this scope"},
		{"(letmac m (1) 2)", "Non-symbol in parameter list"},
		{"(isstr 1 2)", "Native function isstring was given bad syntax. Perhaps it was given the wrong amount of args? Args expected: 1. Got: 2"},
	}
	for _, tc := range tests {
		t.Run(tc.source, func(t *testing.T) {
			_, err := interpreter.New().DoString(tc.source)
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestTailCallsUseConstantStack(t *testing.T) {
	v := eval(t, "(let loop (lambda (n acc) (if (= n 0) acc (loop (- n 1) (+ acc 1))))) (loop 250000 0)")
	n, ok := v.(value.Number)
	if !ok || n.Val != 250000 {
		t.Fatalf("got %v", v)
	}
}

func TestExitIsNotCaught(t *testing.T) {
	_, err := interpreter.New().DoString("(try (exit 3) 42)")
	got, ok := err.(*value.Error)
	if !ok || got.Kind != value.ErrExit || got.Code != 3 {
		t.Fatalf("got %v", err)
	}
}

func TestNumberSpecialValues(t *testing.T) {
	v := eval(t, "(mod 0 0)")
	n, ok := v.(value.Number)
	if !ok || !math.IsNaN(n.Val) {
		t.Fatalf("got %v", v)
	}
}

func BenchmarkInterpreterStartup(b *testing.B) {
	for range b.N {
		_ = interpreter.New()
	}
}

func BenchmarkTailRecursion100k(b *testing.B) {
	for range b.N {
		i := interpreter.New()
		if _, err := i.DoString("(let loop (lambda (n acc) (if (= n 0) acc (loop (- n 1) (+ acc 1))))) (loop 100000 0)"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap1000(b *testing.B) {
	for range b.N {
		i := interpreter.New()
		if _, err := i.DoString("(map (lambda (x) (+ x 1)) (list_init 1000 (lambda (x) x)))"); err != nil {
			b.Fatal(err)
		}
	}
}
