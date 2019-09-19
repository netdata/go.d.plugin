package hdfs

// HDFS Architecture
// https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html#NameNode+and+DataNodes

// Metrics description
// https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-common/Metrics.html

type metrics struct {
	*jvmMetrics           `stm:"jvm"`
	*fsnNameSystemMetrics `stm:"fsn"`
}

type jvmMetrics struct {
	ProcessName string `json:"tag.ProcessName"`
	HostName    string `json:"tag.Hostname"`
	//MemNonHeapUsedM            float64 `stm:"mem_non_heap_used,1000,1"`
	//MemNonHeapCommittedM       float64 `stm:"mem_non_heap_committed,1000,1"`
	//MemNonHeapMaxM             float64 `stm:"mem_non_heap_max"`
	MemHeapUsedM      float64 `stm:"mem_heap_used,1000,1"`
	MemHeapCommittedM float64 `stm:"mem_heap_committed,1000,1"`
	MemHeapMaxM       float64 `stm:"mem_heap_max"`
	//MemMaxM                    float64 `stm:"mem_max"`
	GcCount                    float64 `stm:"gc_count"`
	GcTimeMillis               float64 `stm:"gc_time_millis"`
	GcNumWarnThresholdExceeded float64 `stm:"gc_num_warn_threshold_exceeded"`
	GcNumInfoThresholdExceeded float64 `stm:"gc_num_info_threshold_exceeded"`
	GcTotalExtraSleepTime      float64 `stm:"gc_total_extra_sleep_time"`
	ThreadsNew                 float64 `stm:"threads_new"`
	ThreadsRunnable            float64 `stm:"threads_runnable"`
	ThreadsBlocked             float64 `stm:"threads_blocked"`
	ThreadsWaiting             float64 `stm:"threads_waiting"`
	ThreadsTimedWaiting        float64 `stm:"threads_timed_waiting"`
	ThreadsTerminated          float64 `stm:"threads_terminated"`
	LogFatal                   float64 `stm:"log_fatal"`
	LogError                   float64 `stm:"log_error"`
	LogWarn                    float64 `stm:"log_warn"`
	LogInfo                    float64 `stm:"log_info"`
}

type fsnNameSystemMetrics struct {
	HostName string `json:"tag.Hostname"`
	HAState  string `json:"tag.HAState"`
	//TotalSyncTimes                               float64 `json:"tag.tag.TotalSyncTimes"`
	//MissingBlocks                                float64 `stm:"missing_blocks"`
	//MissingReplOneBlocks                         float64 `stm:"missing_repl_one_blocks"`
	//ExpiredHeartbeats                            float64 `stm:"expired_heartbeats"`
	//TransactionsSinceLastCheckpoint              float64 `stm:"transactions_since_last_checkpoint"`
	//TransactionsSinceLastLogRoll                 float64 `stm:"transactions_since_last_log_roll"`
	//LastWrittenTransactionId                     float64 `stm:"last_written_transaction_id"`
	//LastCheckpointTime                           float64 `stm:"last_checkpoint_time"`
	//CapacityTotal                                float64 `stm:"capacity_total"`
	//CapacityTotalGB                              float64 `stm:"capacity_total_gb"`
	CapacityUsed float64 `stm:"capacity_used"`
	//CapacityUsedGB                               float64 `stm:"capacity_used_gb"`
	CapacityRemaining float64 `stm:"capacity_remaining"`
	//ProvidedCapacityTotal                        float64 `stm:"provided_capacity_total"`
	//CapacityRemainingGB                          float64 `stm:"capacity_remaining_gb"`
	//CapacityUsedNonDFS                           float64 `stm:"capacity_used_non_dfs"`
	TotalLoad float64 `stm:"total_load"`
	//SnapshottableDirectories                     float64 `stm:"snapshottable_directories"`
	//Snapshots                                    float64 `stm:"snapshots"`
	//NumEncryptionZones                           float64 `stm:"num_encryption_zones"`
	//LockQueueLength                              float64 `stm:"lock_queue_length"`
	//BlocksTotal                                  float64 `stm:"blocks_total"`
	//NumFilesUnderConstruction                    float64 `stm:"num_files_under_construction"`
	//NumActiveClients                             float64 `stm:"num_active_clients"`
	//FilesTotal                                   float64 `stm:"files_total"`
	//PendingReplicationBlocks                     float64 `stm:"pending_replication_blocks"`
	//PendingReconstructionBlocks                  float64 `stm:"pending_reconstruction_blocks"`
	//UnderReplicatedBlocks                        float64 `stm:"under_replicated_blocks"`
	//LowRedundancyBlocks                          float64 `stm:"low_redundancy_blocks"`
	//CorruptBlocks                                float64 `stm:"corrupt_blocks"`
	//ScheduledReplicationBlocks                   float64 `stm:"scheduled_replication_blocks"`
	//PendingDeletionBlocks                        float64 `stm:"pending_deletion_blocks"`
	//LowRedundancyReplicatedBlocks                float64 `stm:"low_redundancy_replicated_blocks"`
	//CorruptReplicatedBlocks                      float64 `stm:"corrupt_replicated_blocks"`
	//MissingReplicatedBlocks                      float64 `stm:"missing_replicated_blocks"`
	//MissingReplicationOneBlocks                  float64 `stm:"missing_replication_one_blocks"`
	//HighestPriorityLowRedundancyReplicatedBlocks float64 `stm:"highest_priority_low_redundancy_replicated_blocks"`
	//HighestPriorityLowRedundancyECBlocks         float64 `stm:"highest_priority_low_redundancy_ec_blocks"`
	//BytesInFutureReplicatedBlocks                float64 `stm:"bytes_in_future_replicated_blocks"`
	//PendingDeletionReplicatedBlocks              float64 `stm:"pending_deletion_replicated_blocks"`
	//TotalReplicatedBlocks                        float64 `stm:"total_replicated_blocks"`
	//LowRedundancyECBlockGroups                   float64 `stm:"low_redundancy_ec_block_groups"`
	//CorruptECBlockGroups                         float64 `stm:"corrupt_ec_block_groups"`
	//MissingECBlockGroups                         float64 `stm:"missing_ec_block_groups"`
	//BytesInFutureECBlockGroups                   float64 `stm:"bytes_in_future_ec_block_groups"`
	//PendingDeletionECBlocks                      float64 `stm:"pending_deletion_ec_blocks"`
	//TotalECBlockGroups                           float64 `stm:"total_ec_block_groups"`
	//ExcessBlocks                                 float64 `stm:"excess_blocks"`
	//NumTimedOutPendingReconstructions            float64 `stm:"num_timed_out_pending_reconstructions"`
	//PostponedMisreplicatedBlocks                 float64 `stm:"postponed_misreplicated_blocks"`
	//PendingDataNodeMessageCount                  float64 `stm:"pending_data_node_message_count"`
	//MillisSinceLastLoadedEdits                   float64 `stm:"millis_since_last_loaded_edits"`
	//BlockCapacity                                float64 `stm:"block_capacity"`
	NumLiveDataNodes float64 `stm:"num_live_data_nodes"`
	NumDeadDataNodes float64 `stm:"num_dead_data_nodes"`
	//NumDecomLiveDataNodes                        float64 `stm:"num_decom_live_data_nodes"`
	//NumDecomDeadDataNodes                        float64 `stm:"num_decom_dead_data_nodes"`
	VolumeFailuresTotal float64 `stm:"volume_failures_total"`
	//EstimatedCapacityLostTotal                   float64 `stm:"estimated_capacity_lost_total"`
	//NumDecommissioningDataNodes                  float64 `stm:"num_decommissioning_data_nodes"`
	//StaleDataNodes                               float64 `stm:"stale_data_nodes"`
	//NumStaleStorages                             float64 `stm:"num_stale_storages"`
	//TotalSyncCount                               float64 `stm:"total_sync_count"`
	//NumInMaintenanceLiveDataNodes                float64 `stm:"num_in_maintenance_live_data_nodes"`
	//NumInMaintenanceDeadDataNodes                float64 `stm:"num_in_maintenance_dead_data_nodes"`
	//NumEnteringMaintenanceDataNodes              float64 `stm:"num_entering_maintenance_data_nodes"`
}
