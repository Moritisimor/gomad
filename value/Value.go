package value

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Moritisimor/gomad/expr"
)

type Value interface {
	v()
	String() string
}

type Unit struct{}

func (u Unit) v()             {}
func (u Unit) String() string { return "<UNIT>" }

type String struct{ Val string }

func (s String) v()             {}
func (s String) String() string { return s.Val }

type Number struct{ Val float64 }

func (n Number) v() {}
func (n Number) String() string {
	if math.IsNaN(n.Val) {
		return "nan"
	}
	if math.IsInf(n.Val, 0) {
		if n.Val > 0 {
			return "inf"
		}
		return "-inf"
	}
	if n.Val == float64(int64(n.Val)) && math.Abs(n.Val) < 9.223372036854776e18 {
		return fmt.Sprintf("%d", int64(n.Val))
	}
	return fmt.Sprintf("%.2f", n.Val)
}

type Boolean struct{ Val bool }

func (b Boolean) v()             {}
func (b Boolean) String() string { return fmt.Sprintf("%t", b.Val) }

type List struct {
	head Value
	tail *List
}

var nilList = List{}

func (l List) v() {}

func NewNil() List { return nilList }

func (l List) IsNil() bool { return l.tail == nil }

func (l List) Head() Value { return l.head }

func (l List) Tail() List { return *l.tail }

func Cons(h Value, t List) List {
	return List{head: h, tail: &t}
}

func FromVec(values []Value) List {
	var l List = NewNil()
	for i := len(values) - 1; i >= 0; i-- {
		l = Cons(values[i], l)
	}
	return l
}

func (l List) Len() int {
	n := 0
	for cur := l; !cur.IsNil(); cur = *cur.tail {
		n++
	}
	return n
}

func (l List) Get(i int) (Value, bool) {
	for cur := l; !cur.IsNil(); cur = *cur.tail {
		if i == 0 {
			return cur.head, true
		}
		i--
	}
	return nil, false
}

func (l List) Foreach(f func(Value)) {
	for cur := l; !cur.IsNil(); cur = *cur.tail {
		f(cur.head)
	}
}

func (l List) ToSlice() []Value {
	out := make([]Value, 0, l.Len())
	l.Foreach(func(v Value) { out = append(out, v) })
	return out
}

func (l List) Append(right List) List {
	values := l.ToSlice()
	out := right
	for i := len(values) - 1; i >= 0; i-- {
		out = Cons(values[i], out)
	}
	return out
}

func (l List) String() string {
	var b strings.Builder
	b.WriteByte('(')
	first := true
	for cur := l; !cur.IsNil(); cur = *cur.tail {
		if !first {
			b.WriteString(" ")
		}
		b.WriteString(cur.head.String())
		first = false
	}
	b.WriteByte(')')
	return b.String()
}

type Record struct {
	Val map[string]Value
	mu  sync.RWMutex
}

func NewRecord(initial ...map[string]Value) *Record {
	fields := map[string]Value{}
	if len(initial) > 0 {
		for k, v := range initial[0] {
			fields[k] = v
		}
	}
	return &Record{Val: fields}
}

func (r *Record) v()             {}
func (r *Record) String() string { return "<RECORD>" }

func (r *Record) GetField(name string) (Value, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.Val[name]
	return v, ok
}

func (r *Record) SetField(name string, v Value) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Val[name] = v
}

type NativeFunc func(args []expr.Expr, env *Env) (Value, error)

func (n NativeFunc) v()             {}
func (n NativeFunc) String() string { return "<NATIVEFUNCTION>" }

type NativeFunction struct{ Callback NativeFunc }

func (n NativeFunction) v()             {}
func (n NativeFunction) String() string { return "<NATIVEFUNCTION>" }

type Macro struct {
	Params      []string
	Expressions []expr.Expr
	Id          uint64
}

func (m Macro) v()             {}
func (m Macro) String() string { return "<MACRO>" }

type Lambda struct {
	Params   []string
	Body     expr.Expr
	Captured *Env
	Id       uint64
}

var lambdaIdCounter atomic.Uint64

func NextLambdaID() uint64 { return lambdaIdCounter.Add(1) }

func (l Lambda) v()             {}
func (l Lambda) String() string { return "<FUNCTION>" }

type Env struct {
	Bindings map[string]Value
	Parent   *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{Bindings: map[string]Value{}, Parent: parent}
}

func RootEnv() *Env { return NewEnv(nil) }

func (e *Env) Get(name string) (Value, error) {
	for cur := e; cur != nil; cur = cur.Parent {
		if v, ok := cur.Bindings[name]; ok {
			return v, nil
		}
	}
	return Unit{}, EvalErrf("No such variable: %s", name)
}

func (e *Env) Set(name string, val Value) error {
	if _, ok := e.Bindings[name]; ok {
		return EvalErrf("Cannot bind %s: Already exists in this scope", name)
	}
	e.Bindings[name] = val
	return nil
}

func (e *Env) Mutate(name string, val Value) error {
	for cur := e; cur != nil; cur = cur.Parent {
		if _, ok := cur.Bindings[name]; ok {
			cur.Bindings[name] = val
			return nil
		}
	}
	return EvalErrf("Cannot mutate non-existant binding: %s", name)
}

func (e *Env) Define(name string, val Value) {
	e.Bindings[name] = val
}

func (e *Env) RegisterNative(name string, fn NativeFunc) {
	e.Define(name, fn)
}

func (e *Env) GetBinding(name string) (Value, error)      { return e.Get(name) }
func (e *Env) SetBinding(name string, val Value) error    { return e.Set(name, val) }
func (e *Env) MutateBinding(name string, val Value) error { return e.Mutate(name, val) }

func NewUnit() Unit               { return Unit{} }
func NewNumber(n float64) Number  { return Number{Val: n} }
func NewString(s string) String   { return String{Val: s} }
func NewBool(b bool) Boolean      { return Boolean{Val: b} }
func NewList(values []Value) List { return FromVec(values) }

type ErrorKind int

const (
	ErrEval ErrorKind = iota
	ErrParse
	ErrTokenize
	ErrIo
	ErrExit
)

type Error struct {
	Kind ErrorKind
	Msg  string
	Code int
}

func (e *Error) Error() string {
	if e.Kind == ErrExit {
		return fmt.Sprintf("exit(%d)", e.Code)
	}
	return e.Msg
}

func EvalErr(msg string) *Error {
	return &Error{Kind: ErrEval, Msg: msg}
}

func EvalErrf(format string, args ...any) *Error {
	return &Error{Kind: ErrEval, Msg: fmt.Sprintf(format, args...)}
}

func Equal(a, b Value) bool {
	switch x := a.(type) {
	case Number:
		y, ok := b.(Number)
		return ok && x.Val == y.Val
	case String:
		y, ok := b.(String)
		return ok && x.Val == y.Val
	case Boolean:
		y, ok := b.(Boolean)
		return ok && x.Val == y.Val
	case Unit:
		_, ok := b.(Unit)
		return ok
	case *Record:
		y, ok := b.(*Record)
		return ok && x == y
	case Lambda:
		y, ok := b.(Lambda)
		return ok && x.Id == y.Id
	case Macro:
		y, ok := b.(Macro)
		return ok && x.Id == y.Id
	case NativeFunc:
		y, ok := b.(NativeFunc)
		return ok && reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
	case NativeFunction:
		y, ok := b.(NativeFunction)
		return ok && reflect.ValueOf(x.Callback).Pointer() == reflect.ValueOf(y.Callback).Pointer()
	case List:
		y, ok := b.(List)
		if !ok {
			return false
		}
		ac, bc := x, y
		for !ac.IsNil() && !bc.IsNil() {
			if !Equal(ac.head, bc.head) {
				return false
			}
			ac = *ac.tail
			bc = *bc.tail
		}
		return ac.IsNil() && bc.IsNil()
	}
	return false
}
