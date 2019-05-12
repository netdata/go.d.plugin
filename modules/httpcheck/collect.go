package httpcheck

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

type reqErrCode int

const (
	codeTimeout reqErrCode = iota
	codeDNSLookup
	codeParseAddress
	codeRedirect
	codeNoConnection
)

func (hc *HTTPCheck) collect() (map[string]int64, error) {
	req, err := web.NewHTTPRequest(hc.Request)
	if err != nil {
		return nil, fmt.Errorf("error on creating HTTP requests to %s : %v", hc.Request.UserURL, err)
	}

	var mx metrics

	start := time.Now()
	resp, err := hc.client.Do(req)
	defer closeBody(resp)
	end := time.Since(start)

	if err != nil {
		hc.Warning(err)
		hc.collectErrResponse(&mx, err)
	} else {
		mx.ResponseTime = durationToMs(end)
		hc.collectOKResponse(&mx, resp)
	}

	return stm.ToMap(mx), nil
}

func (hc HTTPCheck) collectErrResponse(mx *metrics, err error) {
	switch code := decodeReqError(err); code {
	default:
		panic(fmt.Sprintf("unknown request error code : %d", code))
	case codeNoConnection:
		mx.Status.NoConnection = true
	case codeDNSLookup:
		mx.Status.DNSLookupError = true
	case codeParseAddress:
		mx.Status.ParseAddressError = true
	case codeRedirect:
		mx.Status.RedirectError = true
	case codeTimeout:
		mx.Status.Timeout = true
	}
}

func (hc HTTPCheck) collectOKResponse(mx *metrics, resp *http.Response) {
	if !hc.acceptedStatuses[resp.StatusCode] {
		mx.Status.BadStatusCode = true
		return
	}

	if hc.reResponse != nil {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mx.Status.BodyReadError = true
			return
		}
		if !hc.reResponse.Match(bs) {
			mx.Status.BadContent = true
			return
		}
	}

	mx.Status.Success = true
}

func decodeReqError(err error) reqErrCode {
	if err == nil {
		panic("nil error")
	}

	netErr, isNetErr := err.(net.Error)
	if isNetErr && netErr.Timeout() {
		return codeTimeout
	}

	urlErr, isURLErr := err.(*url.Error)
	if !isURLErr {
		return codeNoConnection
	}

	if urlErr.Err == web.ErrRedirectAttempted {
		return codeRedirect
	}

	opErr, isOpErr := (urlErr.Err).(*net.OpError)
	if !isOpErr {
		return codeNoConnection
	}

	switch (opErr.Err).(type) {
	case *net.DNSError:
		return codeDNSLookup
	case *net.ParseError:
		return codeParseAddress
	}

	return codeNoConnection
}

func closeBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
}

func durationToMs(duration time.Duration) int64 {
	return int64(duration) / (int64(time.Millisecond) / int64(time.Nanosecond))
}
