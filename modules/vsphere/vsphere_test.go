// SPDX-License-Identifier: GPL-3.0-or-later
package vsphere

import (
	"crypto/tls"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/discover"
	"github.com/netdata/go.d.plugin/modules/vsphere/match"
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/simulator"
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
}

func TestVSphere_Init(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()

	assert.True(t, vSphere.Init())
	assert.NotNil(t, vSphere.discoverer)
	assert.NotNil(t, vSphere.scraper)
	assert.NotNil(t, vSphere.resources)
	assert.NotNil(t, vSphere.discoveryTask)
	assert.True(t, vSphere.discoveryTask.isRunning())
}

func TestVSphere_Init_ReturnsFalseIfURLNotSet(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()
	vSphere.URL = ""

	assert.False(t, vSphere.Init())
}

func TestVSphere_Init_ReturnsFalseIfUsernameNotSet(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()
	vSphere.Username = ""

	assert.False(t, vSphere.Init())
}

func TestVSphere_Init_ReturnsFalseIfPasswordNotSet(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()
	vSphere.Password = ""

	assert.False(t, vSphere.Init())
}

func TestVSphere_Init_ReturnsFalseIfClientWrongTLSCA(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()
	vSphere.Client.TLSConfig.TLSCA = "testdata/tls"

	assert.False(t, vSphere.Init())
}

func TestVSphere_Init_ReturnsFalseIfConnectionRefused(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()
	vSphere.URL = "http://127.0.0.1:32001"

	assert.False(t, vSphere.Init())
}

func TestVSphere_Init_ReturnsFalseIfInvalidHostVMIncludeFormat(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()

	vSphere.HostsInclude = match.HostIncludes{"invalid"}
	assert.False(t, vSphere.Init())

	vSphere.HostsInclude = vSphere.HostsInclude[:0]

	vSphere.VMsInclude = match.VMIncludes{"invalid"}
	assert.False(t, vSphere.Init())
}

func TestVSphere_Check(t *testing.T) {
	assert.NotNil(t, New().Check())
}

func TestVSphere_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestVSphere_Cleanup(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()

	require.True(t, vSphere.Init())

	vSphere.Cleanup()
	time.Sleep(time.Second)
	assert.True(t, vSphere.discoveryTask.isStopped())
	assert.False(t, vSphere.discoveryTask.isRunning())
}

func TestVSphere_Cleanup_NotPanicsIfNotInitialized(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestVSphere_Collect(t *testing.T) {
	vSphere, model, teardown := prepareVSphereSim(t)
	defer teardown()

	require.True(t, vSphere.Init())

	vSphere.scraper = mockScraper{vSphere.scraper}

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
		"host-20_overall.status.gray":         1,
		"host-20_overall.status.green":        0,
		"host-20_overall.status.red":          0,
		"host-20_overall.status.yellow":       0,
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
		"host-32_overall.status.gray":         1,
		"host-32_overall.status.green":        0,
		"host-32_overall.status.red":          0,
		"host-32_overall.status.yellow":       0,
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
		"host-39_overall.status.gray":         1,
		"host-39_overall.status.green":        0,
		"host-39_overall.status.red":          0,
		"host-39_overall.status.yellow":       0,
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
		"host-46_overall.status.gray":         1,
		"host-46_overall.status.green":        0,
		"host-46_overall.status.red":          0,
		"host-46_overall.status.yellow":       0,
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
		"vm-53_overall.status.gray":           0,
		"vm-53_overall.status.green":          1,
		"vm-53_overall.status.red":            0,
		"vm-53_overall.status.yellow":         0,
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
		"vm-56_overall.status.gray":           0,
		"vm-56_overall.status.green":          1,
		"vm-56_overall.status.red":            0,
		"vm-56_overall.status.yellow":         0,
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
		"vm-59_overall.status.gray":           0,
		"vm-59_overall.status.green":          1,
		"vm-59_overall.status.red":            0,
		"vm-59_overall.status.yellow":         0,
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
		"vm-62_overall.status.gray":           0,
		"vm-62_overall.status.green":          1,
		"vm-62_overall.status.red":            0,
		"vm-62_overall.status.yellow":         0,
		"vm-62_sys.uptime.latest":             200,
	}

	collected := vSphere.Collect()
	require.Equal(t, expected, collected)

	count := model.Count()
	assert.Len(t, vSphere.discoveredHosts, count.Host)
	assert.Len(t, vSphere.discoveredVMs, count.Machine)
	assert.Len(t, vSphere.charted, count.Host+count.Machine)

	assert.Len(t, *vSphere.Charts(), count.Host*len(hostChartsTmpl)+count.Machine*len(vmChartsTmpl))
	ensureCollectedHasAllChartsDimsVarsIDs(t, vSphere, collected)
}

func TestVSphere_Collect_RemoveHostsVMsInRuntime(t *testing.T) {
	vSphere, _, teardown := prepareVSphereSim(t)
	defer teardown()

	require.True(t, vSphere.Init())
	require.True(t, vSphere.Check())

	okHostID := "host-46"
	okVMID := "vm-62"
	vSphere.discoverer.(*discover.Discoverer).HostMatcher = mockHostMatcher{okHostID}
	vSphere.discoverer.(*discover.Discoverer).VMMatcher = mockVMMatcher{okVMID}

	require.NoError(t, vSphere.discoverOnce())

	numOfRuns := 5
	for i := 0; i < numOfRuns; i++ {
		vSphere.Collect()
	}

	host := vSphere.resources.Hosts.Get(okHostID)
	for k, v := range vSphere.discoveredHosts {
		if k == host.ID {
			assert.Equal(t, 0, v)
		} else {
			assert.Equal(t, numOfRuns, v)
		}
	}

	vm := vSphere.resources.VMs.Get(okVMID)
	for id, fails := range vSphere.discoveredVMs {
		if id == vm.ID {
			assert.Equal(t, 0, fails)
		} else {
			assert.Equal(t, numOfRuns, fails)
		}

	}

	for i := numOfRuns; i < failedUpdatesLimit; i++ {
		vSphere.Collect()
	}

	assert.Len(t, vSphere.discoveredHosts, 1)
	assert.Len(t, vSphere.discoveredVMs, 1)
	assert.Len(t, vSphere.charted, 2)

	for _, c := range *vSphere.Charts() {
		if strings.HasPrefix(c.ID, okHostID) || strings.HasPrefix(c.ID, okVMID) {
			assert.False(t, c.Obsolete)
		} else {
			assert.True(t, c.Obsolete)
		}
	}
}

func TestVSphere_Collect_Run(t *testing.T) {
	vSphere, model, teardown := prepareVSphereSim(t)
	defer teardown()

	vSphere.DiscoveryInterval.Duration = time.Second * 2
	require.True(t, vSphere.Init())
	require.True(t, vSphere.Check())

	runs := 20
	for i := 0; i < runs; i++ {
		assert.True(t, len(vSphere.Collect()) > 0)
		if i < 6 {
			time.Sleep(time.Second)
		}
	}

	count := model.Count()
	assert.Len(t, vSphere.discoveredHosts, count.Host)
	assert.Len(t, vSphere.discoveredVMs, count.Machine)
	assert.Len(t, vSphere.charted, count.Host+count.Machine)
	assert.Len(t, *vSphere.charts, count.Host*len(hostChartsTmpl)+count.Machine*len(vmChartsTmpl))
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, vSphere *VSphere, collected map[string]int64) {
	for _, chart := range *vSphere.Charts() {
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareVSphereSim(t *testing.T) (vSphere *VSphere, model *simulator.Model, teardown func()) {
	model, srv := createSim(t)
	vSphere = New()
	teardown = func() { model.Remove(); srv.Close(); vSphere.Cleanup() }

	vSphere.Username = "administrator"
	vSphere.Password = "password"
	vSphere.URL = srv.URL.String()
	vSphere.TLSConfig.InsecureSkipVerify = true

	return vSphere, model, teardown
}

func createSim(t *testing.T) (*simulator.Model, *simulator.Server) {
	model := simulator.VPX()
	err := model.Create()
	require.NoError(t, err)
	model.Service.TLS = new(tls.Config)
	return model, model.Service.NewServer()
}

type mockScraper struct {
	scraper
}

func (s mockScraper) ScrapeHosts(hosts rs.Hosts) []performance.EntityMetric {
	ms := s.scraper.ScrapeHosts(hosts)
	return populateMetrics(ms, 100)
}
func (s mockScraper) ScrapeVMs(vms rs.VMs) []performance.EntityMetric {
	ms := s.scraper.ScrapeVMs(vms)
	return populateMetrics(ms, 200)
}

func populateMetrics(ms []performance.EntityMetric, value int64) []performance.EntityMetric {
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
	return ms
}

type mockHostMatcher struct{ name string }
type mockVMMatcher struct{ name string }

func (m mockHostMatcher) Match(host *rs.Host) bool { return m.name == host.ID }
func (m mockVMMatcher) Match(vm *rs.VM) bool       { return m.name == vm.ID }
