package raw

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	varID = iota
	varValue
)

type Variable [2]interface{}

func (v *Variable) IsValid() error {
	if v.ID() == "" {
		return errors.New("id not specified")
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// ID returns 0 element of Variable converted to string (valid types are fmt.Stringer, string).
func (v *Variable) ID() string {
	switch v := v[varID].(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}

// Value returns 1 element of Variable converted to int.
func (v *Variable) Value() int64 {
	switch v := v[varValue].(type) {
	case int:
		return int64(v)
	case string:
		if c, err := strconv.Atoi(v); err == nil {
			return int64(c)
		}
		return 0
	default:
		return 0
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets 0 element of Variable (valid types are fmt.Stringer, string).
func (v *Variable) SetID(id interface{}) *Variable {
	v[varID] = id
	return v
}

// SetValue sets 1 element of Variable.
func (v *Variable) SetValue(val int) *Variable {
	v[varValue] = val
	return v
}

// ---------------------------------------------------------------------------------------------------------------------
