package raw

import (
	"testing"
)

func TestDimension_IsValid(t *testing.T) {
	d := Dimension{}

	if err := d.IsValid(); err == nil {
		t.Error("expected error, but got nil")
	}

	d[IdxDimID] = nil

	if err := d.IsValid(); err == nil {
		t.Error("expected error, but got nil")
	}

	d[IdxDimID] = 33

	if err := d.IsValid(); err == nil {
		t.Error("expected error, but got nil")
	}

	d[IdxDimID] = "dimID"

	if err := d.IsValid(); err != nil {
		t.Error("expected nil, but got error")
	}
}

func TestDimension_ID(t *testing.T) {
	id := "dimID"
	d := Dimension{IdxDimID: id}

	if d.ID() != id {
		t.Errorf("expected %s, but got %s", id, d.ID())
	}

	d[IdxDimID] = nil

	if d.ID() != "" {
		t.Errorf("expected empty string, but got %s", d.ID())
	}
}

func TestDimension_Name(t *testing.T) {
	name := "dimName"
	d := Dimension{IdxDimName: name}

	if d.Name() != name {
		t.Errorf("expected %s, but got %s", name, d.Name())
	}

	d[IdxDimName] = nil

	if d.Name() != "" {
		t.Errorf("expected empty string, but got %s", d.Name())
	}
}

func TestDimension_Algorithm(t *testing.T) {
	d := Dimension{}

	if d.Algorithm() != defaultDimAlgorithm {
		t.Errorf("expected %s, but got %s", Absolute, d.Algorithm())
	}

	d[IdxDimAlgorithm] = Incremental

	if d.Algorithm() != Incremental {
		t.Errorf("expected %s, but got %s", Incremental, d.Algorithm())
	}

	d[IdxDimAlgorithm] = "Algorithm"

	if d.Algorithm() != defaultDimAlgorithm {
		t.Errorf("expected %s, but got %s", Absolute, d.Algorithm())
	}
}

func TestDimension_Multiplier(t *testing.T) {
	d := Dimension{}

	if d.Multiplier() != defaultDimMultiplier {
		t.Errorf("expected %d, but got %d", defaultDimMultiplier, d.Multiplier())
	}

	d[IdxDimMultiplier] = 5

	if d.Multiplier() != 5 {
		t.Errorf("expected 5, but got %d", d.Multiplier())
	}

	d[IdxDimMultiplier] = -5

	if d.Multiplier() != defaultDimMultiplier {
		t.Errorf("expected %d, but got %d", defaultDimMultiplier, d.Multiplier())
	}

	d[IdxDimMultiplier] = "5"

	if d.Multiplier() != defaultDimMultiplier {
		t.Errorf("expected 1, but got %d", d.Multiplier())
	}

	d[IdxDimMultiplier] = 1e7

	if d.Multiplier() != defaultDimMultiplier {
		t.Errorf("expected %d, but got %d", defaultDimMultiplier, d.Multiplier())
	}
}

func TestDimension_Divisor(t *testing.T) {
	d := Dimension{}

	if d.Divisor() != defaultDimDivisor {
		t.Errorf("expected %d, but got %d", defaultDimDivisor, d.Divisor())
	}

	d[IdxDimDivisor] = 5

	if d.Divisor() != 5 {
		t.Errorf("expected 5, but got %d", d.Divisor())
	}

	d[IdxDimDivisor] = -5

	if d.Divisor() != defaultDimDivisor {
		t.Errorf("expected %d, but got %d", defaultDimDivisor, d.Divisor())
	}

	d[IdxDimDivisor] = "5"

	if d.Divisor() != defaultDimDivisor {
		t.Errorf("expected %d, but got %d", defaultDimDivisor, d.Divisor())
	}

	d[IdxDimDivisor] = 1e7

	if d.Divisor() != defaultDimMultiplier {
		t.Errorf("expected %d, but got %d", defaultDimMultiplier, d.Divisor())
	}
}

func TestDimension_Hidden(t *testing.T) {
	d := Dimension{}

	if d.Hidden() != defaultDimHidden {
		t.Errorf("expected %s, but got %s", defaultDimHidden, d.Hidden())
	}

	d[IdxDimHidden] = true

	if d.Hidden() != "hidden" {
		t.Errorf("expected hidden, but got %s", d.Hidden())
	}

	d[IdxDimHidden] = "hidden"

	if d.Hidden() != "hidden" {
		t.Errorf("expected hidden, but got %s", d.Hidden())
	}
}

func TestDimension_SetID(t *testing.T) {
	d := Dimension{}
	id := "newID"

	d.SetID(id)

	if d.ID() != id {
		t.Errorf("expected %s, but got %s", id, d.ID())
	}
}

func TestDimension_SetName(t *testing.T) {
	d := Dimension{}
	name := "newName"

	d.SetName(name)

	if d.Name() != name {
		t.Errorf("expected %s, but got %s", name, d.Name())
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
	d := Dimension{}
	mul := 420

	d.SetMultiplier(mul)
	if d.Multiplier() != mul {
		t.Errorf("expected %d, but got %d", mul, d.Multiplier())
	}
}

func TestDimension_SetDivisor(t *testing.T) {
	d := Dimension{}
	div := 420

	d.SetDivisor(div)

	if d.Divisor() != div {
		t.Errorf("expected %d, but got %d", div, d.Divisor())
	}
}

func TestDimension_SetHidden(t *testing.T) {
	d := Dimension{}

	d.SetHidden(true)

	if d.Hidden() != "hidden" {
		t.Errorf("expected hidden, but got %s", d.Hidden())
	}
}

func TestValidAlgorithm(t *testing.T) {
	for _, v := range []string{
		Absolute,
		Incremental,
		PercentOfAbsolute,
		PercentOfIncremental,
	} {
		if !ValidAlgorithm(v) {
			t.Fatalf("function returned false for correct dimension algorithm")
		}
	}

	for _, v := range []string{
		"",
		"this",
		"is",
		"wrong",
	} {
		if ValidAlgorithm(v) {
			t.Fatalf("function returned true for incorrect dimension algorithm")
		}
	}
}
