package raw

import (
	"strings"
	"testing"
)

func TestChart_IsValid(t *testing.T) {
	if (&Chart{}).IsValid() == nil {
		t.Error("expected error, but got nil")
	}
	if (&Chart{
		ID: "id",
		Options: Options{
			IdxChartTitle: "title",
			IdxChartUnits: "units",
		}}).IsValid() == nil {
		t.Error("expected error, but got nil")
	}
	if (&Chart{
		ID: "id",
		Options: Options{
			IdxChartTitle:  "title",
			IdxChartFamily: "family",
		}}).IsValid() == nil {
		t.Error("expected error, but got nil")
	}
	if (&Chart{
		ID: "id",
		Options: Options{
			IdxChartUnits:  "units",
			IdxChartFamily: "family",
		}}).IsValid() == nil {
		t.Error("expected error, but got nil")
	}
	if (&Chart{
		Options: Options{
			IdxChartTitle:  "title",
			IdxChartUnits:  "units",
			IdxChartFamily: "family",
		}}).IsValid() == nil {
		t.Error("expected error, but got nil")
	}
	if (&Chart{
		ID: "id",
		Options: Options{
			IdxChartTitle:  "title",
			IdxChartUnits:  "units",
			IdxChartFamily: "family",
		}}).IsValid() != nil {
		t.Error("expected nil, but got error")
	}
}

func TestChart_Title(t *testing.T) {
	title := "title"
	c := Chart{
		Options: Options{
			IdxChartTitle: title,
		}}

	if c.Title() != title {
		t.Errorf("expected %s, but got %s", title, c.Title())
	}
}

func TestChart_Units(t *testing.T) {
	units := "units"
	c := Chart{
		Options: Options{
			IdxChartUnits: units,
		}}

	if c.Units() != units {
		t.Errorf("expected %s, but got %s", units, c.Units())
	}
}

func TestChart_Family(t *testing.T) {
	family := "FAMILY"
	c := Chart{
		Options: Options{
			IdxChartFamily: family,
		}}

	if c.Family() != strings.ToLower(family) {
		t.Errorf("expected %s, but got %s", strings.ToLower(family), c.Family())
	}
}

func TestChart_Context(t *testing.T) {
	context := "context"
	c := Chart{
		Options: Options{
			IdxChartContext: context,
		}}

	if c.Context() != context {
		t.Errorf("expected %s, but got %s", context, c.Context())
	}
}

func TestChart_ChartType(t *testing.T) {
	c := Chart{}

	if c.ChartType() != defaultChartType {
		t.Errorf("expected %s, but got %s", defaultChartType, c.ChartType())
	}

	c = Chart{
		Options: Options{
			IdxChartType: Stacked,
		}}

	if c.ChartType() != Stacked {
		t.Errorf("expected %s, but got %s", Stacked, c.ChartType())
	}
}

func TestChart_OverrideID(t *testing.T) {
	overrideID := "id"
	c := Chart{
		Options: Options{
			IdxChartOverrideID: overrideID,
		}}

	if c.OverrideID() != overrideID {
		t.Errorf("expected %s, but got %s", overrideID, c.OverrideID())
	}
}

func TestChart_SetTitle(t *testing.T) {
	c := Chart{}
	newTitle := "title"

	c.SetTitle(newTitle)
	if c.Title() != newTitle {
		t.Errorf("expected %s, but got %s", newTitle, c.Title())
	}
}

func TestChart_SetUnits(t *testing.T) {
	c := Chart{}
	newUnits := "units"

	c.SetUnits(newUnits)
	if c.Units() != newUnits {
		t.Errorf("expected %s, but got %s", newUnits, c.Units())
	}
}

func TestChart_SetFamily(t *testing.T) {
	c := Chart{}
	newFamily := "family"

	c.SetFamily(newFamily)
	if c.Family() != newFamily {
		t.Errorf("expected %s, but got %s", newFamily, c.Family())
	}
}

func TestChart_SetContext(t *testing.T) {
	c := Chart{}
	newContext := "context"

	c.SetContext(newContext)
	if c.Context() != newContext {
		t.Errorf("expected %s, but got %s", newContext, c.Context())
	}
}

func TestChart_SetChartType(t *testing.T) {
	c := Chart{}

	c.SetChartType(Area)
	if c.ChartType() != Area {
		t.Errorf("expected %s, but got %s", Area, c.ChartType())
	}

	c = Chart{}
	c.SetChartType("SecretChartType")
	if c.ChartType() != Line {
		t.Errorf("expected %s, but got %s", Line, c.ChartType())
	}
}

func TestChart_SetOverrideID(t *testing.T) {
	c := Chart{}
	newOverrideID := "id"

	c.SetOverrideID(newOverrideID)
	if c.OverrideID() != newOverrideID {
		t.Errorf("expected %s, but got %s", newOverrideID, c.OverrideID())
	}
}

func TestChart_GetDimByID(t *testing.T) {
	id := "id"
	d := Dimension{id}
	c := Chart{Dimensions: Dimensions{d}}

	if c.GetDimByID(id) == nil {
		t.Fatal("expected dimension, but got nil")
	}

	if c.GetDimByID("id2") != nil {
		t.Fatalf("expected nil, but got %v", c.GetDimByID("id2"))
	}

	if c.GetDimByID(id).ID() != id {
		t.Errorf("expected %s, but got %s", id, c.GetDimByID(id).ID())
	}
}

func TestChart_GetDimByIndex(t *testing.T) {
	id1, id2 := "id1", "id2"
	d1, d2 := Dimension{id1}, Dimension{id2}
	c := Chart{Dimensions: Dimensions{d1, d2}}

	if c.GetDimByIndex(0) == nil || c.GetDimByIndex(1) == nil {
		t.Fatal("expected dimension, but got nil")
	}

	if c.GetDimByIndex(2) != nil {
		t.Fatalf("expected nil, but got %v", c.GetDimByIndex(2))
	}

	if c.GetDimByIndex(0).ID() != id1 {
		t.Errorf("expected %s, but got %s", id1, c.GetDimByIndex(0).ID())
	}

	if c.GetDimByIndex(1).ID() != id2 {
		t.Errorf("expected %s, but got %s", id2, c.GetDimByIndex(1).ID())
	}
}

func TestChart_GetVarByID(t *testing.T) {
	id := "id"
	v := Variable{id}
	c := Chart{Variables: Variables{v}}

	if c.GetVarByID(id) == nil {
		t.Fatal("expected variable, but got nil")
	}

	if c.GetVarByID("id2") != nil {
		t.Fatalf("expected nil, but got %v", c.GetVarByID("id2"))
	}

	if c.GetVarByID(id).ID() != id {
		t.Errorf("expected %s, but got %s", id, c.GetVarByID(id).ID())
	}
}

func TestChart_DeleteDimByID(t *testing.T) {
	id := "id"
	d := Dimension{id}
	c := Chart{Dimensions: Dimensions{d}}

	l := len(c.Dimensions)

	if err := c.DeleteDimByID(id); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if err := c.DeleteDimByID("id2"); err == nil {
		t.Error("expected error, but got nil")
	}

	if !(len(c.Dimensions) < l) {
		t.Errorf("start length %d, length after delete %d", l, len(c.Dimensions))
	}
}

func TestChart_DeleteDimByIndex(t *testing.T) {
	id := "id"
	d := Dimension{id}
	c := Chart{Dimensions: Dimensions{d}}

	l := len(c.Dimensions)

	if err := c.DeleteDimByIndex(0); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if err := c.DeleteDimByIndex(99); err == nil {
		t.Error("expected error, but got nil")
	}

	if !(len(c.Dimensions) < l) {
		t.Errorf("start length %d, length after delete %d", l, len(c.Dimensions))
	}
}

func TestChart_DeleteVarByID(t *testing.T) {
	id := "id"
	v := Variable{id}
	c := Chart{Variables: Variables{v}}

	l := len(c.Variables)

	if err := c.DeleteVarByID(id); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if err := c.DeleteVarByID("id2"); err == nil {
		t.Error("expected error, but got nil")
	}

	if !(len(c.Variables) < l) {
		t.Errorf("start length %d, length after delete %d", l, len(c.Variables))
	}
}

func TestChart_AddDim(t *testing.T) {
	id := "id"
	d := Dimension{id}
	c := Chart{}

	if c.AddDim(Dimension{}) == nil {
		t.Error("expected error, got nil")
	}

	if err := c.AddDim(d); err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if err := c.AddDim(d); err == nil {
		t.Error("expected error, got nil")
	}

	if len(c.Dimensions) != 1 {
		t.Errorf("expected dimension length 1, got %d", len(c.Dimensions))
	}
}

func TestChart_AddVar(t *testing.T) {
	id := "id"
	v := Variable{id}
	c := Chart{}

	if c.AddVar(Variable{}) == nil {
		t.Error("expected error, got nil")
	}

	if err := c.AddVar(v); err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if err := c.AddVar(v); err == nil {
		t.Error("expected error, got nil")
	}

	if len(c.Variables) != 1 {
		t.Errorf("expected dimension length 1, got %d", len(c.Variables))
	}
}

func TestNewChart(t *testing.T) {
	c1 := Chart{
		ID:      "id",
		Options: Options{"Title", "Units", "Family", "Context", "Type", "OverrideID"},
		Dimensions: Dimensions{
			Dimension{"1"},
			Dimension{"2"},
		}}

	if !chartsEqual(c1, NewChart(c1.ID, c1.Options, c1.Dimensions...)) {
		t.Error("expected an equal chart, got not equal")
	}
}

func TestChart_Copy(t *testing.T) {
	chart := c1
	chartCopy := chart.Copy()

	chart.SetTitle("New Title")
	chart.GetDimByIndex(0).SetID(999)

	if chart.Title() == chartCopy.Title() || chart.GetDimByIndex(0) == chartCopy.GetDimByIndex(0) {
		t.Error("copy funcion fails")
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
