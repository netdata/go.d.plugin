package prometheus

type (
	Prometheus struct {
		Scraper

		metrics Metrics
	}
)

func New(URL string) {

}

func (p *Prometheus) GetMetrics() (Metrics, error) {
	p.metrics = p.metrics[:0]
	if err := p.Scraper.Scrape(&p.metrics); err != nil {
		return nil, err
	}
	return p.metrics, nil
}
