package raw

import (
	"strings"
	"testing"
)

var testChart = Chart{}

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
	testChart.Options[IdxChartTitle] = title

	if testChart.Title() != title {
		t.Errorf("expected %s, but got %s", title, testChart.Title())
	}
}

func TestChart_Units(t *testing.T) {
	units := "units"
	testChart.Options[IdxChartUnits] = units

	if testChart.Units() != units {
		t.Errorf("expected %s, but got %s", units, testChart.Units())
	}
}

func TestChart_Family(t *testing.T) {
	family := "family"
	testChart.Options[IdxChartFamily] = family

	if testChart.Family() != strings.ToLower(family) {
		t.Errorf("expected %s, but got %s", strings.ToLower(family), testChart.Family())
	}
}

func TestChart_Context(t *testing.T) {
	context := "context"
	testChart.Options[IdxChartContext] = context

	if testChart.Context() != context {
		t.Errorf("expected %s, but got %s", context, testChart.Context())
	}
}

func TestChart_ChartType(t *testing.T) {
	if testChart.ChartType() != defaultChartType {
		t.Errorf("expected %s, but got %s", defaultChartType, testChart.ChartType())
	}

	testChart.Options[IdxChartType] = Stacked

	if testChart.ChartType() != Stacked {
		t.Errorf("expected %s, but got %s", Stacked, testChart.ChartType())
	}
}

func TestChart_OverrideID(t *testing.T) {
	id := "id"
	testChart.Options[IdxChartOverrideID] = id

	if testChart.OverrideID() != id {
		t.Errorf("expected %s, but got %s", id, testChart.OverrideID())
	}
}

func TestChart_SetTitle(t *testing.T) {
	title := "newTitle"

	testChart.SetTitle(title)

	if testChart.Title() != title {
		t.Errorf("expected %s, but got %s", title, testChart.Title())
	}
}

func TestChart_SetUnits(t *testing.T) {
	units := "newUnits"

	testChart.SetUnits(units)

	if testChart.Units() != units {
		t.Errorf("expected %s, but got %s", units, testChart.Units())
	}
}

func TestChart_SetFamily(t *testing.T) {
	family := "newFamily"

	testChart.SetFamily(family)

	if testChart.Family() != "newfamily" {
		t.Errorf("expected %s, but got %s", family, testChart.Family())
	}
}

func TestChart_SetContext(t *testing.T) {
	context := "newContext"

	testChart.SetContext(context)

	if testChart.Context() != context {
		t.Errorf("expected %s, but got %s", context, testChart.Context())
	}
}

func TestChart_SetChartType(t *testing.T) {
	testChart.SetChartType(Area)

	if testChart.ChartType() != Area {
		t.Errorf("expected %s, but got %s", Area, testChart.ChartType())
	}

	testChart.SetChartType("ChartType")

	if testChart.ChartType() != Line {
		t.Errorf("expected %s, but got %s", Line, testChart.ChartType())
	}
}

func TestChart_SetOverrideID(t *testing.T) {
	id := "id"

	testChart.SetOverrideID(id)

	if testChart.OverrideID() != id {
		t.Errorf("expected %s, but got %s", id, testChart.OverrideID())
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

	if ok := c.DeleteDimByID(id); !ok {
		t.Errorf("expected true, but got false")
	}

	if c.DeleteDimByID("id2") {
		t.Errorf("expected false, but got true")
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

	if !c.DeleteDimByIndex(0) {
		t.Errorf("expected true, but got false")
	}

	if c.DeleteDimByIndex(99) {
		t.Errorf("expected false, but got true")
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

	if !c.DeleteVarByID(id) {
		t.Errorf("expected true, but got false")
	}

	if c.DeleteVarByID("id2") {
		t.Errorf("expected false, but got true")
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

	if !chartsEqual(c1, *NewChart(c1.ID, c1.Options, c1.Dimensions...)) {
		t.Error("expected an equal chart, got not equal")
	}
}

func TestChart_Copy(t *testing.T) {
	chart1 := c1
	chart2 := chart1.Copy()

	chart1.SetTitle("New Title")
	chart1.GetDimByIndex(0).SetID("dimID")

	if chart1.Title() == chart2.Title() || chart1.GetDimByIndex(0) == chart2.GetDimByIndex(0) {
		t.Error("chart2 funcion fails")
	}
}

func TestValidChartType(t *testing.T) {
	for _, v := range []string{
		Line,
		Area,
		Stacked,
	} {
		if !ValidChartType(v) {
			t.Fatalf("function returned false for correct chart type")
		}
	}

	for _, v := range []string{
		"",
		"this",
		"is",
		"wrong",
	} {
		if ValidChartType(v) {
			t.Fatalf("function returned true for incorrect chart type")
		}
	}
}
