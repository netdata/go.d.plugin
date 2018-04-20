package cooked

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/charts/raw"
)

// TODO actually there is no "delete dimension" feature in netdata, we can only hide it :(
// hide dimension after N failed updates in a row.
var dimMaxRetries = 5

func newDimension(d raw.Dimension) (*dimension, error) {
	if err := d.IsValid(); err != nil {
		return nil, err
	}
	return &dimension{
		id:         d.ID(),
		name:       d.Name(),
		algorithm:  d.Algorithm(),
		multiplier: d.Multiplier(),
		divisor:    d.Divisor(),
		hidden:     d.Hidden(),
		flagsDim:   &flagsDim{retriesMax: dimMaxRetries},
	}, nil

}

type dimension struct {
	id         string
	name       string
	algorithm  string
	multiplier int
	divisor    int
	hidden     string
	*flagsDim
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// ID returns dimension id.
func (d *dimension) ID() string {
	return d.id
}

// Name returns dimension name.
func (d *dimension) Name() string {
	return d.name
}

// Algorithm returns dimension algorithm.
func (d *dimension) Algorithm() string {
	return d.algorithm
}

// Multiplier returns dimension multiplier.
func (d *dimension) Multiplier() int {
	return d.multiplier
}

// Divisor returns dimension divisor.
func (d *dimension) Divisor() int {
	return d.divisor
}

// Hidden returns dimension hidden.
func (d *dimension) Hidden() bool {
	return d.hidden == "hidden"
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets dimension id.
func (d *dimension) SetID(id string) *dimension {
	d.id = id
	return d
}

// SetName sets dimension name.
func (d *dimension) SetName(name string) *dimension {
	d.name = name
	return d
}

// SetAlgorithm sets dimension algorithm (only if algorithm is valid).
func (d *dimension) SetAlgorithm(algorithm string) *dimension {
	if raw.ValidAlgorithm(algorithm) {
		d.algorithm = algorithm
	}
	return d
}

// SetMultiplier sets dimension multiplier (only if mul > 0).
func (d *dimension) SetMultiplier(mul int) *dimension {
	if mul > 0 {
		d.multiplier = mul
	}
	return d
}

// SetDivisor sets dimension divisor (only if div > 0).
func (d *dimension) SetDivisor(div int) *dimension {
	if div > 0 {
		d.divisor = div
	}
	return d
}

// SetHidden sets dimension hidden.
func (d *dimension) SetHidden(b bool) *dimension {
	if b {
		d.hidden = "hidden"
	} else {
		d.hidden = ""
	}
	return d
}

// ---------------------------------------------------------------------------------------------------------------------

func (d *dimension) create() string {
	return fmt.Sprintf(formatDimCREATE,
		d.id,
		d.name,
		d.algorithm,
		d.multiplier,
		d.divisor,
		d.hidden)
}

func (d *dimension) get(m *map[string]int64) (int64, bool) {
	v, ok := (*m)[d.id]
	if !ok && d.alive() {
		d.retries++
		if !d.alive() {
			d.SetHidden(true)
			d.push = true
		}
	}
	return v, ok
}

func (d *dimension) set(value int64) string {
	if !d.alive() {
		d.SetHidden(false)
		d.push = true
	}
	d.retries = 0
	return fmt.Sprintf(formatDimSET,
		d.id,
		value)
}
