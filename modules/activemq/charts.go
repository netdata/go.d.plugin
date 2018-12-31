package activemq

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{ID: "messages",
		Title: "Messages",
		Units: "messages/s",
		Fam:   "queues",
		Ctx:   "activemq.messages",
		Dims: Dims{
			{ID: "enqueued", Name: "enqueued", Algo: modules.Incremental},
			{ID: "dequeued", Name: "dequeued", Algo: modules.Incremental},
			{ID: "unprocessed", Name: "unprocessed", Algo: modules.Incremental},
		},
	},
	{ID: "consumers",
		Title: "Consumers",
		Units: "consumers",
		Fam:   "queues",
		Ctx:   "activemq.consumers",
		Dims: Dims{
			{ID: "consumers", Name: "consumers"},
		},
	},
}
