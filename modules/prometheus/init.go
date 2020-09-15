package prometheus

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/prometheus/prometheus/pkg/labels"
)

func (p Prometheus) validateConfig() error {
	if p.URL == "" {
		return errors.New("URL not set")
	}
	return nil
}

func (p Prometheus) initPrometheusClient() (prometheus.Prometheus, error) {
	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client: %v", err)
	}

	req := p.Request.Copy()
	if p.BearerTokenFile != "" {
		token, err := ioutil.ReadFile(p.BearerTokenFile)
		if err != nil {
			return nil, fmt.Errorf("reading bearer token file: %v", err)
		}
		req.Headers["Authorization"] = "Bearer " + string(token)
	}

	sr, err := p.Selector.Parse()
	if err != nil {
		return nil, fmt.Errorf("parsing selector: %v", err)
	}

	if sr != nil {
		return prometheus.NewWithSelector(client, req, sr), nil
	}
	return prometheus.New(client, req), nil
}

func (p Prometheus) initOptionalGrouping() ([]optionalGrouping, error) {
	var optGrps []optionalGrouping
	for _, item := range p.Grouping {
		if item.Selector == "" {
			return nil, errors.New("empty group selector")
		}

		if item.ByLabel == "" {
			return nil, fmt.Errorf("group selector '%s' has no 'by_label'", item.Selector)
		}

		sr, err := selector.Parse(item.Selector)
		if err != nil {
			return nil, fmt.Errorf("parse group selector '%s': %v", item.Selector, err)
		}
		if sr == nil {
			continue
		}

		names := strings.Fields(item.ByLabel)
		fn := selector.Func(func(lbs labels.Labels) bool {
			return lbs.Len() >= 3 && labelsContainsAll(lbs, names...) && sr.Matches(lbs)
		})
		optGrps = append(optGrps, optionalGrouping{
			sr:  fn,
			grp: newGroupingGroupedBy(names...),
		})
	}
	return optGrps, nil
}

func labelsContainsAll(lbs labels.Labels, names ...string) bool {
	switch len(names) {
	case 0:
		return true
	case 1:
		return lbs.Has(names[0])
	default:
		return lbs.Has(names[0]) && labelsContainsAll(lbs, names[1:]...)
	}
}
