package elasticsearch

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

func (es *Elasticsearch) collect() (map[string]int64, error) {
	var mx esMetrics
	es.scrapeElasticsearch(&mx, true)

	return stm.ToMap(mx), nil
}

func (es *Elasticsearch) scrapeElasticsearch(mx *esMetrics, concurrently bool) {
	type scrapeJob func(mx *esMetrics)

	wg := &sync.WaitGroup{}
	wrap := func(job scrapeJob) scrapeJob {
		return func(mx *esMetrics) {
			job(mx)
			wg.Done()
		}
	}

	jobs := []scrapeJob{
		es.scrapeNodeStats,
		es.scrapeClusterHealth,
		es.scrapeClusterStats,
	}

	for _, job := range jobs {
		if !concurrently {
			job(mx)
		} else {
			wg.Add(1)
			job := wrap(job)
			go job(mx)
		}
	}
	wg.Wait()
}

func (es *Elasticsearch) scrapeNodeStats(mx *esMetrics) {

}

func (es *Elasticsearch) scrapeClusterHealth(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = "/_cluster/health"

	var health esClusterHealth
	if err := es.doOKDecode(req, &health); err != nil {
		es.Warning(err)
		return
	}
	mx.ClusterHealth = &health
}

func (es *Elasticsearch) scrapeClusterStats(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = "/_cluster/health"

	var stats esClusterStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	mx.ClusterStats = &stats
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
