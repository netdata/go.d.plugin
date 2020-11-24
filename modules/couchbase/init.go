package couchbase

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (cb *Couchbase) initCharts() (*Charts, error) {
	charts := module.Charts{}

	bucketNames := cb.bucketNames
	if err := newDbPercentChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newOpPerSecChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newDiskFetchesChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newItemCountChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newDiskUsedChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newDataUsedChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newMemUsedChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	if err := newVbActiveNumNonResidentChart(&charts, bucketNames); err != nil {
		return nil, err
	}

	return &charts, nil

}

func newItemCountChart(charts *Charts, buckets []string) error {
	chart := itemCountCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID: fmt.Sprintf("item_count_%s", v),
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newDbPercentChart(charts *Charts, buckets []string) error {
	chart := dbPercentCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID: fmt.Sprintf("quota_used_%s", v),
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newOpPerSecChart(charts *Charts, buckets []string) error {
	chart := opPerSecCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID: fmt.Sprintf("ops_%s", v),
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newDiskFetchesChart(charts *Charts, buckets []string) error {
	chart := diskFetchesCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID: fmt.Sprintf("fetches_%s", v),
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newDataUsedChart(charts *Charts, buckets []string) error {
	chart := dataUsedCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID:  fmt.Sprintf("data_%s", v),
			Div: 1024,
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newDiskUsedChart(charts *Charts, buckets []string) error {
	chart := diskUsedCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID:  fmt.Sprintf("disk_%s", v),
			Div: 1024,
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newMemUsedChart(charts *Charts, buckets []string) error {
	chart := memUsedCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID:  fmt.Sprintf("mem_%s", v),
			Div: 1024,
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}

func newVbActiveNumNonResidentChart(charts *Charts, buckets []string) error {
	chart := vbActiveNumNonResidentCharts.Copy()
	for _, v := range buckets {
		dim := &Dim{
			ID: fmt.Sprintf("num_non_resident_%s", v),
		}
		if err := chart.AddDim(dim); err != nil {
			return err
		}
	}
	if err := charts.Add(chart); err != nil {
		return err
	}
	return nil
}
