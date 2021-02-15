package supervisord

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var summaryCharts = module.Charts{
	{
		ID:    "processes",
		Title: "Processes",
		Units: "processes",
		Fam:   "summary",
		Ctx:   "supervisord.summary_processes",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "running_processes", Name: "running"},
			{ID: "non_running_processes", Name: "non-running"},
		},
	},
}

var (
	groupChartsTmpl = module.Charts{
		groupProcessesChartTmpl.Copy(),
		groupProcessesStateCodeChartTmpl.Copy(),
		groupProcessesExitStatusChartTmpl.Copy(),
		groupProcessesUptimeChartTmpl.Copy(),
		groupProcessesDowntimeChartTmpl.Copy(),
	}

	groupProcessesChartTmpl = module.Chart{
		ID:    "group_%s_processes",
		Title: "Processes",
		Units: "processes",
		Fam:   "group %s",
		Ctx:   "supervisord.processes",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "group_%s_running_processes", Name: "running"},
			{ID: "group_%s_non_running_processes", Name: "non-running"},
		},
	}
	groupProcessesStateCodeChartTmpl = module.Chart{
		ID:    "group_%s_processes_state_code",
		Title: "Processes state code",
		Units: "state",
		Fam:   "group %s",
		Ctx:   "supervisord.process_state_code",
	}
	groupProcessesExitStatusChartTmpl = module.Chart{
		ID:    "group_%s_processes_exit_status",
		Title: "Processes exit status",
		Units: "status",
		Fam:   "group %s",
		Ctx:   "supervisord.process_exit_status",
	}
	groupProcessesUptimeChartTmpl = module.Chart{
		ID:    "group_%s_processes_uptime",
		Title: "Processes uptime",
		Units: "seconds",
		Fam:   "group %s",
		Ctx:   "supervisord.process_uptime",
	}
	groupProcessesDowntimeChartTmpl = module.Chart{
		ID:    "group_%s_processes_downtime",
		Title: "Processes downtime",
		Units: "seconds",
		Fam:   "group %s",
		Ctx:   "supervisord.process_downtime",
	}
)

func newProcGroupCharts(group string) *module.Charts {
	charts := groupChartsTmpl.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, group)
		c.Fam = fmt.Sprintf(c.Fam, group)
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, group)
		}
	}
	return charts
}
