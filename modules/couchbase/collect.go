package couchbase

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

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

		if !cb.collectedBuckets[b.Name] && ms.hasBucketsStats() {
			cb.collectedBuckets[b.Name] = true
			cb.addDimsToCharts(b.Name)
		}

		bs := b.BasicStats
		collected[indexDimID(b.Name, "quota_used")] = int64(bs.QuotaPercentUsed)
		collected[indexDimID(b.Name, "ops")] = int64(bs.OpsPerSec)
		collected[indexDimID(b.Name, "fetches")] = int64(bs.DiskFetches)
		collected[indexDimID(b.Name, "item_count")] = int64(bs.ItemCount)
		collected[indexDimID(b.Name, "disk")] = int64(bs.DiskUsed)
		collected[indexDimID(b.Name, "data")] = int64(bs.DataUsed)
		collected[indexDimID(b.Name, "mem")] = int64(bs.MemUsed)
		collected[indexDimID(b.Name, "num_non_resident")] = int64(bs.VbActiveNumNonResident)
	}

	return nil
}

func (cb *Couchbase) addDimsToCharts(bucket string) error {
	for _, chart := range *cb.Charts() {
		switch chart.ID {
		case "couchbase_quota_percent_used_stats":
			dim := &Dim{
				ID: indexDimID(bucket, "quota_used"),
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_ops_per_sec_stats":
			dim := &Dim{
				ID: indexDimID(bucket, "ops"),
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_disk_fetches_stats":
			dim := &Dim{
				ID: indexDimID(bucket, "fetches"),
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_item_count_stats":
			dim := &Dim{
				ID: indexDimID(bucket, "item_count"),
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_disk_used_stats":
			dim := &Dim{
				ID:  indexDimID(bucket, "disk"),
				Div: 1024,
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_data_used_stats":
			dim := &Dim{
				ID:  indexDimID(bucket, "data"),
				Div: 1024,
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_mem_used_stats":
			dim := &Dim{
				ID:  indexDimID(bucket, "mem"),
				Div: 1024,
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		case "couchbase_vb_active_num_non_resident_stats":
			dim := &Dim{
				ID: indexDimID(bucket, "num_non_resident"),
			}
			if err := chart.AddDim(dim); err != nil {
				return err
			}
		default:
			continue
		}
	}
	return nil
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

	var stats []bucketsStats
	if err := cb.doOKDecode(req, &stats); err != nil {
		cb.Warning(err)
		return
	}
	for _, bucket := range stats {
		ms.BucketsStats = append(ms.BucketsStats, bucket)
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
	return fmt.Sprintf("basic_stats_%s_%s", name, metric)
}
