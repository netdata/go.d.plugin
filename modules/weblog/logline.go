package weblog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// TODO: it is not clear how to handle "-", current handling is not good
// In general it is:
//   - If a field is unused in a particular entry dash "-" marks the omitted field.
// In addition to that "-" is used as zero value in:
//   - apache: %b '-' when no bytes are sent.
//
// Log Format:
//  - CLF: https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format
//  - ELF: https://www.w3.org/TR/WD-logfile.html
//  - Apache CLF: https://httpd.apache.org/docs/trunk/logs.html#common

// Variables:
//  - nginx: http://nginx.org/en/docs/varindex.html
//  - apache: http://httpd.apache.org/docs/current/mod/mod_log_config.html#logformat

/*
| name               | nginx                   | apache    |
|--------------------|-------------------------|-----------|
| vhost              | $host ($http_host)      | %v        | name of the server which accepted a request.
| port               | $server_port            | %p        |
| req_scheme         | $scheme                 | -         | “http” or “https”.
| req_client         | $remote_addr            | %a (%h)   | %h: logs the IP address if HostnameLookups is Off.
| request            | $request                | %r        | req_method + req_uri + req_protocol.
| req_method         | $request_method         | %m        |
| req_url            | $request_uri            | %U        | nginx: w/ queries, apache: w/o.
| req_proto          | $server_protocol        | %H        | usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”.
| resp_status        | $status                 | %s (%>s)  |
| req_size           | $request_length         | $I        | w/ http headers, apache: need "mod_logio".
| resp_size          | $bytes_sent             | %O        | w/ http headers.
| resp_size          | $body_bytes_sent        | %B (%b)   | w/o http headers. %b '-' when no bytes are sent.
| req_proc_time      | $request_time           | %D        | the time taken to serve the request.
| ups_resp_time      | $upstream_response_time | -         | time spent on rec the response from the upstream server.
| ssl_proto          | $ssl_protocol           | -         | protocol of an established SSL connection.
| ssl_cipher_suite   | $ssl_cipher             | -         | string of ciphers used for an established SSL connection.
| custom             | -                       | -         |
*/

/*
Apache:
Since httpd 2.0, unlike 1.3, the %b and %B format strings do not represent the number of bytes sent to the client,
but simply the size in bytes of the HTTP response. It will will differ, for instance, if the connection is aborted,
or if SSL is used.
The %O format provided by mod_logio will log the actual number of bytes sent over the network
*/

var (
	errEmptyLine         = errors.New("empty line")
	errBadVhost          = errors.New("bad vhost")
	errBadVhostPort      = errors.New("bad vhost with port")
	errBadPort           = errors.New("bad port")
	errBadReqScheme      = errors.New("bad req scheme")
	errBadReqClient      = errors.New("bad req client")
	errBadRequest        = errors.New("bad request")
	errBadReqMethod      = errors.New("bad req method")
	errBadReqURL         = errors.New("bad req url")
	errBadReqProto       = errors.New("bad req protocol")
	errBadReqSize        = errors.New("bad req size")
	errBadRespCode       = errors.New("bad resp status code")
	errBadRespSize       = errors.New("bad resp size")
	errBadReqProcTime    = errors.New("bad req processing time")
	errBadUpsRespTime    = errors.New("bad upstream resp time")
	errBadSSLProto       = errors.New("bad ssl protocol")
	errBadSSLCipherSuite = errors.New("bad ssl cipher suite")
)

func newEmptyLogLine() *logLine {
	var l logLine
	l.custom.fields = make(map[string]struct{})
	l.custom.values = make([]customValue, 0, 20)
	l.reset()
	return &l
}

type (
	logLine struct {
		lineWebFields
		custom lineCustomFields
	}
	lineWebFields struct {
		vhost          string
		port           string // apache has no $scheme, this is workaround to collect per scheme requests. Lame.
		reqScheme      string
		reqClient      string
		reqMethod      string
		reqURL         string
		reqProto       string
		reqSize        int
		reqProcTime    float64
		respCode       int
		respSize       int
		upsRespTime    float64
		sslProto       string
		sslCipherSuite string
	}
	lineCustomFields struct {
		fields map[string]struct{}
		values []customValue
	}
	customValue struct {
		name  string
		value string
	}
)

func (l *logLine) Assign(field string, value string) (err error) {
	if value == "" {
		return
	}

	switch field {
	case "host", "http_host", "v":
		err = l.assignVhost(value)
	case "server_port", "p":
		err = l.assignPort(value)
	case "host:$server_port", "v:%p":
		err = l.assignVhostWithPort(value)
	case "scheme":
		err = l.assignReqScheme(value)
	case "remote_addr", "a", "h":
		err = l.assignReqClient(value)
	case "request", "r":
		err = l.assignRequest(value)
	case "request_method", "m":
		err = l.assignReqMethod(value)
	case "request_uri", "U":
		err = l.assignReqURL(value)
	case "server_protocol", "H":
		err = l.assignReqProto(value)
	case "status", "s", ">s":
		err = l.assignRespCode(value)
	case "request_length", "I":
		err = l.assignReqSize(value)
	case "bytes_sent", "body_bytes_sent", "b", "O", "B":
		err = l.assignRespSize(value)
	case "request_time", "D":
		err = l.assignReqProcTime(value)
	case "upstream_response_time":
		err = l.assignUpsRespTime(value)
	case "ssl_protocol":
		err = l.assignSSLProto(value)
	case "ssl_cipher":
		err = l.assignSSLCipherSuite(value)
	default:
		err = l.assignCustom(field, value)
	}
	return err
}

const hyphen = "-"

func (l *logLine) assignCustom(field, value string) error {
	if len(l.custom.fields) == 0 || value == hyphen {
		return nil
	}
	if _, ok := l.custom.fields[field]; ok {
		l.custom.values = append(l.custom.values, customValue{name: field, value: value})
	}
	return nil
}

func (l *logLine) assignVhost(vhost string) error {
	if vhost == hyphen {
		return nil
	}
	// nginx $host and $http_host returns ipv6 in [], apache not
	if idx := strings.IndexByte(vhost, ']'); idx > 0 {
		vhost = vhost[1:idx]
	}
	l.vhost = vhost
	return nil
}

func (l *logLine) assignPort(port string) error {
	if port == hyphen {
		return nil
	}
	if !isValidPort(port) {
		return fmt.Errorf("assign '%s' : %w", port, errBadPort)
	}
	l.port = port
	return nil
}

func (l *logLine) assignVhostWithPort(vhostPort string) error {
	if vhostPort == hyphen {
		return nil
	}
	idx := strings.LastIndexByte(vhostPort, ':')
	if idx == -1 {
		return fmt.Errorf("assign '%s' : %w", vhostPort, errBadVhostPort)
	}
	if err := l.assignPort(vhostPort[idx+1:]); err != nil {
		return fmt.Errorf("assign '%s' : %w", vhostPort, errBadVhostPort)
	}
	if err := l.assignVhost(vhostPort[0:idx]); err != nil {
		return fmt.Errorf("assign '%s' : %w", vhostPort, errBadVhostPort)
	}
	return nil
}

func (l *logLine) assignReqScheme(scheme string) error {
	if scheme == hyphen {
		return nil
	}
	if !isValidScheme(scheme) {
		return fmt.Errorf("assign '%s' : %w", scheme, errBadReqScheme)
	}
	l.reqScheme = scheme
	return nil
}

func (l *logLine) assignReqClient(client string) error {
	if client == hyphen {
		return nil
	}
	l.reqClient = client
	return nil
}

func (l *logLine) assignRequest(request string) error {
	if request == hyphen {
		return nil
	}
	first := strings.IndexByte(request, ' ')
	if first < 0 {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	last := strings.LastIndexByte(request, ' ')
	if first == last {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	proto := request[last+1:]
	url := request[first+1 : last]
	method := request[0:first]
	if err := l.assignReqMethod(method); err != nil {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	if err := l.assignReqURL(url); err != nil {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	if err := l.assignReqProto(proto); err != nil {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	return nil
}

func (l *logLine) assignReqMethod(method string) error {
	if method == hyphen {
		return nil
	}
	if !isValidReqMethod(method) {
		return fmt.Errorf("assign '%s' : %w", method, errBadReqMethod)
	}
	l.reqMethod = method
	return nil
}

func (l *logLine) assignReqURL(url string) error {
	if url == hyphen {
		return nil
	}
	l.reqURL = url
	return nil
}

func (l *logLine) assignReqProto(proto string) error {
	if proto == hyphen {
		return nil
	}
	if !isValidReqProto(proto) {
		return fmt.Errorf("assign '%s': %w", proto, errBadReqProto)
	}
	l.reqProto = proto[5:]
	return nil
}

func (l *logLine) assignRespCode(status string) error {
	if status == hyphen {
		return nil
	}
	v, err := strconv.Atoi(status)
	if err != nil || !isValidRespCode(v) {
		return fmt.Errorf("assign '%s': %w", status, errBadRespCode)
	}
	l.respCode = v
	return nil
}

func (l *logLine) assignReqSize(size string) error {
	// apache: can be "-" according web_log py regexp.
	if size == hyphen {
		l.reqSize = 0
		return nil
	}
	v, err := strconv.Atoi(size)
	if err != nil || !isValidSize(v) {
		return fmt.Errorf("assign '%s': %w", size, errBadReqSize)
	}
	l.reqSize = v
	return nil
}

func (l *logLine) assignRespSize(size string) error {
	// apache: %b. In CLF format, i.e. a '-' rather than a 0 when no bytes are sent.
	if size == hyphen {
		l.respSize = 0
		return nil
	}
	v, err := strconv.Atoi(size)
	if err != nil || !isValidSize(v) {
		return fmt.Errorf("assign '%s': %w", size, errBadRespSize)
	}
	l.respSize = v
	return nil
}

func (l *logLine) assignReqProcTime(time string) error {
	if time == hyphen {
		return nil
	}
	if time == "0.000" {
		l.reqProcTime = 0
		return nil
	}
	v, err := strconv.ParseFloat(time, 64)
	if err != nil || !isValidTime(v) {
		return fmt.Errorf("assign '%s': %w", time, errBadReqProcTime)
	}
	l.reqProcTime = v * respTimeMultiplier(time)
	return nil
}

func (l *logLine) assignUpsRespTime(time string) error {
	if time == hyphen {
		return nil
	}
	// times of several responses are separated by commas and colons.
	if idx := strings.IndexByte(time, ','); idx >= 0 {
		time = time[0:idx]
	}
	v, err := strconv.ParseFloat(time, 64)
	if err != nil || !isValidTime(v) {
		return fmt.Errorf("assign '%s': %w", time, errBadUpsRespTime)
	}
	l.upsRespTime = v * respTimeMultiplier(time)
	return nil
}

func (l *logLine) assignSSLProto(proto string) error {
	if proto == hyphen {
		return nil
	}
	if !isValidSSLProto(proto) {
		return fmt.Errorf("assign '%s': %w", proto, errBadSSLProto)
	}
	l.sslProto = proto
	return nil
}

func (l *logLine) assignSSLCipherSuite(cipher string) error {
	if cipher == hyphen {
		return nil
	}
	if idx := strings.IndexByte(cipher, '-'); idx <= 0 {
		return fmt.Errorf("assign '%s': %w", cipher, errBadSSLCipherSuite)
	}
	l.sslCipherSuite = cipher
	return nil
}

func (l logLine) verify() error {
	if l.isEmpty() {
		return fmt.Errorf("verify: %w", errEmptyLine)
	}
	if l.hasRespCode() && !l.validRespCode() {
		return fmt.Errorf("verify '%d': %w", l.respCode, errBadRespCode)
	}
	if l.hasVhost() && !l.validVhost() {
		return fmt.Errorf("verify '%s': %w", l.vhost, errBadVhost)
	}
	if l.hasPort() && !l.validPort() {
		return fmt.Errorf("verify '%s': %w", l.port, errBadPort)
	}
	if l.hasReqScheme() && !l.validReqScheme() {
		return fmt.Errorf("verify '%s': %w", l.reqScheme, errBadReqScheme)
	}
	if l.hasReqClient() && !l.validReqClient() {
		return fmt.Errorf("verify '%s': %w", l.reqClient, errBadReqClient)
	}
	if l.hasReqMethod() && !l.validReqMethod() {
		return fmt.Errorf("verify '%s': %w", l.reqMethod, errBadReqMethod)
	}
	if l.hasReqURL() && !l.validReqURL() {
		return fmt.Errorf("verify '%s': %w", l.reqURL, errBadReqURL)
	}
	if l.hasReqProto() && !l.validReqProto() {
		return fmt.Errorf("verify '%s': %w", l.reqProto, errBadReqProto)
	}
	if l.hasReqSize() && !l.validReqSize() {
		return fmt.Errorf("verify '%d': %w", l.reqSize, errBadReqSize)
	}
	if l.hasRespSize() && !l.validRespSize() {
		return fmt.Errorf("verify '%d': %w", l.respSize, errBadRespSize)
	}
	if l.hasReqProcTime() && !l.validReqProcTime() {
		return fmt.Errorf("verify '%f': %w", l.reqProcTime, errBadReqProcTime)
	}
	if l.hasUpsRespTime() && !l.validUpsRespTime() {
		return fmt.Errorf("verify '%f': %w", l.upsRespTime, errBadUpsRespTime)
	}
	if l.hasSSLProto() && !l.validSSLProto() {
		return fmt.Errorf("verify '%s': %w", l.sslProto, errBadSSLProto)
	}
	if l.hasSSLCipherSuite() && !l.validSSLCipherSuite() {
		return fmt.Errorf("verify '%s': %w", l.sslCipherSuite, errBadSSLCipherSuite)
	}
	return nil
}

func (l logLine) isEmpty() bool             { return !l.hasWebFields() && !l.hasCustomFields() }
func (l logLine) hasCustomFields() bool     { return len(l.custom.values) > 0 }
func (l logLine) hasWebFields() bool        { return l.lineWebFields != emptyWebFields }
func (l logLine) hasVhost() bool            { return !isEmptyString(l.vhost) }
func (l logLine) hasPort() bool             { return !isEmptyString(l.port) }
func (l logLine) hasReqScheme() bool        { return !isEmptyString(l.reqScheme) }
func (l logLine) hasReqClient() bool        { return !isEmptyString(l.reqClient) }
func (l logLine) hasReqMethod() bool        { return !isEmptyString(l.reqMethod) }
func (l logLine) hasReqURL() bool           { return !isEmptyString(l.reqURL) }
func (l logLine) hasReqProto() bool         { return !isEmptyString(l.reqProto) }
func (l logLine) hasRespCode() bool         { return !isEmptyNumber(l.respCode) }
func (l logLine) hasReqSize() bool          { return !isEmptyNumber(l.reqSize) }
func (l logLine) hasRespSize() bool         { return !isEmptyNumber(l.respSize) }
func (l logLine) hasReqProcTime() bool      { return !isEmptyNumber(int(l.reqProcTime)) }
func (l logLine) hasUpsRespTime() bool      { return !isEmptyNumber(int(l.upsRespTime)) }
func (l logLine) hasSSLProto() bool         { return !isEmptyString(l.sslProto) }
func (l logLine) hasSSLCipherSuite() bool   { return !isEmptyString(l.sslCipherSuite) }
func (l logLine) validVhost() bool          { return reVhost.MatchString(l.vhost) }
func (l logLine) validPort() bool           { return isValidPort(l.port) }
func (l logLine) validReqScheme() bool      { return isValidScheme(l.reqScheme) }
func (l logLine) validReqClient() bool      { return reClient.MatchString(l.reqClient) }
func (l logLine) validReqMethod() bool      { return isValidReqMethod(l.reqMethod) }
func (l logLine) validReqURL() bool         { return isValidURL(l.reqMethod, l.reqURL) }
func (l logLine) validReqProto() bool       { return isValidReqProtoVer(l.reqProto) }
func (l logLine) validRespCode() bool       { return isValidRespCode(l.respCode) }
func (l logLine) validReqSize() bool        { return isValidSize(l.reqSize) }
func (l logLine) validRespSize() bool       { return isValidSize(l.respSize) }
func (l logLine) validReqProcTime() bool    { return isValidTime(l.reqProcTime) }
func (l logLine) validUpsRespTime() bool    { return isValidTime(l.upsRespTime) }
func (l logLine) validSSLProto() bool       { return isValidSSLProto(l.sslProto) }
func (l logLine) validSSLCipherSuite() bool { return reCipherSuite.MatchString(l.sslCipherSuite) }

func (l *logLine) reset() {
	l.lineWebFields = emptyWebFields
	l.custom.values = l.custom.values[:0]
}

var (
	// TODO: reClient doesnt work with %h when HostnameLookups is On.
	reVhost       = regexp.MustCompile(`^[a-zA-Z0-9-:.]+$`)
	reClient      = regexp.MustCompile(`^([\da-f:.]+|localhost)$`)
	reCipherSuite = regexp.MustCompile(`^[A-Z0-9-]+$`) // openssl -v
)

var emptyWebFields = lineWebFields{
	vhost:          emptyString,
	port:           emptyString,
	reqScheme:      emptyString,
	reqClient:      emptyString,
	reqMethod:      emptyString,
	reqURL:         emptyString,
	reqProto:       emptyString,
	reqSize:        emptyNumber,
	reqProcTime:    emptyNumber,
	respCode:       emptyNumber,
	respSize:       emptyNumber,
	upsRespTime:    emptyNumber,
	sslProto:       emptyString,
	sslCipherSuite: emptyString,
}

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

func isValidURL(method, url string) bool {
	// CONNECT www.example.com:443 HTTP/1.1
	if method == "CONNECT" {
		return !hasBadURICharacter(url)
	}
	return url[0] == '/' && !hasBadURICharacter(url)
}

func hasBadURICharacter(s string) bool {
	// A URI specified by RFC1738, relative URIs are specified by RFC1808.
	// URIs cannot by definition include whitespace or ASCII control characters.
	for _, v := range s {
		if unicode.IsSpace(v) || unicode.IsControl(v) {
			return true
		}
	}
	return false
}

func isValidReqMethod(method string) bool {
	if method == "GET" {
		return true
	}
	switch method {
	case "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		return true
	}
	return false
}

func isValidReqProto(proto string) bool {
	return len(proto) >= 6 && strings.HasPrefix(proto, "HTTP/") && isValidReqProtoVer(proto[5:])
}

func isValidReqProtoVer(version string) bool {
	if version == "1.1" {
		return true
	}
	switch version {
	case "1", "1.0", "2", "2.0":
		return true
	}
	return false
}

func isValidPort(port string) bool {
	v, err := strconv.Atoi(port)
	return err == nil && v >= 80 && v <= 49151
}

func isValidScheme(scheme string) bool {
	return scheme == "http" || scheme == "https"
}

func isValidRespCode(code int) bool {
	// rfc7231
	// Informational responses (100–199),
	// Successful responses (200–299),
	// Redirects (300–399),
	// Client errors (400–499),
	// Server errors (500–599).
	return code >= 100 && code <= 600
}

func isValidSize(size int) bool {
	return size >= 0
}

func isValidTime(time float64) bool {
	return time >= 0
}

func isValidSSLProto(proto string) bool {
	if proto == "TLSv1.2" {
		return true
	}
	switch proto {
	case "TLSv1.3", "SSLv2", "SSLv3", "TLSv1", "TLSv1.1":
		return true
	}
	return false
}

func respTimeMultiplier(time string) float64 {
	// Convert to microseconds:
	//   - nginx time is in seconds with a milliseconds resolution.
	if strings.IndexByte(time, '.') > 0 {
		return 1e6
	}
	//   - apache time is in microseconds.
	return 1
}
