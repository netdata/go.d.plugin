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

func TestNewCharts(t *testing.T) {
	charts := NewCharts(createTestChart("1"), createTestChart("2"), createTestChart("1"))
	assert.IsType(t, (*Charts)(nil), charts)
	assert.True(t, len(*charts) == 2)
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
