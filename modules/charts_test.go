package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestChart(id string) *Chart {
	return &Chart{
		ID:    id,
		Title: "Title",
		Units: "units",
		Fam:   "family",
		Ctx:   "context",
		Type:  Line,
		Dims: Dims{
			{ID: "dim1", Algo: Absolute},
		},
		Vars: Vars{
			{ID: "var1", Value: 1},
		},
	}
}

func TestDimAlgo_String(t *testing.T) {
	assert.Equal(t, Line.String(), string(Line))
	assert.Equal(t, Line.String(), "line")
	assert.Equal(t, Area.String(), string(Area))
	assert.Equal(t, Area.String(), "area")
	assert.Equal(t, Stacked.String(), string(Stacked))
	assert.Equal(t, Stacked.String(), "stacked")

	assert.Equal(t, dimAlgo("wrong").String(), "")
}

func TestChartType_String(t *testing.T) {
	assert.Equal(t, Absolute.String(), string(Absolute))
	assert.Equal(t, Absolute.String(), "absolute")
	assert.Equal(t, Incremental.String(), string(Incremental))
	assert.Equal(t, Incremental.String(), "incremental")
	assert.Equal(t, PercentOfAbsolute.String(), string(PercentOfAbsolute))
	assert.Equal(t, PercentOfAbsolute.String(), "percentage-of-absolute-row")
	assert.Equal(t, PercentOfIncremental.String(), string(PercentOfIncremental))
	assert.Equal(t, PercentOfIncremental.String(), "percentage-of-incremental-row")

	assert.Equal(t, chartType("wrong").String(), "")
}

func TestDimHidden_String(t *testing.T) {
	assert.Equal(t, dimHidden(false).String(), "")
	assert.Equal(t, dimHidden(true).String(), "hidden")
}

func TestDimDivMul_String(t *testing.T) {
	assert.Equal(t, dimDivMul(0).String(), "")
	assert.Equal(t, dimDivMul(1).String(), "1")
	assert.Equal(t, dimDivMul(-1).String(), "-1")
}

func TestOpts_String(t *testing.T) {
	assert.Equal(t, Opts{}.String(), "")
	assert.Equal(t, Opts{
		Obsolete:   true,
		Detail:     true,
		StoreFirst: true,
		Hidden:     true,
	}.String(), "obsolete detail store_first hidden")
	assert.Equal(t, Opts{
		Obsolete:   true,
		Detail:     false,
		StoreFirst: false,
		Hidden:     true,
	}.String(), "obsolete hidden")
}

func TestCharts_Copy(t *testing.T) {
	orig := &Charts{
		createTestChart("1"),
		createTestChart("2"),
	}
	copied := orig.Copy()

	assert.False(t, orig == copied, "copied charts points to the same address")
	require.Len(t, *orig, len(*copied))

	for idx := range *orig {
		compareCharts(t, (*orig)[idx], (*copied)[idx])

	}
}

func TestChart_Copy(t *testing.T) {
	orig := createTestChart("1")

	compareCharts(t, orig, orig.Copy())
}

func TestCharts_Add(t *testing.T) {
	charts := Charts{}
	chart1 := createTestChart("1")
	chart2 := createTestChart("2")
	chart3 := createTestChart("")
	charts.Add(
		chart1,
		chart2,
		chart1,
		chart3,
	)
	assert.Len(t, charts, 2)
	assert.True(t, charts[0] == chart1)
	assert.True(t, charts[1] == chart2)
}

func TestCharts_Get(t *testing.T) {
	chart := createTestChart("1")
	charts := Charts{
		chart,
	}
	assert.Nil(t, charts.Get("2"))
	assert.IsType(t, (*Chart)(nil), charts.Get("1"))
	assert.True(t, chart == charts.Get("1"))
}

func TestCharts_Has(t *testing.T) {
	chart := createTestChart("1")
	charts := &Charts{
		chart,
	}

	assert.True(t, charts.Has("1"))
	assert.False(t, charts.Has("2"))
}

func TestCharts_Remove(t *testing.T) {
	chart := createTestChart("1")
	charts := &Charts{
		chart,
	}

	assert.False(t, charts.Remove("2"))
	assert.True(t, charts.Remove("1"))
	assert.Len(t, *charts, 0)
}

func TestChart_AddDim(t *testing.T) {
	chart := createTestChart("1")
	dim := &Dim{ID: "dim2"}

	assert.True(t, chart.AddDim(dim))
	assert.False(t, chart.AddDim(dim))
	assert.Len(t, chart.Dims, 2)
}

func TestChart_AddVar(t *testing.T) {
	chart := createTestChart("1")
	variable := &Var{ID: "var2"}

	assert.True(t, chart.AddVar(variable))
	assert.False(t, chart.AddVar(variable))
	assert.Len(t, chart.Vars, 2)
}

func TestChart_RemoveDim(t *testing.T) {
	chart := createTestChart("1")

	assert.False(t, chart.RemoveDim("dim2"))
	assert.True(t, chart.RemoveDim("dim1"))
	assert.Len(t, chart.Dims, 0)
}

func TestChart_HasDim(t *testing.T) {
	chart := createTestChart("1")

	assert.False(t, chart.HasDim("dim2"))
	assert.True(t, chart.HasDim("dim1"))
}

func TestChart_MarkPush(t *testing.T) {
	chart := createTestChart("1")

	assert.False(t, chart.pushed)
	chart.pushed = true
	chart.MarkPush()
	assert.False(t, chart.pushed)
}

func compareCharts(t *testing.T, orig, copied *Chart) {
	assert.False(t, orig == copied, "copied charts points to the same address")

	require.Len(t, orig.Dims, len(copied.Dims))
	require.Len(t, orig.Vars, len(copied.Vars))

	for idx := range (*orig).Dims {
		assert.False(t, orig.Dims[idx] == copied.Dims[idx], "copied dim points to the same address")
		assert.Equal(t, orig.Dims[idx], copied.Dims[idx], "copied dim isn't equal to orig")
	}

	for idx := range (*orig).Vars {
		assert.False(t, orig.Vars[idx] == copied.Vars[idx], "copied var points to the same address")
		assert.Equal(t, orig.Vars[idx], copied.Vars[idx], "copied var isn't equal to orig")
	}
}
