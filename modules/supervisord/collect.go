package supervisord

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	// http://supervisord.org/subprocess.html#process-states
	// STOPPED  (0)
	// STARTING (10)
	// RUNNING (20)
	// BACKOFF (30)
	// STOPPING (40)
	// EXITED (100)
	// FATAL (200)
	// UNKNOWN (1000)

	stateRunning = 20
)

func (s *Supervisord) collect() (map[string]int64, error) {
	info, err := s.client.getAllProcessInfo()
	if err != nil {
		return nil, err
	}

	ms := make(map[string]int64)
	s.collectAllProcessInfo(ms, info)

	return ms, nil
}

func (s *Supervisord) collectAllProcessInfo(ms map[string]int64, info []processStatus) {
	ms["running_processes"] = 0
	ms["non_running_processes"] = 0
	for _, p := range info {
		if !s.collectedGroups[p.group] {
			s.collectedGroups[p.group] = true
			s.addProcessGroupCharts(p)
		}
		id := procID(p)
		if !s.collectedProcesses[id] {
			s.collectedProcesses[id] = true
			s.addProcessToCharts(p)
		}

		ms[fmt.Sprintf("group_%s_running_processes", p.group)] += 0
		ms[fmt.Sprintf("group_%s_non_running_processes", p.group)] += 0
		if p.state == stateRunning {
			ms["running_processes"] += 1
			ms[fmt.Sprintf("group_%s_running_processes", p.group)] += 1
		} else {
			ms["non_running_processes"] += 1
			ms[fmt.Sprintf("group_%s_non_running_processes", p.group)] += 1
		}
		ms[id+"_state"] = int64(p.state)
		ms[id+"_exit_status"] = int64(p.exitStatus)
		ms[id+"_uptime"] = calcProcessUptime(p)
		ms[id+"_downtime"] = calcProcessDowntime(p)
	}
}

func calcProcessUptime(p processStatus) int64 {
	if p.state != stateRunning {
		return 0
	}
	return int64(p.now - p.start)
}

func calcProcessDowntime(p processStatus) int64 {
	if p.state == stateRunning || p.stop == 0 {
		return 0
	}
	return int64(p.now - p.stop)
}

func (s *Supervisord) addProcessGroupCharts(p processStatus) {
	charts := newProcGroupCharts(p.group)
	if err := s.Charts().Add(*charts...); err != nil {
		s.Warning(err)
	}
}

func (s *Supervisord) addProcessToCharts(p processStatus) {
	id := procID(p)
	for _, c := range *s.Charts() {
		var dimID string
		switch c.ID {
		case fmt.Sprintf(groupProcessesStateCodeChartTmpl.ID, p.group):
			dimID = id + "_state"
		case fmt.Sprintf(groupProcessesExitStatusChartTmpl.ID, p.group):
			dimID = id + "_exit_status"
		case fmt.Sprintf(groupProcessesUptimeChartTmpl.ID, p.group):
			dimID = id + "_uptime"
		case fmt.Sprintf(groupProcessesDowntimeChartTmpl.ID, p.group):
			dimID = id + "_downtime"
		default:
			continue
		}
		dim := &module.Dim{ID: dimID, Name: p.name}
		if err := c.AddDim(dim); err != nil {
			s.Warning(err)
			return
		}
		c.MarkNotCreated()
	}
}

func (s *Supervisord) removeProcessFromCharts(p processStatus) {
	id := procID(p)
	for _, c := range *s.Charts() {
		var dimID string
		switch c.ID {
		case fmt.Sprintf(groupProcessesStateCodeChartTmpl.ID, p.group):
			dimID = id + "_state"
		case fmt.Sprintf(groupProcessesExitStatusChartTmpl.ID, p.group):
			dimID = id + "_exit_status"
		case fmt.Sprintf(groupProcessesUptimeChartTmpl.ID, p.group):
			dimID = id + "_uptime"
		case fmt.Sprintf(groupProcessesDowntimeChartTmpl.ID, p.group):
			dimID = id + "_downtime"
		default:
			continue
		}
		if err := c.MarkDimRemove(dimID, true); err != nil {
			s.Warning(err)
			return
		}
		c.MarkNotCreated()
	}
}

func procID(p processStatus) string {
	return fmt.Sprintf("group_%s_process_%s", p.group, p.name)
}
