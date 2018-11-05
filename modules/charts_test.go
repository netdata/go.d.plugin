package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, dimAlgo("wrong").String(), "")
	assert.Equal(t, Line.String(), string(Line))
	assert.Equal(t, Area.String(), string(Area))
	assert.Equal(t, Stacked.String(), string(Stacked))
}

func TestChartType_String(t *testing.T) {
	assert.Equal(t, chartType("wrong").String(), "")
	assert.Equal(t, Absolute.String(), string(Absolute))
	assert.Equal(t, Incremental.String(), string(Incremental))
	assert.Equal(t, PercentOfAbsolute.String(), string(PercentOfAbsolute))
	assert.Equal(t, PercentOfIncremental.String(), string(PercentOfIncremental))
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
	chart := createTestChart("")
	chartCopy := chart.Copy()
	assert.False(t, chart == chartCopy, "chart points to the same address")

	for idx := range chart.Dims {
		assert.False(t, chart.Dims[idx] == chartCopy.Dims[idx], "chart dimension points to the same address")
	}

	for idx := range chart.Vars {
		assert.False(t, chart.Vars[idx] == chartCopy.Vars[idx], "char var points to the same address")
	}
}

func TestNewCharts(t *testing.T) {
	charts := NewCharts(
		createTestChart("1"),
		createTestChart("2"),
		createTestChart("1"),
		createTestChart(""),
	)
	assert.IsType(t, (*Charts)(nil), charts)
	assert.Len(t, *charts, 2)
}

func TestCharts_Add(t *testing.T) {
	charts := new(Charts)
	chart1 := createTestChart("1")
	chart2 := createTestChart("2")
	chart3 := createTestChart("")
	charts.Add(
		chart1,
		chart2,
		chart1,
		chart3,
	)
	assert.Len(t, *charts, 2)
	assert.True(t, (*charts)[0] == chart1)
	assert.True(t, (*charts)[1] == chart2)
}

func TestCharts_Get(t *testing.T) {
	chart := createTestChart("1")
	charts := &Charts{
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
