package expr

import (
	"fmt"
	"strings"
)

type Expr interface {
	expr()
	String() string
}

type Expression = Expr

type NumLit struct{ Val float64 }

func (n NumLit) expr()          {}
func (n NumLit) String() string { return fmt.Sprintf("Number(%.6f)", n.Val) }

type StringLit struct{ Val string }

func (s StringLit) expr()          {}
func (s StringLit) String() string { return fmt.Sprintf("String(%q)", s.Val) }

type BoolLit struct{ Val bool }

func (b BoolLit) expr()          {}
func (b BoolLit) String() string { return fmt.Sprintf("Bool(%t)", b.Val) }

type UnitLit struct{}

func (u UnitLit) expr()          {}
func (u UnitLit) String() string { return "<UNIT>" }

type Symbol struct{ Val string }

func (s Symbol) expr()          {}
func (s Symbol) String() string { return fmt.Sprintf("Symbol('%s')", s.Val) }

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

type Lambda struct {
	Params []string
	Body   Expr
}

type Number = NumLit
type String = StringLit
type Boolean = BoolLit
type Unit = UnitLit

func (l Lambda) expr()          {}
func (l Lambda) String() string { return fmt.Sprintf("<LAMBDA (%v)>", l.Params) }
