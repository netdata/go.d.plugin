package cooked

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

func newVariable(v raw.Variable) (*variable, error) {
	if err := v.IsValid(); err != nil {
		return nil, err
	}
	return &variable{
		id:    v.ID(),
		value: v.Value(),
	}, nil
}

type variable struct {
	id    string
	value int64
}

// set formats variables SET (ex.: "VARIABLE CHART 'max_conn' = '123'\n").
func (v variable) set(value int64) string {
	return fmt.Sprintf(formatVarSET, v.id, value)
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// ID returns variable id.
func (v variable) ID() string {
	return v.id
}

// Value returns variable value.
func (v variable) Value() int64 {
	return v.value
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetValue sets variable value.
func (v *variable) SetValue(val int64) *variable {
	v.value = val
	return v
}

// ---------------------------------------------------------------------------------------------------------------------
