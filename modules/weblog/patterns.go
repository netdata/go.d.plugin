package weblog

import "sort"

const (
	keyAddress              = "address"                // check
	keyCode                 = "code"                   // check
	keyRequest              = "request"                // no
	keyBytesSent            = "bytes_sent"             // check
	keyResponseTime         = "response_time"          // check
	keyResponseTimeUpstream = "response_time_upstream" // check
	keyResponseLength       = "response_length"        // check
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
		{keyCode, 5},
		{keyRequest, 6},
		{keyBytesSent, 7},
		{keyResponseLength, 8},
		{keyResponseTime, 9},
		{keyResponseTimeUpstream, 10},
	},
	// TODO: add examples
	{
		{keyAddress, 1},
		{keyRequest, 6},
		{keyCode, 7},
		{keyBytesSent, 8},
		{keyResponseLength, 9},
		{keyResponseTime, 10},
		{keyResponseTimeUpstream, 11},
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
