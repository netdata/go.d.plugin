package elasticsearch

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (e *Elasticsearch) collect() (map[string]int64, error) {
	var mx esMetrics
	e.scrapeElasticsearch(&mx, true)

	return stm.ToMap(mx), nil
}

func (e *Elasticsearch) scrapeElasticsearch(mx *esMetrics, concurrently bool) {
	type scrapeJob func(mx *esMetrics)

	wg := &sync.WaitGroup{}
	wrap := func(job scrapeJob) scrapeJob {
		return func(mx *esMetrics) {
			job(mx)
			wg.Done()
		}
	}

	jobs := []scrapeJob{
		e.scrapeNodeStats,
		e.scrapeClusterHealth,
		e.scrapeClusterStats,
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

func (e *Elasticsearch) scrapeNodeStats(mx *esMetrics) {

}

func (e *Elasticsearch) scrapeClusterHealth(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = "/_cluster/health"

	resp, err := e.httpClient.Do(req)
	if err != nil {
		e.Warningf("error on HTTP request '%s': %v", req.URL, err)
		return
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		e.Warningf("'%s' returned HTTP status code: %d", req.URL, resp.StatusCode)
		return
	}

	var health esClusterHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		e.Warningf("decoding response from '%s': %v", req.URL, err)
		return
	}

	mx.ClusterHealth = &health
}

func (e *Elasticsearch) scrapeClusterStats(mx *esMetrics) {

}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
