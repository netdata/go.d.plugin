package cooked

import (
	"testing"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

func TestNewCharts(t *testing.T) {
	if v := NewCharts(testBC{}); v == nil {
		t.Error("expected charts, not nil")
	} else {
		if _, ok := toInterface(v).(*charts); !ok {
			t.Error("expected *charts type, but got another")
		}
	}
}

func TestCharts_AddOne(t *testing.T) {
	ch := NewCharts(testBC{})

	if err := ch.AddOne(&testRawChart); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if err := ch.AddOne(&testRawChart); err != nil {
		t.Errorf("expected nil, but got %s", err)
	}

	if len(ch.charts) != 1 || ch.GetChartByID("chart1") == nil {
		t.Fatal("chart not added")
	}

	if ch.priority != initPriority+1 {
		t.Errorf("expected %d, but got %d", initPriority+1, ch.priority)
	}
}

func TestCharts_AddMany(t *testing.T) {
	r1, r2, rc := testRawChart, testRawChart, raw.Charts{}
	r2.ID = "chart2"
	rc.AddChart(&r1)
	rc.AddChart(&r2)

	ch := NewCharts(testBC{})
	n := ch.AddMany(&rc)

	switch {
	case
		n != 2,
		len(ch.charts) != 2,
		ch.GetChartByID(r1.ID) == nil,
		ch.GetChartByID(r2.ID) == nil:
		t.Error("charts not added")
	}
}

func TestCharts_GetChartByID(t *testing.T) {
	ch := NewCharts(testBC{})
	ch.AddOne(&testRawChart)

	if ch.GetChartByID(testRawChart.ID) == nil {
		t.Error("expected chart, but got nil")
	}

	if _, ok := toInterface(ch.GetChartByID(testRawChart.ID)).(*Chart); !ok {
		t.Error("expected *Chart type, but got another")
	}
}

func toInterface(i interface{}) interface{} {
	return i
}
