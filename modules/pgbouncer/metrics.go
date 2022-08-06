package pgbouncer

type metrics struct {
	dbs map[string]*dbMetrics
}

type dbMetrics struct {
	name      string
	updated   bool
	hasCharts bool

	// command 'SHOW DATABASES;'
	maxConnections     int64
	currentConnections int64
	paused             int64
	disabled           int64

	// command 'SHOW STATS;'
	// v1.17.0: https://github.com/pgbouncer/pgbouncer/blob/9a346b0e451d842d7202abc3eccf0ff5a66b2dd6/src/stats.c#L76
	// v1.7.0: https://github.com/pgbouncer/pgbouncer/blob/b8eab2128d43e895107e7cddcee74f65181d3673/src/stats.c#L58
	totalXactCount  int64 // v1.8+
	totalQueryCount int64 // v1.8+
	totalReceived   int64
	totalSent       int64
	totalXactTime   int64 // v1.8+
	totalQueryTime  int64
	totalWaitTime   int64 // v1.8+
	avgXactTime     int64 // v1.8+
	avgQueryTime    int64

	// command 'SHOW POOLS;'
	clActive    int64
	clWaiting   int64
	clCancelReq int64
	svActive    int64
	svIdle      int64
	svUsed      int64
	svTested    int64
	svLogin     int64
	maxWait     int64
	maxWaitUS   int64
}
