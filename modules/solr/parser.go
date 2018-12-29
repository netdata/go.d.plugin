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
	Count        int64
	MeanRate     float64 `json:"meanRate"`
	MinRate1min  float64 `json:"1minRate"`
	MinRate5min  float64 `json:"5minRate"`
	MinRate15min float64 `json:"15minRate"`
}

type requestTimes struct {
	Count        int64
	MeanRate     float64 `json:"meanRate"`
	MinRate1min  float64 `json:"1minRate"`
	MinRate5min  float64 `json:"5minRate"`
	MinRate15min float64 `json:"15minRate"`
	MinMS        float64 `json:"min_ms"`
	MaxMS        float64 `json:"max_ms"`
	MeanMS       float64 `json:"mean_ms"`
	MedianMS     float64 `json:"median_ms"`
	StdDevMS     float64 `json:"stddev_ms"`
	P75MS        float64 `json:"p75_ms"`
	P95MS        float64 `json:"p95_ms"`
	P99MS        float64 `json:"p99_ms"`
	P999MS       float64 `json:"p999_ms"`
}

type coresMetrics struct {
	Metrics map[string]map[string]json.RawMessage
}

type parser struct {
	simpleCount  int64
	count        count
	value        value
	common       common
	requestTimes requestTimes

	version float64
}

func (v *parser) parse(resp *http.Response) (map[string]int64, error) {
	var m coresMetrics
	parsed := make(map[string]int64)

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	for core, data := range m.Metrics {
		if err := v.parseCore(core[10:], data, parsed); err != nil {
			return nil, err
		}
	}

	return parsed, nil
}

func (p *parser) parseCore(core string, data map[string]json.RawMessage, parsed map[string]int64) error {
	for metric, stats := range data {
		parts := strings.Split(metric, ".")

		if len(parts) != 3 {
			continue
		}

		typ, handler, stat := strings.ToLower(parts[0]), parts[1], parts[2]

		if handler == "updateHandler" {
			//switch stat {
			//case "adds", "autoCommits", "deletesById", "deletesByQuery", "docsPending", "errors", "softAutoCommits":
			//case "commits", "cumulativeAdds", "cumulativeDeletesById", "cumulativeDeletesByQuery", "cumulativeErrors", "expungeDeletes", "merges", "optimizes", "rollbacks", "splits":
			//}
			continue
		}

		switch stat {
		case "clientErrors", "errors", "serverErrors", "timeouts":
			if err := json.Unmarshal(stats, &p.common); err != nil {
				return err
			}
			parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += p.common.Count
		case "requests", "totalTime":
			//
			// 7.0+:
			// "UPDATE./update.requests": 0
			//
			// 6.4, 6.5:
			// "UPDATE./update.requests": { "count": 0 }
			//
			if p.version < 7.0 {
				if err := json.Unmarshal(stats, &p.count); err != nil {
					return err
				}
				parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += p.count.Count
			} else {
				if err := json.Unmarshal(stats, &p.simpleCount); err != nil {
					return err
				}
				parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += p.simpleCount
			}
		case "requestTimes":
			if err := json.Unmarshal(stats, &p.requestTimes); err != nil {
				return err
			}
			parsed[fmt.Sprintf("%s_%s_%s_count", core, typ, stat)] += p.requestTimes.Count
			parsed[fmt.Sprintf("%s_%s_%s_mean_ms", core, typ, stat)] += int64(p.requestTimes.MeanMS * 1e6)
			parsed[fmt.Sprintf("%s_%s_%s_median_ms", core, typ, stat)] += int64(p.requestTimes.MedianMS * 1e6)
			parsed[fmt.Sprintf("%s_%s_%s_p75_ms", core, typ, stat)] += int64(p.requestTimes.P75MS * 1e6)
			parsed[fmt.Sprintf("%s_%s_%s_p95_ms", core, typ, stat)] += int64(p.requestTimes.P95MS * 1e6)
			parsed[fmt.Sprintf("%s_%s_%s_p99_ms", core, typ, stat)] += int64(p.requestTimes.P99MS * 1e6)
			parsed[fmt.Sprintf("%s_%s_%s_p999_ms", core, typ, stat)] += int64(p.requestTimes.P999MS * 1e6)
		}
	}

	return nil
}
