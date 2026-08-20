package lexer

type Token interface{ t() }

type LPAREN struct{}

func (l LPAREN) t() {}

type RPAREN struct{}

func (r RPAREN) t() {}

type NUMLIT struct{ val float64 }

func (n NUMLIT) t() {}

type BOOLLIT struct{ val bool }

func (b BOOLLIT) t() {}

type STRINGLIT struct{ val string }

func (s STRINGLIT) t() {}

type UNITLIT struct{ val string }

func (u UNITLIT) t() {}

type SYMBOL struct{ val string }

func (s SYMBOL) t() {}
