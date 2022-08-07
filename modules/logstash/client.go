// SPDX-License-Identifier: GPL-3.0-or-later

package logstash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	jvmStatusPath = "/_node/stats"
)

type jvmStats struct {
	JVM       jvm                 `json:"jvm" stm:"jvm"`
	Process   process             `json:"process" stm:"process"`
	Event     events              `json:"event" stm:"event"`
	Pipelines map[string]pipeline `json:"pipelines" stm:"pipelines"`
}

type pipeline struct {
	Event events `json:"events" stm:"event"`
}

type events struct {
	In                        int `json:"in" stm:"in"`
	Filtered                  int `json:"filtered" stm:"filtered"`
	Out                       int `json:"out" stm:"out"`
	DurationInMillis          int `json:"duration_in_millis" stm:"duration_in_millis"`
	QueuePushDurationInMillis int `json:"queue_push_duration_in_millis" stm:"queue_push_duration_in_millis"`
}

type process struct {
	OpenFileDescriptors int `json:"open_file_descriptors" stm:"open_file_descriptors"`
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

func newClient(httpClient *http.Client, request web.Request) *client {
	return &client{httpClient: httpClient, request: request}
}

type client struct {
	httpClient *http.Client
	request    web.Request
}

func (a client) jvmStats() (*jvmStats, error) {
	req, err := a.createRequest(jvmStatusPath)
	if err != nil {
		return nil, err
	}

	resp, err := a.doRequestOK(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	var stats jvmStats
	if err = json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (a client) doRequestOK(req *http.Request) (resp *http.Response, err error) {
	if resp, err = a.httpClient.Do(req); err != nil {
		err = fmt.Errorf("error on request to %s : %w", req.URL, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
		return
	}
	return
}

func (a client) createRequest(urlPath string) (*http.Request, error) {
	req := a.request.Copy()
	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, urlPath)
	req.URL = u.String()
	return web.NewHTTPRequest(req)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
