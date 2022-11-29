package prometheus

import (
	"fmt"
	"os"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (p *Prometheus) validateConfig() error {
	return nil
}

func (p *Prometheus) initPrometheusClient() (prometheus.Prometheus, error) {
	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client: %v", err)
	}

	req := p.Request.Copy()
	if p.BearerTokenFile != "" {
		token, err := os.ReadFile(p.BearerTokenFile)
		if err != nil {
			return nil, fmt.Errorf("bearer token file: %v", err)
		}
		req.Headers["Authorization"] = "Bearer " + string(token)
	}

	return prometheus.New(client, req), nil
}
