package raw

import "testing"

func TestDimension_IsValid(t *testing.T) {
	d := Dimension{}

	if err := d.IsValid(); err == nil {
		t.Error("expected error, but got nil")
	}

	d = Dimension{IdxDimID: nil}

	if err := d.IsValid(); err == nil {
		t.Error("expected error, but got nil")
	}

	d = Dimension{IdxDimID: 1}

	if err := d.IsValid(); err != nil {
		t.Error("expected nil, but got error")
	}

	d = Dimension{IdxDimID: "1"}

	if err := d.IsValid(); err != nil {
		t.Error("expected nil, but got error")
	}

	d = Dimension{IdxDimID: stringer{}}

	if err := d.IsValid(); err != nil {
		t.Error("expected nil, but got error")
	}
}

func TestDimension_ID(t *testing.T) {
	d := Dimension{IdxDimID: "1"}

	if d.ID() != "1" {
		t.Errorf("expected 1, but got %s", d.ID())
	}

	d = Dimension{IdxDimID: stringer{}}

	if d.ID() != (stringer{}).String() {
		t.Errorf("expected %s, but got %s", (stringer{}).String(), d.ID())
	}

	d = Dimension{IdxDimID: nil}

	if d.ID() != "" {
		t.Errorf("expected \"\", but got %s", d.ID())
	}

}

func TestDimension_Name(t *testing.T) {
	d := Dimension{IdxDimName: "1"}

	if d.Name() != "1" {
		t.Errorf("expected 1, but got %s", d.Name())
	}

	d = Dimension{IdxDimName: stringer{}}

	if d.Name() != (stringer{}).String() {
		t.Errorf("expected %s, but got %s", (stringer{}).String(), d.Name())
	}

	d = Dimension{IdxDimName: nil}

	if d.Name() != "" {
		t.Errorf("expected \"\", but got %s", d.Name())
	}

}

func TestDimension_Algorithm(t *testing.T) {
	d := Dimension{}

	if d.Algorithm() != defaultDimAlgorithm {
		t.Errorf("expected %s, but got %s", Absolute, d.Algorithm())
	}

	d = Dimension{IdxDimAlgorithm: Incremental}

	if d.Algorithm() != Incremental {
		t.Errorf("expected %s, but got %s", Incremental, d.Algorithm())
	}

	d = Dimension{IdxDimAlgorithm: "SuperAlgorithm"}

	if d.Algorithm() != defaultDimAlgorithm {
		t.Errorf("expected %s, but got %s", Absolute, d.Algorithm())
	}

}

func TestDimension_Multiplier(t *testing.T) {
	d := Dimension{}

	if d.Multiplier() != defaultDimMultiplier {
		t.Errorf("expected 1, but got %d", d.Multiplier())
	}

	d = Dimension{IdxDimMultiplier: 5}

	if d.Multiplier() != 5 {
		t.Errorf("expected 5, but got %d", d.Multiplier())
	}

	d = Dimension{IdxDimMultiplier: -5}

	if d.Multiplier() != -5 {
		t.Errorf("expected -5, but got %d", d.Multiplier())
	}

	d = Dimension{IdxDimMultiplier: "5"}

	if d.Multiplier() != 5 {
		t.Errorf("expected 5, but got %d", d.Multiplier())
	}

	d = Dimension{IdxDimMultiplier: 1e6}

	if d.Multiplier() != 1000000 {
		t.Errorf("expected 1000000, but got %d", d.Multiplier())
	}
}

func TestDimension_Divisor(t *testing.T) {
	d := Dimension{}

	if d.Divisor() != defaultDimDivisor {
		t.Errorf("expected 1, but got %d", d.Divisor())
	}

	d = Dimension{IdxDimDivisor: 5}

	if d.Divisor() != 5 {
		t.Errorf("expected 5, but got %d", d.Divisor())
	}

	d = Dimension{IdxDimDivisor: -5}

	if d.Divisor() != -5 {
		t.Errorf("expected -5, but got %d", d.Divisor())
	}

	d = Dimension{IdxDimDivisor: "5"}

	if d.Divisor() != 5 {
		t.Errorf("expected 5, but got %d", d.Divisor())
	}

	d = Dimension{IdxDimDivisor: 1e6}

	if d.Divisor() != 1000000 {
		t.Errorf("expected 1000000, but got %d", d.Divisor())
	}

}

func TestDimension_Hidden(t *testing.T) {
	d := Dimension{}

	if d.Hidden() != defaultDimHidden {
		t.Errorf("expected %s, but got %s", defaultDimHidden, d.Hidden())
	}

	d = Dimension{IdxDimHidden: true}

	if d.Hidden() != "hidden" {
		t.Errorf("expected hidden, but got %s", d.Hidden())
	}

}

func TestDimension_SetID(t *testing.T) {
	d := Dimension{}
	newID := "newID"

	d.SetID(newID)
	if d.ID() != newID {
		t.Errorf("expected %s, but got %s", newID, d.ID())
	}
}

func TestDimension_SetName(t *testing.T) {
	d := Dimension{}
	newName := "newID"

	d.SetName(newName)
	if d.Name() != newName {
		t.Errorf("expected %s, but got %s", newName, d.Name())
	}

}

func TestDimension_SetAlgorithm(t *testing.T) {
	d := Dimension{}

	d.SetAlgorithm(PercentOfIncremental)
	if d.Algorithm() != PercentOfIncremental {
		t.Errorf("expected %s, but got %s", PercentOfIncremental, d.Algorithm())
	}
}

func TestDimension_SetMultiplier(t *testing.T) {
	d := Dimension{IdxDimMultiplier: 1}
	newMul := 420

	d.SetMultiplier(newMul)
	if d.Multiplier() != newMul {
		t.Errorf("expected %d, but got %d", newMul, d.Multiplier())
	}
}

func TestDimension_SetDivisor(t *testing.T) {
	d := Dimension{IdxDimDivisor: 1}
	newDiv := 420

	d.SetDivisor(newDiv)
	if d.Divisor() != newDiv {
		t.Errorf("expected %d, but got %d", newDiv, d.Divisor())
	}

}

func TestDimension_SetHidden(t *testing.T) {
	d := Dimension{IdxDimHidden: false}

	d.SetHidden(true)
	if d.Hidden() != "hidden" {
		t.Errorf("expected hidden, but got %s", d.Hidden())
	}
}

func TestValidAlgorithm(t *testing.T) {
	for _, v := range []string{Absolute, Incremental, PercentOfAbsolute, PercentOfIncremental} {
		if !ValidAlgorithm(v) {
			t.Fatalf("function returned false for correct dimension algorithm")
		}
	}

	for _, v := range []string{"", "this", "is", "wrong"} {
		if ValidAlgorithm(v) {
			t.Fatalf("function returned true for incorrect dimension algorithm")
		}
	}

}
