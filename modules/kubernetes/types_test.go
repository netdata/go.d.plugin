package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

var testKubernetesSumStatsData, _ = ioutil.ReadFile("testdata/stats_summary.json")

func BenchmarkParser_JSON(b *testing.B) {
	var sum Summary
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(testKubernetesSumStatsData, &sum)
	}
}

func BenchmarkParser_EasyJSON(b *testing.B) {
	var sum Summary
	for i := 0; i < b.N; i++ {
		_ = easyjson.Unmarshal(testKubernetesSumStatsData, &sum)
	}
}

func TestUnmarshalEquality(t *testing.T) {
	var sumJSON, sumEasyJSON Summary

	assert.NoError(t, json.Unmarshal(testKubernetesSumStatsData, &sumJSON))
	assert.NoError(t, easyjson.Unmarshal(testKubernetesSumStatsData, &sumEasyJSON))
	assert.Equal(t, sumJSON, sumEasyJSON)
}
