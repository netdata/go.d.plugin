package raw

const (
	IdxDimID = iota
	IdxDimName
	IdxDimAlgorithm
	IdxDimMultiplier
	IdxDimDivisor
	IdxDimHidden
)

const (
	Absolute             = "absolute"
	Incremental          = "incremental"
	PercentOfAbsolute    = "percentage-of-absolute-row"
	PercentOfIncremental = "percentage-of-incremental-row"
)

var (
	defaultDimMultiplier = 1
	defaultDimDivisor    = 1
	defaultDimAlgorithm  = Absolute
	defaultDimHidden     = ""
)

type Dimension [6]interface{}

func (d Dimension) IsValid() error {
	if d.ID() != "" {
		return nil
	}
	return errNoID
}

// ---------------------------------------------------------------------------------------------------------------------_

// FIELD GETTER

// ID returns 0 element of Dimension converted to string.
func (d Dimension) ID() string {
	id, ok := d[IdxDimID].(string)
	if ok {
		return id
	}
	return ""
}

// Name returns 1 element of Dimension converted to string.
func (d Dimension) Name() string {
	name, ok := d[IdxDimName].(string)
	if ok {
		return name
	}
	return ""
}

// Algorithm returns 2 element of Dimension converted to string.
func (d Dimension) Algorithm() string {
	algorithm, ok := d[IdxDimAlgorithm].(string)
	if !ok {
		return defaultDimAlgorithm
	}
	if !ValidAlgorithm(algorithm) {
		return defaultDimAlgorithm
	}
	return algorithm
}

// Multiplier returns 3 element of Dimension converted to int.
func (d Dimension) Multiplier() int {
	mul, ok := d[IdxDimMultiplier].(int)
	if ok && mul > 0 {
		return mul
	}
	return defaultDimMultiplier
}

// Divisor returns 4 element of Dimension converted to int.
func (d Dimension) Divisor() int {
	div, ok := d[IdxDimDivisor].(int)
	if ok && div > 0 {
		return div
	}
	return defaultDimDivisor
}

// Hidden returns 5 element of Dimension converted to string.
func (d Dimension) Hidden() string {
	switch v := d[IdxDimHidden].(type) {
	case string:
		if v == "hidden" {
			return "hidden"
		}
	case bool:
		if v {
			return "hidden"
		}
	}
	return defaultDimHidden
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets 0 element of Dimension.
func (d *Dimension) SetID(s string) *Dimension {
	d[IdxDimID] = s
	return d
}

// SetName sets 1 element of Dimension.
func (d *Dimension) SetName(s string) *Dimension {
	d[IdxDimName] = s
	return d
}

// SetAlgorithm sets 2 element of Dimension.
func (d *Dimension) SetAlgorithm(s string) *Dimension {
	d[IdxDimAlgorithm] = s
	return d
}

// SetMultiplier sets 3 element of Dimension.
func (d *Dimension) SetMultiplier(i int) *Dimension {
	d[IdxDimMultiplier] = i
	return d
}

// SetDivisor sets 4 element of Dimension.
func (d *Dimension) SetDivisor(i int) *Dimension {
	d[IdxDimDivisor] = i
	return d
}

// SetHidden sets 5 element of Dimension.
func (d *Dimension) SetHidden(b bool) *Dimension {
	d[IdxDimHidden] = b
	return d
}

// ---------------------------------------------------------------------------------------------------------------------

// ValidAlgorithm returns whether the dimension algorithm is valid.
// Valid algorithms: "absolute", "incremental", "percentage-of-absolute-row", "percentage-of-incremental-row".
func ValidAlgorithm(s string) bool {
	switch s {
	case Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental:
		return true
	}
	return false
}
