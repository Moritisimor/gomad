package value

func NewUnit() Unit {
	return Unit{}
}

func NewNumber(n float64) Number {
	return Number{Val: n}
}

func NewString(s string) String {
	return String{Val: s}
}

func NewBool(b bool) Boolean {
	return Boolean{Val: b}
}

func NewList(l []Value) List {
	return List{Val: l}
}

func NewRecord(r map[string]Value) Record {
	return Record{Val: r}
}
