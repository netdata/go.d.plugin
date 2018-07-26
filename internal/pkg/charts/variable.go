package charts

type Var struct {
	ID    string
	Value int64
}

func (v Var) copy() Var {
	return v
}
