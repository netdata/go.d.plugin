package apache

type (
	serverStatus struct {
		Total struct {
			// A total number of accesses
			Accesses *int64 `stm:"accesses"`
			// A total number of byte count served.
			KBytes *int64 `stm:"kBytes"`
		} `stm:"total"`
		Averages struct {
			//The average number of requests per second
			ReqPerSec *float64 `stm:"req_per_sec,1,10000"`
			// The average number of bytes served per second
			BytesPerSec *float64 `stm:"bytes_per_sec,1,10000"`
			// The average number of bytes per request
			BytesPerReq *float64 `stm:"bytes_per_req,1,10000"`
		} `stm:""`
		Workers struct {
			// The number of worker serving requests.
			Busy *int64 `stm:"busy_workers"`
			// The number of idle worker.
			Idle *int64 `stm:"idle_workers"`
		} `stm:""`
		Connections struct {
			Total          *int64 `stm:"total"`
			AsyncWriting   *int64 `stm:"async_writing"`
			AsyncKeepAlive *int64 `stm:"async_keep_alive"`
			AsyncClosing   *int64 `stm:"async_closing"`
		} `stm:"conns"`
		Scoreboard *scoreboard `stm:"scoreboard"`
		Uptime     *int64      `stm:"uptime"`
	}
	scoreboard struct {
		Waiting     int64 `stm:"waiting"`
		Starting    int64 `stm:"starting"`
		Reading     int64 `stm:"reading"`
		Sending     int64 `stm:"sending"`
		KeepAlive   int64 `stm:"keepalive"`
		DNSLookup   int64 `stm:"dns_lookup"`
		Closing     int64 `stm:"closing"`
		Logging     int64 `stm:"logging"`
		Finishing   int64 `stm:"finishing"`
		IdleCleanup int64 `stm:"idle_cleanup"`
		Open        int64 `stm:"open"`
	}
)
