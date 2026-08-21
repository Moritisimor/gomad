package expr

type Expression interface { e() }

type Lambda struct {
	Params []string
	Body Expression
}

func (l Lambda) e() {}

type Symbol struct{ Val string }

func (s Symbol) e() {}

type Number struct{ Val float64 }

func (n Number) e() {}

type String struct{ Val string }

func (s String) e() {}

type Boolean struct{ Val bool }

func (b Boolean) e() {}

type Unit struct{}

func (u Unit) e() {}

type List struct{ Val []Expression }

func (l List) e() {}
