package httpcheck

type metrics struct {
	Status         status `stm:""`
	InState        int    `stm:"in_state"`
	ResponseTime   int    `stm:"time"`
	ResponseLength int    `stm:"length"`
}

type status struct {
	Success bool `stm:"success"` // No error on request, body reading and checking its content
	Timeout bool `stm:"timeout"`
	//DNSLookupError    bool `stm:"dns_lookup_error"`
	//ParseAddressError bool `stm:"address_parse_error"`
	//RedirectError     bool `stm:"redirect_error"`
	//BodyReadError     bool `stm:"body_read_error"`
	BadContent    bool `stm:"bad_content"`
	BadStatusCode bool `stm:"bad_status"`
	NoConnection  bool `stm:"no_connection"` // All other errors basically
}
