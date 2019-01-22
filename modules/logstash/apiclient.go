package logstash

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

type jvmStats struct {
	JVM jvm `stm:"jvm"`
}

type jvm struct {
	Mem jvmMem `stm:"mem"`
	GC  jvmGC  `stm:"gc"`
}

type jvmMem struct {
	HeapUsedPercent      int `json:"heap_used_percent",stm:"heap_used_percent"`
	HeapCommittedInBytes int `json:"heap_committed_in_bytes",stm:"heap_committed_in_bytes"`
	HeapUsedInBytes      int `json:"heap_used_in_bytes",stm:"heap_used_in_bytes"`
	Pools                struct {
		Survivor jvmPool `stm:"survivor"`
		Old      jvmPool `stm:"old"`
		Young    jvmPool `stm:"young"`
	} `stm:"pools"`
}

type jvmPool struct {
	UsedInBytes      int `json:"used_in_bytes",stm:"used_in_bytes"`
	CommittedInBytes int `json:"used_in_bytes",stm:"used_in_bytes"`
}

type jvmGC struct {
	Collectors struct {
		Old   gcCollector `stm:"old"`
		Young gcCollector `stm:"young"`
	} `stm:"collectors"`
}

type gcCollector struct {
	CollectionTimeInMillis int `json:"collection_time_in_millis",stm:"collection_time_in_millis"`
	CollectionCount        int `json:"collection_count",stm:"collection_count"`
}

type apiClient struct {
	req        web.Request
	httpClient *http.Client
}

func (a apiClient) jvmStats() (*jvmStats, error) {
	var stats jvmStats

	req, err := a.createRequest("/_node/stats/jvm")

	if err != nil {
		return nil, err
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (a apiClient) doRequest(req *http.Request) (*http.Response, error) {
	return a.httpClient.Do(req)
}

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	if resp, err = a.doRequest(req); err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)

	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (a apiClient) createRequest(uri string) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	a.req.URI = uri

	if req, err = web.NewHTTPRequest(a.req); err != nil {
		return nil, err
	}

	return req, nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
