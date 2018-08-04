package charts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testChart = &Chart{
	ID:   "test1",
	Opts: Opts{Title: "Test1", Units: "test"},
	Dims: Dims{
		{ID: "dim1"},
	},
	Vars: Vars{
		{ID: "var1"},
	},
}

func TestChart_Copy(t *testing.T) {
	chart1 := &Chart{
		ID:   "test1",
		Opts: Opts{Title: "Test1"},
		Dims: Dims{
			{ID: "dim1"},
		},
	}

	chart2 := chart1.Copy()

	assert.Equal(t, chart1, chart2)

	chart2.Dims = append(chart2.Dims, &Dim{ID: "dim3"})

	assert.NotEqual(t, chart1, chart2)
}

func TestChart_AddDim(t *testing.T) {
	chart := testChart.Copy()

	chart.AddDim(&Dim{ID: "dim1"})

	assert.Equal(t, len(chart.Dims), 1)

	chart.AddDim(&Dim{ID: "dim2"})

	assert.Equal(t, len(chart.Dims), 2)
}

func TestChart_AddVar(t *testing.T) {
	chart := testChart.Copy()

	chart.AddVar(&Var{ID: "var1"})

	assert.Equal(t, len(chart.Vars), 1)

	chart.AddVar(&Var{ID: "var2"})

	assert.Equal(t, len(chart.Vars), 2)
}

func TestChart_DeleteDimByID(t *testing.T) {
	chart := testChart.Copy()

	chart.DeleteDimByID("dim2")

	assert.Equal(t, len(chart.Dims), 1)

	chart.DeleteDimByID("dim1")

	assert.Equal(t, len(chart.Dims), 0)
}

func TestChart_GetDimByID(t *testing.T) {
	chart := testChart.Copy()
	dim := chart.Dims[0]

	v := chart.GetDimByID("dim1")

	assert.Equal(t, v, dim)
	assert.IsType(t, (*Dim)(nil), v)
	assert.Nil(t, chart.GetDimByID("dim2"))
}

func TestChart_LookupDimByID(t *testing.T) {
	chart := testChart.Copy()
	dim := chart.Dims[0]

	v, ok := chart.LookupDimByID("dim1")

	assert.True(t, ok)
	assert.Equal(t, v, dim)
	assert.IsType(t, (*Dim)(nil), dim)

	v, ok = chart.LookupDimByID("dim2")

	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestChart_Refresh(t *testing.T) {
	testChart.Refresh()
}

func Test_ChartType_Algorithm_Hidden(t *testing.T) {
	assert.Equal(t, Line.String(), "line")
	assert.Equal(t, Area.String(), "area")
	assert.Equal(t, Stacked.String(), "stacked")
	assert.Equal(t, Absolute.String(), "absolute")
	assert.Equal(t, Incremental.String(), "incremental")
	assert.Equal(t, PercentOfAbsolute.String(), "percentage-of-absolute-row")
	assert.Equal(t, PercentOfIncremental.String(), "percentage-of-incremental-row")
	assert.Equal(t, Hidden.String(), "hidden")
	assert.Equal(t, NotHidden.String(), "")
}
