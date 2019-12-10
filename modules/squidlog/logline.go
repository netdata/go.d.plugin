package squidlog

import (
	"errors"
)

// https://wiki.squid-cache.org/Features/LogFormat
// http://www.squid-cache.org/Doc/config/logformat/
// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Squid_result_codes

type (
	logLine struct {
		conn  connFields
		time  timeFields
		http  httpFields
		squid squidFields
	}
	connFields struct {
		clientIP string
	}
	timeFields struct {
		respTime int
	}
	httpFields struct {
		reqMethod    string
		respSentCode int
		mimeType     string
		respSize     int
	}
	squidFields struct {
		cacheCode string
		hierCode  string
	}
)

/*
#logformat squid %ts.%03tu %6tr %>a %Ss/%03>Hs %<st %rm %ru %un %Sh/%<A %mt
#logformat squidmime %ts.%03tu %6tr %>a %Ss/%03>Hs %<st %rm %ru %un %Sh/%<A %mt [%>h] [%<h]
#logformat common %>a %ui %un [%tl] "%rm %ru HTTP/%rv" %>Hs %<st %Ss:%Sh
#logformat combined %>a %ui %un [%tl] "%rm %ru HTTP/%rv" %>Hs %<st "%{Referer}>h" "%{User-Agent}>h" %Ss:%Sh
*/
func (l *logLine) Assign(field string, value string) (err error) {
	if value == "" {
		return
	}

	switch field {
	case "tr":
		// Response time (milliseconds)
		err = l.assignRespTime(value)
	case ">a":
		// Client source IP address
		err = l.assignClientIP(value)
	case "Ss":
		// Squid request status
		err = l.assignCacheCode(value)
	case ">Hs":
		// HTTP status code sent to the client
		err = l.assignHTTPSentCode(value)
	case "<st":
		// Total size of reply sent to client (after adaptation)
		err = l.assignRespSize(value)
	case "rm":
		// Request method
		err = l.assignReqMethod(value)
	case "Sh":
		// Squid hierarchy status
		err = l.assignHierCode(value)
	case "mt":
		// MIME content type
		err = l.assignMimeType(value)
	}
	return err
}

func (l *logLine) assignRespTime(duration string) error { return nil }
func (l *logLine) assignClientIP(address string) error  { return nil }
func (l *logLine) assignCacheCode(code string) error    { return nil }
func (l *logLine) assignHTTPSentCode(code string) error { return nil }
func (l *logLine) assignRespSize(bytes string) error    { return nil }
func (l *logLine) assignReqMethod(method string) error  { return nil }
func (l *logLine) assignHierCode(code string) error     { return nil }
func (l *logLine) assignMimeType(mime string) error     { return nil }

func (l logLine) verify() error {
	if l.empty() {
		return errors.New("")
	}
	if l.hasDuration() && !l.isDurationValid() {
		return errors.New("")
	}
	if l.hasClientIP() && !l.isClientIPValid() {
		return errors.New("")
	}
	if l.hasCacheCode() && !l.isCacheCodeValid() {
		return errors.New("")
	}
	if l.hasHTTPSentCode() && !l.isHTTPSentCodeValid() {
		return errors.New("")
	}
	if l.hasRespSize() && !l.isRespSizeValid() {
		return errors.New("")
	}
	if l.hasReqMethod() && !l.isReqMethodValid() {
		return errors.New("")
	}
	if l.hasHierCode() && !l.isHierCodeValid() {
		return errors.New("")
	}
	if l.hasMimeType() && !l.isMimeTypeValid() {
		return errors.New("")
	}
	return nil
}

func (l logLine) empty() bool           { return l == emptyLogLine }
func (l logLine) hasDuration() bool     { return !isEmptyNumber(l.time.respTime) }
func (l logLine) hasClientIP() bool     { return !isEmptyString(l.conn.clientIP) }
func (l logLine) hasCacheCode() bool    { return !isEmptyString(l.squid.cacheCode) }
func (l logLine) hasHTTPSentCode() bool { return !isEmptyNumber(l.http.respSentCode) }
func (l logLine) hasRespSize() bool     { return !isEmptyNumber(l.http.respSize) }
func (l logLine) hasReqMethod() bool    { return !isEmptyString(l.http.reqMethod) }
func (l logLine) hasHierCode() bool     { return !isEmptyString(l.squid.hierCode) }
func (l logLine) hasMimeType() bool     { return !isEmptyString(l.http.mimeType) }

func (l logLine) isDurationValid() bool     { return false }
func (l logLine) isClientIPValid() bool     { return false }
func (l logLine) isCacheCodeValid() bool    { return false }
func (l logLine) isHTTPSentCodeValid() bool { return false }
func (l logLine) isRespSizeValid() bool     { return false }
func (l logLine) isReqMethodValid() bool    { return false }
func (l logLine) isHierCodeValid() bool     { return false }
func (l logLine) isMimeTypeValid() bool     { return false }

func (l *logLine) reset() {

}

func newEmptyLogLine() *logLine {
	return &logLine{
		conn: connFields{
			clientIP: emptyString,
		},
		time: timeFields{
			respTime: emptyNumber,
		},
		http: httpFields{
			reqMethod:    emptyString,
			respSentCode: emptyNumber,
			mimeType:     emptyString,
			respSize:     emptyNumber,
		},
		squid: squidFields{
			cacheCode: emptyString,
			hierCode:  emptyString,
		},
	}
}

var emptyLogLine = *newEmptyLogLine()

const (
	emptyString = "__empty_string__"
	emptyNumber = -9999
)

func isEmptyString(s string) bool {
	return s == emptyString || s == ""
}

func isEmptyNumber(n int) bool {
	return n == emptyNumber
}
