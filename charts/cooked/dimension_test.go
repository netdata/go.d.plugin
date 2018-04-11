package cooked

import (
	"testing"

	"github.com/l2isbad/go.d.plugin/charts/raw"
)

func TestDimension_ID(t *testing.T) {
	id := "id"
	d := dimension{id: id}

	if d.ID() != id {
		t.Errorf("expected %s, but got %s", id, d.ID())
	}
}

func TestDimension_Name(t *testing.T) {
	n := "name"
	d := dimension{name: n}

	if d.Name() != n {
		t.Errorf("expected %s, but got %s", n, d.Name())
	}
}

func TestDimension_Algorithm(t *testing.T) {
	a := "a"
	d := dimension{algorithm: a}

	if d.Algorithm() != a {
		t.Errorf("expected %s, but got %s", a, d.Algorithm())
	}
}

func TestDimension_Multiplier(t *testing.T) {
	m := 2
	d := dimension{multiplier: m}

	if d.Multiplier() != m {
		t.Errorf("expected %d, but got %d", m, d.Multiplier())
	}
}

func TestDimension_Divisor(t *testing.T) {
	m := 2
	d := dimension{divisor: m}

	if d.Divisor() != m {
		t.Errorf("expected %d, but got %d", m, d.Divisor())
	}

}

func TestDimension_Hidden(t *testing.T) {
	if (&dimension{}).Hidden() {
		t.Error("expected false, but got true")
	}

	if !(&dimension{hidden: "hidden"}).Hidden() {
		t.Error("expected true, but got false")
	}
}

func TestDimension_SetID(t *testing.T) {
	id := "id"
	d := dimension{}

	if d.SetID(id); d.ID() != id {
		t.Errorf("expected %s, but got %s", id, d.ID())
	}
}

func TestDimension_SetName(t *testing.T) {
	n := "n"
	d := dimension{}

	if d.SetName(n); d.Name() != n {
		t.Errorf("expected %s, but got %s", n, d.Name())
	}
}

func TestDimension_SetAlgorithm(t *testing.T) {
	d := dimension{}

	if d.SetAlgorithm(raw.Incremental); d.Algorithm() != raw.Incremental {
		t.Errorf("expected %s, but got %s", raw.Incremental, d.Algorithm())
	}

	if d.SetAlgorithm("SecretAlgorithm"); d.Algorithm() != raw.Incremental {
		t.Errorf("expected %s, but got %s", raw.Incremental, d.Algorithm())
	}
}

func TestDimension_SetMultiplier(t *testing.T) {
	m := 2
	d := dimension{}

	if d.SetMultiplier(m); d.Multiplier() != m {
		t.Errorf("expected %d, but got %d", m, d.Multiplier())
	}

	if d.SetMultiplier(-1); d.Multiplier() != m {
		t.Errorf("expected %d, but got %d", m, d.Multiplier())
	}

}

func TestDimension_SetDivisor(t *testing.T) {
	m := 2
	d := dimension{}

	if d.SetDivisor(m); d.Divisor() != m {
		t.Errorf("expected %d, but got %d", m, d.Divisor())
	}

	if d.SetDivisor(-1); d.Divisor() != m {
		t.Errorf("expected %d, but got %d", m, d.Divisor())
	}
}

func TestDimension_SetHidden(t *testing.T) {
	d := dimension{}

	if d.SetHidden(false); d.Hidden() != false {
		t.Error("expected false, but got true")
	}
	if d.SetHidden(true); d.Hidden() != true {
		t.Error("expected true, but got false")
	}
}

func TestDimension_MultipleSet(t *testing.T) {
	d := dimension{}
	id, name, mul, div := "id", "name", 10, 100

	d.SetID(id).SetName(name).SetAlgorithm(raw.Absolute).SetMultiplier(mul).SetDivisor(div)

	if d.ID() != id || d.Name() != name || d.Algorithm() != raw.Absolute || d.Multiplier() != mul || d.Divisor() != div {
		t.Error("multiple set failed")
	}
}
