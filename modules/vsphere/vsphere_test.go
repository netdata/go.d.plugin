package vsphere

import (
	"crypto/tls"
	"testing"
	"time"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/simulator"
)

func newTestJob(vCenterURL string) *VSphere {
	job := New()
	job.Username = "administrator"
	job.Password = "password"
	job.UserURL = vCenterURL
	return job
}

func createSim() (*simulator.Model, *simulator.Server, error) {
	model := simulator.VPX()

	err := model.Create()
	if err != nil {
		return nil, nil, err
	}

	model.Service.TLS = new(tls.Config)

	s := model.Service.NewServer()
	return model, s, nil
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultURL, job.UserURL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.Equal(t, defaultDiscoveryInterval, job.DiscoveryInterval.Duration)
}

func TestVSphere_Init(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()
	job := newTestJob(srv.URL.String())
	defer job.Cleanup()

	assert.True(t, job.Init())
	assert.NotNil(t, job.discoverer)
	assert.NotNil(t, job.metricScraper)
	assert.NotNil(t, job.resources)
	assert.NotNil(t, job.discoveryTask)
	assert.True(t, job.discoveryTask.isRunning())
}

func TestVSphere_Init_NG(t *testing.T) {
	job := newTestJob("http://127.0.0.1:32001")

	assert.False(t, job.Init())
}

func TestVSphere_Check(t *testing.T) {
	assert.NotNil(t, New().Check())
}

func TestVSphere_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestVSphere_Cleanup(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()
	job := newTestJob(srv.URL.String())
	require.True(t, job.Init())

	job.Cleanup()
	time.Sleep(time.Second)
	assert.True(t, job.discoveryTask.isStopped())
	assert.False(t, job.discoveryTask.isRunning())
}

func TestVSphere_Cleanup_NotInited(t *testing.T) {
	New().Cleanup()
}

func TestVSphere_Collect(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	job := newTestJob(srv.URL.String())
	defer job.Cleanup()
	require.True(t, job.Init())
	require.True(t, job.Check())
	job.metricScraper = testMetricScraper{job.metricScraper}

	expected := map[string]int64{
		"host-20_cpu.usage.average":           100,
		"host-20_disk.maxTotalLatency.latest": 100,
		"host-20_disk.read.average":           100,
		"host-20_disk.write.average":          100,
		"host-20_mem.active.average":          100,
		"host-20_mem.consumed.average":        100,
		"host-20_mem.granted.average":         100,
		"host-20_mem.shared.average":          100,
		"host-20_mem.sharedcommon.average":    100,
		"host-20_mem.swapinRate.average":      100,
		"host-20_mem.swapoutRate.average":     100,
		"host-20_mem.usage.average":           100,
		"host-20_net.bytesRx.average":         100,
		"host-20_net.bytesTx.average":         100,
		"host-20_net.droppedRx.summation":     100,
		"host-20_net.droppedTx.summation":     100,
		"host-20_net.errorsRx.summation":      100,
		"host-20_net.errorsTx.summation":      100,
		"host-20_net.packetsRx.summation":     100,
		"host-20_net.packetsTx.summation":     100,
		"host-20_sys.uptime.latest":           100,
		"host-32_cpu.usage.average":           100,
		"host-32_disk.maxTotalLatency.latest": 100,
		"host-32_disk.read.average":           100,
		"host-32_disk.write.average":          100,
		"host-32_mem.active.average":          100,
		"host-32_mem.consumed.average":        100,
		"host-32_mem.granted.average":         100,
		"host-32_mem.shared.average":          100,
		"host-32_mem.sharedcommon.average":    100,
		"host-32_mem.swapinRate.average":      100,
		"host-32_mem.swapoutRate.average":     100,
		"host-32_mem.usage.average":           100,
		"host-32_net.bytesRx.average":         100,
		"host-32_net.bytesTx.average":         100,
		"host-32_net.droppedRx.summation":     100,
		"host-32_net.droppedTx.summation":     100,
		"host-32_net.errorsRx.summation":      100,
		"host-32_net.errorsTx.summation":      100,
		"host-32_net.packetsRx.summation":     100,
		"host-32_net.packetsTx.summation":     100,
		"host-32_sys.uptime.latest":           100,
		"host-39_cpu.usage.average":           100,
		"host-39_disk.maxTotalLatency.latest": 100,
		"host-39_disk.read.average":           100,
		"host-39_disk.write.average":          100,
		"host-39_mem.active.average":          100,
		"host-39_mem.consumed.average":        100,
		"host-39_mem.granted.average":         100,
		"host-39_mem.shared.average":          100,
		"host-39_mem.sharedcommon.average":    100,
		"host-39_mem.swapinRate.average":      100,
		"host-39_mem.swapoutRate.average":     100,
		"host-39_mem.usage.average":           100,
		"host-39_net.bytesRx.average":         100,
		"host-39_net.bytesTx.average":         100,
		"host-39_net.droppedRx.summation":     100,
		"host-39_net.droppedTx.summation":     100,
		"host-39_net.errorsRx.summation":      100,
		"host-39_net.errorsTx.summation":      100,
		"host-39_net.packetsRx.summation":     100,
		"host-39_net.packetsTx.summation":     100,
		"host-39_sys.uptime.latest":           100,
		"host-46_cpu.usage.average":           100,
		"host-46_disk.maxTotalLatency.latest": 100,
		"host-46_disk.read.average":           100,
		"host-46_disk.write.average":          100,
		"host-46_mem.active.average":          100,
		"host-46_mem.consumed.average":        100,
		"host-46_mem.granted.average":         100,
		"host-46_mem.shared.average":          100,
		"host-46_mem.sharedcommon.average":    100,
		"host-46_mem.swapinRate.average":      100,
		"host-46_mem.swapoutRate.average":     100,
		"host-46_mem.usage.average":           100,
		"host-46_net.bytesRx.average":         100,
		"host-46_net.bytesTx.average":         100,
		"host-46_net.droppedRx.summation":     100,
		"host-46_net.droppedTx.summation":     100,
		"host-46_net.errorsRx.summation":      100,
		"host-46_net.errorsTx.summation":      100,
		"host-46_net.packetsRx.summation":     100,
		"host-46_net.packetsTx.summation":     100,
		"host-46_sys.uptime.latest":           100,
		"vm-53_cpu.usage.average":             200,
		"vm-53_disk.maxTotalLatency.latest":   200,
		"vm-53_disk.read.average":             200,
		"vm-53_disk.write.average":            200,
		"vm-53_mem.active.average":            200,
		"vm-53_mem.consumed.average":          200,
		"vm-53_mem.granted.average":           200,
		"vm-53_mem.shared.average":            200,
		"vm-53_mem.swapinRate.average":        200,
		"vm-53_mem.swapoutRate.average":       200,
		"vm-53_mem.swapped.average":           200,
		"vm-53_mem.usage.average":             200,
		"vm-53_net.bytesRx.average":           200,
		"vm-53_net.bytesTx.average":           200,
		"vm-53_net.droppedRx.summation":       200,
		"vm-53_net.droppedTx.summation":       200,
		"vm-53_net.packetsRx.summation":       200,
		"vm-53_net.packetsTx.summation":       200,
		"vm-53_sys.uptime.latest":             200,
		"vm-56_cpu.usage.average":             200,
		"vm-56_disk.maxTotalLatency.latest":   200,
		"vm-56_disk.read.average":             200,
		"vm-56_disk.write.average":            200,
		"vm-56_mem.active.average":            200,
		"vm-56_mem.consumed.average":          200,
		"vm-56_mem.granted.average":           200,
		"vm-56_mem.shared.average":            200,
		"vm-56_mem.swapinRate.average":        200,
		"vm-56_mem.swapoutRate.average":       200,
		"vm-56_mem.swapped.average":           200,
		"vm-56_mem.usage.average":             200,
		"vm-56_net.bytesRx.average":           200,
		"vm-56_net.bytesTx.average":           200,
		"vm-56_net.droppedRx.summation":       200,
		"vm-56_net.droppedTx.summation":       200,
		"vm-56_net.packetsRx.summation":       200,
		"vm-56_net.packetsTx.summation":       200,
		"vm-56_sys.uptime.latest":             200,
		"vm-59_cpu.usage.average":             200,
		"vm-59_disk.maxTotalLatency.latest":   200,
		"vm-59_disk.read.average":             200,
		"vm-59_disk.write.average":            200,
		"vm-59_mem.active.average":            200,
		"vm-59_mem.consumed.average":          200,
		"vm-59_mem.granted.average":           200,
		"vm-59_mem.shared.average":            200,
		"vm-59_mem.swapinRate.average":        200,
		"vm-59_mem.swapoutRate.average":       200,
		"vm-59_mem.swapped.average":           200,
		"vm-59_mem.usage.average":             200,
		"vm-59_net.bytesRx.average":           200,
		"vm-59_net.bytesTx.average":           200,
		"vm-59_net.droppedRx.summation":       200,
		"vm-59_net.droppedTx.summation":       200,
		"vm-59_net.packetsRx.summation":       200,
		"vm-59_net.packetsTx.summation":       200,
		"vm-59_sys.uptime.latest":             200,
		"vm-62_cpu.usage.average":             200,
		"vm-62_disk.maxTotalLatency.latest":   200,
		"vm-62_disk.read.average":             200,
		"vm-62_disk.write.average":            200,
		"vm-62_mem.active.average":            200,
		"vm-62_mem.consumed.average":          200,
		"vm-62_mem.granted.average":           200,
		"vm-62_mem.shared.average":            200,
		"vm-62_mem.swapinRate.average":        200,
		"vm-62_mem.swapoutRate.average":       200,
		"vm-62_mem.swapped.average":           200,
		"vm-62_mem.usage.average":             200,
		"vm-62_net.bytesRx.average":           200,
		"vm-62_net.bytesTx.average":           200,
		"vm-62_net.droppedRx.summation":       200,
		"vm-62_net.droppedTx.summation":       200,
		"vm-62_net.packetsRx.summation":       200,
		"vm-62_net.packetsTx.summation":       200,
		"vm-62_sys.uptime.latest":             200,
	}
	assert.Equal(t, expected, job.Collect())

	count := model.Count()
	assert.Len(t, job.discoveredHosts, count.Host)
	assert.Len(t, job.discoveredVMs, count.Machine)
	assert.Len(t, job.charted, count.Host+count.Machine)
}

type testMetricScraper struct {
	metricScraper
}

func (s testMetricScraper) ScrapeHostsMetrics(hosts rs.Hosts) []performance.EntityMetric {
	ms := s.metricScraper.ScrapeHostsMetrics(hosts)
	setValueInMetrics(ms, 100)
	return ms
}
func (s testMetricScraper) ScrapeVMsMetrics(vms rs.VMs) []performance.EntityMetric {
	ms := s.metricScraper.ScrapeVMsMetrics(vms)
	setValueInMetrics(ms, 200)
	return ms
}

func setValueInMetrics(ms []performance.EntityMetric, value int64) {
	for i := range ms {
		for ii := range ms[i].Value {
			v := &ms[i].Value[ii].Value
			if *v == nil {
				*v = append(*v, value)
			} else {
				(*v)[0] = value
			}
		}
	}
}
