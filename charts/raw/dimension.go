package raw

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	idxDimID = iota
	idxDimName
	idxDimAlgorithm
	idxDimMultiplier
	idxDimDivisor
	idxDimHidden
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

func (d *Dimension) IsValid() error {
	if d.ID() != "" {
		return nil
	}
	return errors.New("id not specified")
}

// ---------------------------------------------------------------------------------------------------------------------_

// FIELD GETTER

// ID returns 0 element of Dimension converted to string (valid types are fmt.Stringer, string, int).
func (d *Dimension) ID() string {
	switch v := d[idxDimID].(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// Name returns 1 element of Dimension converted to string (valid types are fmt.Stringer, string, int).
func (d *Dimension) Name() string {
	switch v := d[idxDimName].(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// Algorithm returns 2 element of Dimension converted to string.
func (d *Dimension) Algorithm() string {
	switch v := d[idxDimAlgorithm].(type) {
	case string:
		if ValidAlgorithm(v) {
			return v
		}
		return defaultDimAlgorithm
	default:
		return defaultDimAlgorithm
	}
}

// Multiplier returns 3 element of Dimension converted to int.
func (d *Dimension) Multiplier() int {
	switch v := d[idxDimMultiplier].(type) {
	case string:
		if val, err := strconv.Atoi(v); err != nil || val <= 0 {
			return defaultDimMultiplier
		} else {
			return val
		}
	case int:
		if v <= 0 {
			return defaultDimMultiplier
		}
		return v
	case float64:
		if v <= 0 {
			return defaultDimMultiplier
		}
		return int(v)
	case float32:
		if v <= 0 {
			return defaultDimMultiplier
		}
		return int(v)
	default:
		return defaultDimMultiplier
	}
}

// Divisor returns 4 element of Dimension converted to int.
func (d *Dimension) Divisor() int {
	switch v := d[idxDimDivisor].(type) {
	case string:
		if val, err := strconv.Atoi(v); err != nil || val <= 0 {
			return defaultDimDivisor
		} else {
			return val
		}
	case int:
		if v <= 0 {
			return defaultDimDivisor
		}
		return v
	case float64:
		if v <= 0 {
			return defaultDimDivisor
		}
		return int(v)
	case float32:
		if v <= 0 {
			return defaultDimDivisor
		}
		return int(v)
	default:
		return defaultDimDivisor
	}
}

// Hidden returns 5 element of Dimension converted to string.
func (d *Dimension) Hidden() string {
	switch v := d[idxDimHidden].(type) {
	case string:
		if v == "hidden" {
			return "hidden"
		}
		return defaultDimHidden
	case bool:
		if v {
			return "hidden"
		}
		return defaultDimHidden
	default:
		return defaultDimHidden
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets 0 element of Dimension (valid types are fmt.Stringer, string, int).
func (d *Dimension) SetID(a interface{}) *Dimension {
	d[idxDimID] = a
	return d
}

// SetName sets 1 element of Dimension (valid types are fmt.Stringer, string, int).
func (d *Dimension) SetName(a interface{}) *Dimension {
	d[idxDimName] = a
	return d
}

// SetAlgorithm sets 2 element of Dimension.
func (d *Dimension) SetAlgorithm(a string) *Dimension {
	d[idxDimAlgorithm] = a
	return d
}

// SetMultiplier sets 3 element of Dimension.
func (d *Dimension) SetMultiplier(a int) *Dimension {
	d[idxDimMultiplier] = a
	return d
}

// SetDivisor sets 4 element of Dimension.
func (d *Dimension) SetDivisor(a int) *Dimension {
	d[idxDimDivisor] = a
	return d
}

// SetHidden sets 5 element of Dimension.
func (d *Dimension) SetHidden(a bool) *Dimension {
	d[idxDimHidden] = a
	return d
}

// ---------------------------------------------------------------------------------------------------------------------

// ValidAlgorithm returns whether the dimension algorithm is valid.
// Valid algorithms: "absolute", "incremental", "percentage-of-absolute-row", "percentage-of-incremental-row".
func ValidAlgorithm(a string) bool {
	switch a {
	case Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental:
		return true
	}
	return false
}
