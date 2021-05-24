package redis

import "github.com/netdata/go.d.plugin/agent/module"

const (
	prioConnections = module.Priority + iota
	prioClients

	prioMemory
	prioMemoryFragmentationRatio

	prioNet

	prioConnectedReplicas

	prioPersistenceRDBChanges
	prioPersistenceRDBBgSaveNow
	prioPersistenceRDBBgSaveHealth
	prioPersistenceAOFSize

	prioCommands
	prioCommandsCalls
	prioCommandsUsec
	prioCommandsUsecPerSec

	prioKeyLookupHitRate
	prioKeyEviction
	prioKeyExpiration
	prioKeys
	prioExpiresKeys

	prioUptime
)

var redisCharts = module.Charts{
	chartConnections.Copy(),
	chartClients.Copy(),

	chartMemory.Copy(),
	chartMemoryFragmentationRatio.Copy(),

	chartNet.Copy(),

	chartConnectedReplicas.Copy(),

	chartPersistenceRDBChanges.Copy(),
	chartPersistenceRDBBgSaveNow.Copy(),
	chartPersistenceRDBBgSaveHealth.Copy(),

	chartCommands.Copy(),
	chartCommandsCalls.Copy(),
	chartCommandsUsec.Copy(),
	chartCommandsUsecPerSec.Copy(),

	chartKeyLookupHitRate.Copy(),
	chartKeyEviction.Copy(),
	chartKeyExpiration.Copy(),
	chartKeys.Copy(),
	chartExpiresKeys.Copy(),

	chartUptime.Copy(),
}

var (
	chartConnections = module.Chart{
		ID:       "connections",
		Title:    "Accepted and rejected (maxclients limit) connections",
		Units:    "connections/s",
		Fam:      "connections",
		Ctx:      "redis.connections",
		Priority: prioConnections,
		Dims: module.Dims{
			{ID: "total_connections_received", Name: "accepted", Algo: module.Incremental},
			{ID: "rejected_connections", Name: "rejected", Algo: module.Incremental},
		},
	}
	chartClients = module.Chart{
		ID:       "clients",
		Title:    "Clients",
		Units:    "clients",
		Fam:      "connections",
		Ctx:      "redis.clients",
		Priority: prioClients,
		Dims: module.Dims{
			{ID: "connected_clients", Name: "connected"},
			{ID: "blocked_clients", Name: "blocked"},
			{ID: "tracking_clients", Name: "tracking"},
			{ID: "clients_in_timeout_table", Name: "in_timeout_table"},
		},
	}
)

var (
	chartMemory = module.Chart{
		ID:       "memory",
		Title:    "Memory usage",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "redis.memory",
		Type:     module.Area,
		Priority: prioMemory,
		Dims: module.Dims{
			{ID: "maxmemory", Name: "max"},
			{ID: "used_memory", Name: "used"},
			{ID: "used_memory_rss", Name: "rss"},
			{ID: "used_memory_peak", Name: "peak"},
			{ID: "used_memory_dataset", Name: "dataset"},
			{ID: "used_memory_lua", Name: "lua"},
			{ID: "used_memory_scripts", Name: "scripts"},
		},
	}
	chartMemoryFragmentationRatio = module.Chart{
		ID:       "mem_fragmentation_ratio",
		Title:    "Ratio between used_memory_rss and used_memory",
		Units:    "ratio",
		Fam:      "memory",
		Ctx:      "redis.mem_fragmentation_ratio",
		Priority: prioMemoryFragmentationRatio,
		Dims: module.Dims{
			{ID: "mem_fragmentation_ratio", Name: "mem_fragmentation", Div: precision},
		},
	}
)

var (
	chartNet = module.Chart{
		ID:       "net",
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "network",
		Ctx:      "redis.net",
		Type:     module.Area,
		Priority: prioNet,
		Dims: module.Dims{
			{ID: "total_net_input_bytes", Name: "received", Mul: 8, Div: 1024, Algo: module.Incremental},
			{ID: "total_net_output_bytes", Name: "sent", Mul: -8, Div: 1024, Algo: module.Incremental},
		},
	}
)

var (
	chartPersistenceRDBChanges = module.Chart{
		ID:       "persistence",
		Title:    "Operations that produced changes since the last SAVE or BGSAVE",
		Units:    "operations",
		Fam:      "persistence rdb",
		Ctx:      "redis.rdb_changes",
		Priority: prioPersistenceRDBChanges,
		Dims: module.Dims{
			{ID: "rdb_changes_since_last_save", Name: "changes"},
		},
	}
	chartPersistenceRDBBgSaveNow = module.Chart{
		ID:       "bgsave_now",
		Title:    "Duration of the on-going RDB save operation if any",
		Units:    "seconds",
		Fam:      "persistence rdb",
		Ctx:      "redis.bgsave_now",
		Priority: prioPersistenceRDBBgSaveNow,
		Dims: module.Dims{
			{ID: "rdb_current_bgsave_time_sec", Name: "current_bgsave_time"},
		},
	}
	chartPersistenceRDBBgSaveHealth = module.Chart{
		ID:       "bgsave_health",
		Title:    "Status of the last RDB save operation (0: ok, 1: err)",
		Units:    "status",
		Fam:      "persistence rdb",
		Ctx:      "redis.bgsave_health",
		Priority: prioPersistenceRDBBgSaveHealth,
		Dims: module.Dims{
			{ID: "rdb_last_bgsave_status", Name: "last_bgsave"},
		},
	}

	chartPersistenceAOFSize = module.Chart{
		ID:       "persistence_aof_size",
		Title:    "AOF file size",
		Units:    "bytes",
		Fam:      "persistence aof",
		Ctx:      "redis.aof_file_size",
		Priority: prioPersistenceAOFSize,
		Dims: module.Dims{
			{ID: "aof_current_size", Name: "current"},
			{ID: "aof_base_size", Name: "base"},
		},
	}
)

var (
	chartCommands = module.Chart{
		ID:       "commands",
		Title:    "Processed commands",
		Units:    "commands/s",
		Fam:      "commands",
		Ctx:      "redis.commands",
		Priority: prioCommands,
		Dims: module.Dims{
			{ID: "total_commands_processed", Name: "processed", Algo: module.Incremental},
		},
	}
	chartCommandsCalls = module.Chart{
		ID:       "commands_calls",
		Title:    "Calls per command",
		Units:    "calls/s",
		Fam:      "commands",
		Ctx:      "redis.commands_calls",
		Type:     module.Stacked,
		Priority: prioCommandsCalls,
	}
	chartCommandsUsec = module.Chart{
		ID:       "commands_usec",
		Title:    "Total CPU time consumed by the commands",
		Units:    "microseconds",
		Fam:      "commands",
		Ctx:      "redis.commands_usec",
		Type:     module.Stacked,
		Priority: prioCommandsUsec,
	}
	chartCommandsUsecPerSec = module.Chart{
		ID:       "commands_usec_per_sec",
		Title:    "Average CPU consumed per command execution",
		Units:    "microseconds/s",
		Fam:      "commands",
		Ctx:      "redis.commands_usec_per_sec",
		Priority: prioCommandsUsecPerSec,
	}
)

var (
	chartKeyLookupHitRate = module.Chart{
		ID:       "key_lookup_hit_rate",
		Title:    "Keys lookup hit rate",
		Units:    "percentage",
		Fam:      "keyspace",
		Ctx:      "redis.keyspace_lookup_hit_rate",
		Priority: prioKeyLookupHitRate,
		Dims: module.Dims{
			{ID: "keyspace_hit_rate", Name: "lookup_hit_rate", Div: precision},
		},
	}
	chartKeyEviction = module.Chart{
		ID:       "key_eviction_events",
		Title:    "Evicted keys due to maxmemory limit",
		Units:    "keys/s",
		Fam:      "keyspace",
		Ctx:      "redis.key_eviction_events",
		Priority: prioKeyEviction,
		Dims: module.Dims{
			{ID: "evicted_keys", Name: "evicted", Algo: module.Incremental},
		},
	}
	chartKeyExpiration = module.Chart{
		ID:       "key_expiration_events",
		Title:    "Expired keys",
		Units:    "keys/s",
		Fam:      "keyspace",
		Ctx:      "redis.key_expiration_events",
		Priority: prioKeyExpiration,
		Dims: module.Dims{
			{ID: "expired_keys", Name: "expired", Algo: module.Incremental},
		},
	}
	chartKeys = module.Chart{
		ID:       "keys",
		Title:    "Keys per database",
		Units:    "keys",
		Fam:      "keyspace",
		Ctx:      "redis.database_keys",
		Type:     module.Stacked,
		Priority: prioKeys,
	}
	chartExpiresKeys = module.Chart{
		ID:       "expires_keys",
		Title:    "Keys with an expiration per database",
		Units:    "keys",
		Fam:      "keyspace",
		Ctx:      "redis.database_expires_keys",
		Type:     module.Stacked,
		Priority: prioExpiresKeys,
	}
)

var (
	chartConnectedReplicas = module.Chart{
		ID:       "connected_replicas",
		Title:    "Connected replicas",
		Units:    "replicas",
		Fam:      "replication",
		Ctx:      "redis.connected_replicas",
		Priority: prioConnectedReplicas,
		Dims: module.Dims{
			{ID: "connected_slaves", Name: "connected"},
		},
	}
)

var (
	chartUptime = module.Chart{
		ID:       "uptime",
		Title:    "Uptime",
		Units:    "seconds",
		Fam:      "uptime",
		Ctx:      "redis.uptime",
		Priority: prioUptime,
		Dims: module.Dims{
			{ID: "uptime_in_seconds", Name: "uptime"},
		},
	}
)
