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
	Dimension struct {
		ID         string
		Name       string
		Algorithm  algorithm
		Multiplier int
		Divisor    int
		Hidden     hidden
	}
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

func (d Dimension) copy() Dimension {
	return d
}
