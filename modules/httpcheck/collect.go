// SPDX-License-Identifier: GPL-3.0-or-later

package httpcheck

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

type reqErrCode int

const (
	codeTimeout reqErrCode = iota
	//codeDNSLookup
	//codeParseAddress
	//codeRedirect
	codeNoConnection
)

func (hc *HTTPCheck) collect() (map[string]int64, error) {
	req, err := web.NewHTTPRequest(hc.Request)
	if err != nil {
		return nil, fmt.Errorf("error on creating HTTP requests to %s : %v", hc.Request.URL, err)
	}

	var mx metrics

	start := time.Now()
	resp, err := hc.client.Do(req)
	dur := time.Since(start)
	defer closeBody(resp)

	if err != nil {
		hc.Warning(err)
		hc.collectErrResponse(&mx, err)
	} else {
		mx.ResponseTime = durationToMs(dur)
		hc.collectOKResponse(&mx, resp)
	}

	changed := hc.metrics.Status != mx.Status
	if changed {
		mx.InState = hc.UpdateEvery
	} else {
		mx.InState = hc.metrics.InState + hc.UpdateEvery
	}
	hc.metrics = mx

	//if err == nil || mx.Status.RedirectError {
	//	mx.ResponseTime = durationToMs(end)
	//}

	return stm.ToMap(mx), nil
}

func (hc *HTTPCheck) collectErrResponse(mx *metrics, err error) {
	switch code := decodeReqError(err); code {
	default:
		panic(fmt.Sprintf("unknown request error code : %d", code))
	case codeNoConnection:
		mx.Status.NoConnection = true
	//case codeDNSLookup:
	//	mx.Status.DNSLookupError = true
	//case codeParseAddress:
	//	mx.Status.ParseAddressError = true
	//case codeRedirect:
	//	mx.Status.RedirectError = true
	case codeTimeout:
		mx.Status.Timeout = true
	}
}

func (hc *HTTPCheck) collectOKResponse(mx *metrics, resp *http.Response) {
	hc.Debugf("endpoint '%s' returned %d (%s) HTTP status code", hc.URL, resp.StatusCode, resp.Status)

	if !hc.acceptedStatuses[resp.StatusCode] {
		mx.Status.BadStatusCode = true
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		hc.Warningf("error on reading body : %v", err)
		mx.Status.BadContent = true
		return
	}

	mx.ResponseLength = len(bs)

	if hc.reResponse != nil && !hc.reResponse.Match(bs) {
		mx.Status.BadContent = true
		return
	}

	mx.Status.Success = true
}

func decodeReqError(err error) reqErrCode {
	if err == nil {
		panic("nil error")
	}
	if v, ok := err.(net.Error); ok && v.Timeout() {
		return codeTimeout
	}
	return codeNoConnection
	//
	//netErr, isNetErr := err.(net.Error)
	//if isNetErr && netErr.Timeout() {
	//	return codeTimeout
	//}
	//
	//urlErr, isURLErr := err.(*url.Error)
	//if !isURLErr {
	//	return codeNoConnection
	//}
	//
	//if urlErr.Err == web.ErrRedirectAttempted {
	//	return codeRedirect
	//}
	//
	//opErr, isOpErr := (urlErr.Err).(*net.OpError)
	//if !isOpErr {
	//	return codeNoConnection
	//}
	//
	//switch (opErr.Err).(type) {
	//case *net.DNSError:
	//	return codeDNSLookup
	//case *net.ParseError:
	//	return codeParseAddress
	//}
	//
	//return codeNoConnection
}

func closeBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

func durationToMs(duration time.Duration) int {
	return int(duration) / (int(time.Millisecond) / int(time.Nanosecond))
}
