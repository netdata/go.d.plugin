package couchbase

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathBucketsStats = "/pools/default/buckets"
)

func (cb *Couchbase) collect() (map[string]int64, error) {

	ms := cb.scrapeCouchbase()
	if ms.empty() {
		return nil, nil
	}
	collected := make(map[string]int64)
	err := cb.collectBasicStats(collected, ms)
	if err != nil {
		return nil, fmt.Errorf("error on creating a connection: %v", err)
	}

	return collected, nil
}

func (cb Couchbase) collectBasicStats(collected map[string]int64, ms *cbMetrics) error {
	for _, b := range ms.BucketsStats {

		if !cb.collectedBuckets[b.Name] {
			cb.collectedBuckets[b.Name] = true
			cb.addBucketToCharts(b.Name)
		}

		bs := b.BasicStats
		collected[indexDimID(b.Name, "quota_percent_used")] = int64(bs.QuotaPercentUsed)
		collected[indexDimID(b.Name, "ops_per_sec")] = int64(bs.OpsPerSec)
		collected[indexDimID(b.Name, "disk_fetches")] = int64(bs.DiskFetches)
		collected[indexDimID(b.Name, "item_count")] = int64(bs.ItemCount)
		collected[indexDimID(b.Name, "disk_used")] = int64(bs.DiskUsed)
		collected[indexDimID(b.Name, "data_used")] = int64(bs.DataUsed)
		collected[indexDimID(b.Name, "mem_used")] = int64(bs.MemUsed)
		collected[indexDimID(b.Name, "vb_active_num_non_resident")] = int64(bs.VbActiveNumNonResident)
	}

	return nil
}

func (cb *Couchbase) addBucketToCharts(bucket string) {

	cb.addDimToChart(dbPercentCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "quota_percent_used"),
	})

	cb.addDimToChart(opsPerSecCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "ops_per_sec"),
	})

	cb.addDimToChart(diskFetchesCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "disk_fetches"),
	})

	cb.addDimToChart(itemCountCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "item_count"),
	})

	cb.addDimToChart(diskUsedCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "disk_used"),
	})

	cb.addDimToChart(dataUsedCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "data_used"),
	})

	cb.addDimToChart(memUsedCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "mem_used"),
	})

	cb.addDimToChart(vbActiveNumNonResidentCharts.ID, &module.Dim{
		ID: indexDimID(bucket, "vb_active_num_non_resident"),
	})
}

func (cb *Couchbase) addDimToChart(chartID string, dim *module.Dim) {
	chart := cb.Charts().Get(chartID)
	if chart == nil {
		cb.Warningf("error on adding '%s' dimension: can not find '%s' chart", dim.ID, chartID)
		return
	}
	if err := chart.AddDim(dim); err != nil {
		cb.Warning(err)
		return
	}
	chart.MarkNotCreated()
}

func (cb Couchbase) scrapeCouchbase() *cbMetrics {
	ms := &cbMetrics{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() { defer wg.Done(); cb.scrapeBucketsStats(ms) }()

	wg.Wait()
	return ms
}

func (cb Couchbase) scrapeBucketsStats(ms *cbMetrics) {
	req, _ := web.NewHTTPRequest(cb.Request)
	req.URL.Path = urlPathBucketsStats

	if err := cb.doOKDecode(req, &ms.BucketsStats); err != nil {
		cb.Warning(err)
		return
	}
}

func (cb Couchbase) pingCouchbase() error {
	req, _ := web.NewHTTPRequest(cb.Request)
	req.URL.Path = urlPathBucketsStats

	var stats []bucketsStats
	return cb.doOKDecode(req, &stats)
}

func (cb Couchbase) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := cb.httpClient.Do(req)
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

func indexDimID(name, metric string) string {
	return fmt.Sprintf("bucket_%s_%s", name, metric)
}
