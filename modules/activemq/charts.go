package activemq

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "%s_%s_messages",
		Title: "%s Messages",
		Units: "messages/s",
		Fam:   "",
		Ctx:   "activemq.messages",
		Dims: Dims{
			{ID: "%s_%s_enqueued", Name: "enqueued", Algo: modules.Incremental},
			{ID: "%s_%s_dequeued", Name: "dequeued", Algo: modules.Incremental},
			{ID: "%s_%s_unprocessed", Name: "unprocessed", Algo: modules.Incremental},
		},
	},
	{
		ID:    "%s_%s_consumers",
		Title: "%s Consumers",
		Units: "consumers",
		Fam:   "",
		Ctx:   "activemq.consumers",
		Dims: Dims{
			{ID: "%s_%s_consumers", Name: "consumers"},
		},
	},
}
