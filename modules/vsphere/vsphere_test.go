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
		"host-20_DC0_H0_cpu.usage.average":                 100,
		"host-20_DC0_H0_disk.maxTotalLatency.latest":       100,
		"host-20_DC0_H0_disk.read.average":                 100,
		"host-20_DC0_H0_disk.write.average":                100,
		"host-20_DC0_H0_mem.active.average":                100,
		"host-20_DC0_H0_mem.consumed.average":              100,
		"host-20_DC0_H0_mem.granted.average":               100,
		"host-20_DC0_H0_mem.shared.average":                100,
		"host-20_DC0_H0_mem.sharedcommon.average":          100,
		"host-20_DC0_H0_mem.swapinRate.average":            100,
		"host-20_DC0_H0_mem.swapoutRate.average":           100,
		"host-20_DC0_H0_mem.usage.average":                 100,
		"host-20_DC0_H0_net.bytesRx.average":               100,
		"host-20_DC0_H0_net.bytesTx.average":               100,
		"host-20_DC0_H0_net.droppedRx.summation":           100,
		"host-20_DC0_H0_net.droppedTx.summation":           100,
		"host-20_DC0_H0_net.errorsRx.summation":            100,
		"host-20_DC0_H0_net.errorsTx.summation":            100,
		"host-20_DC0_H0_net.packetsRx.summation":           100,
		"host-20_DC0_H0_net.packetsTx.summation":           100,
		"host-20_DC0_H0_sys.uptime.latest":                 100,
		"host-32_DC0_C0_H0_cpu.usage.average":              100,
		"host-32_DC0_C0_H0_disk.maxTotalLatency.latest":    100,
		"host-32_DC0_C0_H0_disk.read.average":              100,
		"host-32_DC0_C0_H0_disk.write.average":             100,
		"host-32_DC0_C0_H0_mem.active.average":             100,
		"host-32_DC0_C0_H0_mem.consumed.average":           100,
		"host-32_DC0_C0_H0_mem.granted.average":            100,
		"host-32_DC0_C0_H0_mem.shared.average":             100,
		"host-32_DC0_C0_H0_mem.sharedcommon.average":       100,
		"host-32_DC0_C0_H0_mem.swapinRate.average":         100,
		"host-32_DC0_C0_H0_mem.swapoutRate.average":        100,
		"host-32_DC0_C0_H0_mem.usage.average":              100,
		"host-32_DC0_C0_H0_net.bytesRx.average":            100,
		"host-32_DC0_C0_H0_net.bytesTx.average":            100,
		"host-32_DC0_C0_H0_net.droppedRx.summation":        100,
		"host-32_DC0_C0_H0_net.droppedTx.summation":        100,
		"host-32_DC0_C0_H0_net.errorsRx.summation":         100,
		"host-32_DC0_C0_H0_net.errorsTx.summation":         100,
		"host-32_DC0_C0_H0_net.packetsRx.summation":        100,
		"host-32_DC0_C0_H0_net.packetsTx.summation":        100,
		"host-32_DC0_C0_H0_sys.uptime.latest":              100,
		"host-39_DC0_C0_H1_cpu.usage.average":              100,
		"host-39_DC0_C0_H1_disk.maxTotalLatency.latest":    100,
		"host-39_DC0_C0_H1_disk.read.average":              100,
		"host-39_DC0_C0_H1_disk.write.average":             100,
		"host-39_DC0_C0_H1_mem.active.average":             100,
		"host-39_DC0_C0_H1_mem.consumed.average":           100,
		"host-39_DC0_C0_H1_mem.granted.average":            100,
		"host-39_DC0_C0_H1_mem.shared.average":             100,
		"host-39_DC0_C0_H1_mem.sharedcommon.average":       100,
		"host-39_DC0_C0_H1_mem.swapinRate.average":         100,
		"host-39_DC0_C0_H1_mem.swapoutRate.average":        100,
		"host-39_DC0_C0_H1_mem.usage.average":              100,
		"host-39_DC0_C0_H1_net.bytesRx.average":            100,
		"host-39_DC0_C0_H1_net.bytesTx.average":            100,
		"host-39_DC0_C0_H1_net.droppedRx.summation":        100,
		"host-39_DC0_C0_H1_net.droppedTx.summation":        100,
		"host-39_DC0_C0_H1_net.errorsRx.summation":         100,
		"host-39_DC0_C0_H1_net.errorsTx.summation":         100,
		"host-39_DC0_C0_H1_net.packetsRx.summation":        100,
		"host-39_DC0_C0_H1_net.packetsTx.summation":        100,
		"host-39_DC0_C0_H1_sys.uptime.latest":              100,
		"host-46_DC0_C0_H2_cpu.usage.average":              100,
		"host-46_DC0_C0_H2_disk.maxTotalLatency.latest":    100,
		"host-46_DC0_C0_H2_disk.read.average":              100,
		"host-46_DC0_C0_H2_disk.write.average":             100,
		"host-46_DC0_C0_H2_mem.active.average":             100,
		"host-46_DC0_C0_H2_mem.consumed.average":           100,
		"host-46_DC0_C0_H2_mem.granted.average":            100,
		"host-46_DC0_C0_H2_mem.shared.average":             100,
		"host-46_DC0_C0_H2_mem.sharedcommon.average":       100,
		"host-46_DC0_C0_H2_mem.swapinRate.average":         100,
		"host-46_DC0_C0_H2_mem.swapoutRate.average":        100,
		"host-46_DC0_C0_H2_mem.usage.average":              100,
		"host-46_DC0_C0_H2_net.bytesRx.average":            100,
		"host-46_DC0_C0_H2_net.bytesTx.average":            100,
		"host-46_DC0_C0_H2_net.droppedRx.summation":        100,
		"host-46_DC0_C0_H2_net.droppedTx.summation":        100,
		"host-46_DC0_C0_H2_net.errorsRx.summation":         100,
		"host-46_DC0_C0_H2_net.errorsTx.summation":         100,
		"host-46_DC0_C0_H2_net.packetsRx.summation":        100,
		"host-46_DC0_C0_H2_net.packetsTx.summation":        100,
		"host-46_DC0_C0_H2_sys.uptime.latest":              100,
		"vm-53_DC0_H0_VM0_cpu.usage.average":               200,
		"vm-53_DC0_H0_VM0_disk.maxTotalLatency.latest":     200,
		"vm-53_DC0_H0_VM0_disk.read.average":               200,
		"vm-53_DC0_H0_VM0_disk.write.average":              200,
		"vm-53_DC0_H0_VM0_mem.active.average":              200,
		"vm-53_DC0_H0_VM0_mem.consumed.average":            200,
		"vm-53_DC0_H0_VM0_mem.granted.average":             200,
		"vm-53_DC0_H0_VM0_mem.shared.average":              200,
		"vm-53_DC0_H0_VM0_mem.swapinRate.average":          200,
		"vm-53_DC0_H0_VM0_mem.swapoutRate.average":         200,
		"vm-53_DC0_H0_VM0_mem.swapped.average":             200,
		"vm-53_DC0_H0_VM0_mem.usage.average":               200,
		"vm-53_DC0_H0_VM0_net.bytesRx.average":             200,
		"vm-53_DC0_H0_VM0_net.bytesTx.average":             200,
		"vm-53_DC0_H0_VM0_net.droppedRx.summation":         200,
		"vm-53_DC0_H0_VM0_net.droppedTx.summation":         200,
		"vm-53_DC0_H0_VM0_net.packetsRx.summation":         200,
		"vm-53_DC0_H0_VM0_net.packetsTx.summation":         200,
		"vm-53_DC0_H0_VM0_sys.uptime.latest":               200,
		"vm-56_DC0_H0_VM1_cpu.usage.average":               200,
		"vm-56_DC0_H0_VM1_disk.maxTotalLatency.latest":     200,
		"vm-56_DC0_H0_VM1_disk.read.average":               200,
		"vm-56_DC0_H0_VM1_disk.write.average":              200,
		"vm-56_DC0_H0_VM1_mem.active.average":              200,
		"vm-56_DC0_H0_VM1_mem.consumed.average":            200,
		"vm-56_DC0_H0_VM1_mem.granted.average":             200,
		"vm-56_DC0_H0_VM1_mem.shared.average":              200,
		"vm-56_DC0_H0_VM1_mem.swapinRate.average":          200,
		"vm-56_DC0_H0_VM1_mem.swapoutRate.average":         200,
		"vm-56_DC0_H0_VM1_mem.swapped.average":             200,
		"vm-56_DC0_H0_VM1_mem.usage.average":               200,
		"vm-56_DC0_H0_VM1_net.bytesRx.average":             200,
		"vm-56_DC0_H0_VM1_net.bytesTx.average":             200,
		"vm-56_DC0_H0_VM1_net.droppedRx.summation":         200,
		"vm-56_DC0_H0_VM1_net.droppedTx.summation":         200,
		"vm-56_DC0_H0_VM1_net.packetsRx.summation":         200,
		"vm-56_DC0_H0_VM1_net.packetsTx.summation":         200,
		"vm-56_DC0_H0_VM1_sys.uptime.latest":               200,
		"vm-59_DC0_C0_RP0_VM0_cpu.usage.average":           200,
		"vm-59_DC0_C0_RP0_VM0_disk.maxTotalLatency.latest": 200,
		"vm-59_DC0_C0_RP0_VM0_disk.read.average":           200,
		"vm-59_DC0_C0_RP0_VM0_disk.write.average":          200,
		"vm-59_DC0_C0_RP0_VM0_mem.active.average":          200,
		"vm-59_DC0_C0_RP0_VM0_mem.consumed.average":        200,
		"vm-59_DC0_C0_RP0_VM0_mem.granted.average":         200,
		"vm-59_DC0_C0_RP0_VM0_mem.shared.average":          200,
		"vm-59_DC0_C0_RP0_VM0_mem.swapinRate.average":      200,
		"vm-59_DC0_C0_RP0_VM0_mem.swapoutRate.average":     200,
		"vm-59_DC0_C0_RP0_VM0_mem.swapped.average":         200,
		"vm-59_DC0_C0_RP0_VM0_mem.usage.average":           200,
		"vm-59_DC0_C0_RP0_VM0_net.bytesRx.average":         200,
		"vm-59_DC0_C0_RP0_VM0_net.bytesTx.average":         200,
		"vm-59_DC0_C0_RP0_VM0_net.droppedRx.summation":     200,
		"vm-59_DC0_C0_RP0_VM0_net.droppedTx.summation":     200,
		"vm-59_DC0_C0_RP0_VM0_net.packetsRx.summation":     200,
		"vm-59_DC0_C0_RP0_VM0_net.packetsTx.summation":     200,
		"vm-59_DC0_C0_RP0_VM0_sys.uptime.latest":           200,
		"vm-62_DC0_C0_RP0_VM1_cpu.usage.average":           200,
		"vm-62_DC0_C0_RP0_VM1_disk.maxTotalLatency.latest": 200,
		"vm-62_DC0_C0_RP0_VM1_disk.read.average":           200,
		"vm-62_DC0_C0_RP0_VM1_disk.write.average":          200,
		"vm-62_DC0_C0_RP0_VM1_mem.active.average":          200,
		"vm-62_DC0_C0_RP0_VM1_mem.consumed.average":        200,
		"vm-62_DC0_C0_RP0_VM1_mem.granted.average":         200,
		"vm-62_DC0_C0_RP0_VM1_mem.shared.average":          200,
		"vm-62_DC0_C0_RP0_VM1_mem.swapinRate.average":      200,
		"vm-62_DC0_C0_RP0_VM1_mem.swapoutRate.average":     200,
		"vm-62_DC0_C0_RP0_VM1_mem.swapped.average":         200,
		"vm-62_DC0_C0_RP0_VM1_mem.usage.average":           200,
		"vm-62_DC0_C0_RP0_VM1_net.bytesRx.average":         200,
		"vm-62_DC0_C0_RP0_VM1_net.bytesTx.average":         200,
		"vm-62_DC0_C0_RP0_VM1_net.droppedRx.summation":     200,
		"vm-62_DC0_C0_RP0_VM1_net.droppedTx.summation":     200,
		"vm-62_DC0_C0_RP0_VM1_net.packetsRx.summation":     200,
		"vm-62_DC0_C0_RP0_VM1_net.packetsTx.summation":     200,
		"vm-62_DC0_C0_RP0_VM1_sys.uptime.latest":           200,
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
