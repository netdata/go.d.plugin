package charts

var (
	Absolute             = algorithm{"absolute"}
	Incremental          = algorithm{"incremental"}
	PercentOfAbsolute    = algorithm{"percentage-of-absolute-row"}
	PercentOfIncremental = algorithm{"percentage-of-incremental-row"}

	Hidden    = hidden{"hidden"}
	NotHidden = hidden{}
)

type (
	algorithm struct {
		a string
	}
	hidden struct {
		h string
	}
)

func (a algorithm) String() string {
	return a.a
}

func (h hidden) String() string {
	return h.h
}

type (
	Dim struct {
		ID     string
		Name   string
		Algo   algorithm
		Mul    int
		Div    int
		Hidden hidden
	}
)

func (d Dim) copy() *Dim {
	return &d
}

func (d Dim) isValid() bool {
	return d.ID != ""
}
