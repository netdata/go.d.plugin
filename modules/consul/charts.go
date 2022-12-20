// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import (
	"fmt"

	"github.com/blang/semver/v4"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	_ = module.Priority + iota
	prioKVSApplyTime
	prioKVSApplyOperations
	prioTXNApplyTime
	prioTXNApplyOperations
	prioRaftCommitTime
	prioRaftCommitsRate

	prioServerLeadershipStatus
	prioRaftLeaderLastContactTime
	prioRaftLeaderElections
	prioRaftLeadershipTransitions

	prioAutopilotHealthStatus
	prioAutopilotFailureTolerance

	prioRPCRequests
	prioRPCRequestsExceeded
	prioRPCRequestsFailed

	prioRaftThreadMainSaturation
	prioRaftThreadFSMSaturation

	prioRaftFSMLastRestoreDuration
	prioRaftLeaderOldestLogAge
	prioRaftRPCInstallSnapshotTime

	prioBoltDBFreelistBytes
	prioBoltDBLogsPerBatch
	prioBoltDBStoreLogsTime

	prioMemoryAllocated
	prioMemorySys
	prioGCPauseTime

	prioServiceHealthCheckStatus
	prioNodeHealthCheckStatus
)

var (
	clientCharts = module.Charts{
		clientRPCRequestsRateChart.Copy(),
		clientRPCRequestsExceededRateChart.Copy(),
		clientRPCRequestsFailedRateChart.Copy(),

		memoryAllocatedChart.Copy(),
		memorySysChart.Copy(),
		gcPauseTimeChart.Copy(),
	}

	serverLeaderCharts = module.Charts{
		raftCommitTimeChart.Copy(),
		raftLeaderLastContactTimeChart.Copy(),
		raftCommitsRateChart.Copy(),
		raftLeaderOldestLogAgeChart.Copy(),
	}
	serverFollowerCharts = module.Charts{
		raftRPCInstallSnapshotTimeChart.Copy(),
	}
	serverCommonCharts = module.Charts{
		kvsApplyTimeChart.Copy(),
		kvsApplyOperationsRateChart.Copy(),
		txnApplyTimeChart.Copy(),
		txnApplyOperationsRateChart.Copy(),

		autopilotHealthStatusChart.Copy(),
		autopilotFailureTolerance.Copy(),

		raftLeaderElectionsRateChart.Copy(),
		raftLeadershipTransitionsRateChart.Copy(),
		serverLeadershipStatusChart.Copy(),

		clientRPCRequestsRateChart.Copy(),
		clientRPCRequestsExceededRateChart.Copy(),
		clientRPCRequestsFailedRateChart.Copy(),

		raftThreadMainSaturationPercChart.Copy(),
		raftThreadFSMSaturationPercChart.Copy(),

		raftFSMLastRestoreDurationChart.Copy(),

		raftBoltDBFreelistBytesChart.Copy(),
		raftBoltDBLogsPerBatchChart.Copy(),
		raftBoltDBStoreLogsTimeChart.Copy(),

		memoryAllocatedChart.Copy(),
		memorySysChart.Copy(),
		gcPauseTimeChart.Copy(),
	}

	kvsApplyTimeChart = module.Chart{
		ID:       "kvs_apply_time",
		Title:    "KVS apply time",
		Units:    "ms",
		Fam:      "transaction timing",
		Ctx:      "consul.kvs_apply_time",
		Priority: prioKVSApplyTime,
		Dims: module.Dims{
			{ID: "kvs_apply_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "kvs_apply_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "kvs_apply_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}
	kvsApplyOperationsRateChart = module.Chart{
		ID:       "kvs_apply_operations_rate",
		Title:    "KVS apply operations",
		Units:    "ops/s",
		Fam:      "transaction timing",
		Ctx:      "consul.kvs_apply_operations_rate",
		Priority: prioKVSApplyOperations,
		Dims: module.Dims{
			{ID: "kvs_apply_count", Name: "kvs_apply"},
		},
	}
	txnApplyTimeChart = module.Chart{
		ID:       "txn_apply_time",
		Title:    "Transaction apply time",
		Units:    "ms",
		Fam:      "transaction timing",
		Ctx:      "consul.txn_apply_time",
		Priority: prioTXNApplyTime,
		Dims: module.Dims{
			{ID: "txn_apply_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "txn_apply_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "txn_apply_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}
	txnApplyOperationsRateChart = module.Chart{
		ID:       "txn_apply_operations_rate",
		Title:    "Transaction apply operations",
		Units:    "ops/s",
		Fam:      "transaction timing",
		Ctx:      "consul.txn_apply_operations_rate",
		Priority: prioTXNApplyOperations,
		Dims: module.Dims{
			{ID: "txn_apply_count", Name: "kvs_apply"},
		},
	}

	raftCommitTimeChart = module.Chart{
		ID:       "raft_commit_time",
		Title:    "Raft commit time",
		Units:    "ms",
		Fam:      "transaction timing",
		Ctx:      "consul.raft_commit_time",
		Priority: prioRaftCommitTime,
		Dims: module.Dims{
			{ID: "raft_commitTime_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "raft_commitTime_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "raft_commitTime_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}
	raftCommitsRateChart = module.Chart{
		ID:       "raft_commits_rate",
		Title:    "Raft commits rate",
		Units:    "commits/s",
		Fam:      "transaction timing",
		Ctx:      "consul.raft_commits_rate",
		Priority: prioRaftCommitsRate,
		Dims: module.Dims{
			{ID: "raft_apply", Name: "commits", Div: precision, Algo: module.Incremental},
		},
	}

	autopilotHealthStatusChart = module.Chart{
		ID:       "autopilot_health_status",
		Title:    "Autopilot health status",
		Units:    "status",
		Fam:      "autopilot",
		Ctx:      "consul.autopilot_health_status",
		Priority: prioAutopilotHealthStatus,
		Dims: module.Dims{
			{ID: "autopilot_healthy_yes", Name: "healthy"},
			{ID: "autopilot_healthy_no", Name: "unhealthy"},
		},
	}
	autopilotFailureTolerance = module.Chart{
		ID:       "autopilot_failure_tolerance",
		Title:    "Autopilot failure tolerance",
		Units:    "servers",
		Fam:      "autopilot",
		Ctx:      "consul.autopilot_failure_tolerance",
		Priority: prioAutopilotFailureTolerance,
		Dims: module.Dims{
			{ID: "autopilot_failure_tolerance", Name: "failure_tolerance"},
		},
	}

	raftLeaderLastContactTimeChart = module.Chart{
		ID:       "raft_leader_last_contact_time",
		Title:    "Raft leader last contact time",
		Units:    "ms",
		Fam:      "leadership changes",
		Ctx:      "consul.raft_leader_last_contact_time",
		Priority: prioRaftLeaderLastContactTime,
		Dims: module.Dims{
			{ID: "raft_leader_lastContact_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "raft_leader_lastContact_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "raft_leader_lastContact_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}
	raftLeaderElectionsRateChart = module.Chart{
		ID:       "raft_leader_elections_rate",
		Title:    "Raft leader elections rate",
		Units:    "elections/s",
		Fam:      "leadership changes",
		Ctx:      "consul.raft_leader_elections_rate",
		Priority: prioRaftLeaderElections,
		Dims: module.Dims{
			{ID: "raft_state_candidate", Name: "leader", Algo: module.Incremental},
		},
	}
	raftLeadershipTransitionsRateChart = module.Chart{
		ID:       "raft_leadership_transitions_rate",
		Title:    "Raft leadership transitions rate",
		Units:    "transitions/s",
		Fam:      "leadership changes",
		Ctx:      "consul.raft_leadership_transitions_rate",
		Priority: prioRaftLeadershipTransitions,
		Dims: module.Dims{
			{ID: "raft_state_leader", Name: "leadership", Algo: module.Incremental},
		},
	}
	serverLeadershipStatusChart = module.Chart{
		ID:       "server_leadership_status",
		Title:    "Server leadership status",
		Units:    "status",
		Fam:      "leadership changes",
		Ctx:      "consul.server_leadership_status",
		Priority: prioServerLeadershipStatus,
		Dims: module.Dims{
			{ID: "server_isLeader_yes", Name: "leader"},
			{ID: "server_isLeader_no", Name: "not_leader"},
		},
	}

	clientRPCRequestsRateChart = module.Chart{
		ID:       "client_rpc_requests_rate",
		Title:    "Client RPC requests",
		Units:    "requests/s",
		Fam:      "rpc network activity",
		Ctx:      "consul.client_rpc_requests_rate",
		Priority: prioRPCRequests,
		Dims: module.Dims{
			{ID: "client_rpc", Name: "rpc", Algo: module.Incremental},
		},
	}
	clientRPCRequestsExceededRateChart = module.Chart{
		ID:       "client_rpc_requests_exceeded_rate",
		Title:    "Client rate-limited RPC requests",
		Units:    "requests/s",
		Fam:      "rpc network activity",
		Ctx:      "consul.client_rpc_requests_exceeded_rate",
		Priority: prioRPCRequestsExceeded,
		Dims: module.Dims{
			{ID: "client_rpc_exceeded", Name: "exceeded", Algo: module.Incremental},
		},
	}
	clientRPCRequestsFailedRateChart = module.Chart{
		ID:       "client_rpc_requests_failed_rate",
		Title:    "Client failed RPC requests",
		Units:    "requests/s",
		Fam:      "rpc network activity",
		Ctx:      "consul.client_rpc_requests_failed_rate",
		Priority: prioRPCRequestsFailed,
		Dims: module.Dims{
			{ID: "client_rpc_failed", Name: "failed", Algo: module.Incremental},
		},
	}

	raftThreadMainSaturationPercChart = module.Chart{
		ID:       "raft_thread_main_saturation_perc",
		Title:    "Raft main thread saturation",
		Units:    "percentage",
		Fam:      "raft saturation",
		Ctx:      "consul.raft_thread_main_saturation_perc",
		Priority: prioRaftThreadMainSaturation,
		Dims: module.Dims{
			{ID: "raft_thread_main_saturation_sum", Name: "saturation", Algo: module.Incremental, Div: precision},
		},
	}
	raftThreadFSMSaturationPercChart = module.Chart{
		ID:       "raft_thread_fsm_saturation_perc",
		Title:    "Raft FSM thread saturation",
		Units:    "percentage",
		Fam:      "raft saturation",
		Ctx:      "consul.raft_thread_fsm_saturation_perc",
		Priority: prioRaftThreadFSMSaturation,
		Dims: module.Dims{
			{ID: "raft_thread_fsm_saturation_sum", Name: "saturation", Algo: module.Incremental, Div: precision},
		},
	}

	raftFSMLastRestoreDurationChart = module.Chart{
		ID:       "raft_fsm_last_restore_duration",
		Title:    "Raft last restore duration",
		Units:    "ms",
		Fam:      "raft replication capacity",
		Ctx:      "consul.raft_fsm_last_restore_duration",
		Priority: prioRaftFSMLastRestoreDuration,
		Dims: module.Dims{
			{ID: "raft_fsm_lastRestoreDuration", Name: "last_restore_duration"},
		},
	}
	raftLeaderOldestLogAgeChart = module.Chart{
		ID:       "raft_leader_oldest_log_age",
		Title:    "Raft leader oldest log age",
		Units:    "seconds",
		Fam:      "raft replication capacity",
		Ctx:      "consul.raft_leader_oldest_log_age",
		Priority: prioRaftLeaderOldestLogAge,
		Dims: module.Dims{
			{ID: "raft_leader_oldestLogAge", Name: "oldest_log_age", Div: 1000},
		},
	}
	raftRPCInstallSnapshotTimeChart = module.Chart{
		ID:       "raft_rpc_install_snapshot_time",
		Title:    "Raft RPC install snapshot time",
		Units:    "ms",
		Fam:      "raft replication capacity",
		Ctx:      "consul.raft_rpc_install_snapshot_time",
		Priority: prioRaftRPCInstallSnapshotTime,
		Dims: module.Dims{
			{ID: "raft_rpc_installSnapshot_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "raft_rpc_installSnapshot_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "raft_rpc_installSnapshot_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}

	raftBoltDBFreelistBytesChart = module.Chart{
		ID:       "raft_boltdb_freelist_bytes",
		Title:    "Raft BoltDB freelist",
		Units:    "bytes",
		Fam:      "boltdb performance",
		Ctx:      "consul.raft_boltdb_freelist_bytes",
		Priority: prioBoltDBFreelistBytes,
		Dims: module.Dims{
			{ID: "raft_boltdb_freelistBytes", Name: "freelist"},
		},
	}
	raftBoltDBLogsPerBatchChart = module.Chart{
		ID:       "raft_boltdb_logs_per_batch_rate",
		Title:    "Raft BoltDB logs written per batch",
		Units:    "logs/s",
		Fam:      "boltdb performance",
		Ctx:      "consul.raft_boltdb_logs_per_batch_rate",
		Priority: prioBoltDBLogsPerBatch,
		Dims: module.Dims{
			{ID: "raft_boltdb_logsPerBatch_sum", Name: "written", Algo: module.Incremental},
		},
	}

	raftBoltDBStoreLogsTimeChart = module.Chart{
		ID:       "raft_boltdb_store_logs_time",
		Title:    "Raft BoltDB store logs time",
		Units:    "ms",
		Fam:      "boltdb performance",
		Ctx:      "consul.raft_boltdb_store_logs_time",
		Priority: prioBoltDBStoreLogsTime,
		Dims: module.Dims{
			{ID: "raft_boltdb_storeLogs_quantile=0.5", Name: "quantile_0.5", Div: precision * precision},
			{ID: "raft_boltdb_storeLogs_quantile=0.9", Name: "quantile_0.9", Div: precision * precision},
			{ID: "raft_boltdb_storeLogs_quantile=0.99", Name: "quantile_0.99", Div: precision * precision},
		},
	}

	memoryAllocatedChart = module.Chart{
		ID:       "memory_allocated",
		Title:    "Memory allocated by the Consul process",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "consul.memory_allocated",
		Priority: prioMemoryAllocated,
		Dims: module.Dims{
			{ID: "runtime_alloc_bytes", Name: "allocated"},
		},
	}
	memorySysChart = module.Chart{
		ID:       "memory_sys",
		Title:    "Memory obtained from the OS",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "consul.memory_sys",
		Priority: prioMemorySys,
		Dims: module.Dims{
			{ID: "runtime_sys_bytes", Name: "sys"},
		},
	}

	gcPauseTimeChart = module.Chart{
		ID:       "gc_pause_time",
		Title:    "Garbage collection stop-the-world pause time",
		Units:    "seconds",
		Fam:      "garbage collection",
		Ctx:      "consul.gc_pause_time",
		Priority: prioGCPauseTime,
		Dims: module.Dims{
			{ID: "runtime_total_gc_pause_ns", Name: "gc_pause", Algo: module.Incremental, Div: 1e9},
		},
	}
)

var (
	serviceHealthCheckStatusChartTmpl = module.Chart{
		ID:       "health_check_%s_status",
		Title:    "Service health check status",
		Units:    "status",
		Fam:      "service health checks",
		Ctx:      "consul.service_health_check_status",
		Priority: prioServiceHealthCheckStatus,
		Dims: module.Dims{
			{ID: "health_check_%s_passing_status", Name: "passing"},
			{ID: "health_check_%s_critical_status", Name: "critical"},
			{ID: "health_check_%s_maintenance_status", Name: "maintenance"},
			{ID: "health_check_%s_warning_status", Name: "warning"},
		},
	}
	nodeHealthCheckStatusChartTmpl = module.Chart{
		ID:       "health_check_%s_status",
		Title:    "Node health check status",
		Units:    "status",
		Fam:      "node health checks",
		Ctx:      "consul.node_health_check_status",
		Priority: prioNodeHealthCheckStatus,
		Dims: module.Dims{
			{ID: "health_check_%s_passing_status", Name: "passing"},
			{ID: "health_check_%s_critical_status", Name: "critical"},
			{ID: "health_check_%s_maintenance_status", Name: "maintenance"},
			{ID: "health_check_%s_warning_status", Name: "warning"},
		},
	}
)

func (c *Consul) addGlobalCharts() {
	if !c.isTelemetryPrometheusEnabled() {
		return
	}

	var charts *module.Charts

	if !c.cfg.Config.Server {
		charts = clientCharts.Copy()
	} else {
		charts = serverCommonCharts.Copy()

		// can't really rely on checking if a response contains a metric due to retention of some metrics
		// https://github.com/hashicorp/go-metrics/blob/b6d5c860c07ef6eeec89f4a662c7b452dd4d0c93/prometheus/prometheus.go#L75-L76
		if c.version != nil {
			if c.version.LT(semver.Version{Major: 1, Minor: 13, Patch: 0}) {
				_ = charts.Remove(raftThreadMainSaturationPercChart.ID)
				_ = charts.Remove(raftThreadFSMSaturationPercChart.ID)
			}
			if c.version.LT(semver.Version{Major: 1, Minor: 11, Patch: 0}) {
				_ = charts.Remove(kvsApplyTimeChart.ID)
				_ = charts.Remove(kvsApplyOperationsRateChart.ID)
				_ = charts.Remove(txnApplyTimeChart.ID)
				_ = charts.Remove(txnApplyOperationsRateChart.ID)
				_ = charts.Remove(raftBoltDBFreelistBytesChart.ID)
			}
		}
	}

	for _, chart := range *charts {
		chart.Labels = []module.Label{
			{Key: "datacenter", Value: c.cfg.Config.Datacenter},
			{Key: "node_name", Value: c.cfg.Config.NodeName},
		}
	}

	if err := c.Charts().Add(*charts.Copy()...); err != nil {
		c.Warning(err)
	}
}

func newServiceHealthCheckChart(check *agentCheck) *module.Chart {
	chart := serviceHealthCheckStatusChartTmpl.Copy()
	chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
	chart.Labels = []module.Label{
		{Key: "node_name", Value: check.Node},
		{Key: "check_name", Value: check.Name},
		{Key: "service_name", Value: check.ServiceName},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, check.CheckID)
	}
	return chart
}

func newNodeHealthCheckChart(check *agentCheck) *module.Chart {
	chart := nodeHealthCheckStatusChartTmpl.Copy()
	chart.ID = fmt.Sprintf(chart.ID, check.CheckID)
	chart.Labels = []module.Label{
		{Key: "node_name", Value: check.Node},
		{Key: "check_name", Value: check.Name},
	}
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, check.CheckID)
	}
	return chart
}

func (c *Consul) addHealthCheckCharts(check *agentCheck) {
	var chart *module.Chart

	if check.ServiceName != "" {
		chart = newServiceHealthCheckChart(check)
	} else {
		chart = newNodeHealthCheckChart(check)
	}

	chart.Labels = append(chart.Labels, module.Label{
		Key:   "datacenter",
		Value: c.cfg.Config.Datacenter,
	})

	if err := c.Charts().Add(chart); err != nil {
		c.Warning(err)
	}
}

func (c *Consul) removeHealthCheckCharts(checkID string) {
	id := fmt.Sprintf("health_check_%s_status", checkID)

	chart := c.Charts().Get(id)
	if chart == nil {
		c.Warningf("failed to remove '%s' chart: the chart does not exist", id)
		return
	}

	chart.MarkRemove()
	chart.MarkNotCreated()
}

func (c *Consul) addLeaderCharts() {
	charts := serverLeaderCharts.Copy()

	if err := c.Charts().Add(*charts...); err != nil {
		c.Warning(err)
	}
}

func (c *Consul) removeLeaderCharts() {
	s := make(map[string]bool)
	for _, v := range serverLeaderCharts {
		s[v.ID] = true
	}

	for _, v := range *c.Charts() {
		if s[v.ID] {
			v.MarkRemove()
			v.MarkNotCreated()
		}
	}
}

func (c *Consul) addFollowerCharts() {
	charts := serverFollowerCharts.Copy()

	if err := c.Charts().Add(*charts...); err != nil {
		c.Warning(err)
	}
}

func (c *Consul) removeFollowerCharts() {
	s := make(map[string]bool)
	for _, v := range serverFollowerCharts {
		s[v.ID] = true
	}

	for _, v := range *c.Charts() {
		if s[v.ID] {
			v.MarkRemove()
			v.MarkNotCreated()
		}
	}
}
