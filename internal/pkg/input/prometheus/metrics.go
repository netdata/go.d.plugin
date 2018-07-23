package prometheus

import "github.com/prometheus/prometheus/pkg/labels"

type (
	// Metric is a pair of label set and value
	Metric struct {
		labels.Labels
		Value float64
	}

	// Metrics is a list of Metric
	Metrics []Metric
)
