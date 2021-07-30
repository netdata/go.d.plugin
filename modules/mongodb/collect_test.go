package mongo

import (
	"fmt"
	"testing"

	v5_0_0 "github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongo_serverStatusCollect(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.V500ServerStatus), true, &status)
	require.NoError(t, err)

	ms := map[string]int64{}
	for _, chart := range serverStatusCharts {
		for _, dim := range chart.Dims {
			addIfExists(status, dim.ID, ms)
		}
	}

	for _, chart := range serverStatusCharts {
		for _, dim := range chart.Dims {
			_, ok := ms[toID(dim.ID)]
			if contains(v5_0_0.V500ServerStatusMissingValues, toID(dim.ID)) {
				assert.False(t, ok, fmt.Sprintf("value for dim.ID:%s not found", toID(dim.ID)))
			} else {
				assert.True(t, ok, fmt.Sprintf("value for dim.ID:%s not found", toID(dim.ID)))
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
