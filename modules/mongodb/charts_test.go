package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeChart(t *testing.T) {
	assert.Len(t, serverStatusCharts, 26)
}
