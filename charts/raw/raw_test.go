package raw

import "testing"

func TestVariable(t *testing.T) {
	if len(Variable{}) != 2 {
		t.Errorf("expected 2 length, but got %d", len(Variable{}))
	}

	v := Variable{""}

	if e := v.IsValid(); e == nil {
		t.Error("id = \"\", expected error after IsValid(), but got nil")
	}

	v = Variable{"var1"}

	if e := v.IsValid(); e != nil {
		t.Error("id = \"var1\", expected nil after IsValid(), but got error")
	}

	v = Variable{"var1", "33"}

	if v.Value() != 33 {
		t.Errorf("val = \"33\", expected 33 after Value(), but got %d", v.Value())
	}

	v = Variable{"var1", 33}

	if v.Value() != 33 {
		t.Errorf("val = 33, expected 33 after Value(), but got %d", v.Value())
	}

	newID := "var2"
	v.SetID(newID)

	if v.ID() != newID {
		t.Errorf("expected \"%s\" after setID(), but got %s", newID, v.ID())
	}

	newValue := 99
	v.SetValue(newValue)

	if v.Value() != int64(newValue) {
		t.Errorf("expected %d after setValue(), but got %d", newValue, v.Value())
	}

	v.SetID("var3").SetValue(111)

	if v.ID() != "var3" || v.Value() != 111 {
		t.Errorf("expected \"var3\" 111 after SetID().setValue(), but got %s %d", v.ID(), v.Value())
	}
}
