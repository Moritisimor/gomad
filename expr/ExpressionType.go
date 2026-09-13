package expr

import (
	"fmt"
	"math"
	"strings"
)

type Expr interface {
	expr()
	String() string
	Source() string
}

type Expression = Expr

type NumLit struct{ Val float64 }

func (n NumLit) expr()          {}
func (n NumLit) String() string { return fmt.Sprintf("Number(%.6f)", n.Val) }
func (n NumLit) Source() string {
	if (math.Mod(n.Val, 1) == 0) {
		return fmt.Sprintf("%d", int(n.Val))
	}

	return fmt.Sprintf("%.2f", n.Val)
}

type StringLit struct{ Val string }

func (s StringLit) expr()          {}
func (s StringLit) String() string { return fmt.Sprintf("String(%q)", s.Val) }
func (s StringLit) Source() string { return fmt.Sprintf("\"%s\"", s.Val) }

type BoolLit struct{ Val bool }

func (b BoolLit) expr()          {}
func (b BoolLit) String() string { return fmt.Sprintf("Bool(%t)", b.Val) }
func (b BoolLit) Source() string { return fmt.Sprintf("%t", b.Val) }

type UnitLit struct{}

func (u UnitLit) expr()          {}
func (u UnitLit) String() string { return "<UNIT>" }
func (u UnitLit) Source() string { return "unit" }

type Symbol struct{ Val string }

func (s Symbol) expr()          {}
func (s Symbol) String() string { return fmt.Sprintf("Symbol('%s')", s.Val) }
func (s Symbol) Source() string { return s.Val }

type List struct{ Val []Expr }

func (l List) expr() {}
func (l List) String() string {
	var b strings.Builder
	b.WriteString("List(")
	for i, e := range l.Val {
		if i > 0 {
			b.WriteString(" ")
		}

		b.WriteString(e.String())
	}

	b.WriteByte(')')
	return b.String()
}

func (l List) Source() string {
	var b strings.Builder
	b.WriteByte('(')
	for i, e := range l.Val {
		if i > 0 {
			b.WriteByte(' ')
		}

		b.WriteString(e.Source())
	}

	b.WriteByte(')')
	return b.String()
}

type Lambda struct {
	Params []string
	Body   Expr
}

func (l Lambda) expr()          {}
func (l Lambda) String() string { return fmt.Sprintf("<LAMBDA (%v)>", l.Params) }
func (l Lambda) Source() string {
	var b strings.Builder
	b.WriteString("(lambda (")
	for i, e := range l.Params {
		if i > 0 {
			b.WriteByte(' ')
		}

		b.WriteString(e)
	}

	b.WriteString(l.Body.Source())
	b.WriteByte(')')
	return b.String()
}

type Number = NumLit
type String = StringLit
type Boolean = BoolLit
type Unit = UnitLit
