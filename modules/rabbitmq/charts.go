package rabbitmq

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

func charts() *Charts {
	c := Charts{}
	panicIfErr(c.Add(*overviewCharts.Copy()...))
	panicIfErr(c.Add(*nodeCharts.Copy()...))
	return &c
}

var (
	overviewCharts = Charts{
		{
			ID:    "queued_messages",
			Title: "Queued Messages",
			Units: "messages",
			Fam:   "overview",
			Ctx:   "rabbitmq.queued_messages",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "queue_totals_messages_ready", Name: "ready"},
				{ID: "queue_totals_messages_unacknowledged", Name: "unacknowledged"},
			},
		},
		{
			ID:    "message_rates",
			Title: "Messages",
			Units: "messages/s",
			Fam:   "overview",
			Ctx:   "rabbitmq.message_rates",
			Dims: Dims{
				{ID: "message_stats_ack", Name: "ack", Algo: module.Incremental},
				{ID: "message_stats_publish", Name: "publish", Algo: module.Incremental},
				{ID: "message_stats_publish_in", Name: "publish in", Algo: module.Incremental},
				{ID: "message_stats_publish_out", Name: "publish out", Algo: module.Incremental},
				{ID: "message_stats_confirm", Name: "confirm", Algo: module.Incremental},
				{ID: "message_stats_deliver", Name: "deliver", Algo: module.Incremental},
				{ID: "message_stats_deliver_no_ack", Name: "deliver no ack", Algo: module.Incremental},
				{ID: "message_stats_get", Name: "get", Algo: module.Incremental},
				{ID: "message_stats_get_no_ack", Name: "get no ack", Algo: module.Incremental},
				{ID: "message_stats_deliver_get", Name: "deliver get", Algo: module.Incremental},
				{ID: "message_stats_redeliver", Name: "redeliver", Algo: module.Incremental},
				{ID: "message_stats_return_unroutable", Name: "return unroutable", Algo: module.Incremental},
			},
		},
		{
			ID:    "global_counts",
			Title: "Global Counts",
			Units: "counts",
			Fam:   "overview",
			Ctx:   "rabbitmq.global_counts",
			Dims: Dims{
				{ID: "object_totals_channels", Name: "channels"},
				{ID: "object_totals_consumers", Name: "consumers"},
				{ID: "object_totals_connections", Name: "connections"},
				{ID: "object_totals_queues", Name: "queues"},
				{ID: "object_totals_exchanges", Name: "exchanges"},
			},
		},
	}

	nodeCharts = Charts{
		{
			ID:    "file_descriptors",
			Title: "File Descriptors",
			Units: "descriptors",
			Fam:   "overview",
			Ctx:   "rabbitmq.file_descriptors",
			Dims: Dims{
				{ID: "fd_used", Name: "used"},
			},
		},
		{
			ID:    "socket_descriptors",
			Title: "Socket Descriptors",
			Units: "descriptors",
			Fam:   "overview",
			Ctx:   "rabbitmq.sockets",
			Dims: Dims{
				{ID: "sockets_used", Name: "used"},
			},
		},
		{
			ID:    "erlang_processes",
			Title: "Erlang Processes",
			Units: "processes",
			Fam:   "overview",
			Ctx:   "rabbitmq.processes",
			Dims: Dims{
				{ID: "proc_used", Name: "used"},
			},
		},
		{
			ID:    "erlang_run_queue",
			Title: "Erlang Run Queue",
			Units: "processes",
			Fam:   "overview",
			Ctx:   "rabbitmq.erlang_run_queue",
			Dims: Dims{
				{ID: "run_queue", Name: "length"},
			},
		},
		{
			ID:    "memory",
			Title: "Memory",
			Units: "MiB",
			Fam:   "overview",
			Ctx:   "rabbitmq.memory",
			Dims: Dims{
				{ID: "mem_used", Name: "used", Div: 1024 << 10},
			},
		},
		{
			ID:    "disk_space",
			Title: "Disk Space",
			Units: "GiB",
			Fam:   "overview",
			Ctx:   "rabbitmq.disk_space",
			Type:  module.Area,
			Dims: Dims{
				{ID: "disk_free", Name: "free", Div: 1024 * 1024 * 1024},
			},
		},
	}

	vhostMessagesChart = Chart{
		ID:    "vhost_%s_message_stats",
		Title: "Vhost \"%s\" Messages",
		Units: "messages/s",
		Fam:   "vhost %s",
		Ctx:   "rabbitmq.vhost_messages",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "vhost_%s_message_stats_ack", Name: "ack", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_confirm", Name: "confirm", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_deliver", Name: "deliver", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_get", Name: "get", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_get_no_ack", Name: "get_no_ack", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_publish", Name: "publish", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_redeliver", Name: "redeliver", Algo: module.Incremental},
			{ID: "vhost_%s_message_stats_return_unroutable", Name: "return_unroutable", Algo: module.Incremental},
		},
	}
)

func (r *RabbitMQ) updateCharts(mx *metrics) {
	r.updateVhostsCharts(mx)
}

func (r *RabbitMQ) updateVhostsCharts(mx *metrics) {
	for _, v := range mx.vhosts {
		if v.MessageStats == nil {
			continue
		}
		if r.collectedVhosts[v.Name] {
			continue
		}
		r.collectedVhosts[v.Name] = true
		r.addVhostCharts(v.Name)
	}
}

func (r *RabbitMQ) addVhostCharts(name string) {
	chart := vhostMessagesChart.Copy()
	chart.ID = fmt.Sprintf(chart.ID, name)
	chart.Title = fmt.Sprintf(chart.Title, name)
	chart.Fam = fmt.Sprintf(chart.Fam, name)

	for _, dim := range chart.Dims {
		dim.ID = fmt.Sprintf(dim.ID, name)
	}

	err := r.charts.Add(chart)
	if err != nil {
		r.Warningf("error on adding '%s' chart : %v", chart.ID, err)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
