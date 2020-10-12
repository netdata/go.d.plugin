package couchdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathActiveTasks = "/_active_tasks"

	httpStatusCodePrefix    = "couchdb_httpd_status_codes_"
	httpStatusCodePrefixLen = len(httpStatusCodePrefix)
)

func (cdb *CouchDB) collect() (map[string]int64, error) {
	ms := cdb.scrapeCouchDB()
	if ms.empty() {
		return nil, nil
	}

	collected := make(map[string]int64)
	cdb.collectNodeStats(collected, ms)
	cdb.collectSystemStats(collected, ms)
	cdb.collectActiveTasks(collected, ms)
	cdb.collectDBStats(collected, ms)

	return collected, nil
}

func (CouchDB) collectNodeStats(collected map[string]int64, ms *cdbMetrics) {
	if !ms.hasNodeStats() {
		return
	}

	for metric, value := range stm.ToMap(ms.NodeStats) {
		if strings.HasPrefix(metric, httpStatusCodePrefix) {
			aggregateHTTPStatusCodes(collected, metric, value)
		} else {
			collected[metric] = value
		}
	}
}

func (CouchDB) collectSystemStats(collected map[string]int64, ms *cdbMetrics) {
	if !ms.hasNodeSystem() {
		return
	}

	for metric, value := range stm.ToMap(ms.NodeSystem) {
		collected[metric] = value
	}

	collected["peak_msg_queue"] = findMaxMQSize(ms.NodeSystem.MessageQueues)
}

func (CouchDB) collectActiveTasks(collected map[string]int64, ms *cdbMetrics) {
	collected["active_tasks_indexer"] = 0
	collected["active_tasks_database_compaction"] = 0
	collected["active_tasks_replication"] = 0
	collected["active_tasks_view_compaction"] = 0

	if !ms.hasActiveTasks() {
		return
	}

	for _, task := range ms.ActiveTasks {
		collected["active_tasks_"+task.Type]++
	}
}

func (CouchDB) collectDBStats(collected map[string]int64, ms *cdbMetrics) {
	if !ms.hasDBStats() {
		return
	}

	for metric, value := range stm.ToMap(ms.DBStats) {
		collected[metric] = value
	}
}

func (cdb CouchDB) scrapeCouchDB() *cdbMetrics {
	ms := &cdbMetrics{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeNodeStats(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeSystemStats(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeActiveTasks(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeDBStats(ms) }()

	wg.Wait()
	return ms
}

func (cdb CouchDB) scrapeNodeStats(ms *cdbMetrics) {
	req, _ := web.NewHTTPRequest(cdb.Request)
	req.URL.Path = cdb.urlPathOverviewStats

	var stats cdbNodeStats
	if err := cdb.doOKDecode(req, &stats); err != nil {
		cdb.Warning(err)
		return
	}
	ms.NodeStats = &stats
}

func (cdb *CouchDB) scrapeSystemStats(ms *cdbMetrics) {
	req, _ := web.NewHTTPRequest(cdb.Request)
	req.URL.Path = cdb.urlPathSystemStats

	var stats cdbNodeSystem
	if err := cdb.doOKDecode(req, &stats); err != nil {
		cdb.Warning(err)
		return
	}
	ms.NodeSystem = &stats
}

func (cdb *CouchDB) scrapeActiveTasks(ms *cdbMetrics) {
	req, _ := web.NewHTTPRequest(cdb.Request)
	req.URL.Path = urlPathActiveTasks

	var stats []cdbActiveTask
	if err := cdb.doOKDecode(req, &stats); err != nil {
		cdb.Warning(err)
		return
	}
	ms.ActiveTasks = stats
}

func (cdb *CouchDB) scrapeDBStats(ms *cdbMetrics) {
	wg := &sync.WaitGroup{}

	ms.DBStats = make(map[string]*cdbDBStats)
	var DBStatsMutex sync.Mutex

	for _, db := range cdb.databases {
		req, _ := web.NewHTTPRequest(cdb.Request)
		req.URL.Path = "/" + db

		wg.Add(1)
		go func(db string) {
			defer wg.Done()

			var stats cdbDBStats
			if err := cdb.doOKDecode(req, &stats); err != nil {
				cdb.Warning(err)
				return
			}
			DBStatsMutex.Lock()
			ms.DBStats["db_"+db] = &stats
			DBStatsMutex.Unlock()
		}(db)
	}

	wg.Wait()
}

func aggregateHTTPStatusCodes(collected map[string]int64, metric string, value int64) {
	code := metric[httpStatusCodePrefixLen:]

	switch {
	case code == "200" || code == "201" || code == "202":
		collected[metric] = value
	case strings.HasPrefix(code, "2"):
		collected["couchdb_httpd_status_codes_2xx"] += value
	case strings.HasPrefix(code, "3"):
		collected["couchdb_httpd_status_codes_3xx"] += value
	case strings.HasPrefix(code, "4"):
		collected["couchdb_httpd_status_codes_4xx"] += value
	case strings.HasPrefix(code, "5"):
		collected["couchdb_httpd_status_codes_5xx"] += value
	default:
		collected[metric] = value
	}
}

func findMaxMQSize(MessageQueues map[string]interface{}) int64 {
	var max float64
	for _, mq := range MessageQueues {
		switch mqSize := mq.(type) {
		case float64:
			max = math.Max(max, mqSize)
		case map[string]interface{}:
			if v, ok := mqSize["count"].(float64); ok {
				max = math.Max(max, v)
			}
		}
	}
	return int64(max)
}

func (cdb CouchDB) pingCouchDB() error {
	req, _ := web.NewHTTPRequest(cdb.Request)

	var info struct{ Couchdb string }
	if err := cdb.doOKDecode(req, &info); err != nil {
		return err
	}

	if info.Couchdb != "Welcome" {
		return errors.New("not a CouchDB endpoint")
	}

	return nil
}

func (cdb CouchDB) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := cdb.httpClient.Do(req)
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

func merge(dst, src map[string]int64, prefix string) {
	for k, v := range src {
		dst[prefix+"_"+k] = v
	}
}
