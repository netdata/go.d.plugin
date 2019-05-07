package logstash

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	jvmStatusPath = "/_node/stats/jvm"
)

type jvmStats struct {
	JVM jvm `stm:"jvm"`
}

type jvm struct {
	Threads struct {
		Count int `stm:"count"`
	} `stm:"threads"`
	Mem            jvmMem `stm:"mem"`
	GC             jvmGC  `stm:"gc"`
	UptimeInMillis int    `json:"uptime_in_millis" stm:"uptime_in_millis"`
}

type jvmMem struct {
	HeapUsedPercent      int `json:"heap_used_percent" stm:"heap_used_percent"`
	HeapCommittedInBytes int `json:"heap_committed_in_bytes" stm:"heap_committed_in_bytes"`
	HeapUsedInBytes      int `json:"heap_used_in_bytes" stm:"heap_used_in_bytes"`
	Pools                struct {
		Survivor jvmPool `stm:"survivor"`
		Old      jvmPool `stm:"old"`
		Young    jvmPool `stm:"eden"`
	} `stm:"pools"`
}

type jvmPool struct {
	UsedInBytes      int `json:"used_in_bytes" stm:"used_in_bytes"`
	CommittedInBytes int `json:"committed_in_bytes" stm:"committed_in_bytes"`
}

type jvmGC struct {
	Collectors struct {
		Old   gcCollector `stm:"old"`
		Young gcCollector `stm:"eden"`
	} `stm:"collectors"`
}

type gcCollector struct {
	CollectionTimeInMillis int `json:"collection_time_in_millis" stm:"collection_time_in_millis"`
	CollectionCount        int `json:"collection_count" stm:"collection_count"`
}

func newAPIClient(client *http.Client, request web.Request) *apiClient {
	return &apiClient{httpClient: client, request: request}
}

type apiClient struct {
	httpClient *http.Client
	request    web.Request
}

func (a apiClient) jvmStats() (*jvmStats, error) {
	var stats jvmStats

	req, err := a.createRequest(jvmStatusPath)

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

func (a apiClient) doRequestOK(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	if resp, err = a.httpClient.Do(req); err != nil {
		return resp, fmt.Errorf("error on request to %s : %v", req.URL, err)

	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return resp, err
}

func (a apiClient) createRequest(urlPath string) (*http.Request, error) {
	req := a.request.Copy()
	req.URL.Path = urlPath
	return web.NewHTTPRequest(req)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
