package weblog

//
//const (
//	keyVhost   = "vhost"   // check
//	keyAddress = "address" // check
//	keyCode    = "code"    // check
//	//keyRequest          = "request"            // check
//	//keyBytesSent        = "bytes_sent"         // check
//	//keyRespTime         = "resp_time"          // check
//	//keyRespTimeUpstream = "resp_time_upstream" // check
//	//keyReqLength       = "resp_length"        // check
//	//keyUserDefined      = "user_defined"
//	keyMethod  = "http_method"  // check, parsed request fieldID
//	keyVersion = "http_version" // check, parsed request fieldID
//	keyURL     = "url"          // parsed request fieldID
//
//	keyRespTimeHistogram         = "resp_time_histogram"
//	keyRespTimeUpstreamHistogram = "resp_time_upstream_histogram"
//)
//
//type (
//	csvPattern []csvField
//	csvField   struct {
//		Key   string
//		Index int
//	}
//)
//
//func (c csvPattern) max() int {
//	return c[len(c)-1].Index
//}
//
//func (c csvPattern) isSorted() bool {
//	return sort.SliceIsSorted(c, func(i, j int) bool {
//		return c[i].Index < c[j].Index
//	})
//}
//
//func (c csvPattern) isValid() bool {
//	set := make(map[int]bool)
//
//	for _, p := range c {
//		if !(p.Key != "" && !set[p.Index]) {
//			return false
//		}
//		set[p.Index] = true
//	}
//	return true
//}
//
//var (
//	logFormatNetdata = csvPattern{
//		{keyAddress, 0},
//		{keyRequest, 5},
//		{keyCode, 6},
//		{keyBytesSent, 7},
//		{keyReqLength, 8},
//		{keyRespTime, 9},
//		{keyRespTimeUpstream, 10},
//	}
//	logFormatNetdataVhost = csvPattern{
//		{keyVhost, 0},
//		{keyAddress, 1},
//		{keyRequest, 6},
//		{keyCode, 7},
//		{keyBytesSent, 8},
//		{keyReqLength, 9},
//		{keyRespTime, 10},
//		{keyRespTimeUpstream, 11},
//	}
//	logFormatDefault = csvPattern{
//		{keyAddress, 0},
//		{keyRequest, 5},
//		{keyCode, 6},
//		{keyBytesSent, 7},
//	}
//	logFormatDefaultVhost = csvPattern{
//		{keyVhost, 0},
//		{keyAddress, 1},
//		{keyRequest, 6},
//		{keyCode, 7},
//		{keyBytesSent, 8},
//	}
//
//	csvDefaultPatterns = []csvPattern{
//		logFormatNetdata,
//		logFormatNetdataVhost,
//		logFormatDefault,
//		logFormatDefaultVhost,
//	}
//)
