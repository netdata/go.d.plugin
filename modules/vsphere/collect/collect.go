package collect

import (
	"sync"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go-orchestrator/logger"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
)

type APIClient interface {
	Version() string
	PerformanceMetrics([]types.PerfQuerySpec) ([]performance.EntityMetric, error)
}

func NewVSphereMetricCollector(client APIClient) *VSphereMetricCollector {
	return &VSphereMetricCollector{
		APIClient: client,
	}
}

type VSphereMetricCollector struct {
	*logger.Logger
	APIClient
}

func (c VSphereMetricCollector) CollectHostsMetrics(hosts rs.Hosts) []performance.EntityMetric {
	pqs := newHostsPerfQuerySpecs(hosts)
	return c.collectMetrics(pqs)
}

func (c VSphereMetricCollector) CollectVMsMetrics(vms rs.VMs) []performance.EntityMetric {
	pqs := newVMsPerfQuerySpecs(vms)
	return c.collectMetrics(pqs)
}

func (c VSphereMetricCollector) collectMetrics(pqs []types.PerfQuerySpec) []performance.EntityMetric {
	// TODO: hardcoded
	chunks := chunkify(pqs, 256)
	tc := newThrottledCaller(5)

	var ms []performance.EntityMetric
	lock := &sync.Mutex{}

	for _, chunk := range chunks {
		pqs := chunk
		job := func() {
			c.collect(&ms, lock, pqs)
		}
		tc.call(job)
	}
	tc.wait()

	return ms
}

func (c VSphereMetricCollector) collect(metrics *[]performance.EntityMetric, lock *sync.Mutex, pqs []types.PerfQuerySpec) {
	m, err := c.PerformanceMetrics(pqs)
	if err != nil {
		c.Error(err)
		return
	}

	lock.Lock()
	*metrics = append(*metrics, m...)
	lock.Unlock()
}

func chunkify(pqs []types.PerfQuerySpec, chunkSize int) (chunks [][]types.PerfQuerySpec) {
	for i := 0; i < len(pqs); i += chunkSize {
		end := i + chunkSize
		if end > len(pqs) {
			end = len(pqs)
		}
		chunks = append(chunks, pqs[i:end])
	}
	return chunks
}

const (
	pqsMaxSample  = 1
	pqsIntervalID = 20
	pqsFormat     = "normal"
)

func newHostsPerfQuerySpecs(hosts rs.Hosts) []types.PerfQuerySpec {
	var pqs []types.PerfQuerySpec
	for _, host := range hosts {
		pq := types.PerfQuerySpec{
			Entity:     host.Ref,
			MaxSample:  pqsMaxSample,
			MetricId:   host.MetricList,
			IntervalId: pqsIntervalID,
			Format:     pqsFormat,
		}
		pqs = append(pqs, pq)
	}
	return pqs
}

func newVMsPerfQuerySpecs(vms rs.VMs) []types.PerfQuerySpec {
	var pqs []types.PerfQuerySpec
	for _, vm := range vms {
		pq := types.PerfQuerySpec{
			Entity:     vm.Ref,
			MaxSample:  pqsMaxSample,
			MetricId:   vm.MetricList,
			IntervalId: pqsIntervalID,
			Format:     pqsFormat,
		}
		pqs = append(pqs, pq)
	}
	return pqs
}
