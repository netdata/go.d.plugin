package raw

import "testing"

func TestVariable_IsValid(t *testing.T) {
	v := Variable{""}

	if e := v.IsValid(); e == nil {
		t.Error("expected error , but got nil")
	}

	v = Variable{1}

	if e := v.IsValid(); e == nil {
		t.Error("expected error, but got nil")
	}

	v = Variable{"varID"}

	if e := v.IsValid(); e != nil {
		t.Error("expected nil, but got error")
	}

}

func TestVariable_ID(t *testing.T) {
	id := "varID"
	v := Variable{id}

	if v.ID() != id {
		t.Errorf("expected %s, but got %s", id, v.ID())
	}
}

func TestVariable_Value(t *testing.T) {
	v := Variable{1: 33}

	if v.Value() != 33 {
		t.Errorf("expected 33, but got %d", v.Value())
	}

	v = Variable{1: "33"}

	if v.Value() != 0 {
		t.Errorf("expected 0, but got %d", v.Value())
	}

}

func TestVariable_SetID(t *testing.T) {
	v := Variable{"varID"}

	newID := "newID"
	v.SetID(newID)

	if v.ID() != newID {
		t.Errorf("expected %s, but got %s", newID, v.ID())
	}

}

func TestVariable_SetValue(t *testing.T) {
	v := Variable{}

	val := 2
	v.SetValue(val)

	if v.Value() != int64(val) {
		t.Errorf("expected %d, but got %d", val, v.Value())
	}

}

func TestVariable_ChainSet(t *testing.T) {
	v := Variable{"varID", 33}
	id, val := "2", 2

	v.SetID(id).SetValue(val)

	if v.ID() != id || v.Value() != int64(val) {
		t.Errorf("expected %s %d, but got %s %d", id, val, v.ID(), v.Value())
	}
}
