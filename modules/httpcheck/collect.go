package httpcheck

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (hc *HTTPCheck) collect() (map[string]int64, error) {
	req, err := web.NewHTTPRequest(hc.Request)
	if err != nil {
		return nil, fmt.Errorf("error on creating HTTP requests to %s : %v", hc.Request.UserURL, err)
	}

	var mx metrics

	before := time.Now()
	resp, err := hc.client.Do(req)
	closeBody(resp)
	mx.Response.Time = time.Now().Sub(before).Nanoseconds()

	if err != nil {
		hc.Debug(err)
		hc.collectErrResponse(&mx, err)
	} else {
		hc.collectOKResponse(&mx, resp)
	}

	return stm.ToMap(mx), nil
}

func (hc HTTPCheck) collectErrResponse(mx *metrics, err error) {
	if v, ok := err.(net.Error); ok && v.Timeout() {
		mx.Request.Status.Timeout = true
	} else {
		mx.Request.Status.Failed = true
	}
}

func (hc HTTPCheck) collectOKResponse(mx *metrics, resp *http.Response) {
	mx.Request.Status.Success = true
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	mx.Response.Length = len(bodyBytes)

	if !hc.acceptedStatuses[resp.StatusCode] {
		mx.Response.BadStatusCode = true
	}

	if hc.reResponse != nil && !hc.reResponse.Match(bodyBytes) {
		mx.Response.BadContent = true
	}
}

func closeBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
}
