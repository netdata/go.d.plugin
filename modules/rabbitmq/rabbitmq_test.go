package rabbitmq

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "queued_messages",
		Title: "Queued Messages",
		Units: "messages",
		Fam:   "overview",
		Ctx:   "rabbitmq.queued_messages",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "queue_totals_messages_ready", Name: "ready"},
			{ID: "queue_totals_messages_unacknowledged", Name: "unacknowledged"},
		},
	},
	{
		ID:    "message_rates",
		Title: "Message Rates",
		Units: "messages/s",
		Fam:   "overview",
		Ctx:   "rabbitmq.message_rates",
		Dims: Dims{
			{ID: "message_stats_ack", Name: "ack", Algo: modules.Incremental},
			{ID: "message_stats_redeliver", Name: "redeliver", Algo: modules.Incremental},
			{ID: "message_stats_deliver", Name: "deliver", Algo: modules.Incremental},
			{ID: "message_stats_publish", Name: "publish", Algo: modules.Incremental},
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
		Type:  modules.Area,
		Dims: Dims{
			{ID: "disk_free", Name: "free", Div: 1024 * 1024 * 1024},
		},
	},
}
