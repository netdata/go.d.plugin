package prometheus

import (
	"errors"
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/prometheus/prometheus/pkg/labels"
)

func (p Prometheus) validateConfig() error {
	if p.UserURL == "" {
		return errors.New("URL not set")
	}
	return nil
}

func (p Prometheus) initPrometheusClient() (prometheus.Prometheus, error) {
	client, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client: %v", err)
	}

	sr, err := p.Selector.Parse()
	if err != nil {
		return nil, fmt.Errorf("parsing selector: %v", err)
	}

	if sr != nil {
		return prometheus.NewWithSelector(client, p.Request, sr), nil
	}
	return prometheus.New(client, p.Request), nil
}

func (p Prometheus) initOptionalGrouping() ([]optionalGrouping, error) {
	var optGrps []optionalGrouping
	for _, item := range p.Grouping {
		if item.Selector == "" {
			return nil, errors.New("empty grouping selector")
		}

		if item.ByLabel == "" {
			return nil, fmt.Errorf("grouping selector '%s' has no 'by_label'", item.Selector)
		}

		sr, err := selector.Parse(item.Selector)
		if err != nil {
			return nil, fmt.Errorf("parsing grouping selector '%s': %v", item.Selector, err)
		}
		if sr == nil {
			continue
		}

		byLabels := strings.Fields(item.ByLabel)
		fn := selector.Func(func(lbs labels.Labels) bool {
			return lbs.Len() >= 3 && labelsContains(lbs, byLabels...) && sr.Matches(lbs)
		})
		optGrps = append(optGrps, optionalGrouping{
			sr:  fn,
			grp: newGroupingGroupedBy(byLabels...),
		})
	}
	return optGrps, nil
}

func labelsContains(lbs labels.Labels, names ...string) bool {
	switch len(names) {
	case 0:
		return true
	case 1:
		return lbs.Has(names[0])
	default:
		return lbs.Has(names[0]) && labelsContains(lbs, names[1:]...)
	}
}
