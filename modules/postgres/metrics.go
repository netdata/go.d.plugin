package postgres

type pgMetrics struct {
	srvMetrics
	dbs       map[string]*dbMetrics
	tables    map[string]*tableMetrics
	replApps  map[string]*replStandbyAppMetrics
	replSlots map[string]*replSlotMetrics
}

type srvMetrics struct {
	maxConnections int64

	uptime int64

	relkindOrdinaryTable        int64
	relkindIndex                int64
	relkindSequence             int64
	relkindTOASTTable           int64
	relkindView                 int64
	relkindMatView              int64
	relkindCompositeType        int64
	relkindForeignTable         int64
	relkindPartitionedTable     int64
	relkindPartitionedIndex     int64
	relkindOrdinaryTableSize    int64
	relkindIndexSize            int64
	relkindSequenceSize         int64
	relkindTOASTTableSize       int64
	relkindViewSize             int64
	relkindMatViewSize          int64
	relkindCompositeTypeSize    int64
	relkindForeignTableSize     int64
	relkindPartitionedTableSize int64
	relkindPartitionedIndexSize int64

	connUsed                      int64
	connStateActive               int64
	connStateIdle                 int64
	connStateIdleInTrans          int64
	connStateIdleInTransAborted   int64
	connStateFastpathFunctionCall int64
	connStateDisabled             int64

	checkpointsTimed    int64
	checkpointsReq      int64
	checkpointWriteTime int64
	checkpointSyncTime  int64
	buffersCheckpoint   int64
	buffersClean        int64
	maxwrittenClean     int64
	buffersBackend      int64
	buffersBackendFsync int64
	buffersAlloc        int64

	oldestXID                         int64
	percentTowardsWraparound          int64
	percentTowardsEmergencyAutovacuum int64

	walWrites            int64
	walRecycledFiles     int64
	walWrittenFiles      int64
	walArchiveFilesReady int64
	walArchiveFilesDone  int64

	autovacuumWorkersAnalyze       int64
	autovacuumWorkersVacuumAnalyze int64
	autovacuumWorkersVacuum        int64
	autovacuumWorkersVacuumFreeze  int64
	autovacuumWorkersBrinSummarize int64
}

type dbMetrics struct {
	name string

	updated   bool
	hasCharts bool

	numBackends  int64
	datConnLimit int64
	xactCommit   int64
	xactRollback int64
	blksRead     int64
	blksHit      int64
	tupReturned  int64
	tupFetched   int64
	tupInserted  int64
	tupUpdated   int64
	tupDeleted   int64
	conflicts    int64
	size         int64
	tempFiles    int64
	tempBytes    int64
	deadlocks    int64

	conflTablespace int64
	conflLock       int64
	conflSnapshot   int64
	conflBufferpin  int64
	conflDeadlock   int64

	accessShareLockHeld             int64
	rowShareLockHeld                int64
	rowExclusiveLockHeld            int64
	shareUpdateExclusiveLockHeld    int64
	shareLockHeld                   int64
	shareRowExclusiveLockHeld       int64
	exclusiveLockHeld               int64
	accessExclusiveLockHeld         int64
	accessShareLockAwaited          int64
	rowShareLockAwaited             int64
	rowExclusiveLockAwaited         int64
	shareUpdateExclusiveLockAwaited int64
	shareLockAwaited                int64
	shareRowExclusiveLockAwaited    int64
	exclusiveLockAwaited            int64
	accessExclusiveLockAwaited      int64
}

type replStandbyAppMetrics struct {
	name string

	updated   bool
	hasCharts bool

	walSentDelta   int64
	walWriteDelta  int64
	walFlushDelta  int64
	walReplayDelta int64

	walWriteLag  int64
	walFlushLag  int64
	walReplayLag int64
}

type replSlotMetrics struct {
	name string

	updated   bool
	hasCharts bool

	walKeep int64
	files   int64
}

type tableMetrics struct {
	name   string
	schema string
	db     string

	updated                 bool
	hasCharts               bool
	hasLastAutoVacuumChart  bool
	hasLastVacuumChart      bool
	hasLastAutoAnalyzeChart bool
	hasLastAnalyzeChart     bool

	seqScan            int64
	seqTupRead         int64
	idxScan            int64
	idxTupFetch        int64
	nTupIns            int64
	nTupUpd            int64
	nTupDel            int64
	nTupHotUpd         int64
	nLiveTup           int64
	nDeadTup           int64
	lastVacuumAgo      int64
	lastAutoVacuumAgo  int64
	lastAnalyzeAgo     int64
	lastAutoAnalyzeAgo int64
	vacuumCount        int64
	autovacuumCount    int64
	analyzeCount       int64
	autoAnalyzeCount   int64

	totalSize int64
}
