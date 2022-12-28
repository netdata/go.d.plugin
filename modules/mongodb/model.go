// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import "time"

const (
	mongos = "mongos"
)

type serverStatus struct {
	// available in many versions and builds of mongo
	Opcounters  Opcounters   `bson:"opcounters" stm:"operations"`
	OpLatencies *OpLatencies `bson:"opLatencies" stm:"operations_latency" `
	Connections Connections  `bson:"connections" stm:"connections"`
	Network     Network      `bson:"network" stm:"network"`
	ExtraInfo   ExtraInfo    `bson:"extra_info" stm:"extra_info"`
	Asserts     Asserts      `bson:"asserts" stm:"asserts"`

	// available in newer or specific builds of mongo
	// for example, in he hosted version of mongoDB(atlas)
	// these are not available
	Transactions *Transactions         `bson:"transactions" stm:"transactions"`
	GlobalLock   *GlobalLock           `bson:"globalLock" stm:"glock"`
	Tcmalloc     *ServerStatusTcmalloc `bson:"tcmalloc" stm:"tcmalloc"`
	Locks        *Locks                `bson:"locks" stm:"locks"`
	FlowControl  *FlowControl          `bson:"flowControl" stm:"flow"`
	WiredTiger   *WiredTiger           `bson:"wiredTiger" stm:"wiredtiger"`
	Repl         interface{}           `bson:"repl"`
	Process      string                `bson:"process"` // mongod|mongos
}

type Opcounters struct {
	Insert  *int64 `bson:"insert" stm:"insert"`
	Query   *int64 `bson:"query" stm:"query"`
	Update  *int64 `bson:"update" stm:"update"`
	Delete  *int64 `bson:"delete" stm:"delete"`
	Getmore *int64 `bson:"getmore" stm:"getmore"`
	Command *int64 `bson:"command" stm:"command"`
}

type OpLatencies struct {
	Reads *struct {
		Latency *int64 `bson:"latency" stm:""`
	} `bson:"reads" stm:"read"`
	Writes *struct {
		Latency *int64 `bson:"latency" stm:""`
	} `bson:"writes" stm:"write"`
	Commands *struct {
		Latency *int64 `bson:"latency" stm:""`
	} `bson:"commands" stm:"command"`
}

type Connections struct {
	Current                 *int64 `bson:"current" stm:"current"`
	Available               *int64 `bson:"available" stm:"available"`
	TotalCreated            *int64 `bson:"totalCreated" stm:"total_created"`
	Active                  *int64 `bson:"active" stm:"active"`
	Threaded                *int64 `bson:"threaded" stm:"threaded"`
	ExhaustIsMaster         *int64 `bson:"exhaustIsMaster" stm:"exhaustIsMaster"`
	ExhaustHello            *int64 `bson:"exhaustHello" stm:"exhaustHello"`
	AwaitingTopologyChanges *int64 `bson:"awaitingTopologyChanges" stm:"awaitingTopologyChanges"`
}

type Network struct {
	BytesIn     *int64 `bson:"bytesIn" stm:"bytes_in"`
	BytesOut    *int64 `bson:"bytesOut" stm:"bytes_out"`
	NumRequests *int64 `bson:"numRequests" stm:"requests"`
}

type ExtraInfo struct {
	PageFaults *int64 `bson:"page_faults" stm:"page_faults"`
}

type Asserts struct {
	Regular   *int64 `bson:"regular" stm:"regular"`
	Warning   *int64 `bson:"warning" stm:"warning"`
	Msg       *int64 `bson:"msg" stm:"msg"`
	User      *int64 `bson:"user" stm:"user"`
	Tripwire  *int64 `bson:"tripwire" stm:"tripwire"`
	Rollovers *int64 `bson:"rollovers" stm:"rollovers"`
}

type Transactions struct {
	CurrentActive   *int64       `bson:"currentActive" stm:"active"`
	CurrentInactive *int64       `bson:"currentInactive" stm:"inactive"`
	CurrentOpen     *int64       `bson:"currentOpen" stm:"open"`
	CurrentPrepared *int64       `bson:"currentPrepared" stm:"prepared"`
	CommitTypes     *CommitTypes `bson:"commitTypes" stm:"commit_types"`
}

type CommitTypes struct {
	NoShards         CommitType `bson:"noShards" stm:"no_shards"`
	SingleShard      CommitType `bson:"singleShard" stm:"single_shard"`
	SingleWriteShard CommitType `bson:"singleWriteShard" stm:"single_write_shard"`
	TwoPhaseCommit   CommitType `bson:"twoPhaseCommit" stm:"two_phase"`
}

type CommitType struct {
	Initiated  int64 `json:"initiated" stm:"initiated"`
	Successful int64 `json:"successful" stm:"successful"`
}

type GlobalLock struct {
	ActiveClients *struct {
		Readers *int64 `bson:"readers" stm:"readers"`
		Writers *int64 `bson:"writers" stm:"writers"`
	} `bson:"activeClients" stm:"active_clients"`
	CurrentQueue *struct {
		Readers *int64 `bson:"readers" stm:"readers"`
		Writers *int64 `bson:"writers" stm:"writers"`
	} `bson:"currentQueue" stm:"current_queue"`
}

type ServerStatusTcmalloc struct {
	Generic  *Generic          `bson:"generic" stm:"generic"`
	Tcmalloc *TcmallocTcmalloc `bson:"tcmalloc" stm:"tcmalloc"`
}

type Generic struct {
	CurrentAllocatedBytes *int64 `bson:"current_allocated_bytes" stm:"current_allocated"`
	HeapSize              *int64 `bson:"heap_size" stm:"heap_size"`
}

type TcmallocTcmalloc struct {
	PageheapFreeBytes          *int64 `bson:"pageheap_free_bytes" stm:"pageheap_free"`
	PageheapUnmappedBytes      *int64 `bson:"pageheap_unmapped_bytes" stm:"pageheap_unmapped"`
	MaxTotalThreadCacheBytes   *int64 `bson:"max_total_thread_cache_bytes" stm:"max_total_thread_cache"`
	TotalFreeBytes             *int64 `bson:"total_free_bytes" stm:"total_free"`
	PageheapCommittedBytes     *int64 `bson:"pageheap_committed_bytes" stm:"pageheap_committed"`
	PageheapTotalCommitBytes   *int64 `bson:"pageheap_total_commit_bytes" stm:"pageheap_total_commit"`
	PageheapTotalDecommitBytes *int64 `bson:"pageheap_total_decommit_bytes" stm:"pageheap_total_decommit"`
	PageheapTotalReserveBytes  *int64 `bson:"pageheap_total_reserve_bytes" stm:"pageheap_total_reserve"`
}

type Locks struct {
	Global *struct {
		AcquireCount struct {
			R *int64 `bson:"r" stm:"read"`
			W *int64 `bson:"W" stm:"write"`
		} `bson:"acquireCount" stm:""`
	} `bson:"Global" stm:"global"`
	Database *struct {
		AcquireCount struct {
			R *int64 `bson:"r" stm:"read"`
			W *int64 `bson:"W" stm:"write"`
		} `bson:"acquireCount" stm:""`
	} `bson:"Database" stm:"database"`
	Collection *struct {
		AcquireCount struct {
			R *int64 `bson:"r" stm:"read"`
			W *int64 `bson:"W" stm:"write"`
		} `bson:"acquireCount" stm:""`
	} `bson:"Collection" stm:"collection"`
}

type FlowControl struct {
	TargetRateLimit     *int64 `bson:"targetRateLimit" stm:"target_rate_limit"`
	TimeAcquiringMicros *int64 `bson:"timeAcquiringMicros" stm:"time_acquiring_micros"`
}

type WiredTiger struct {
	BlockManager *struct {
		BytesRead                    int `bson:"bytes read" stm:"read"`
		BytesReadViaMemoryMapAPI     int `bson:"bytes read via memory map API" stm:"read_via_memory"`
		BytesReadViaSystemCallAPI    int `bson:"bytes read via system call API" stm:"read_via_system_api"`
		BytesWritten                 int `bson:"bytes written" stm:"written"`
		BytesWrittenForCheckpoint    int `bson:"bytes written for checkpoint" stm:"written_for_checkpoint"`
		BytesWrittenViaMemoryMapAPI  int `bson:"bytes written via memory map API" stm:"written_via_memory"`
		BytesWrittenViaSystemCallAPI int `bson:"bytes written via system call API" stm:"written_via_system_api"`
	} `bson:"block-manager" json:"block-manager" stm:"block_manager"`
	Cache *struct {
		BytesAllocatedForUpdates int `bson:"bytes allocated for updates" stm:"alloccated"`
		BytesReadIntoCache       int `bson:"bytes read into cache" stm:"read"`
		BytesWrittenFromCache    int `bson:"bytes written from cache" stm:"write"`
	} `bson:"cache" stm:"cache"`
	Capacity *struct {
		TimeWaitingDueToTotalCapacityUsecs int `bson:"time waiting due to total capacity (usecs)" stm:"wait_capacity"`
		TimeWaitingDuringCheckpointUsecs   int `bson:"time waiting during checkpoint (usecs)" stm:"wait_checkpoint"`
		TimeWaitingDuringEvictionUsecs     int `bson:"time waiting during eviction (usecs)" stm:"wait_eviction"`
		TimeWaitingDuringLoggingUsecs      int `bson:"time waiting during logging (usecs)" stm:"wait_logging"`
		TimeWaitingDuringReadUsecs         int `bson:"time waiting during read (usecs)" stm:"wait_read"`
	} `bson:"capacity" stm:"capacity"`
	Connection *struct {
		MemoryAllocations   int `bson:"memory allocations" stm:"allocations"`
		MemoryFrees         int `bson:"memory frees" stm:"frees"`
		MemoryReAllocations int `bson:"memory re-allocations" stm:"reallocations"`
	} `bson:"connection" stm:"connection"`
	Cursor *struct {
		CachedCursorCount                 int `bson:"cached cursor count" stm:"count"`
		CursorBulkLoadedCursorInsertCalls int `bson:"cursor bulk loaded cursor insert calls" stm:"bulk"`
		CursorCloseCallsThatResultInCache int `bson:"cursor close calls that result in cache" stm:"close"`
		CursorCreateCalls                 int `bson:"cursor create calls" stm:"create"`
		CursorInsertCalls                 int `bson:"cursor insert calls" stm:"insert"`
		CursorModifyCalls                 int `bson:"cursor modify calls" stm:"modify"`
		CursorNextCalls                   int `bson:"cursor next calls" stm:"next"`
		CursorOperationRestarted          int `bson:"cursor operation restarted" stm:"restarted"`
		CursorPrevCalls                   int `bson:"cursor prev calls" stm:"prev"`
		CursorRemoveCalls                 int `bson:"cursor remove calls" stm:"remove"`
		CursorReserveCalls                int `bson:"cursor reserve calls" stm:"reserve"`
		CursorResetCalls                  int `bson:"cursor reset calls" stm:"reset"`
		CursorSearchCalls                 int `bson:"cursor search calls" stm:"search"`
		CursorSearchHistoryStoreCalls     int `bson:"cursor search history store calls" stm:"search_history"`
		CursorSearchNearCalls             int `bson:"cursor search near calls" stm:"search_near"`
		CursorSweepBuckets                int `bson:"cursor sweep buckets" stm:"sweep_buckets"`
		CursorSweepCursorsClosed          int `bson:"cursor sweep cursors closed" stm:"sweep_cursors"`
		CursorSweepCursorsExamined        int `bson:"cursor sweep cursors examined" stm:"sweep_examined"`
		CursorSweeps                      int `bson:"cursor sweeps" stm:"sweeps"`
		CursorTruncateCalls               int `bson:"cursor truncate calls" stm:"truncate"`
		CursorUpdateCalls                 int `bson:"cursor update calls" stm:"update"`
		CursorUpdateValueSizeChange       int `bson:"cursor update value size change" stm:"update_value"`
	} `bson:"cursor" stm:"cursor"`
	Lock *struct {
		CheckpointLockAcquisitions                 int `bson:"checkpoint lock acquisitions" stm:"checkpoint_acquisitions"`
		DhandleReadLockAcquisitions                int `bson:"dhandle read lock acquisitions" stm:"read_acquisitions"`
		DhandleWriteLockAcquisitions               int `bson:"dhandle write lock acquisitions" stm:"write_acquisitions"`
		DurableTimestampQueueReadLockAcquisitions  int `bson:"durable timestamp queue read lock acquisitions" stm:"durable_timestamp_queue_read_acquisitions"`
		DurableTimestampQueueWriteLockAcquisitions int `bson:"durable timestamp queue write lock acquisitions" stm:"durable_timestamp_queue_write_acquisitions"`
		MetadataLockAcquisitions                   int `bson:"metadata lock acquisitions" stm:"metadata_acquisitions"`
		ReadTimestampQueueReadLockAcquisitions     int `bson:"read timestamp queue read lock acquisitions" stm:"read_timestamp_queue_read_acquisitions"`
		ReadTimestampQueueWriteLockAcquisitions    int `bson:"read timestamp queue write lock acquisitions" stm:"read_timestamp_queue_write_acquisitions"`
		SchemaLockAcquisitions                     int `bson:"schema lock acquisitions" stm:"schema_acquisitions"`
		TableReadLockAcquisitions                  int `bson:"table read lock acquisitions" stm:"table_read_acquisitions"`
		TableWriteLockAcquisitions                 int `bson:"table write lock acquisitions" stm:"table_write_acquisitions"`
		TxnGlobalReadLockAcquisitions              int `bson:"txn global read lock acquisitions" stm:"txn_global_read_acquisitions"`

		CheckpointLockApplicationThreadWaitTimeUsecs               int `bson:"checkpoint lock application thread wait time (usecs)" stm:"checkpoint_wait_time"`
		CheckpointLockInternalThreadWaitTimeUsecs                  int `bson:"checkpoint lock internal thread wait time (usecs)" stm:"checkpoint_internal_thread_wait_time"`
		DhandleLockApplicationThreadTimeWaitingUsecs               int `bson:"dhandle lock application thread time waiting (usecs)" stm:"application_thread_time_waiting"`
		DhandleLockInternalThreadTimeWaitingUsecs                  int `bson:"dhandle lock internal thread time waiting (usecs)" stm:"internal_thread_time_waiting"`
		DurableTimestampQueueLockApplicationThreadTimeWaitingUsecs int `bson:"durable timestamp queue lock application thread time waiting (usecs)" stm:"durable_timestamp_queue_application_thread_time_waiting"`
		DurableTimestampQueueLockInternalThreadTimeWaitingUsecs    int `bson:"durable timestamp queue lock internal thread time waiting (usecs)" stm:"durable_timestamp_queue_internal_thread_time_waiting"`
		MetadataLockApplicationThreadWaitTimeUsecs                 int `bson:"metadata lock application thread wait time (usecs)" stm:"metadata_application_thread_wait_time"`
		MetadataLockInternalThreadWaitTimeUsecs                    int `bson:"metadata lock internal thread wait time (usecs)" stm:"metadata_internal_thread_wait_time"`
		ReadTimestampQueueLockApplicationThreadTimeWaitingUsecs    int `bson:"read timestamp queue lock application thread time waiting (usecs)" stm:"read_timestamp_queue_application_thread_time_waiting"`
		ReadTimestampQueueLockInternalThreadTimeWaitingUsecs       int `bson:"read timestamp queue lock internal thread time waiting (usecs)" stm:"read_timestamp_queue_internal_thread_time_waiting"`
		SchemaLockApplicationThreadWaitTimeUsecs                   int `bson:"schema lock application thread wait time (usecs)" stm:"schema_application_thread_wait_time"`
		SchemaLockInternalThreadWaitTimeUsecs                      int `bson:"schema lock internal thread wait time (usecs)" stm:"schema_internal_thread_wait_time"`
	} `bson:"lock" stm:"lock"`
	Log *struct {
		LogFlushOperations             int `bson:"log flush operations" stm:"flush"`
		LogForceWriteOperations        int `bson:"log force write operations" stm:"force_write"`
		LogForceWriteOperationsSkipped int `bson:"log force write operations skipped" stm:"write_skip"`
		LogScanOperations              int `bson:"log scan operations" stm:"scan"`
		LogSyncOperations              int `bson:"log sync operations" stm:"sync"`
		LogSyncDirOperations           int `bson:"log sync_dir operations" stm:"sync_dir"`
		LogWriteOperations             int `bson:"log write operations" stm:"write"`

		LogBytesOfPayloadData    int `bson:"log bytes of payload data" stm:"payload"`
		LogBytesWritten          int `bson:"log bytes written" stm:"written"`
		LoggingBytesConsolidated int `bson:"logging bytes consolidated" stm:"consolidated"`
		TotalLogBufferSize       int `bson:"total log buffer size" stm:"buffer_size"`
	} `bson:"log" stm:"log"`
	Transaction *struct {
		PreparedTransactions   int `bson:"prepared transactions" stm:"prepare"`
		QueryTimestampCalls    int `bson:"query timestamp calls" stm:"query"`
		RollbackToStableCalls  int `bson:"rollback to stable calls" stm:"rollback"`
		SetTimestampCalls      int `bson:"set timestamp calls" stm:"set_timestamp"`
		TransactionBegins      int `bson:"transaction begins" stm:"begin"`
		TransactionSyncCalls   int `bson:"transaction sync calls" stm:"sync"`
		TransactionsCommitted  int `bson:"transactions committed" stm:"committed"`
		TransactionsRolledBack int `bson:"transactions rolled back" stm:"rolled_back"`
	} `bson:"transaction" stm:"transaction"`
}

type dbStats struct {
	Collections int64 `bson:"collections"`
	Views       int64 `bson:"views"`
	Indexes     int64 `bson:"indexes"`
	Objects     int64 `bson:"objects"`
	DataSize    int64 `bson:"dataSize"`
	IndexSize   int64 `bson:"indexSize"`
	StorageSize int64 `bson:"storageSize"`
}

type replSetStatus struct {
	Date    time.Time       `bson:"date"`
	Members []replSetMember `bson:"members"`
}

type replSetMember struct {
	Name              string     `bson:"name"`
	Self              *bool      `bson:"self"`
	State             int        `bson:"state"`
	Health            int        `bson:"health"`
	OptimeDate        time.Time  `bson:"optimeDate"`
	LastHeartbeat     *time.Time `bson:"lastHeartbeat"`
	LastHeartbeatRecv *time.Time `bson:"lastHeartbeatRecv"`
	PingMs            *int64     `bson:"pingMs"`
	Uptime            int64      `bson:"uptime"`
}

type aggrResults struct {
	Bool  bool  `bson:"_id"`
	Count int64 `bson:"count"`
}

type aggrResult struct {
	True  int64
	False int64
}

type partitionedResult struct {
	Partitioned   int64
	UnPartitioned int64
}

type shardNodesResult struct {
	ShardAware   int64
	ShardUnaware int64
}
