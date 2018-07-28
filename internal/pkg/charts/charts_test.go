package charts

import "testing"

var testCharts = &Charts{
	testChart.Copy(),
}

func TestNew(t *testing.T) {
	if _, ok := interface{}(New()).(*Charts); !ok {
		t.Error("excpected *Charts")
	}
}

func TestCharts_Copy(t *testing.T) {
	c := testCharts.Copy()

	c[0].ID = "test2"
	c[0].Dims = append(c[0].Dims, &Dim{ID: "dim2"})
	c = append(c, testChart.Copy())

	if len(c) == len(*testCharts) || c[0].ID == (*testCharts)[0].ID || len(c[0].Dims) == len((*testCharts)[0].Dims) {
		t.Error("expected full copy, but got partial")
	}

}

func TestCharts_AddChart(t *testing.T) {
	c := testCharts.Copy()

	c.AddChart(testChart.Copy())

	if len(c) != 2 {
		t.Errorf("excpected 2 charts, but got %d", len(c))
	}
}

func TestCharts_DeleteChart(t *testing.T) {
	c := testCharts.Copy()

	c.DeleteChart("test1")

	if len(c) != 0 {
		t.Errorf("excpected 0 charts, but got %d", len(c))
	}
}

func TestCharts_GetChart(t *testing.T) {
	c := testCharts.Copy()

	if c.GetChart("test1") == nil {
		t.Error("expected not nil")
	}

	if c.GetChart("test2") != nil {
		t.Error("expected nil")
	}

	if _, ok := interface{}(c.GetChart("test1")).(*Chart); !ok {
		t.Error("excpected *Chart")
	}
}

func TestCharts_LookupChart(t *testing.T) {
	c := testCharts.Copy()

	v, ok := c.LookupChart("test1")

	if !ok {
		t.Error("expected true")
	}

	if _, ok := interface{}(v).(*Chart); !ok {
		t.Error("excpected *Chart")
	}

	if _, ok := c.LookupChart("test2"); ok {
		t.Error("expected false")
	}
}
