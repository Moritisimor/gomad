package eval

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

func RegisterNatives(env *value.Env) {
	regs := coreNatives()
	for _, c := range CoreFormBindings() {
		env.Define(c.Name, c.Fn)
	}
	for _, r := range regs {
		env.Define(r.name, r.fn)
	}
}

type nativeReg struct {
	name string
	fn   value.NativeFunc
}

func reg(name string, fn value.NativeFunc) nativeReg {
	return nativeReg{name: name, fn: fn}
}

func coreNatives() []nativeReg {
	return []nativeReg{
		reg("throw", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("throw", 1, len(args))
			}
			v, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			if s, ok := v.(value.String); ok {
				return nil, value.EvalErr(s.Val)
			}
			return nil, value.EvalErrf("Cannot throw non-string: %s", v)
		}),

		reg("letmac", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) < 3 {
				return nil, arityErr("letmac", 3, len(args))
			}
			nameSym, ok := args[0].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("Macro name was expected to be a symbol")
			}
			paramsList, ok := args[1].(expr.List)
			if !ok {
				return nil, value.EvalErr("Macro params were expected to be a list")
			}
			params, err := symbolParams(paramsList.Val, true)
			if err != nil {
				return nil, err
			}
			if err := env.Set(nameSym.Val, value.Macro{Params: params, Expressions: args[2:], Id: value.NextLambdaID()}); err != nil {
				return nil, err
			}
			return value.Unit{}, nil
		}),

		reg("let", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("let", 2, len(args))
			}
			nameSym, ok := args[0].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("Name of let-binding was expected to be a symbol")
			}
			v, err := Eval(args[1], env)
			if err != nil {
				return nil, err
			}
			if err := env.Set(nameSym.Val, v); err != nil {
				return nil, value.EvalErr(err.Error())
			}
			return value.Unit{}, nil
		}),

		reg("letfun", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 3 {
				return nil, arityErr("letfun", 3, len(args))
			}
			nameSym, ok := args[0].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("Function name was expected to be a symbol")
			}
			paramsList, ok := args[1].(expr.List)
			if !ok {
				return nil, value.EvalErr("Function params were expected to be a list")
			}
			params, err := symbolParams(paramsList.Val, false)
			if err != nil {
				return nil, err
			}
			if err := env.Set(nameSym.Val, value.Lambda{
				Params:   params,
				Body:     args[2],
				Captured: env,
				Id:       value.NextLambdaID(),
			}); err != nil {
				return nil, err
			}
			return value.Unit{}, nil
		}),

		reg("mut", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("mut", 2, len(args))
			}
			nameSym, ok := args[0].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("First argument to mut was expected to be a symbol")
			}
			v, err := Eval(args[1], env)
			if err != nil {
				return nil, err
			}
			if err := env.Mutate(nameSym.Val, v); err != nil {
				return nil, value.EvalErr(err.Error())
			}
			return value.Unit{}, nil
		}),

		reg("lambda", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("lambda", 2, len(args))
			}
			paramsList, ok := args[0].(expr.List)
			if !ok {
				return nil, value.EvalErr("Expected parameter list after lambda")
			}
			params, err := symbolParams(paramsList.Val, false)
			if err != nil {
				return nil, err
			}
			return value.Lambda{
				Params:   params,
				Body:     args[1],
				Captured: env,
				Id:       value.NextLambdaID(),
			}, nil
		}),

		reg("record", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			rec := value.NewRecord()
			for _, field := range args {
				items, ok := field.(expr.List)
				if !ok || len(items.Val) != 2 {
					return nil, value.EvalErr("Record field has bad syntax")
				}
				nameSym, ok := items.Val[0].(expr.Symbol)
				if !ok {
					return nil, value.EvalErr("Record field has bad syntax")
				}
				v, err := Eval(items.Val[1], env)
				if err != nil {
					return nil, err
				}
				rec.SetField(nameSym.Val, v)
			}
			return rec, nil
		}),

		reg(".", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr(".", 2, len(args))
			}
			nameSym, ok := args[1].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("Field name was expected to be a symbol")
			}
			v, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			rec, ok := v.(*value.Record)
			if !ok {
				return nil, value.EvalErrf(
					"Attempt to access field of non-record expression: %s", v)
			}
			fv, ok := rec.GetField(nameSym.Val)
			if !ok {
				return nil, value.EvalErrf(
					"Attempt to access non-existant field of record: %s", nameSym.Val)
			}
			return fv, nil
		}),

		reg("record_mut", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 3 {
				return nil, arityErr("record_mut", 3, len(args))
			}
			nameSym, ok := args[1].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("Field name was expected to be a symbol")
			}
			v, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			rec, ok := v.(*value.Record)
			if !ok {
				return nil, value.EvalErrf(
					"Attempt to mutate field of non-record expression: %s", v)
			}
			if _, ok := rec.GetField(nameSym.Val); !ok {
				return nil, value.EvalErrf(
					"Cannot mutate non-existant field: %s", nameSym.Val)
			}
			nv, err := Eval(args[2], env)
			if err != nil {
				return nil, err
			}
			rec.SetField(nameSym.Val, nv)
			return value.Unit{}, nil
		}),

		reg("+", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("+", 2, len(args))
			}
			x, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := Eval(args[1], env)
			if err != nil {
				return nil, err
			}
			switch a := x.(type) {
			case value.Number:
				if b, ok := y.(value.Number); ok {
					return value.Number{Val: a.Val + b.Val}, nil
				}
			case value.String:
				if b, ok := y.(value.String); ok {
					return value.String{Val: a.Val + b.Val}, nil
				}
			}
			return nil, value.EvalErrf("Cannot add these expressions: %s and %s", x, y)
		}),
		reg("-", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("-", 2, len(args))
			}
			x, err := GetNumber(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := GetNumber(args[1], env)
			if err != nil {
				return nil, err
			}
			return value.Number{Val: x - y}, nil
		}),
		reg("*", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("*", 2, len(args))
			}
			x, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := Eval(args[1], env)
			if err != nil {
				return nil, err
			}
			switch a := x.(type) {
			case value.Number:
				if b, ok := y.(value.Number); ok {
					return value.Number{Val: a.Val * b.Val}, nil
				}
				if s, ok := y.(value.String); ok {
					return value.String{Val: mulString(s.Val, a.Val)}, nil
				}
			case value.String:
				if b, ok := y.(value.Number); ok {
					return value.String{Val: mulString(a.Val, b.Val)}, nil
				}
			}
			return nil, value.EvalErrf("Cannot multiply these expressions: %s and %s", x, y)
		}),
		reg("/", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("/", 2, len(args))
			}
			x, err := GetNumber(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := GetNumber(args[1], env)
			if err != nil {
				return nil, err
			}
			if y == 0 {
				return nil, value.EvalErr("Attempt to divide by 0")
			}
			return value.Number{Val: x / y}, nil
		}),
		reg("mod", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("mod", 2, len(args))
			}
			x, err := GetNumber(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := GetNumber(args[1], env)
			if err != nil {
				return nil, err
			}
			return value.Number{Val: math.Mod(x, y)}, nil
		}),

		reg("=", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("=", 2, len(args))
			}
			x, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := Eval(args[1], env)
			if err != nil {
				return nil, err
			}
			return value.Boolean{Val: value.Equal(x, y)}, nil
		}),
		reg(">", cmpNative(">", func(a, b float64) bool { return a > b })),
		reg(">=", cmpNative(">=", func(a, b float64) bool { return a >= b })),
		reg("<", cmpNative("<", func(a, b float64) bool { return a < b })),
		reg("<=", cmpNative("<=", func(a, b float64) bool { return a <= b })),
		reg("or", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("or", 2, len(args))
			}
			a, err := GetBool(args[0], env)
			if err != nil {
				return nil, err
			}
			if a {
				return value.Boolean{Val: true}, nil
			}
			b, err := GetBool(args[1], env)
			if err != nil {
				return nil, err
			}
			return value.Boolean{Val: b}, nil
		}),
		reg("and", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("and", 2, len(args))
			}
			a, err := GetBool(args[0], env)
			if err != nil {
				return nil, err
			}
			if !a {
				return value.Boolean{Val: false}, nil
			}
			b, err := GetBool(args[1], env)
			if err != nil {
				return nil, err
			}
			return value.Boolean{Val: b}, nil
		}),

		reg("list", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			values := make([]value.Value, len(args))
			for i, p := range args {
				v, err := Eval(p, env)
				if err != nil {
					return nil, err
				}
				values[i] = v
			}
			return value.FromVec(values), nil
		}),
		reg("append", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("append", 2, len(args))
			}
			x, err := GetList(args[0], env)
			if err != nil {
				return nil, err
			}
			y, err := GetList(args[1], env)
			if err != nil {
				return nil, err
			}
			return x.Append(y), nil
		}),
		reg("car", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("car", 1, len(args))
			}
			l, err := GetList(args[0], env)
			if err != nil {
				return nil, err
			}
			if l.IsNil() {
				return value.Unit{}, nil
			}
			return l.Head(), nil
		}),
		reg("cdr", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("cdr", 1, len(args))
			}
			l, err := GetList(args[0], env)
			if err != nil {
				return nil, err
			}
			if l.IsNil() {
				return value.Unit{}, nil
			}
			return l.Tail(), nil
		}),
		reg("cons", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("cons", 2, len(args))
			}
			t, err := GetList(args[1], env)
			if err != nil {
				return nil, err
			}
			h, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return value.Cons(h, t), nil
		}),

		reg("sprint", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			var b strings.Builder
			for _, p := range args {
				v, err := Eval(p, env)
				if err != nil {
					return nil, err
				}
				b.WriteString(v.String())
			}
			return value.String{Val: b.String()}, nil
		}),
		reg("print", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			var b strings.Builder
			for _, p := range args {
				v, err := Eval(p, env)
				if err != nil {
					return nil, err
				}
				b.WriteString(v.String())
			}
			fmt.Print(b.String())
			return value.Unit{}, nil
		}),
		reg("println", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			var b strings.Builder
			for _, p := range args {
				v, err := Eval(p, env)
				if err != nil {
					return nil, err
				}
				b.WriteString(v.String())
			}
			fmt.Println(b.String())
			return value.Unit{}, nil
		}),
		reg("readln", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("readln", 1, len(args))
			}
			prompt, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			fmt.Print(prompt)
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return nil, value.EvalErrf("readln: %v", err)
			}
			return value.String{Val: strings.TrimRight(line, "\r\n")}, nil
		}),
		reg("chars", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("chars", 1, len(args))
			}
			s, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			runes := []rune(s)
			vals := make([]value.Value, len(runes))
			for i, c := range runes {
				vals[i] = value.String{Val: string(c)}
			}
			return value.FromVec(vals), nil
		}),
		reg("lower", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("lower", 1, len(args))
			}
			s, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			return value.String{Val: lowerASCII(s)}, nil
		}),
		reg("trim", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("trim", 1, len(args))
			}
			v, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			s, ok := v.(value.String)
			if !ok {
				return nil, value.EvalErrf("Cannot apply trim-operation on non-string expression: %s", v)
			}
			return value.String{Val: strings.Trim(s.Val, " \t\n\r\v\f")}, nil
		}),
		reg("splitws", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("splitws", 1, len(args))
			}
			s, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			fields := splitNomadWhitespace(s)
			vals := make([]value.Value, len(fields))
			for i, f := range fields {
				vals[i] = value.String{Val: f}
			}
			return value.FromVec(vals), nil
		}),
		reg("to_string", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("to_string", 1, len(args))
			}
			v, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return value.String{Val: v.String()}, nil
		}),
		reg("string_to_num", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("string_to_num", 1, len(args))
			}
			s, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			n, err := parseNomadFloat(s)
			if err != nil {
				return nil, value.EvalErrf("Cannot parse this string to a number: %s", s)
			}
			return value.Number{Val: n}, nil
		}),

		reg("isunit", predicate("isunit", func(v value.Value) bool {
			_, ok := v.(value.Unit)
			return ok
		})),
		reg("isstr", predicate("isstring", func(v value.Value) bool {
			_, ok := v.(value.String)
			return ok
		})),
		reg("isnum", predicate("isnum", func(v value.Value) bool {
			_, ok := v.(value.Number)
			return ok
		})),
		reg("islist", predicate("islist", func(v value.Value) bool {
			_, ok := v.(value.List)
			return ok
		})),
		reg("isfun", predicate("islambda", func(v value.Value) bool {
			_, ok := v.(value.Lambda)
			return ok
		})),
		reg("isnative", predicate("isnative", func(v value.Value) bool {
			switch v.(type) {
			case value.NativeFunc, value.NativeFunction:
				return true
			}
			return false
		})),
		reg("ismac", predicate("isnative", func(v value.Value) bool {
			_, ok := v.(value.Macro)
			return ok
		})),
		reg("isbool", predicate("isbool", func(v value.Value) bool {
			_, ok := v.(value.Boolean)
			return ok
		})),
		reg("isrecord", predicate("isrecord", func(v value.Value) bool {
			_, ok := v.(*value.Record)
			return ok
		})),
	}
}

func cmpNative(name string, f func(a, b float64) bool) value.NativeFunc {
	return func(args []expr.Expr, env *value.Env) (value.Value, error) {
		if len(args) != 2 {
			return nil, arityErr(name, 2, len(args))
		}
		x, err := GetNumber(args[0], env)
		if err != nil {
			return nil, err
		}
		y, err := GetNumber(args[1], env)
		if err != nil {
			return nil, err
		}
		return value.Boolean{Val: f(x, y)}, nil
	}
}

func predicate(name string, f func(value.Value) bool) value.NativeFunc {
	return func(args []expr.Expr, env *value.Env) (value.Value, error) {
		if len(args) != 1 {
			return nil, arityErr(name, 1, len(args))
		}
		v, err := Eval(args[0], env)
		if err != nil {
			if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrEval {
				return value.Boolean{Val: false}, nil
			}
			return nil, err
		}
		return value.Boolean{Val: f(v)}, nil
	}
}

func symbolParams(params []expr.Expr, lowercaseErr bool) ([]string, error) {
	out := make([]string, 0, len(params))
	for _, p := range params {
		s, ok := p.(expr.Symbol)
		if !ok {
			if lowercaseErr {
				return nil, value.EvalErr("Non-symbol in parameter list")
			}
			return nil, value.EvalErr("Non-Symbol in parameter list")
		}
		out = append(out, s.Val)
	}
	return out, nil
}

func lowerASCII(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return r
	}, s)
}

func splitNomadWhitespace(s string) []string {
	parts := strings.Split(s, " ")
	out := parts[:0]
	for _, p := range parts {
		if strings.Trim(p, " \t\n\r\v\f") != "" {
			out = append(out, p)
		}
	}
	return out
}

func mulString(s string, factor float64) string {
	n := int64(factor)
	if math.IsNaN(factor) {
		n = 0
	} else if math.IsInf(factor, 0) {
		n = math.MaxInt64
	}
	if n < 1 {
		return ""
	}
	if n > 1000000 {
		n = 1000000
	}
	return strings.Repeat(s, int(n))
}

func parseNomadFloat(s string) (float64, error) {
	clean := strings.ReplaceAll(s, "_", "")
	if f, err := strconv.ParseFloat(clean, 64); err == nil {
		return f, nil
	}
	return parseHexFloat(clean)
}

func parseHexFloat(s string) (float64, error) {
	body := s
	negative := false
	if strings.HasPrefix(body, "+") || strings.HasPrefix(body, "-") {
		negative = strings.HasPrefix(body, "-")
		body = body[1:]
	}
	if !strings.HasPrefix(body, "0x") && !strings.HasPrefix(body, "0X") {
		return 0, fmt.Errorf("not a hex float")
	}
	body = body[2:]

	var mantissaStr string
	var exp float64
	var err error
	if lo := strings.IndexAny(body, "pP"); lo >= 0 {
		mantissaStr = body[:lo]
		var e int
		e, err = strconv.Atoi(body[lo+1:])
		if err != nil {
			return 0, err
		}
		exp = float64(e)
	} else {
		mantissaStr = body
	}
	if strings.IndexByte(mantissaStr, '.') < 0 {
		mantissaStr += ".0"
	}
	if lo := strings.IndexByte(mantissaStr, '.'); lo >= 0 {
		intPart, fracPart := mantissaStr[:lo], mantissaStr[lo+1:]
		var mantissa float64
		for _, c := range intPart {
			d, ok := hexDigit(c)
			if !ok {
				return 0, fmt.Errorf("bad hex digit")
			}
			mantissa = mantissa*16 + d
		}
		scale := 1.0 / 16.0
		for _, c := range fracPart {
			d, ok := hexDigit(c)
			if !ok {
				return 0, fmt.Errorf("bad hex digit")
			}
			mantissa += d * scale
			scale /= 16
		}
		v := mantissa
		if exp < 0 {
			v /= math.Pow(2, -exp)
		} else {
			v *= math.Pow(2, exp)
		}
		if negative {
			v = -v
		}
		return v, nil
	}
	return 0, fmt.Errorf("malformed hex float")
}

func hexDigit(c rune) (float64, bool) {
	switch {
	case c >= '0' && c <= '9':
		return float64(c - '0'), true
	case c >= 'a' && c <= 'f':
		return float64(c-'a') + 10, true
	case c >= 'A' && c <= 'F':
		return float64(c-'A') + 10, true
	}
	return 0, false
}
