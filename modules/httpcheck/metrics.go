package httpcheck

type metrics struct {
	Request struct {
		Status struct {
			Success bool `stm:"success"`
			Failed  bool `stm:"failed"`
			Timeout bool `stm:"timeout"`
		} `stm:""`
	} `stm:""`
	Response struct {
		BadContent    bool  `stm:"bad_content"`
		BadStatusCode bool  `stm:"bad_status"`
		Time          int64 `stm:"response_time"`
		Length        int   `stm:"response_length"`
	} `stm:""`
}
