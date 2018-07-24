package raw

const (
	varID = iota
	varValue
)

type Variable [2]interface{}

func (v Variable) IsValid() error {
	if v.ID() == "" {
		return errNoID
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD GETTER

// ID returns 0 element of Variable converted to string.
func (v Variable) ID() string {
	id, ok := v[varID].(string)
	if ok {
		return id
	}
	return ""

}

// Value returns 1 element of Variable converted to int.
func (v Variable) Value() int64 {
	val, ok := v[varValue].(int)
	if ok {
		return int64(val)
	}
	return 0
}

// ---------------------------------------------------------------------------------------------------------------------

// FIELD SETTER

// SetID sets 0 element of Variable.
func (v *Variable) SetID(s string) *Variable {
	v[varID] = s
	return v
}

// SetValue sets 1 element of Variable.
func (v *Variable) SetValue(i int) *Variable {
	v[varValue] = i
	return v
}

// ---------------------------------------------------------------------------------------------------------------------
