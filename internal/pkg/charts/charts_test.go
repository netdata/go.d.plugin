package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCharts = &Charts{
	testChart.Copy(),
}

func TestNewCharts(t *testing.T) {
	assert.IsType(t, (*Charts)(nil), NewCharts())
}

func TestCharts_Copy(t *testing.T) {
	charts := testCharts.Copy()

	charts[0].ID = "test2"
	charts[0].Dims = append(charts[0].Dims, &Dim{ID: "dim2"})
	charts = append(charts, testChart.Copy())

	assert.NotEqual(t, len(charts), len(*testCharts))
	assert.NotEqual(t, len(charts[0].Dims), len((*testCharts)[0].Dims))
	assert.NotEqual(t, charts[0].ID, (*testCharts)[0].ID)
}

func TestCharts_Add(t *testing.T) {
	charts := testCharts.Copy()
	chart := testChart.Copy()

	charts.Add(chart)

	assert.Equal(t, len(charts), 1)

	chart.ID = "test2"

	charts.Add(chart)

	assert.Equal(t, len(charts), 2)
}

func TestCharts_Delete(t *testing.T) {
	charts := testCharts.Copy()

	charts.Delete("test99")

	assert.Equal(t, len(charts), 1)

	charts.Delete("test1")

	assert.Equal(t, len(charts), 0)
}

func TestCharts_Get(t *testing.T) {
	charts := testCharts.Copy()

	assert.NotNil(t, charts.Get("test1"))

	assert.Nil(t, charts.Get("test2"))

	assert.IsType(t, (*Chart)(nil), charts.Get("test1"))
}

func TestCharts_Lookup(t *testing.T) {
	charts := testCharts.Copy()

	v, ok := charts.Lookup("test1")

	assert.True(t, ok)
	assert.IsType(t, (*Chart)(nil), v)

	_, ok = charts.Lookup("test2")

	assert.False(t, ok)
}


func TestCharts_AddAfter(t *testing.T) {
	chart1 := testChart.Copy()
	chart2 := testChart.Copy()
	chart3 := testChart.Copy()
	chart2.ID = "test2"
	chart3.ID = "test3"

	charts := NewCharts(chart1, chart2, chart3)

	charts.AddAfter("test2", testChart.Copy())

	assert.Equal(t, len(*charts), 3)

	chart4 := testChart.Copy()
	chart5 := testChart.Copy()
	chart4.ID = "test4"
	chart5.ID = "test5"

	charts.AddAfter("test2", chart4, chart5)

	assert.Equal(t, len(*charts), 5)

	assert.Equal(t, (*charts)[2], chart4)
	assert.Equal(t, (*charts)[3], chart5)
}

func TestCharts_AddBefore(t *testing.T) {
	chart1 := testChart.Copy()
	chart2 := testChart.Copy()
	chart3 := testChart.Copy()
	chart2.ID = "test2"
	chart3.ID = "test3"

	charts := NewCharts(chart1, chart2, chart3)

	charts.AddBefore("test2", testChart.Copy())

	assert.Equal(t, len(*charts), 3)

	chart4 := testChart.Copy()
	chart5 := testChart.Copy()
	chart4.ID = "test4"
	chart5.ID = "test5"

	charts.AddBefore("test2", chart4, chart5)

	assert.Equal(t, len(*charts), 5)

	assert.Equal(t, (*charts)[1], chart4)
	assert.Equal(t, (*charts)[2], chart5)
}