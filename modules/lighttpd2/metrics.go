package lighttpd2

type serverStatus struct {
	Requests struct {
		Total *int64 `stm:"abs"`
	} `stm:"requests"`
	Responses struct {
		Status struct {
			Codes1xx *int64 `stm:"1xx"`
			Codes2xx *int64 `stm:"2xx"`
			Codes3xx *int64 `stm:"3xx"`
			Codes4xx *int64 `stm:"4xx"`
			Codes5xx *int64 `stm:"5xx"`
		} `stm:"status"`
	} `stm:""`
	Traffic struct {
		In  *int64 `stm:"in_abs"`
		Out *int64 `stm:"out_abs"`
	} `stm:"traffic"`
	Connection struct {
		Total *int64 `stm:"abs"`
		State struct {
			Start         *int64 `stm:"start"`
			ReadHeader    *int64 `stm:"read_header"`
			HandleRequest *int64 `stm:"handle_request"`
			WriteResponse *int64 `stm:"write_response"`
			KeepAlive     *int64 `stm:"keepalive"`
			Upgraded      *int64 `stm:"upgraded"`
		} `stm:"state"`
	} `stm:"connection"`
	Memory struct {
		Usage *int64 `stm:"usage"`
	} `stm:"memory"`
	Uptime *int64 `stm:"uptime"`
}
