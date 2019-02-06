package bind

import (
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts.
	Charts = module.Charts
	// Chart is an alias for module.Chart.
	Chart = module.Chart
	// Dims is an alias for module.Dims.
	Dims = module.Dims
	// Dim is an alias for module.Dim.
	Dim = module.Dim
)

const (
	keyReceivedRequests    = "received_requests"
	keyQueriesSuccess      = "queries_success"
	keyRecursiveClients    = "recursive_clients"
	keyProtocolsQueries    = "protocols_queries"
	keyQueriesAnalysis     = "queries_analysis"
	keyReceivedUpdates     = "received_updates"
	keyQueryFailures       = "query_failures"
	keyQueryFailuresDetail = "query_failures_detail"
	keyNSStats             = "nsstats"
	keyInOpCodes           = "in_opcodes"
	keyInQTypes            = "in_qtypes"
	keyInSockStats         = "in_sockstats"
)

var charts = map[string]Chart{
	keyReceivedRequests: {
		ID:    keyReceivedRequests,
		Title: "Global Received Requests by IP version",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "bind.requests",
		Type:  module.Stacked,
	},
	keyQueriesSuccess: {
		ID:    keyQueriesSuccess,
		Title: "Global Successful Queries",
		Units: "queries/s",
		Fam:   "queries",
		Ctx:   "bind.queries_success",
	},
	keyRecursiveClients: {
		ID:    keyRecursiveClients,
		Title: "Global Recursive Clients",
		Units: "clients",
		Fam:   "clients",
		Ctx:   "bind.recursive_clients",
	},
	keyProtocolsQueries: {
		ID:    keyProtocolsQueries,
		Title: "Global Queries by IP Protocol",
		Units: "queries/s",
		Fam:   "queries",
		Ctx:   "bind.protocol_queries",
		Type:  module.Stacked,
	},
	keyQueriesAnalysis: {
		ID:    keyQueriesAnalysis,
		Title: "Global Queries Analysis",
		Units: "queries/s",
		Fam:   "queries",
		Ctx:   "bind.global_queries",
		Type:  module.Stacked,
	},
	keyReceivedUpdates: {
		ID:    keyReceivedUpdates,
		Title: "Global Received Updates",
		Units: "updates/s",
		Fam:   "updates",
		Ctx:   "bind.global_updates",
		Type:  module.Stacked,
	},
	keyQueryFailures: {
		ID:    keyQueryFailures,
		Title: "Global Query Failures",
		Units: "failures/s",
		Fam:   "failures",
		Ctx:   "bind.global_failures",
	},
	keyQueryFailuresDetail: {
		ID:    keyQueryFailuresDetail,
		Title: "Global Query Failures Analysis",
		Units: "failures/s",
		Fam:   "failures",
		Ctx:   "bind.global_failures_detail",
		Type:  module.Stacked,
	},
	keyNSStats: {
		ID:    keyNSStats,
		Title: "Global Server Statistics",
		Units: "operations/s",
		Fam:   "other",
		Ctx:   "bind.nsstats",
	},
	keyInOpCodes: {
		ID:    keyInOpCodes,
		Title: "Incoming Requests by OpCode",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "bind.in_opcodes",
		Type:  module.Stacked,
	},
	keyInQTypes: {
		ID:    keyInQTypes,
		Title: "Incoming Requests by Query Type",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "bind.in_qtypes",
		Type:  module.Stacked,
	},
	keyInSockStats: {
		ID:    keyInSockStats,
		Title: "Socket Statistics",
		Units: "operations/s",
		Fam:   "sockets",
		Ctx:   "bind.in_sockstats",
	},
}
