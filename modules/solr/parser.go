package solr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type count struct {
	Count int64
}

type value struct {
	Value int64
}

type common struct {
	Count         int64
	MeanRate      int64
	MeanRate1min  int64 `json:"1minRate"`
	MeanRate5min  int64 `json:"5minRate"`
	MeanRate15min int64 `json:"15minRate"`
}

type requestTimes struct {
	Count         int64
	MeanRate      int64
	MeanRate1min  int64 `json:"1minRate"`
	MeanRate5min  int64 `json:"5minRate"`
	MeanRate15min int64 `json:"15minRate"`
	MinMS         int64 `json:"min_ms"`
	MaxMS         int64 `json:"max_ms"`
	MeanMS        int64 `json:"mean_ms"`
	MedianMS      int64 `json:"median_ms"`
	StdDevMS      int64 `json:"stddev_ms"`
	P75MS         int64 `json:"p75_ms"`
	P95MS         int64 `json:"p95_ms"`
	P99MS         int64 `json:"p99_ms"`
	P999MS        int64 `json:"p999_ms"`
}

type metrics struct {
	Metrics map[string]map[string]json.RawMessage
}

type V6Parser struct {
	count        count
	value        value
	common       common
	requestTimes requestTimes

	parsed map[string]int64
}

func (v *V6Parser) parse(resp *http.Response) (map[string]int64, error) {
	var m metrics

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	for core, data := range m.Metrics {
		if err := v.parseCore(core[10:], data); err != nil {
			return nil, err
		}
	}

	return v.parsed, nil
}

func (v *V6Parser) parseCore(core string, data map[string]json.RawMessage) error {
	for metric, stats := range data {
		parts := strings.Split(metric, ".")

		if len(parts) != 3 {
			continue
		}

		typ, handler, stat := strings.ToLower(parts[0]), parts[1], parts[2]

		if handler == "updateHandler" {
			//switch stat {
			//case "adds", "autoCommits", "deletesById", "deletesByQuery", "docsPending", "errors", "softAutoCommits":
			//	if err := json.Unmarshal(stats, &v.value); err != nil {
			//		return err
			//	}
			//default:
			//	if err := json.Unmarshal(stats, &v.common); err != nil {
			//		return err
			//	}
			//}
			continue
		}

		switch stat {
		case "clientErrors", "errors", "serverErrors", "timeouts":
			if err := json.Unmarshal(stats, &v.common); err != nil {
				return err
			}
			v.parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += v.common.Count
		case "requests", "totalTime":
			if err := json.Unmarshal(stats, &v.count); err != nil {
				return err
			}
			v.parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += v.count.Count
		case "requestTimes":
			if err := json.Unmarshal(stats, &v.requestTimes); err != nil {
				return err
			}
			v.parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += v.requestTimes.Count
			v.parsed[fmt.Sprintf("%s_%s_%s_mean_ms", core, typ, stat)] += v.requestTimes.MeanMS
			v.parsed[fmt.Sprintf("%s_%s_%s_median_ms", core, typ, stat)] += v.requestTimes.MedianMS
			v.parsed[fmt.Sprintf("%s_%s_%s_p75_ms", core, typ, stat)] += v.requestTimes.P75MS
			v.parsed[fmt.Sprintf("%s_%s_%s_p95_ms", core, typ, stat)] += v.requestTimes.P95MS
			v.parsed[fmt.Sprintf("%s_%s_%s_p99_ms", core, typ, stat)] += v.requestTimes.P99MS
			v.parsed[fmt.Sprintf("%s_%s_%s_p999_ms", core, typ, stat)] += v.requestTimes.P999MS
		}
	}

	return nil
}
