package charts

import (
	"testing"
)

var testCharts = &Charts{
	testChart.Copy(),
}

func TestNewCharts(t *testing.T) {
	if _, ok := interface{}(NewCharts()).(*Charts); !ok {
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

func TestCharts_Add(t *testing.T) {
	ch := testCharts.Copy()
	c := testChart.Copy()

	ch.Add(c)

	if len(ch) != 1 {
		t.Errorf("excpected 1 charts, but got %d", len(ch))
	}

	c.ID = "test2"

	ch.Add(c)

	if len(ch) != 2 {
		t.Errorf("excpected 2 charts, but got %d", len(ch))
	}
}

func TestCharts_Delete(t *testing.T) {
	c := testCharts.Copy()

	c.Delete("test1")

	if len(c) != 0 {
		t.Errorf("excpected 0 charts, but got %d", len(c))
	}
}

func TestCharts_Get(t *testing.T) {
	c := testCharts.Copy()

	if c.Get("test1") == nil {
		t.Error("expected not nil")
	}

	if c.Get("test2") != nil {
		t.Error("expected nil")
	}

	if _, ok := interface{}(c.Get("test1")).(*Chart); !ok {
		t.Error("excpected *Chart")
	}
}

func TestCharts_Lookup(t *testing.T) {
	c := testCharts.Copy()

	v, ok := c.Lookup("test1")

	if !ok {
		t.Error("expected true")
	}

	if _, ok := interface{}(v).(*Chart); !ok {
		t.Error("excpected *Chart")
	}

	if _, ok := c.Lookup("test2"); ok {
		t.Error("expected false")
	}
}


func TestCharts_AddAfter(t *testing.T) {
	c1 := testChart.Copy()
	c2 := testChart.Copy()
	c3 := testChart.Copy()
	c2.ID = "test2"
	c3.ID = "test3"

	ch := NewCharts(c1, c2, c3)

	ch.AddAfter("test2", testChart.Copy())

	if len(*ch) != 3 {
		t.Errorf("expected 3, but got %d", len(*ch))
	}

	c4 := testChart.Copy()
	c5 := testChart.Copy()
	c4.ID = "test4"
	c5.ID = "test5"

	ch.AddAfter("test2", c4, c5)

	if len(*ch) != 5 {
		t.Errorf("expected 5, but got %d", len(*ch))
	}

	if (*ch)[2].ID != c4.ID || (*ch)[3].ID != c5.ID {
		t.Error("insertion order wrong")
	}
}

func TestCharts_AddBefore(t *testing.T) {
	c1 := testChart.Copy()
	c2 := testChart.Copy()
	c3 := testChart.Copy()
	c2.ID = "test2"
	c3.ID = "test3"

	ch := NewCharts(c1, c2, c3)

	ch.AddBefore("test2", testChart.Copy())

	if len(*ch) != 3 {
		t.Errorf("expected 3, but got %d", len(*ch))
	}

	c4 := testChart.Copy()
	c5 := testChart.Copy()
	c4.ID = "test4"
	c5.ID = "test5"

	ch.AddBefore("test2", c4, c5)

	if len(*ch) != 5 {
		t.Errorf("expected 5, but got %d", len(*ch))
	}

	if (*ch)[1].ID != c4.ID || (*ch)[2].ID != c5.ID {
		t.Error("insertion order wrong")
	}
}