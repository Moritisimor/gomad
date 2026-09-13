package expr

func NewNumLit(n float64) NumLit {
	return NumLit{n}
}

func NewStringLit(s string) StringLit {
	return StringLit{s}
}

func NewBoolLit(b bool) BoolLit {
	return BoolLit{b}
}

func NewSymbol(s string) Symbol {
	return Symbol{s}
}

func NewUnitLit() UnitLit {
	return UnitLit{}
}

func NewList(elems... Expr) List {
	return List{elems}
}

func NewLambda(params []string, body Expr) Lambda {
	return Lambda{
		params, 
		body,
	}
}
