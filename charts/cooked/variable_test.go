package cooked

import (
	"testing"
)

func TestVariable_ID(t *testing.T) {
	id := "id"
	v := variable{id:id}

	if v.ID() != id {
		t.Errorf("expected %s, but got %s", id, v.ID())
	}
}

func TestVariable_Value(t *testing.T) {
	value := int64(1)
	v := variable{value:value}

	if v.Value() != value {
		t.Errorf("expected %d, but got %d", value, v.Value())
	}
}

func TestVariable_SetValue(t *testing.T) {
	v := variable{}
	value := int64(2)

	v.SetValue(value)
	if v.Value() != value {
		t.Errorf("expected %d, but got %d", value, v.Value())
	}
}
