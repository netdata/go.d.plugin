package weblog

import "fmt"

type (
	fieldID    int
	logPattern struct {
		Mapping  [fieldMax]int
		maxIndex int
	}
)

const (
	fieldRemoteAddr fieldID = iota
	fieldRequest
	fieldStatus
	fieldBytesSent
	fieldHost
	fieldRespTime
	fieldRespTimeUpstream
	fieldReqLength
	fieldUserDefined

	// last index, keep it at the end
	fieldMax
)

const (
	keyRemoteAddr       = "remote_addr"
	keyRequest          = "request"
	keyStatus           = "status"
	keyBytesSent        = "bytes_sent"
	keyHost             = "host"
	keyRespTime         = "respTime"
	keyRespTimeUpstream = "resp_time_upstream"
	keyReqLength        = "request_length"
	keyUserDefined      = "user_defined"
)

func NewLogPattern(mapping map[string]int) (logPattern, error) {
	var pattern logPattern

	for i := fieldID(0); i < fieldMax; i++ {
		pattern.Mapping[i] = -1
	}
	for key, idx := range mapping {
		f, err := ParseField(key)
		if err != nil {
			return pattern, err
		}
		pattern.Mapping[f] = idx
		if pattern.maxIndex < idx {
			pattern.maxIndex = idx
		}
	}
	return pattern, nil
}

func (p logPattern) MaxIndex() int {
	return p.maxIndex
}

func (f fieldID) String() string {
	switch f {
	case fieldRemoteAddr:
		return keyRemoteAddr
	case fieldRequest:
		return keyRequest
	case fieldStatus:
		return keyStatus
	case fieldBytesSent:
		return keyBytesSent
	case fieldHost:
		return keyHost
	case fieldRespTime:
		return keyRespTime
	case fieldRespTimeUpstream:
		return keyRespTimeUpstream
	case fieldReqLength:
		return keyReqLength
	case fieldUserDefined:
		return keyUserDefined
	default:
		panic("invalid field ID")
	}
}

func ParseField(str string) (fieldID, error) {
	switch str {
	case keyRemoteAddr:
		return fieldRemoteAddr, nil
	case keyRequest:
		return fieldRequest, nil
	case keyStatus:
		return fieldStatus, nil
	case keyBytesSent:
		return fieldBytesSent, nil
	case keyHost:
		return fieldHost, nil
	case keyRespTime:
		return fieldRespTime, nil
	case keyRespTimeUpstream:
		return fieldRespTimeUpstream, nil
	case keyReqLength:
		return fieldReqLength, nil
	case keyUserDefined:
		return fieldUserDefined, nil
	default:
		return -1, fmt.Errorf("unknown field: %s", str)
	}
}

var (
	logFmtNetdata, _ = NewLogPattern(map[string]int{
		keyRemoteAddr:       0,
		keyRequest:          5,
		keyStatus:           6,
		keyBytesSent:        7,
		keyReqLength:        8,
		keyRespTime:         9,
		keyRespTimeUpstream: 10,
	})

	logFmtNetdataVhost, _ = NewLogPattern(map[string]int{
		//keyVhost:            0,
		keyRemoteAddr:       1,
		keyRequest:          6,
		keyStatus:           7,
		keyBytesSent:        8,
		keyReqLength:        9,
		keyRespTime:         10,
		keyRespTimeUpstream: 11,
	})
	logFmtCommon, _ = NewLogPattern(map[string]int{
		keyRemoteAddr: 0,
		keyRequest:    5,
		keyStatus:     6,
		keyBytesSent:  7,
	})
	logFmtDefaultVhost, _ = NewLogPattern(map[string]int{
		//keyVhost:      0,
		keyRemoteAddr: 1,
		keyRequest:    6,
		keyStatus:     7,
		keyBytesSent:  8,
	})

	defaultLogFmtPatterns = []logPattern{
		logFmtNetdata,
		logFmtNetdataVhost,
		logFmtCommon,
		logFmtDefaultVhost,
	}
)
