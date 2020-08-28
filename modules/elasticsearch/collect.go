package elasticsearch

import "sync"

func (e *Elasticsearch) collect() (map[string]int64, error) {
	mx := &esMetrics{}
	e.scrapeElasticsearch(mx, true)

	return nil, nil
}

func (e *Elasticsearch) scrapeElasticsearch(mx *esMetrics, doConcurrently bool) {
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
		if !doConcurrently {
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

}

func (e *Elasticsearch) scrapeClusterStats(mx *esMetrics) {

}
