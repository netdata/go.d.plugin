package squidlog

import (
	"errors"
)

// https://wiki.squid-cache.org/Features/LogFormat
// http://www.squid-cache.org/Doc/config/logformat/
// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Squid_result_codes
// https://www.websense.com/content/support/library/web/v773/wcg_help/squid.aspx

/*
4.6.1:
logformat squid      %ts.%03tu %6tr %>a %Ss/%03>Hs %<st %rm %ru %[un %Sh/%<a %mt
logformat common     %>a %[ui %[un [%tl] "%rm %ru HTTP/%rv" %>Hs %<st %Ss:%Sh
logformat combined   %>a %[ui %[un [%tl] "%rm %ru HTTP/%rv" %>Hs %<st "%{Referer}>h" "%{User-Agent}>h" %Ss:%Sh
logformat referrer   %ts.%03tu %>a %{Referer}>h %ru
logformat useragent  %>a [%tl] "%{User-Agent}>h"
logformat icap_squid %ts.%03tu %6icap::tr %>A %icap::to/%03icap::Hs %icap::<st %icap::rm %icap::ru %un -/%icap::<A -
*/

/*
Valid Capture Name: [A-Za-z0-9_]+
// TODO: namings

| local                   | squid format code | description                                                            |
|-------------------------|-------------------|------------------------------------------------------------------------|
| resp_time               | %tr               | Response time (milliseconds).
| client_ip               | %>a               | Client source IP address.
| cache_code              | %Ss               | Squid request status (TCP_MISS etc).
| http_code               | %>Hs              | The HTTP response status code from Content Gateway to client.
| resp_size               | %<st              | Total size of reply sent to client (after adaptation).
| req_method              | %rm               | Request method (GET/POST etc).
| hier_code               | %Sh               | Squid hierarchy status (DEFAULT_PARENT etc).
| server_ip               | %<a               | Server IP address of the last server or peer connection.
| mime_type               | %mt               | MIME content type.
*/

func newEmptyLogLine() *logLine {
	var l logLine
	l.reset()
	return &l
}

type (
	logLine struct {
		conn  connFields
		time  timeFields
		http  httpFields
		squid squidFields
	}
	connFields struct {
		clientIP string
		serverIP string
	}
	timeFields struct {
		respTime int
	}
	httpFields struct {
		reqMethod string
		respCode  int
		mimeType  string
		respSize  int
	}
	squidFields struct {
		cacheCode string
		hierCode  string
	}
)

func (l *logLine) Assign(field string, value string) (err error) {
	if value == "" {
		return
	}

	switch field {
	case "resp_time":
		err = l.assignRespTime(value)
	case "client_ip":
		err = l.assignClientIP(value)
	case "cache_code":
		err = l.assignCacheCode(value)
	case "http_code":
		err = l.assignHTTPCode(value)
	case "resp_size":
		err = l.assignRespSize(value)
	case "req_method":
		err = l.assignReqMethod(value)
	case "hier_code":
		err = l.assignHierCode(value)
	case "mime_type":
		err = l.assignMimeType(value)
	case "server_ip":
		err = l.assignServerIP(value)
	}
	return err
}

func (l *logLine) assignRespTime(duration string) error { return nil }
func (l *logLine) assignClientIP(address string) error  { return nil }
func (l *logLine) assignCacheCode(code string) error    { return nil }
func (l *logLine) assignHTTPCode(code string) error     { return nil }
func (l *logLine) assignRespSize(bytes string) error    { return nil }
func (l *logLine) assignReqMethod(method string) error  { return nil }
func (l *logLine) assignHierCode(code string) error     { return nil }
func (l *logLine) assignMimeType(mime string) error     { return nil }
func (l *logLine) assignServerIP(address string) error  { return nil }

func (l logLine) verify() error {
	if l.empty() {
		return errors.New("")
	}
	if l.hasRespTime() && !l.isRespTimeValid() {
		return errors.New("")
	}
	if l.hasClientIP() && !l.isClientIPValid() {
		return errors.New("")
	}
	if l.hasCacheCode() && !l.isCacheCodeValid() {
		return errors.New("")
	}
	if l.hasHTTPCode() && !l.isHTTPCodeValid() {
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
	if l.hasMimeType() && !l.isMimeTypeValid() {
		return errors.New("")
	}
	if l.hasServerIP() && !l.isServerIPValid() {
		return errors.New("")
	}
	return nil
}

func (l logLine) empty() bool        { return l == emptyLogLine }
func (l logLine) hasRespTime() bool  { return !isEmptyNumber(l.time.respTime) }
func (l logLine) hasClientIP() bool  { return !isEmptyString(l.conn.clientIP) }
func (l logLine) hasCacheCode() bool { return !isEmptyString(l.squid.cacheCode) }
func (l logLine) hasHTTPCode() bool  { return !isEmptyNumber(l.http.respCode) }
func (l logLine) hasRespSize() bool  { return !isEmptyNumber(l.http.respSize) }
func (l logLine) hasReqMethod() bool { return !isEmptyString(l.http.reqMethod) }
func (l logLine) hasHierCode() bool  { return !isEmptyString(l.squid.hierCode) }
func (l logLine) hasMimeType() bool  { return !isEmptyString(l.http.mimeType) }
func (l logLine) hasServerIP() bool  { return !isEmptyString(l.conn.serverIP) }

func (l logLine) isRespTimeValid() bool  { return false }
func (l logLine) isClientIPValid() bool  { return false }
func (l logLine) isCacheCodeValid() bool { return false }
func (l logLine) isHTTPCodeValid() bool  { return false }
func (l logLine) isRespSizeValid() bool  { return false }
func (l logLine) isReqMethodValid() bool { return false }
func (l logLine) isHierCodeValid() bool  { return false }
func (l logLine) isMimeTypeValid() bool  { return false }
func (l logLine) isServerIPValid() bool  { return false }

func (l *logLine) reset() {
	l.conn.clientIP = emptyString
	l.conn.serverIP = emptyString
	l.time.respTime = emptyNumber
	l.http.reqMethod = emptyString
	l.http.respCode = emptyNumber
	l.http.respSize = emptyNumber
	l.http.mimeType = emptyString
	l.squid.cacheCode = emptyString
	l.squid.hierCode = emptyString
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
