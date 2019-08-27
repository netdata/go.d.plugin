package redis

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

const (
	keysRedisChartId = "keys_redis"
)

var charts = Charts{
	{
		ID:    "operations",
		Title: "Operations",
		Units: "operations/s",
		Fam:   "operations",
		Ctx:   "redis.operations",
		Type:  module.Line,
		Dims: Dims{
			{
				ID:   "total_commands_processed",
				Name: "commands",
				Algo: module.Incremental,
			},
			{ID: "instantaneous_ops_per_sec", Name: "operations"},
		},
	},
	{
		ID:    "hit_rate",
		Title: "Hit rate",
		Units: "percentage",
		Fam:   "hits",
		Ctx:   "redis.hit_rate",
		Type:  module.Line,
		Dims: Dims{
			{ID: "hit_rate", Name: "rate"},
		},
	},
	{
		ID:    "memory",
		Title: "Memory utilization",
		Units: "KiB",
		Fam:   "memory",
		Ctx:   "redis.memory",
		Type:  module.Line,
		Dims: Dims{
			{
				ID:   "used_memory",
				Name: "total",
				Div:  1024,
			},
			{
				ID:   "used_memory_lua",
				Name: "lua",
				Div:  1024,
			},
		},
	},
	{
		ID:    "net",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "network",
		Ctx:   "redis.net",
		Type:  module.Area,
		Dims: Dims{
			{ID: "total_net_input_bytes", Name: "in", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "total_net_output_bytes", Name: "out", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	},
	{
		ID:    keysRedisChartId,
		Title: "Keys per Database",
		Units: "keys",
		Fam:   "keys",
		Ctx:   "redis.keys",
		Type:  module.Line,
	},
	{
		ID:    "keys_pika",
		Title: "Keys",
		Units: "keys",
		Fam:   "keys",
		Ctx:   "redis.keys",
		Type:  module.Line,
		Dims: Dims{
			{ID: "kv_keys", Name: "kv"},
			{ID: "hash_keys", Name: "hash"},
			{ID: "list_keys", Name: "list"},
			{ID: "zset_keys", Name: "zset"},
			{ID: "set_keys", Name: "set"},
		},
	},
	{
		ID:    "eviction",
		Title: "Evicted Keys",
		Units: "keys",
		Fam:   "keys",
		Ctx:   "redis.eviction",
		Type:  module.Line,
		Dims: Dims{
			{ID: "evicted_keys", Name: "evicted"},
		},
	},
	{
		ID:    "connections",
		Title: "Connections",
		Units: "connections/s",
		Fam:   "connections",
		Ctx:   "redis.connections",
		Type:  module.Line,
		Dims: Dims{
			{ID: "total_connections_received", Name: "received", Algo: module.Incremental, Mul: 1},
			{ID: "rejected_connections", Name: "rejected", Algo: module.Incremental, Mul: -1},
		},
	},
	{
		ID:    "clients",
		Title: "Clients",
		Units: "clients",
		Fam:   "connections",
		Ctx:   "redis.clients",
		Type:  module.Line,
		Dims: Dims{
			{ID: "connected_clients", Name: "connected", Mul: 1},
			{ID: "blocked_clients", Name: "blocked", Mul: -1},
		},
	},
	{
		ID:    "slaves",
		Title: "Slaves",
		Units: "slaves",
		Fam:   "replication",
		Ctx:   "redis.slaves",
		Type:  module.Line,
		Dims: Dims{
			{ID: "connected_slaves", Name: "connected"},
		},
	},
	{
		ID:    "persistence",
		Title: "Persistence Changes Since Last Save",
		Units: "changes",
		Fam:   "persistence",
		Ctx:   "redis.rdb_changes",
		Type:  module.Line,
		Dims: Dims{
			{ID: "rdb_changes_since_last_save", Name: "changes"},
		},
	},
	{
		ID:    "bgsave_now",
		Title: "Duration of the RDB Save Operation",
		Units: "seconds",
		Fam:   "persistence",
		Ctx:   "redis.bgsave_now",
		// @TODO the "Type" on python was "absolute", couldn't find an "absolute" chart so I'm putting "Line"
		Type: module.Line,
		Dims: Dims{
			{ID: "rdb_changes_since_last_save", Name: "rdb save"},
		},
	},
	{
		ID:    "bgsave_health",
		Title: "Status of the Last RDB Save Operation",
		Units: "status",
		Fam:   "persistence",
		Ctx:   "redis.bgsave_health",
		Type:  module.Line,
		Dims: Dims{
			{ID: "rdb_last_bgsave_status", Name: "rdb save"},
		},
	},
	{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "redis.uptime",
		Type:  module.Line,
		Dims: Dims{
			{ID: "uptime_in_seconds", Name: "uptime"},
		},
	},
}
