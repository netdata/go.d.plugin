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
		IsBad struct {
			Content    bool `stm:"bad_content"`
			StatusCode bool `stm:"bad_status"`
		} `stm:""`
		Time   int `stm:"response_time"`
		Length int `stm:"response_length"`
	} `stm:""`
}
