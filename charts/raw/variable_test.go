package raw

import "testing"

type stringer struct{}

func (stringer) String() string {
	return "stringer"
}

func TestVariable_IsValid(t *testing.T) {
	v := Variable{""}
	if e := v.IsValid(); e == nil {
		t.Error("expected error , but got nil")
	}

	v = Variable{"1"}
	if e := v.IsValid(); e != nil {
		t.Error("expected nil, but got error")
	}

}

func TestVariable_ID(t *testing.T) {
	id := "1"
	v := Variable{id}

	if val := v.ID(); val != id {
		t.Errorf("expected %s, but got %s", id, val)
	}

	v = Variable{stringer{}}
	if val := v.ID(); val != (stringer{}).String() {
		t.Errorf("expected %s, but got %s", (stringer{}).String(), val)
	}
}

func TestVariable_Value(t *testing.T) {
	v := Variable{"1", "33"}

	if v.Value() != 33 {
		t.Errorf("expected 33, but got %d", v.Value())
	}

	v = Variable{"1", 33}

	if v.Value() != 33 {
		t.Errorf("expected 33, but got %d", v.Value())
	}

}

func TestVariable_SetID(t *testing.T) {
	v := Variable{"1"}

	newID := "2"
	v.SetID(newID)

	if v.ID() != newID {
		t.Errorf("expected %s, but got %s", newID, v.ID())
	}

}

func TestVariable_SetValue(t *testing.T) {
	v := Variable{"1"}
	newValue := 2
	v.SetValue(newValue)

	if v.Value() != int64(newValue) {
		t.Errorf("expected %d, but got %d", newValue, v.Value())
	}

}

func TestVariable_SetIDSetValue(t *testing.T) {
	v := Variable{"1", 1}
	newID, newValue := "2", 2

	v.SetID(newID).SetValue(newValue)

	if v.ID() != newID || v.Value() != int64(newValue) {
		t.Errorf("expected %s %d, but got %s %d", newID, newValue, v.ID(), v.Value())
	}

}

func TestValidChartType(t *testing.T) {

	for _, v := range []string{Line, Area, Stacked} {
		if !ValidChartType(v) {
			t.Fatalf("function returned false for correct chart type")
		}
	}

	for _, v := range []string{"", "this", "is", "wrong"} {
		if ValidChartType(v) {
			t.Fatalf("function returned true for incorrect chart type")
		}
	}

}
