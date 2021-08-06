package mongo

type serverStatus struct {
	//default charts
	Opcounters  Opcounters   `bson:"opcounters,omitempty"`
	OpLatencies *OpLatencies `bson:"opLatencies,omitempty"`
	Connections Connections  `bson:"connections,omitempty"`
	Network     Network      `bson:"network,omitempty"`
	ExtraInfo   ExtraInfo    `bson:"extra_info,omitempty"`
	Asserts     Asserts      `bson:"asserts,omitempty"`

	// optional charts
	Transactions *Transactions         `bson:"transactions,omitempty"`
	GlobalLock   *GlobalLock           `bson:"globalLock,omitempty"`
	Tcmalloc     *ServerStatusTcmalloc `bson:"tcmalloc,omitempty"`
	Locks        *Locks                `bson:"locks,omitempty"`
	FlowControl  *FlowControl          `bson:"flowControl,omitempty"`
	WiredTiger   *WiredTiger           `bson:"wiredTiger,omitempty"`
}

type Opcounters struct {
	Insert  *int64 `bson:"insert,omitempty" stm:"operations_insert"`
	Query   *int64 `bson:"query,omitempty" stm:"operations_query"`
	Update  *int64 `bson:"update,omitempty" stm:"operations_update"`
	Delete  *int64 `bson:"delete,omitempty" stm:"operations_delete"`
	Getmore *int64 `bson:"getmore,omitempty" stm:"operations_getmore"`
	Command *int64 `bson:"command,omitempty" stm:"operations_command"`
}

type OpLatencies struct {
	Reads *struct {
		Latency *int64 `bson:"latency,omitempty" stm:"operations_latency_read"`
	} `bson:"reads,omitempty"`
	Writes *struct {
		Latency *int64 `bson:"latency,omitempty" stm:"operations_latency_write"`
	} `bson:"writes,omitempty"`
	Commands *struct {
		Latency *int64 `bson:"latency,omitempty" stm:"operations_latency_command"`
	} `bson:"commands,omitempty"`
}

type Connections struct {
	Current                 *int64 `bson:"current,omitempty" stm:"connections_current"`
	Available               *int64 `bson:"available,omitempty" stm:"connections_available"`
	TotalCreated            *int64 `bson:"totalCreated,omitempty" stm:"connections_total_created"`
	Active                  *int64 `bson:"active,omitempty" stm:"connections_active"`
	Threaded                *int64 `bson:"threaded,omitempty" stm:"connections_threaded"`
	ExhaustIsMaster         *int64 `bson:"exhaustIsMaster,omitempty" stm:"connections_exhaustIsMaster"`
	ExhaustHello            *int64 `bson:"exhaustHello,omitempty" stm:"connections_exhaustHello"`
	AwaitingTopologyChanges *int64 `bson:"awaitingTopologyChanges,omitempty" stm:"connections_awaitingTopologyChanges"`
}

type Network struct {
	BytesIn     *int64 `bson:"bytesIn,omitempty" stm:"network_bytes_in"`
	BytesOut    *int64 `bson:"bytesOut,omitempty" stm:"network_bytes_out"`
	NumRequests *int64 `bson:"numRequests,omitempty" stm:"network_requests"`
}

type ExtraInfo struct {
	PageFaults *int64 `bson:"page_faults,omitempty" stm:"page_faults"`
}

type Asserts struct {
	Regular   *int64 `bson:"regular,omitempty" stm:"asserts_regular"`
	Warning   *int64 `bson:"warning,omitempty" stm:"asserts_warning"`
	Msg       *int64 `bson:"msg,omitempty" stm:"asserts_msg"`
	User      *int64 `bson:"user,omitempty" stm:"asserts_user"`
	Tripwire  *int64 `bson:"tripwire,omitempty" stm:"asserts_tripwire"`
	Rollovers *int64 `bson:"rollovers,omitempty" stm:"asserts_rollovers"`
}

type Transactions struct {
	CurrentActive   *int64 `bson:"currentActive,omitempty" stm:"transactions_active"`
	CurrentInactive *int64 `bson:"currentInactive,omitempty" stm:"transactions_inactive"`
	CurrentOpen     *int64 `bson:"currentOpen,omitempty" stm:"transactions_open"`
	CurrentPrepared *int64 `bson:"currentPrepared,omitempty" stm:"transactions_prepared"`
}

type GlobalLock struct {
	ActiveClients *struct {
		Readers *int64 `bson:"readers,omitempty" stm:"active_clients_readers"`
		Writers *int64 `bson:"writers,omitempty" stm:"active_clients_writers"`
	} `bson:"activeClients,omitempty"`
	CurrentQueue *struct {
		Readers *int64 `bson:"readers,omitempty" stm:"current_queue_readers"`
		Writers *int64 `bson:"writers,omitempty" stm:"current_queue_writers"`
	} `bson:"currentQueue,omitempty"`
}

type ServerStatusTcmalloc struct {
	Generic  *Generic          `bson:"generic,omitempty"`
	Tcmalloc *TcmallocTcmalloc `bson:"tcmalloc,omitempty"`
}

type Generic struct {
	CurrentAllocatedBytes *int64 `bson:"current_allocated_bytes,omitempty" stm:"tcmalloc_current_allocated"`
	HeapSize              *int64 `bson:"heap_size,omitempty" stm:"tcmalloc_heap_size"`
}

type TcmallocTcmalloc struct {
	PageheapFreeBytes          *int64 `bson:"pageheap_free_bytes,omitempty" stm:"tcmalloc_pageheap_free"`
	PageheapUnmappedBytes      *int64 `bson:"pageheap_unmapped_bytes,omitempty" stm:"tcmalloc_pageheap_unmapped"`
	MaxTotalThreadCacheBytes   *int64 `bson:"max_total_thread_cache_bytes,omitempty" stm:"tcmalloc_max_total_thread_cache"`
	TotalFreeBytes             *int64 `bson:"total_free_bytes,omitempty" stm:"tcmalloc_total_free"`
	PageheapCommittedBytes     *int64 `bson:"pageheap_committed_bytes,omitempty" stm:"tcmalloc_pageheap_committed"`
	PageheapTotalCommitBytes   *int64 `bson:"pageheap_total_commit_bytes,omitempty" stm:"tcmalloc_pageheap_total_commit"`
	PageheapTotalDecommitBytes *int64 `bson:"pageheap_total_decommit_bytes,omitempty" stm:"tcmalloc_pageheap_total_decommit"`
	PageheapTotalReserveBytes  *int64 `bson:"pageheap_total_reserve_bytes,omitempty" stm:"tcmalloc_pageheap_total_reserve"`
}

type Locks struct {
	Global *struct {
		R *int64 `bson:"r,omitempty" stm:"locks_global_read"`
		W *int64 `bson:"W,omitempty" stm:"locks_global_write"`
	} `bson:"Global,omitempty"`
	Database *struct {
		R *int64 `bson:"r,omitempty" stm:"locks_database_read"`
		W *int64 `bson:"W,omitempty" stm:"locks_database_write"`
	} `bson:"Database,omitempty"`
	Collection *struct {
		R *int64 `bson:"r,omitempty" stm:"locks_collection_read"`
		W *int64 `bson:"W,omitempty" stm:"locks_collection_write"`
	} `bson:"Collection,omitempty"`
}

type FlowControl struct {
	TargetRateLimit     *int64 `bson:"targetRateLimit,omitempty" stm:"target_rate_limit"`
	TimeAcquiringMicros *int64 `bson:"timeAcquiringMicros,omitempty" stm:"time_acquiring_micros"`
}

type WiredTiger struct {
	BlockManager *struct {
		BytesRead                    int `bson:"bytes read" stm:"wiredtiger_block_manager_read"`
		BytesReadViaMemoryMapAPI     int `bson:"bytes read via memory map API" stm:"wiredtiger_block_manager_read_via_memory"`
		BytesReadViaSystemCallAPI    int `bson:"bytes read via system call API" stm:"wiredtiger_block_manager_read_via_system_api"`
		BytesWritten                 int `bson:"bytes written" stm:"wiredtiger_block_manager_written"`
		BytesWrittenForCheckpoint    int `bson:"bytes written for checkpoint" stm:"wiredtiger_block_manager_written_for_checkpoint"`
		BytesWrittenViaMemoryMapAPI  int `bson:"bytes written via memory map API" stm:"wiredtiger_block_manager_written_via_memory"`
		BytesWrittenViaSystemCallAPI int `bson:"bytes written via system call API" stm:"wiredtiger_block_manager_written_via_system_api"`
	} `bson:"block-manager" json:"block-manager"`
	Cache *struct {
		BytesAllocatedForUpdates int `bson:"bytes allocated for updates" stm:"wiredtiger_cache_alloccated"`
		BytesReadIntoCache       int `bson:"bytes read into cache" stm:"wiredtiger_cache_read"`
		BytesWrittenFromCache    int `bson:"bytes written from cache" stm:"wiredtiger_cache_write"`
	} `bson:"cache"`
	Capacity *struct {
		TimeWaitingDueToTotalCapacityUsecs int `bson:"time waiting due to total capacity (usecs)" stm:"wiredtiger_cache_wait_capacity"`
		TimeWaitingDuringCheckpointUsecs   int `bson:"time waiting during checkpoint (usecs)" stm:"wiredtiger_cache_wait_checkpoint"`
		TimeWaitingDuringEvictionUsecs     int `bson:"time waiting during eviction (usecs)" stm:"wiredtiger_cache_wait_eviction"`
		TimeWaitingDuringLoggingUsecs      int `bson:"time waiting during logging (usecs)" stm:"wiredtiger_cache_wait_logging"`
		TimeWaitingDuringReadUsecs         int `bson:"time waiting during read (usecs)" stm:"wiredtiger_cache_wait_read"`
	} `bson:"capacity"`
	Connection *struct {
		MemoryAllocations   int `bson:"memory allocations" stm:"wiredtiger_connection_allocations"`
		MemoryFrees         int `bson:"memory frees" stm:"wiredtiger_connection_frees"`
		MemoryReAllocations int `bson:"memory re-allocations" stm:"wiredtiger_connection_reallocations"`
	} `bson:"connection"`
	Cursor *struct {
		CachedCursorCount                 int `bson:"cached cursor count" stm:"wiredtiger_cursor_count"`
		CursorBulkLoadedCursorInsertCalls int `bson:"cursor bulk loaded cursor insert calls" stm:"wiredtiger_cursor_bulk"`
		CursorCloseCallsThatResultInCache int `bson:"cursor close calls that result in cache" stm:"wiredtiger_cursor_close"`
		CursorCreateCalls                 int `bson:"cursor create calls" stm:"wiredtiger_cursor_create"`
		CursorInsertCalls                 int `bson:"cursor insert calls" stm:"wiredtiger_cursor_insert"`
		CursorModifyCalls                 int `bson:"cursor modify calls" stm:"wiredtiger_cursor_modify"`
		CursorNextCalls                   int `bson:"cursor next calls" stm:"wiredtiger_cursor_next"`
		CursorOperationRestarted          int `bson:"cursor operation restarted" stm:"wiredtiger_cursor_restarted"`
		CursorPrevCalls                   int `bson:"cursor prev calls" stm:"wiredtiger_cursor_prev"`
		CursorRemoveCalls                 int `bson:"cursor remove calls" stm:"wiredtiger_cursor_remove"`
		CursorReserveCalls                int `bson:"cursor reserve calls" stm:"wiredtiger_cursor_reserve"`
		CursorResetCalls                  int `bson:"cursor reset calls" stm:"wiredtiger_cursor_reset"`
		CursorSearchCalls                 int `bson:"cursor search calls" stm:"wiredtiger_cursor_search"`
		CursorSearchHistoryStoreCalls     int `bson:"cursor search history store calls" stm:"wiredtiger_cursor_search_history"`
		CursorSearchNearCalls             int `bson:"cursor search near calls" stm:"wiredtiger_cursor_search_near"`
		CursorSweepBuckets                int `bson:"cursor sweep buckets" stm:"wiredtiger_cursor_sweep_buckets"`
		CursorSweepCursorsClosed          int `bson:"cursor sweep cursors closed" stm:"wiredtiger_cursor_sweep_cursors"`
		CursorSweepCursorsExamined        int `bson:"cursor sweep cursors examined" stm:"wiredtiger_cursor_sweep_examined"`
		CursorSweeps                      int `bson:"cursor sweeps" stm:"wiredtiger_cursor_sweeps"`
		CursorTruncateCalls               int `bson:"cursor truncate calls" stm:"wiredtiger_cursor_truncate"`
		CursorUpdateCalls                 int `bson:"cursor update calls" stm:"wiredtiger_cursor_update"`
		CursorUpdateValueSizeChange       int `bson:"cursor update value size change" stm:"wiredtiger_cursor_update_value"`
	} `bson:"cursor"`
	Lock *struct {
		CheckpointLockAcquisitions                 int `bson:"checkpoint lock acquisitions" stm:"wiredtiger_lock_checkpoint_acquisitions"`
		DhandleReadLockAcquisitions                int `bson:"dhandle read lock acquisitions" stm:"wiredtiger_lock_read_acquisitions"`
		DhandleWriteLockAcquisitions               int `bson:"dhandle write lock acquisitions" stm:"wiredtiger_lock_write_acquisitions"`
		DurableTimestampQueueReadLockAcquisitions  int `bson:"durable timestamp queue read lock acquisitions" stm:"wiredtiger_lock_durable_timestamp_queue_read_acquisitions"`
		DurableTimestampQueueWriteLockAcquisitions int `bson:"durable timestamp queue write lock acquisitions" stm:"wiredtiger_lock_durable timestamp_queue_write_acquisitions"`
		MetadataLockAcquisitions                   int `bson:"metadata lock acquisitions" stm:"wiredtiger_lock_metadata_acquisitions"`
		ReadTimestampQueueReadLockAcquisitions     int `bson:"read timestamp queue read lock acquisitions" stm:"wiredtiger_lock_read_timestamp_queue_read_acquisitions"`
		ReadTimestampQueueWriteLockAcquisitions    int `bson:"read timestamp queue write lock acquisitions" stm:"wiredtiger_lock_read timestamp_queue_write_acquisitions"`
		SchemaLockAcquisitions                     int `bson:"schema lock acquisitions" stm:"wiredtiger_lock_schema_acquisitions"`
		TableReadLockAcquisitions                  int `bson:"table read lock acquisitions" stm:"wiredtiger_lock_table_read_acquisitions"`
		TableWriteLockAcquisitions                 int `bson:"table write lock acquisitions" stm:"wiredtiger_lock_table_write_acquisitions"`
		TxnGlobalReadLockAcquisitions              int `bson:"txn global read lock acquisitions" stm:"wiredtiger_lock_txn_global_read_acquisitions"`

		CheckpointLockApplicationThreadWaitTimeUsecs               int `bson:"checkpoint lock application thread wait time (usecs)" stm:"wiredtiger_lock_checkpoint_wait_time"`
		CheckpointLockInternalThreadWaitTimeUsecs                  int `bson:"checkpoint lock internal thread wait time (usecs)" stm:"wiredtiger_lock_checkpoint_internal_thread_wait_time"`
		DhandleLockApplicationThreadTimeWaitingUsecs               int `bson:"dhandle lock application thread time waiting (usecs)" stm:"wiredtiger_lock_application_thread_time_waiting"`
		DhandleLockInternalThreadTimeWaitingUsecs                  int `bson:"dhandle lock internal thread time waiting (usecs)" stm:"wiredtiger_lock_internal_thread_time_waiting"`
		DurableTimestampQueueLockApplicationThreadTimeWaitingUsecs int `bson:"durable timestamp queue lock application thread time waiting (usecs)" stm:"wiredtiger_lock_durable_timestamp_queue_application_thread_time_waiting"`
		DurableTimestampQueueLockInternalThreadTimeWaitingUsecs    int `bson:"durable timestamp queue lock internal thread time waiting (usecs)" stm:"wiredtiger_lock_durable_timestamp_queue_internal_thread_time_waiting"`
		MetadataLockApplicationThreadWaitTimeUsecs                 int `bson:"metadata lock application thread wait time (usecs)" stm:"wiredtiger_lock_metadata_application_thread_wait_time"`
		MetadataLockInternalThreadWaitTimeUsecs                    int `bson:"metadata lock internal thread wait time (usecs)" stm:"wiredtiger_lock_metadata_internal_thread_wait_time"`
		ReadTimestampQueueLockApplicationThreadTimeWaitingUsecs    int `bson:"read timestamp queue lock application thread time waiting (usecs)" stm:"wiredtiger_lock_read_timestamp_queue_application_thread_time_waiting"`
		ReadTimestampQueueLockInternalThreadTimeWaitingUsecs       int `bson:"read timestamp queue lock internal thread time waiting (usecs)" stm:"wiredtiger_lock_read_timestamp_queue_internal_thread_time_waiting"`
		SchemaLockApplicationThreadWaitTimeUsecs                   int `bson:"schema lock application thread wait time (usecs)" stm:"wiredtiger_lock_schema_application_thread_wait_time"`
		SchemaLockInternalThreadWaitTimeUsecs                      int `bson:"schema lock internal thread wait time (usecs)" stm:"wiredtiger_lock_schema_internal_thread_wait_time"`
	} `bson:"lock"`
	Log *struct {
		LogFlushOperations             int `bson:"log flush operations" stm:"wiredtiger_log_fluh"`
		LogForceWriteOperations        int `bson:"log force write operations" stm:"wiredtiger_log_force_write"`
		LogForceWriteOperationsSkipped int `bson:"log force write operations skipped" stm:"wiredtiger_log_write_skip"`
		LogScanOperations              int `bson:"log scan operations" stm:"wiredtiger_log_scan"`
		LogSyncOperations              int `bson:"log sync operations" stm:"wiredtiger_log_sync"`
		LogSyncDirOperations           int `bson:"log sync_dir operations" stm:"wiredtiger_log_sync_dir"`
		LogWriteOperations             int `bson:"log write operations" stm:"wiredtiger_log_write"`

		LogBytesOfPayloadData    int `bson:"log bytes of payload data" stm:"wiredtiger_log_payload"`
		LogBytesWritten          int `bson:"log bytes written" stm:"wiredtiger_log_written"`
		LoggingBytesConsolidated int `bson:"logging bytes consolidated" stm:"wiredtiger_log_consolidated"`
		TotalLogBufferSize       int `bson:"total log buffer size" stm:"wiredtiger_log_buffer_size"`
	} `bson:"log"`
	Transaction *struct {
		PreparedTransactions   int `bson:"prepared transactions" stm:"wiredtiger_transaction_prepare"`
		QueryTimestampCalls    int `bson:"query timestamp calls" stm:"wiredtiger_transaction_query"`
		RollbackToStableCalls  int `bson:"rollback to stable calls" stm:"wiredtiger_transaction_rollback"`
		SetTimestampCalls      int `bson:"set timestamp calls" stm:"wiredtiger_transaction_set_timestamp"`
		TransactionBegins      int `bson:"transaction begins" stm:"wiredtiger_transaction_begin"`
		TransactionSyncCalls   int `bson:"transaction sync calls" stm:"wiredtiger_transaction_sync"`
		TransactionsCommitted  int `bson:"transactions committed" stm:"wiredtiger_transaction_committed"`
		TransactionsRolledBack int `bson:"transactions rolled back" stm:"wiredtiger_transaction_rolled_back"`
	} `bson:"transaction"`
}
