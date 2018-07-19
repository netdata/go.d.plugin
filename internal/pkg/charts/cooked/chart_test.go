package cooked

import (
	"fmt"
	"testing"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

var testRawChart = raw.Chart{
	ID:      "chart1",
	Options: raw.Options{"Title", "Units", "Family", "Context", "Type", "OverrideID"},
	Dimensions: raw.Dimensions{
		raw.Dimension{"dim1"},
		raw.Dimension{"dim2"},
	},
	Variables: raw.Variables{
		raw.Variable{"var1"},
		raw.Variable{"var2"},
	},
}

type testBC struct{}

func (b testBC) GetModuleName() string { return "module" }

func (b testBC) GetJobName() string { return "job" }

func (b testBC) GetFullName() string { return "full" }

func (b testBC) GetUpdateEvery() int { return 1 }

func TestChart_ID(t *testing.T) {
	id := "id"
	c := Chart{id: id}

	if c.ID() != id {
		t.Errorf("expected %s, but got %s", id, c.ID())
	}
}

func TestChart_OverrideID(t *testing.T) {
	id := "id"
	c := Chart{overrideID: id}

	if c.OverrideID() != id {
		t.Errorf("expected %s, but got %s", id, c.OverrideID())
	}
}

func TestChart_Title(t *testing.T) {
	title := "title"
	c := Chart{title: title}

	if c.Title() != title {
		t.Errorf("expected %s, but got %s", title, c.Title())
	}
}

func TestChart_Units(t *testing.T) {
	u := "units"
	c := Chart{units: u}

	if c.Units() != u {
		t.Errorf("expected %s, but got %s", u, c.Units())
	}
}

func TestChart_Family(t *testing.T) {
	f := "family"
	c := Chart{family: f}

	if c.Family() != f {
		t.Errorf("expected %s, but got %s", f, c.Family())
	}
}

func TestChart_Context(t *testing.T) {
	ctx := "context"
	c := Chart{id: "id", context: ctx, bc: testBC{}}

	if c.Context() != ctx {
		t.Errorf("expected %s, but got %s", ctx, c.Context())
	}

	c.context = ""
	ctx = fmt.Sprintf("%s.%s", c.bc.ModuleName(), c.id)
	if c.Context() != ctx {
		t.Errorf("expected %s, but got %s", ctx, c.Context())
	}
}

func TestChart_ChartType(t *testing.T) {
	ct := "type"
	c := Chart{chartType: ct}

	if c.ChartType() != ct {
		t.Errorf("expected %s, but got %s", ct, c.ChartType())
	}
}

func TestChart_SetID(t *testing.T) {
	id := "id"
	c := Chart{}

	if c.SetID(id); c.ID() != id {
		t.Errorf("expected %s, but got %s", id, c.ID())
	}
}

func TestChart_SetOverrideID(t *testing.T) {
	id := "id"
	c := Chart{}

	if c.SetOverrideID(id); c.OverrideID() != id {
		t.Errorf("expected %s, but got %s", id, c.OverrideID())
	}
}

func TestChart_SetTitle(t *testing.T) {
	title := "title"
	c := Chart{}

	if c.SetTitle(title); c.Title() != title {
		t.Errorf("expected %s, but got %s", title, c.Title())
	}
}

func TestChart_SetUnits(t *testing.T) {
	u := "units"
	c := Chart{}

	if c.SetUnits(u); c.Units() != u {
		t.Errorf("expected %s, but got %s", u, c.Units())
	}
}

func TestChart_SetFamily(t *testing.T) {
	f := "family"
	c := Chart{}

	if c.SetFamily(f); c.Family() != f {
		t.Errorf("expected %s, but got %s", f, c.Family())
	}
}

func TestChart_SetContext(t *testing.T) {
	ctx := "context"
	c := Chart{}

	if c.SetContext(ctx); c.Context() != ctx {
		t.Errorf("expected %s, but got %s", ctx, c.Context())
	}
}

func TestChart_SetChartType(t *testing.T) {
	c := Chart{}

	if c.SetChartType(raw.Stacked); c.ChartType() != raw.Stacked {
		t.Errorf("expected %s, but got %s", raw.Stacked, c.ChartType())
	}

	if c.SetChartType("SecretChartType"); c.ChartType() != raw.Stacked {
		t.Errorf("expected %s, but got %s", raw.Stacked, c.ChartType())
	}
}

func TestChart_GetDimByID(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)

	if d := c.GetDimByID("dim0"); d != nil {
		t.Errorf("expected nil, but got %v", d)
	}

	if d := c.GetDimByID("dim1"); d == nil {
		t.Error("expected dimension, but got nil")
	} else {
		if _, ok := toInterface(d).(*dimension); !ok {
			t.Error("expected *dimension type, but got another")
		}
	}
}

func TestChart_GetDimByIndex(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)

	if d := c.GetDimByIndex(-1); d != nil {
		t.Errorf("expected nil, but got %v", d)
	}

	if d := c.GetDimByIndex(1); d == nil {
		t.Error("expected dimension, but got nil")
	} else {
		if _, ok := toInterface(d).(*dimension); !ok {
			t.Error("expected *dimension type, but got another")
		}
	}
}

func TestChart_GetVarByID(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)

	if v := c.GetVarByID("var0"); v != nil {
		t.Errorf("expected nil, but got %v", v)
	}

	if v := c.GetVarByID("var1"); v == nil {
		t.Error("expected dimension, but got nil")
	} else {
		if _, ok := toInterface(v).(*variable); !ok {
			t.Error("expected *variable type, but got another")
		}
	}
}

func TestChart_AddDim(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)
	c.setPush(false)
	c.setObsoleted(true)
	c.FailedUpdates = 1

	d := raw.Dimension{"dim3"}
	if err := c.AddDim(d); err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}

	if c.GetDimByID("dim1") == nil {
		t.Fatal("dimension not added")
	}

	if !c.isPush() {
		t.Error("dimension push flag was not setted to true")
	}

	if c.IsObsoleted() || c.FailedUpdates != 0 {
		t.Error("dimension obsolete flag is not false or failed updates counter not reseted")
	}

	if err := c.AddDim(d); err == nil {
		t.Error("expected error, but got nil")
	}
}

func TestChart_AddVar(t *testing.T) {
	c := Chart{variables: make(map[string]*variable)}
	v := raw.Variable{"var1"}

	if err := c.AddVar(v); err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}

	if err := c.AddVar(v); err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}

	if c.GetVarByID("var1") == nil {
		t.Fatal("variable not added")
	}
}

func TestChart_Refresh(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)
	c.setPush(false)
	c.setObsoleted(true)
	c.setCreated(true)
	c.FailedUpdates = 1

	if c.Refresh(); !c.isPush() || c.IsObsoleted() || c.isCreated() || c.FailedUpdates != 0 {
		t.Error("not all flags were reseted")
	}
}

func TestChart_CanBeUpdated(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)

	if !c.CanBeUpdated(map[string]int64{"dim1": 1}) {
		t.Error("expected true, but got false")
	}

	if c.CanBeUpdated(map[string]int64{"dim3": 1}) {
		t.Error("expected false, but got true")
	}
}

func TestChart_Obsolete(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)

	if c.Obsolete(); !c.IsObsoleted() {
		t.Error("expected true, but got false")
	}

}

func TestChart_Update(t *testing.T) {
	c, _ := newChart(&testRawChart, testBC{}, 1)
	c.setPush(false)
	c.setUpdated(true)

	if ok := c.Update(map[string]int64{"dim3": 1}, 0); ok {
		t.Error("expected false, but got true")
	}

	if c.FailedUpdates != 1 {
		t.Errorf("expected 1, but got %d", c.FailedUpdates)
	}

	if c.isUpdated() {
		t.Error("expected false, but got true")
	}

	if ok := c.Update(map[string]int64{"dim1": 1}, 0); !ok {
		t.Error("expected true, but got false")
	}

	if c.FailedUpdates != 0 {
		t.Errorf("expected 0, but got %d", c.FailedUpdates)
	}

	if !c.isUpdated() {
		t.Error("expected true, but got false")
	}
}
