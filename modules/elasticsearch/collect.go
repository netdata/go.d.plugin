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
	urlPathNodesLocalStats = "/_nodes/_local/stats"
	urlPathClusterHealth   = "/_cluster/health"
	urlPathClusterStats    = "/_cluster/stats"
	urlPathCatIndices      = "/_cat/indices?format=json"
)

func (es *Elasticsearch) collect() (map[string]int64, error) {
	mx := es.scrapeElasticsearch()

	return stm.ToMap(mx), nil
}

func (es *Elasticsearch) scrapeElasticsearch() *esMetrics {
	tasks := []func(metrics *esMetrics){
		es.scrapeLocalNodeStats,
		es.scrapeClusterHealth,
		es.scrapeClusterStats,
		es.scrapeIndicesStats,
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

func (es *Elasticsearch) scrapeLocalNodeStats(metrics *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathNodesLocalStats

	var stats struct {
		Nodes map[string]esNodeStats
	}
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	for _, node := range stats.Nodes {
		metrics.LocalNodeStats = &node
		break
	}
}

func (es *Elasticsearch) scrapeClusterHealth(metrics *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterHealth

	var health esClusterHealth
	if err := es.doOKDecode(req, &health); err != nil {
		es.Warning(err)
		return
	}
	metrics.ClusterHealth = &health
}

func (es *Elasticsearch) scrapeClusterStats(metrics *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterStats

	var stats esClusterStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	metrics.ClusterStats = &stats
}

func (es *Elasticsearch) scrapeIndicesStats(metrics *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathCatIndices

	var stats []esIndexStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	metrics.IndicesStats = stats
	//var i int
	//for _, index := range stats {
	//	if !strings.HasPrefix(index.Index, ".") {
	//		stats[i] = index
	//		i++
	//	}
	//}
}

func (es *Elasticsearch) doOKDecode(req *http.Request, in interface{}) error {
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
