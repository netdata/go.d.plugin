package weblog

import "sort"

const (
	keyAddress          = "address"            // check
	keyCode             = "code"               // check
	keyRequest          = "request"            // no
	keyBytesSent        = "bytes_sent"         // check
	keyRespTime         = "resp_time"          // check
	keyRespTimeUpstream = "resp_time_upstream" // check
	keyRespLength       = "resp_length"        // check
	keyUserDefined      = "user_defined"

	keyMethod                    = "method"
	keyVersion                   = "version"
	keyURL                       = "url"
	keyRespTimeHistogram         = "resp_time_histogram"
	keyRespTimeUpstreamHistogram = "resp_time_upstream_histogram"
)

type (
	csvPattern []csvField
	csvField   struct {
		Name  string
		Index int
	}
)

func (c csvPattern) max() int {
	return c[len(c)-1].Index
}

func (c csvPattern) isSorted() bool {
	return sort.SliceIsSorted(c, func(i, j int) bool {
		return c[i].Index < c[j].Index
	})
}

var csvDefaultPatterns = []csvPattern{
	// TODO: add examples
	{
		{keyAddress, 0},
		{keyRequest, 5},
		{keyCode, 6},
		{keyBytesSent, 7},
		{keyRespLength, 8},
		{keyRespTime, 9},
		{keyRespTimeUpstream, 10},
	},
	// TODO: add examples
	{
		{keyAddress, 1},
		{keyRequest, 6},
		{keyCode, 7},
		{keyBytesSent, 8},
		{keyRespLength, 9},
		{keyRespTime, 10},
		{keyRespTimeUpstream, 11},
	},
	// TODO: add examples
	{
		{keyAddress, 0},
		{keyRequest, 5},
		{keyCode, 6},
		{keyBytesSent, 7},
	},
	// TODO: add examples
	{
		{keyAddress, 1},
		{keyRequest, 6},
		{keyCode, 7},
		{keyBytesSent, 8},
	},
}
