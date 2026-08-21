package lexer

type Token interface{ t() }

type LPAREN struct{}

func (l LPAREN) t() {}

type RPAREN struct{}

func (r RPAREN) t() {}

type NUMLIT struct{ Val float64 }

func (n NUMLIT) t() {}

type BOOLLIT struct{ Val bool }

func (b BOOLLIT) t() {}

type STRINGLIT struct{ Val string }

func (s STRINGLIT) t() {}

type UNITLIT struct{ Val string }

func (u UNITLIT) t() {}

type SYMBOL struct{ Val string }

func (s SYMBOL) t() {}
