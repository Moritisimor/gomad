package expr

import (
	"fmt"
	"strings"
)

type Expression interface {
	e()
	String() string
}

type Lambda struct {
	Params []string
	Body   Expression
}

func (l Lambda) e()             {}
func (l Lambda) String() string { return "<LAMBDA>" }

type Symbol struct{ Val string }

func (s Symbol) e()             {}
func (s Symbol) String() string { return fmt.Sprintf("Symbol(\"%s\")", s.Val) }

type Number struct{ Val float64 }

func (n Number) e()             {}
func (n Number) String() string { return fmt.Sprintf("Number(%.2f)", n.Val) }

type String struct{ Val string }

func (s String) e()             {}
func (s String) String() string { return fmt.Sprintf("String(\"%s\")", s.Val) }

type Boolean struct{ Val bool }

func (b Boolean) e()             {}
func (b Boolean) String() string { return fmt.Sprintf("Boolean(%t)", b.Val) }

type Unit struct{}

func (u Unit) e()             {}
func (u Unit) String() string { return "Unit" }

type List struct{ Val []Expression }

func (l List) e() {}
func (l List) String() string {
	acc := strings.Builder{}
	acc.WriteString("List(")

	for i, elem := range l.Val {
		acc.WriteString(elem.String())
		if i != len(l.Val)-1 {
			acc.WriteString(", ")
		}
	}

	acc.WriteString(")")
	return acc.String()
}
