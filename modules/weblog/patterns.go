package weblog

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
	keyAddress              = "address"                // check
	keyCode                 = "code"                   // check
	keyRequest              = "request"                // makes no sense
	keyUserDefined          = "user_defined"           // can't
	keyBytesSent            = "bytes_sent"             // check
	keyResponseTime         = "response_time"          // check
	keyResponseTimeUpstream = "response_time_upstream" // check
	keyResponseLength       = "response_length"        // check

	keyURL         = "url"
	keyHTTPMethod  = "http_method"
	keyHTTPVersion = "http_version"

	keyResponseTimeHistogram         = "response_time_histogram"
	keyResponseTimeUpstreamHistogram = "response_time_histogram_upstream"
)

type (
	pattern []field
	field   struct {
		name  string
		index int
	}
)

func (p pattern) hasField(name string) bool {
	for idx := range p {
		if p[idx].name == name {
			return true
		}
	}
	return false
}

var defaultPatterns = []pattern{
	// TODO: add examples
	{
		{keyAddress, 0},
		{keyRequest, 5},
		{keyCode, 6},
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
