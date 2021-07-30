package mongo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeChart(t *testing.T) {
	assert.Len(t, serverStatusCharts, 26)
}

func validId(id string) bool {
	return !strings.ContainsAny(id, "- ")
}

func Test_validIds(t *testing.T) {
	for _, chart := range serverStatusCharts {
		t.Run(chart.Title, func(t *testing.T) {
			assert.True(t, validId(chart.ID), fmt.Sprintf("invalid ID: %s", chart.ID))
			assert.True(t, validId(chart.OverID), fmt.Sprintf("invalid OverID: %s", chart.OverID))
			for _, dim := range chart.Dims {
				assert.True(t, validId(dim.ID), fmt.Sprintf("invalid dim ID: %s", dim.ID))
			}
		})
	}
}
