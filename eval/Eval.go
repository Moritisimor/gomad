package eval

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/value"
)

type handler struct {
	cursor expr.Expr
	scope  *value.Env
}

type evaluator struct {
	cursor   expr.Expr
	scope    *value.Env
	handlers []handler
}

func (ev *evaluator) bail(err error) (restart bool, out error) {
	if ve, ok := err.(*value.Error); ok && ve.Kind == value.ErrEval {
		if n := len(ev.handlers); n > 0 {
			h := ev.handlers[n-1]
			ev.handlers = ev.handlers[:n-1]
			ev.cursor = h.cursor
			ev.scope = h.scope
			return true, nil
		}
	}
	return false, err
}

func (ev *evaluator) step(e expr.Expr) (v value.Value, restart bool, err error) {
	v, err = Eval(e, ev.scope)
	if err != nil {
		restart, err = ev.bail(err)
	}
	return v, restart, err
}

func Eval(expression expr.Expr, env *value.Env) (value.Value, error) {
	ev := &evaluator{cursor: expression, scope: env}

eval_loop:
	for {
		switch e := ev.cursor.(type) {
		case expr.NumLit:
			return value.Number{Val: e.Val}, nil
		case expr.StringLit:
			return value.String{Val: e.Val}, nil
		case expr.BoolLit:
			return value.Boolean{Val: e.Val}, nil
		case expr.UnitLit:
			return value.Unit{}, nil

		case expr.Symbol:
			return ev.scope.Get(e.Val)

		case expr.Lambda:
			id := value.NextLambdaID()
			return value.Lambda{
				Params:   e.Params,
				Body:     e.Body,
				Captured: ev.scope,
				Id:       id,
			}, nil

		case expr.List:
			if len(e.Val) == 0 {
				return value.NewNil(), nil
			}

			funExpr := e.Val[0]
			args := e.Val[1:]

			funv, restart, err := ev.step(funExpr)
			if restart {
				continue eval_loop
			}
			if err != nil {
				return nil, err
			}

			switch fn := funv.(type) {
			case value.Lambda:
				if len(fn.Params) != len(args) {
					return nil, value.EvalErrf(
						"Attempted to invoke lambda with wrong amount of params. Expected: %d got: %d",
						len(fn.Params), len(args),
					)
				}
				thisEnv := value.NewEnv(fn.Captured)
				for i := range args {
					argVal, restart, err := ev.step(args[i])
					if restart {
						continue eval_loop
					}
					if err != nil {
						return nil, err
					}
					if err := thisEnv.Set(fn.Params[i], argVal); err != nil {
						return nil, err
					}
				}
				ev.cursor = fn.Body
				ev.scope = thisEnv
				continue eval_loop

			case value.Macro:
				if len(fn.Params) != len(args) {
					return nil, value.EvalErrf(
						"Attempted to invoke macro with wrong amount of params. Expected: %d got: %d",
						len(fn.Params), len(args),
					)
				}
				table := make(map[string]expr.Expr, len(args))
				for i := range args {
					table[fn.Params[i]] = args[i]
				}
				substituted := expr.List{Val: substituteAll(fn.Expressions, table)}
				ev.cursor = substituted
				continue eval_loop

			case value.NativeFunc:
				if sym, ok := funExpr.(expr.Symbol); ok {
					if isCoreFormName(sym.Val) && coreFormName(fn) == sym.Val {
						switch sym.Val {
						case "if":
							if len(args) != 3 {
								return nil, arityErr("if", 3, len(args))
							}
							cond, restart, err := ev.step(args[0])
							if restart {
								continue eval_loop
							}
							if err != nil {
								return nil, err
							}
							switch c := cond.(type) {
							case value.Boolean:
								if c.Val {
									ev.cursor = args[1]
								} else {
									ev.cursor = args[2]
								}
								continue eval_loop
							default:
								return nil, value.EvalErrf(
									"Condition of if-construct does not evaluate to a bool: %s", cond,
								)
							}

						case "do":
							if len(args) == 0 {
								return value.Unit{}, nil
							}
							for i := 0; i < len(args)-1; i++ {
								_, restart, err := ev.step(args[i])
								if restart {
									continue eval_loop
								}
								if err != nil {
									return nil, err
								}
							}
							ev.cursor = args[len(args)-1]
							continue eval_loop

						case "switch":
							if len(args) == 0 {
								return nil, arityErr("switch", 2, 0)
							}
							scrutinee, restart, err := ev.step(args[0])
							if restart {
								continue eval_loop
							}
							if err != nil {
								return nil, err
							}
							var taken expr.Expr
							found := false
							for _, caseExpr := range args[1:] {
								items, ok := caseExpr.(expr.List)
								if !ok || len(items.Val) != 2 {
									return nil, value.EvalErr("Malformed switch-arm syntax")
								}
								matcher, onMatch := items.Val[0], items.Val[1]
								if m, ok := matcher.(expr.Symbol); ok && m.Val == "_" {
									taken, found = onMatch, true
									break
								}
								matcherVal, restart, err := ev.step(matcher)
								if restart {
									continue eval_loop
								}
								if err != nil {
									return nil, err
								}
								if value.Equal(matcherVal, scrutinee) {
									taken, found = onMatch, true
									break
								}
							}
							if found {
								ev.cursor = taken
								continue eval_loop
							}
							return value.Unit{}, nil

						case "scoped":
							if len(args) != 2 {
								return nil, arityErr("scoped", 2, len(args))
							}
							bindingPairs, ok := args[0].(expr.List)
							if !ok {
								return nil, value.EvalErr(
									"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')",
								)
							}
							thisEnv := value.NewEnv(ev.scope)
							for _, pair := range bindingPairs.Val {
								items, ok := pair.(expr.List)
								if !ok || len(items.Val) != 2 {
									return nil, value.EvalErr(
										"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')",
									)
								}
								nameSym, ok := items.Val[0].(expr.Symbol)
								if !ok {
									return nil, value.EvalErr(
										"Bad Syntax! The binding list is in the wrong form! (Expected '(name value)')",
									)
								}
								val, restart, err := ev.step(items.Val[1])
								if restart {
									continue eval_loop
								}
								if err != nil {
									return nil, err
								}
								if err := thisEnv.Set(nameSym.Val, val); err != nil {
									return nil, err
								}
							}
							ev.scope = thisEnv
							ev.cursor = args[1]
							continue eval_loop

						case "try":
							if len(args) != 2 {
								return nil, arityErr("try", 2, len(args))
							}
							ev.handlers = append(ev.handlers, handler{
								cursor: args[1],
								scope:  ev.scope,
							})
							ev.cursor = args[0]
							continue eval_loop
						}
					}
				}
				v, err := fn(args, ev.scope)
				if err != nil {
					restart, err2 := ev.bail(err)
					if restart {
						continue eval_loop
					}
					return nil, err2
				}
				return v, nil

			case value.NativeFunction:
				v, err := fn.Callback(args, ev.scope)
				if err != nil {
					restart, err2 := ev.bail(err)
					if restart {
						continue eval_loop
					}
					return nil, err2
				}
				return v, nil

			default:
				return nil, value.EvalErrf(
					"Attempt to invoke non-function/non-macro: %s (%s)", funExpr, funv,
				)
			}
		}
	}
}

func arityErr(name string, expected, got int) *value.Error {
	return value.EvalErrf(
		"Native function %s was given bad syntax. Perhaps it was given the wrong amount of args? Args expected: %d. Got: %d",
		name, expected, got,
	)
}

func EvalSeq(expressions []expr.Expr, env *value.Env) (value.Value, error) {
	var last value.Value = value.Unit{}
	for _, e := range expressions {
		v, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		last = v
	}
	return last, nil
}

func substituteAll(body []expr.Expr, table map[string]expr.Expr) []expr.Expr {
	out := make([]expr.Expr, len(body))
	for i, e := range body {
		out[i] = substitute(e, table)
	}
	return out
}

func substitute(e expr.Expr, table map[string]expr.Expr) expr.Expr {
	switch x := e.(type) {
	case expr.Symbol:
		if rep, ok := table[x.Val]; ok {
			return rep
		}
	case expr.List:
		return expr.List{Val: substituteAll(x.Val, table)}
	}
	return e
}

var _ = fmt.Sprintf
