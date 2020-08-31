package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathLocalNodeStats = "/_nodes/_local/stats"
	urlPathIndicesStats   = "/_cat/indices"
	urlPathClusterHealth  = "/_cluster/health"
	urlPathClusterStats   = "/_cluster/stats"
)

func (es *Elasticsearch) collect() (map[string]int64, error) {
	ms := es.scrapeElasticsearch()
	if ms.empty() {
		return nil, nil
	}

	collected := make(map[string]int64)
	es.collectLocalNodeStats(collected, ms)
	es.collectClusterHealth(collected, ms)
	es.collectClusterStats(collected, ms)
	es.collectLocalIndicesStats(collected, ms)

	return collected, nil
}

func (Elasticsearch) collectLocalNodeStats(collected map[string]int64, ms *esMetrics) {
	if !ms.hasLocalNodeStats() {
		return
	}
	merge(collected, stm.ToMap(ms.LocalNodeStats), "node")
}

func (Elasticsearch) collectClusterHealth(collected map[string]int64, ms *esMetrics) {
	if !ms.hasClusterHealth() {
		return
	}
	merge(collected, stm.ToMap(ms.ClusterHealth), "cluster")
	collected["cluster_status"] = convertHealthStatus(ms.ClusterHealth.Status)
}

func (Elasticsearch) collectClusterStats(collected map[string]int64, ms *esMetrics) {
	if !ms.hasClusterStats() {
		return
	}
	merge(collected, stm.ToMap(ms.ClusterStats), "cluster")
}

func (es *Elasticsearch) collectLocalIndicesStats(mx map[string]int64, ms *esMetrics) {
	if !ms.hasLocalIndicesStats() {
		return
	}
	key := func(name, metric string) string {
		return fmt.Sprintf("node_indices_stats_%s_index_%s", name, metric)
	}
	for _, index := range ms.LocalIndicesStats {
		mx[key(index.Index, "health")] = convertHealthStatus(index.Health)
		mx[key(index.Index, "shards_count")] = strToInt(index.Rep)
		mx[key(index.Index, "docs_count")] = strToInt(index.DocsCount)
		mx[key(index.Index, "store_size_in_bytes")] = convertIndexStoreSizeToBytes(index.StoreSize)
	}
}

func (es *Elasticsearch) addIndexToCharts(index string) {}

func convertHealthStatus(status string) int64 {
	switch status {
	case "green":
		return 0
	case "yellow":
		return 1
	case "red":
		return 2
	default:
		return 2
	}
}

func convertIndexStoreSizeToBytes(size string) int64 {
	var num float64
	switch {
	case strings.HasSuffix(size, "kb"):
		num, _ = strconv.ParseFloat(size[:len(size)-2], 64)
		num *= math.Pow(1024, 1)
	case strings.HasSuffix(size, "mb"):
		num, _ = strconv.ParseFloat(size[:len(size)-2], 64)
		num *= math.Pow(1024, 2)
	case strings.HasSuffix(size, "gb"):
		num, _ = strconv.ParseFloat(size[:len(size)-2], 64)
		num *= math.Pow(1024, 3)
	case strings.HasSuffix(size, "tb"):
		num, _ = strconv.ParseFloat(size[:len(size)-2], 64)
		num *= math.Pow(1024, 4)
	case strings.HasSuffix(size, "b"):
		num, _ = strconv.ParseFloat(size[:len(size)-1], 64)
	}
	return int64(num)
}

func strToInt(s string) int64 {
	v, _ := strconv.Atoi(s)
	return int64(v)
}

func (es Elasticsearch) scrapeElasticsearch() *esMetrics {
	tasks := []func(ms *esMetrics){
		es.scrapeLocalNodeStats,
		es.scrapeClusterHealth,
		es.scrapeClusterStats,
		es.scrapeLocalIndicesStats,
	}

	var metrics esMetrics
	wg := &sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		task := task
		go func() { defer wg.Done(); task(&metrics) }()
	}
	wg.Wait()
	return &metrics
}

func (es Elasticsearch) scrapeLocalNodeStats(ms *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathLocalNodeStats

	var stats struct {
		Nodes map[string]esNodeStats
	}
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	for _, node := range stats.Nodes {
		ms.LocalNodeStats = &node
		break
	}
}

func (es Elasticsearch) scrapeClusterHealth(ms *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterHealth

	var health esClusterHealth
	if err := es.doOKDecode(req, &health); err != nil {
		es.Warning(err)
		return
	}
	ms.ClusterHealth = &health
}

func (es Elasticsearch) scrapeClusterStats(ms *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterStats

	var stats esClusterStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	ms.ClusterStats = &stats
}

func (es *Elasticsearch) scrapeLocalIndicesStats(ms *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathIndicesStats
	req.URL.RawQuery = "local=true&format=json"

	var stats []esIndexStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	ms.LocalIndicesStats = removeSystemIndices(stats)
}

func (es Elasticsearch) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on HTTP request '%s': %v", req.URL, err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("'%s' returned HTTP status code: %d", req.URL, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(in); err != nil {
		return fmt.Errorf("error on decoding response from '%s': %v", req.URL, err)
	}
	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

func removeSystemIndices(indices []esIndexStats) []esIndexStats {
	var i int
	for _, index := range indices {
		if strings.HasPrefix(index.Index, ".") {
			continue
		}
		indices[i] = index
		i++
	}
	return indices[:i]
}

func merge(dst, src map[string]int64, prefix string) {
	for k, v := range src {
		if prefix == "" {
			dst[k] = v
		} else {
			dst[prefix+"_"+k] = v
		}
	}
}
