package charts

type Variable struct {
	ID    string
	Value int64
}

func (v Variable) copy() Variable {
	return v
}
