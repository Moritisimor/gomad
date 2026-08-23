package value

import (
	"fmt"
	"math"
	"strings"

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
	if math.Mod(n.Val, 1.0) == 0 {
		return fmt.Sprintf("%d", int64(n.Val))
	}

	return fmt.Sprintf("%.2f", n.Val)
}

type Boolean struct{ Val bool }

func (b Boolean) v()             {}
func (b Boolean) String() string { return fmt.Sprintf("%t", b.Val) }

type List struct{ Val []Value }

func (l List) v() {}
func (l List) String() string {
	acc := strings.Builder{}
	acc.WriteByte('(')

	for i, v := range l.Val {
		acc.WriteString(v.String())
		if i != len(l.Val)-1 {
			acc.WriteString(", ")
		}
	}

	acc.WriteByte(')')
	return acc.String()
}

type Record struct{ Val map[string]Value }

func (r Record) v() {}
func (r Record) String() string {
	acc := strings.Builder{}
	acc.WriteString("{ ")
	for k, v := range r.Val {
		acc.WriteByte('(')
		acc.WriteString(k)
		acc.WriteString(": ")
		acc.WriteString(v.String())
		acc.WriteString(") ")
	}

	acc.WriteString("}")
	return acc.String()
}

type NativeFunction struct {
	Callback func(e []expr.Expression, env *Env) (Value, error)
}

func (n NativeFunction) v()             {}
func (n NativeFunction) String() string { return "<NATIVEFUNCTION>" }

type Macro struct {
	Params      []string
	Expressions expr.Expression
}

func (m Macro) v()             {}
func (m Macro) String() string { return "<MACRO>" }

type Lambda struct {
	Params   []string
	Body     expr.Expression
	Captured *Env
}

func (l Lambda) v()             {}
func (l Lambda) String() string { return "<LAMBDA>" }
