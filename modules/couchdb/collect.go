package couchdb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathOverviewStats = "/_node/nonode@nohost/_stats"
	urlPathSystemStats   = "/_node/nonode@nohost/_system"
	urlPathActiveTasks   = "/_active_tasks"
)

func (cdb *CouchDB) collect() (map[string]int64, error) {
	ms := cdb.scrapeCouchDB()
	if ms.empty() {
		return nil, nil
	}

	collected := make(map[string]int64)
	cdb.collectNodeStats(collected, ms)

	return collected, nil
}

func (CouchDB) collectNodeStats(collected map[string]int64, ms *cdbMetrics) {
	if !ms.hasNodeStats() {
		return
	}
	merge(collected, stm.ToMap(ms.NodeStats), "node")
}

func (cdb CouchDB) scrapeCouchDB() *cdbMetrics {
	ms := &cdbMetrics{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeNodeStats(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); cdb.scrapeSystemStats(ms) }()

	wg.Wait()
	return ms
}

func (cdb CouchDB) scrapeNodeStats(ms *cdbMetrics) {
	req, _ := web.NewHTTPRequest(cdb.Request)
	req.URL.Path = urlPathOverviewStats

	var stats cdbNodeStats
	if err := cdb.doOKDecode(req, &stats); err != nil {
		cdb.Warning(err)
		return
	}
	ms.NodeStats = &stats
}

func (cdb *CouchDB) scrapeSystemStats(ms *cdbMetrics) {
	req, _ := web.NewHTTPRequest(cdb.Request)
	req.URL.Path = urlPathSystemStats

	var stats cdbNodeSystem
	if err := cdb.doOKDecode(req, &stats); err != nil {
		cdb.Warning(err)
		return
	}
	ms.NodeSystem = &stats
}

func (cdb CouchDB) pingCouchDB() error {
	req, _ := web.NewHTTPRequest(cdb.Request)

	var info struct{ Name string }
	return cdb.doOKDecode(req, &info)
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
