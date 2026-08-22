package lexer

import "fmt"

type Token interface {
	t()
	String() string
}

type LPAREN struct{}

func (l LPAREN) t()             {}
func (l LPAREN) String() string { return "LPAREN" }

type RPAREN struct{}

func (r RPAREN) t()             {}
func (r RPAREN) String() string { return "RPAREN" }

type NUMLIT struct{ Val float64 }

func (n NUMLIT) t()             {}
func (n NUMLIT) String() string { return fmt.Sprintf("NUMLIT(%.2f)", n.Val) }

type BOOLLIT struct{ Val bool }

func (b BOOLLIT) t()             {}
func (b BOOLLIT) String() string { return fmt.Sprintf("BOOLLIT(%t)", b.Val) }

type STRINGLIT struct{ Val string }

func (s STRINGLIT) t()             {}
func (s STRINGLIT) String() string { return fmt.Sprintf("STRINGLIT(\"%s\")", s.Val) }

type UNITLIT struct{ Val string }

func (u UNITLIT) t()             {}
func (u UNITLIT) String() string { return "UNITLIT" }

type SYMBOL struct{ Val string }

func (s SYMBOL) t()             {}
func (s SYMBOL) String() string { return fmt.Sprintf("SYMBOL('%s')", s.Val) }
