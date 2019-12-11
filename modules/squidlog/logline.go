package squidlog

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
| client_address          | %>a               | Client source IP address.
| client_address          | %>A               | Client FQDN.
| cache_code              | %Ss               | Squid request status (TCP_MISS etc).
| http_code               | %>Hs              | The HTTP response status code from Content Gateway to client.
| resp_size               | %<st              | Total size of reply sent to client (after adaptation).
| req_method              | %rm               | Request method (GET/POST etc).
| hier_code               | %Sh               | Squid hierarchy status (DEFAULT_PARENT etc).
| server_address          | %<a               | Server IP address of the last server or peer connection.
| server_address          | %<A               | Server FQDN or peer name.
| mime_type               | %mt               | MIME content type.

Notes:
- %<a: older versions of Squid would put the origin server hostname here.
*/

var (
	errEmptyLine     = errors.New("empty line")
	errBadRespTime   = errors.New("bad response time")
	errBadClientAddr = errors.New("bad client address")
	errBadCacheCode  = errors.New("bad cache code")
	errBadHTTPCode   = errors.New("bad http code")
	errBadRespSize   = errors.New("bad response size")
	errBadReqMethod  = errors.New("bad request method")
	errBadHierCode   = errors.New("bad hier code")
	errBadServerAddr = errors.New("bad server address")
	errBadMimeType   = errors.New("bad mime type")
)

func newEmptyLogLine() *logLine {
	var l logLine
	l.reset()
	return &l
}

type (
	logLine struct {
		clientAddr string
		serverAddr string

		respTime int
		respSize int
		httpCode int

		reqMethod string
		mimeType  string

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
	case "client_address":
		err = l.assignClientAddress(value)
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
	case "server_address":
		err = l.assignServerAddress(value)
	}
	return err
}

func (l *logLine) assignRespTime(time string) error {
	v, err := strconv.Atoi(time)
	if err != nil || !isRespTimeValid(v) {
		return fmt.Errorf("assign '%s': %w", time, errBadRespTime)
	}
	l.respTime = v
	return nil
}
func (l *logLine) assignClientAddress(address string) error {
	return nil
}

func (l *logLine) assignCacheCode(code string) error {
	if !isCacheCodeValid(code) {
		return fmt.Errorf("assign '%s': %w", code, errBadCacheCode)
	}
	l.cacheCode = code
	return nil
}

func (l *logLine) assignHTTPCode(code string) error {
	v, err := strconv.Atoi(code)
	if err != nil || !isHTTPCodeValid(v) {
		return fmt.Errorf("assign '%s': %w", code, errBadHTTPCode)
	}
	l.httpCode = v
	return nil
}

func (l *logLine) assignRespSize(size string) error {
	v, err := strconv.Atoi(size)
	if err != nil || !isRespSizeValid(v) {
		return fmt.Errorf("assign '%s': %w", size, errBadRespSize)
	}
	l.respSize = v
	return nil
}

func (l *logLine) assignReqMethod(method string) error {
	if !isReqMethodValid(method) {
		return fmt.Errorf("assign '%s': %w", method, errBadReqMethod)
	}
	l.reqMethod = method
	return nil
}

func (l *logLine) assignHierCode(code string) error {
	if !isHierCodeValid(code) {
		return fmt.Errorf("assign '%s': %w", code, errBadHierCode)
	}
	l.hierCode = code
	return nil
}

func (l *logLine) assignMimeType(mime string) error {
	// ICP exchanges usually don't have any content type, and thus are logged "-".
	//Also, some weird replies have content types ":" or even empty ones.
	if mime == "-" || mime == ":" {
		return nil
	}
	l.mimeType = mime
	return nil
}

func (l *logLine) assignServerAddress(address string) error {
	// Logged as "-" if there is no hierarchy information.
	// For TCP HIT, TCP failures, cachemgr requests and all UDP requests, there is no hierarchy information.
	if address == "-" {
		return nil
	}
	l.serverAddr = address
	return nil
}

func (l logLine) verify() error {
	if l.empty() {
		return fmt.Errorf("verify: %w", errEmptyLine)
	}
	if l.hasRespTime() && !l.isRespTimeValid() {
		return fmt.Errorf("verify '%d': %w", l.respTime, errBadRespTime)
	}
	if l.hasClientAddress() && !l.isClientAddressValid() {
		return fmt.Errorf("verify '%s': %w", l.clientAddr, errBadClientAddr)
	}
	if l.hasCacheCode() && !l.isCacheCodeValid() {
		return fmt.Errorf("verify '%s': %w", l.cacheCode, errBadCacheCode)
	}
	if l.hasHTTPCode() && !l.isHTTPCodeValid() {
		return fmt.Errorf("verify '%d': %w", l.httpCode, errBadHTTPCode)
	}
	if l.hasRespSize() && !l.isRespSizeValid() {
		return fmt.Errorf("verify '%d': %w", l.respSize, errBadRespSize)
	}
	if l.hasReqMethod() && !l.isReqMethodValid() {
		return fmt.Errorf("verify '%s': %w", l.reqMethod, errBadReqMethod)
	}
	if l.hasHierCode() && !l.isHierCodeValid() {
		return fmt.Errorf("verify '%s': %w", l.hierCode, errBadHierCode)
	}
	if l.hasMimeType() && !l.isMimeTypeValid() {
		return fmt.Errorf("verify '%s': %w", l.mimeType, errBadMimeType)
	}
	if l.hasServerAddress() && !l.isServerAddressValid() {
		return fmt.Errorf("verify '%s': %w", l.serverAddr, errBadServerAddr)
	}
	return nil
}

func (l logLine) empty() bool                { return l == emptyLogLine }
func (l logLine) hasRespTime() bool          { return !isEmptyNumber(l.respTime) }
func (l logLine) hasClientAddress() bool     { return !isEmptyString(l.clientAddr) }
func (l logLine) hasCacheCode() bool         { return !isEmptyString(l.cacheCode) }
func (l logLine) hasHTTPCode() bool          { return !isEmptyNumber(l.httpCode) }
func (l logLine) hasRespSize() bool          { return !isEmptyNumber(l.respSize) }
func (l logLine) hasReqMethod() bool         { return !isEmptyString(l.reqMethod) }
func (l logLine) hasHierCode() bool          { return !isEmptyString(l.hierCode) }
func (l logLine) hasMimeType() bool          { return !isEmptyString(l.mimeType) }
func (l logLine) hasServerAddress() bool     { return !isEmptyString(l.serverAddr) }
func (l logLine) isRespTimeValid() bool      { return isRespTimeValid(l.respTime) }
func (l logLine) isClientAddressValid() bool { return false }
func (l logLine) isCacheCodeValid() bool     { return isCacheCodeValid(l.cacheCode) }
func (l logLine) isHTTPCodeValid() bool      { return isHTTPCodeValid(l.httpCode) }
func (l logLine) isRespSizeValid() bool      { return isRespSizeValid(l.respSize) }
func (l logLine) isReqMethodValid() bool     { return isReqMethodValid(l.reqMethod) }
func (l logLine) isHierCodeValid() bool      { return isHierCodeValid(l.hierCode) }
func (l logLine) isMimeTypeValid() bool      { return false }
func (l logLine) isServerAddressValid() bool { return false }

func (l *logLine) reset() {
	l.clientAddr = emptyString
	l.serverAddr = emptyString
	l.respTime = emptyNumber
	l.reqMethod = emptyString
	l.httpCode = emptyNumber
	l.respSize = emptyNumber
	l.mimeType = emptyString
	l.cacheCode = emptyString
	l.hierCode = emptyString
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

// isCacheCodeValid does not guarantee cache result code is valid, but it is very likely.
func isCacheCodeValid(code string) bool {
	// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Squid_result_codes
	if i := strings.IndexByte(code, '_'); i <= 0 {
		// TODO: is at least 1 '_' required?
		return false
	} else {
		code = code[:i]
	}
	return code == "TCP" || code == "UDP" || code == "NONE"
}

// isHierCodeValid does not guarantee hierarchy code is valid, but it is very likely.
func isHierCodeValid(code string) bool {
	// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Hierarchy_Codes
	if i := strings.IndexByte(code, '_'); i <= 0 {
		return false
	} else {
		code = code[:i]
	}
	return code == "HIER" || code == "TIMEOUT"
}

func isHTTPCodeValid(code int) bool {
	// rfc7231
	// Informational responses (100–199),
	// Successful responses (200–299),
	// Redirects (300–399),
	// Client errors (400–499),
	// Server errors (500–599).
	return code >= 100 && code <= 600
}

func isRespTimeValid(time int) bool {
	return time >= 0
}

func isRespSizeValid(size int) bool {
	return size >= 0
}

func isReqMethodValid(method string) bool {
	// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Request_methods
	if method == "GET" {
		return true
	}
	switch method {
	case "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		return true
	case "ICP_QUERY", "PURGE":
		return true
	case "PROPFIND", "PROPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK": // rfc2518
		return true
	}
	return false
}
