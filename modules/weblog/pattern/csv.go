package pattern

import (
	"regexp"
)

var (
	reAddress        = regexp.MustCompile(`[\da-f.:]+|localhost`)
	reCode           = regexp.MustCompile(`[1-9]\d{2}`)
	reBytesSent      = regexp.MustCompile(`\d+|-`)
	reResponseLength = regexp.MustCompile(`\d+|-`)
	reResponseTime   = regexp.MustCompile(`\d+|\d+\.\d+`)
)

const (
	address = "address" // check
	code    = "code"    // check
	request = "request" // makes no sense
	// userDefined          = "user_defined"           // can't
	bytesSent            = "bytes_sent"             // check
	responseTime         = "response_time"          // check
	responseTimeUpstream = "response_time_upstream" // check
	responseLength       = "response_length"        // check

	//keyURL         = "url"
	//keyHTTPMethod  = "http_method"
	//keyHTTPVersion = "http_version"
	//
	//keyResponseTimeHistogram         = "response_time_histogram"
	//keyResponseTimeUpstreamHistogram = "response_time_histogram_upstream"
)

type (
	CSVPattern []CSVField
	CSVField   struct {
		Name  string
		Index int
	}
)

//func (p CSVPattern) hasField(name string) bool {
//	for idx := range p {
//		if p[idx].Name == name {
//			return true
//		}
//	}
//	return false
//}

var CSVDefaultPatterns = []CSVPattern{
	// TODO: add examples
	{
		{address, 0},
		{request, 5},
		{code, 6},
		{bytesSent, 7},
		{responseLength, 8},
		{responseTime, 9},
		{responseTimeUpstream, 10},
	},
	// TODO: add examples
	{
		{address, 1},
		{request, 6},
		{code, 7},
		{bytesSent, 8},
		{responseLength, 9},
		{responseTime, 10},
		{responseTimeUpstream, 11},
	},
	// TODO: add examples
	{
		{address, 0},
		{request, 5},
		{code, 6},
		{bytesSent, 7},
	},
	// TODO: add examples
	{
		{address, 1},
		{request, 6},
		{code, 7},
		{bytesSent, 8},
	},
}
