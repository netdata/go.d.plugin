package apache

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	busyServers = "BusyServers"
	idleServers = "IdleServers"

	busyWorkers         = "BusyWorkers"
	idleWorkers         = "IdleWorkers"
	connsTotal          = "ConnsTotal"
	connsAsyncWriting   = "ConnsAsyncWriting"
	connsAsyncKeepAlive = "ConnsAsyncKeepAlive"
	connsAsyncClosing   = "ConnsAsyncClosing"
	totalAccesses       = "Total Accesses"
	totalKBytes         = "Total kBytes"
	uptime              = "Uptime"
	reqPerSec           = "ReqPerSec"
	bytesPerSec         = "BytesPerSec"
	bytesPerReq         = "BytesPerReq"
	scoreBoard          = "Scoreboard"
)

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{httpClient: client, request: request}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
}

func (a apiClient) getServerStatus() (*serverStatus, error) {
	req, err := web.NewHTTPRequest(a.request)
	if err != nil {
		return nil, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)
	defer closeBody(resp)
	if err != nil {
		return nil, err
	}

	status, err := parseResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on parsing response from %s : %v", req.URL, err)
	}

	return status, nil
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on request : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}
	return resp, nil
}

func parseResponse(r io.Reader) (*serverStatus, error) {
	s := bufio.NewScanner(r)
	var status serverStatus

	for s.Scan() {
		parts := strings.Split(s.Text(), ":")
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		default:
		case busyServers, idleServers:
			return nil, fmt.Errorf("found '%s', lighttpd data", key)
		case busyWorkers:
			status.Workers.Busy = mustParseInt(value)
		case idleWorkers:
			status.Workers.Idle = mustParseInt(value)
		case connsTotal:
			status.Connections.Total = mustParseInt(value)
		case connsAsyncWriting:
			status.Connections.Async.Writing = mustParseInt(value)
		case connsAsyncKeepAlive:
			status.Connections.Async.KeepAlive = mustParseInt(value)
		case connsAsyncClosing:
			status.Connections.Async.Closing = mustParseInt(value)
		case totalAccesses:
			status.Total.Accesses = mustParseInt(value)
		case totalKBytes:
			status.Total.KBytes = mustParseInt(value)
		case uptime:
			status.Uptime = mustParseInt(value)
		case reqPerSec:
			status.Averages.ReqPerSec = mustParseFloat(value)
		case bytesPerSec:
			status.Averages.BytesPerSec = mustParseFloat(value)
		case bytesPerReq:
			status.Averages.BytesPerReq = mustParseFloat(value)
		case scoreBoard:
			status.Scoreboard = parseScoreboard(value)
		}
	}

	return &status, nil
}

func parseScoreboard(value string) *scoreboard {
	//  “_” Waiting for Connection
	// “S” Starting up
	// “R” Reading Request
	// “W” Sending Reply
	// “K” Keepalive (read)
	// “D” DNS Lookup
	// “C” Closing connection
	// “L” Logging
	// “G” Gracefully finishing
	// “I” Idle cleanup of worker
	// “.” Open slot with no current process
	var sb scoreboard
	for _, s := range strings.Split(value, "") {

		switch s {
		case "_":
			sb.Waiting++
		case "S":
			sb.Starting++
		case "R":
			sb.Reading++
		case "W":
			sb.Sending++
		case "K":
			sb.KeepAlive++
		case "D":
			sb.DNSLookup++
		case "C":
			sb.Closing++
		case "L":
			sb.Logging++
		case "G":
			sb.Finishing++
		case "I":
			sb.IdleCleanup++
		case ".":
			sb.Open++
		}
	}

	return &sb
}

func mustParseInt(value string) *int64 {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return &v
}

func mustParseFloat(value string) *float64 {
	v, err := strconv.ParseFloat(value, 10)
	if err != nil {
		panic(err)
	}
	return &v
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
