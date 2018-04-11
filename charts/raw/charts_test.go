package raw

import (
	"testing"
)

var (
	c1 = Chart{
		ID:      "id1",
		Options: Options{"Title", "Units", "Family", "Context", "Type", "OverrideID"},
		Dimensions: Dimensions{
			Dimension{"1"},
			Dimension{"2"},
		},
		Variables: Variables{
			Variable{"11"},
			Variable{"22"},
		}}
	c2 = Chart{
		ID:      "",
		Options: Options{"Title", "Units", "Family", "Context", "Type", "OverrideID"},
		Dimensions: Dimensions{
			Dimension{"3"},
			Dimension{"4"},
		}}
)

func chartsEqual(c1, c2 Chart) bool {
	var dimEq, varEq = true, true
	for _, d := range c1.Dimensions {
		if c2.GetDimByID(d.ID()) == nil || d.ID() != c2.GetDimByID(d.ID()).ID() {
			dimEq = false
			break
		}
	}
	for _, v := range c1.Variables {
		if c2.GetVarByID(v.ID()) == nil || v.ID() != c2.GetVarByID(v.ID()).ID() {
			varEq = false
			break
		}
	}
	switch {
	case
		c1.ID != c2.ID,
		c1.Options != c2.Options,
		len(c1.Dimensions) != len(c2.Dimensions),
		len(c1.Variables) != len(c2.Variables),
		!dimEq,
		!varEq:
		return false
	default:
		return true
	}
}

func TestCharts_GetChartByID(t *testing.T) {
	ch := Charts{Order: Order{c1.ID}, Definitions: Definitions{c1}}

	if ch.GetChartByID("id3") != nil {
		t.Errorf("expected nil, but got %v", ch.GetChartByID("id3"))
	}

	if c := ch.GetChartByID(c1.ID); c == nil {
		t.Fatalf("expected %s chart, but got nil", c1.ID)
	}

	if c := ch.GetChartByID(c1.ID); !chartsEqual(c1, *c) {
		t.Error("expected an equal chart, but got not equal")
	}
}

func TestCharts_GetChartByIndex(t *testing.T) {
	ch := Charts{Order: Order{c1.ID}, Definitions: Definitions{c1}}

	if ch.GetChartByIndex(10) != nil {
		t.Errorf("expected nil, but got %v", ch.GetChartByIndex(10))
	}

	if c := ch.GetChartByIndex(0); c == nil {
		t.Fatalf("expected %s chart, but got nil", c1.ID)
	}

	if c := ch.GetChartByIndex(0); !chartsEqual(c1, *c) {
		t.Error("expected an equal chart, but got not equal")
	}

}

func TestCharts_DeleteChartByID(t *testing.T) {
	ch := Charts{Order: Order{c1.ID}, Definitions: Definitions{c1}}

	lenOrd, lenDef := len(ch.Order), len(ch.Definitions)

	if ch.DeleteChartByID("id3") == nil {
		t.Error("expected error, but got nil")
	}

	if err := ch.DeleteChartByID(c1.ID); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if ch.GetChartByID(c1.ID) != nil || len(ch.Order) >= lenOrd || len(ch.Definitions) >= lenDef {
		t.Errorf("chart %s was not deleted", c1.ID)
	}
}

func TestCharts_DeleteChartByIndex(t *testing.T) {
	ch := Charts{Order: Order{c1.ID}, Definitions: Definitions{c1}}

	lenOrd, lenDef := len(ch.Order), len(ch.Definitions)

	if ch.DeleteChartByIndex(99) == nil {
		t.Error("expected error, but got nil")
	}

	if err := ch.DeleteChartByIndex(0); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if ch.GetChartByID(c1.ID) != nil || len(ch.Order) >= lenOrd || len(ch.Definitions) >= lenDef {
		t.Errorf("chart %s was not deleted", c1.ID)
	}
}

func TestCharts_AddChart(t *testing.T) {
	ch := Charts{}

	if ch.AddChart(c2, true) == nil {
		t.Error("expected error, but got nil")
	}

	if err := ch.AddChart(c1, true); err != nil {
		t.Fatalf("expected nil, but got %d", err)
	}

	if ch.GetChartByID(c1.ID) == nil {
		t.Fatalf("chart %s was not added", c1.ID)
	}

	if ch.AddChart(c1, true) == nil {
		t.Fatal("expected error, but got nil")
	}

	if len(ch.Order) == 0 {
		t.Fatalf("chart %s was not added to Order", c1.ID)
	}

	ch.DeleteChartByID(c1.ID)

	ch.AddChart(c1, false)
	if len(ch.Order) != 0 {
		t.Fatalf("chart %s was added to Order", c1.ID)
	}
}

func TestCharts_Copy(t *testing.T) {
	ch1 := Charts{Order: Order{}, Definitions: Definitions{}}
	ch2 := ch1.Copy()
	ch2.Order.Append(c1.ID)

	if len(ch1.Order) != 0 {
		t.Fatal("charts Order was modified by copy")
	}
	ch2.AddChart(c1, false)

	if len(ch1.Definitions) != 0 {
		t.Fatal("charts Definitions was modified by copy")
	}
	ch1 = *(ch2.Copy())

	ch2.GetChartByID(c1.ID).AddDim(Dimension{"newDim"})

	if len(ch1.GetChartByID(c1.ID).Dimensions) == len(ch2.GetChartByID(c1.ID).Dimensions) {
		t.Error("expected different number of dimensions, got equal")
	}

	ch2.GetChartByID(c1.ID).AddVar(Variable{"newVar"})

	if len(ch1.GetChartByID(c1.ID).Variables) == len(ch2.GetChartByID(c1.ID).Variables) {
		t.Error("expected different number of variables, got equal")
	}
}
