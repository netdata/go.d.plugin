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

	"github.com/netdata/go.d.plugin/agent/module"
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
	seen := make(map[string]struct{})
	for _, index := range ms.LocalIndicesStats {
		seen[index.Index] = struct{}{}
		if !es.collectedIndices[index.Index] {
			es.collectedIndices[index.Index] = true
			es.addIndexToCharts(index.Index)
		}
		mx[indexDimID(index.Index, "health")] = convertHealthStatus(index.Health)
		mx[indexDimID(index.Index, "shards_count")] = strToInt(index.Rep)
		mx[indexDimID(index.Index, "docs_count")] = strToInt(index.DocsCount)
		mx[indexDimID(index.Index, "store_size_in_bytes")] = convertIndexStoreSizeToBytes(index.StoreSize)
	}
	for index := range es.collectedIndices {
		if _, ok := seen[index]; !ok {
			delete(es.collectedIndices, index)
			es.removeIndexFromCharts(index)
		}
	}
}

func (es *Elasticsearch) addIndexToCharts(index string) {
	for _, chart := range *es.Charts() {
		dim := module.Dim{Name: index}
		switch chart.ID {
		case "node_index_health":
			dim.ID = indexDimID(index, "health")
		case "node_index_shards_count":
			dim.ID = indexDimID(index, "shards_count")
		case "node_index_docs_count":
			dim.ID = indexDimID(index, "docs_count")
		case "node_index_store_size":
			dim.ID = indexDimID(index, "store_size_in_bytes")
		default:
			continue
		}
		if err := chart.AddDim(&dim); err != nil {
			es.Warningf("add index '%s': %v", index, err)
			continue
		}
		chart.MarkNotCreated()
	}
}

func (es *Elasticsearch) removeIndexFromCharts(index string) {
	for _, chart := range *es.Charts() {
		var id string
		switch chart.ID {
		case "node_index_health":
			id = indexDimID(index, "health")
		case "node_index_shards_count":
			id = indexDimID(index, "shards_count")
		case "node_index_docs_count":
			id = indexDimID(index, "docs_count")
		case "node_index_store_size":
			id = indexDimID(index, "store_size_in_bytes")
		default:
			continue
		}
		if err := chart.MarkDimRemove(id, false); err != nil {
			es.Warningf("remove index '%s': %v", index, err)
			continue
		}
		chart.MarkNotCreated()
	}
}

func indexDimID(name, metric string) string {
	return fmt.Sprintf("node_indices_stats_%s_index_%s", name, metric)
}

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
	ms := &esMetrics{}
	wg := &sync.WaitGroup{}

	if es.DoNodeStats {
		wg.Add(1)
		go func() { defer wg.Done(); es.scrapeLocalNodeStats(ms) }()
	}
	if es.DoClusterHealth {
		wg.Add(1)
		go func() { defer wg.Done(); es.scrapeClusterHealth(ms) }()
	}
	if es.DoClusterStats {
		wg.Add(1)
		go func() { defer wg.Done(); es.scrapeClusterStats(ms) }()
	}
	if es.DoIndicesStats {
		wg.Add(1)
		go func() { defer wg.Done(); es.scrapeLocalIndicesStats(ms) }()
	}
	wg.Wait()
	return ms
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

func (es Elasticsearch) pingElasticsearch() error {
	req, _ := web.NewHTTPRequest(es.Request)

	var info struct{ Name string }
	return es.doOKDecode(req, &info)
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
		dst[prefix+"_"+k] = v
	}
}
