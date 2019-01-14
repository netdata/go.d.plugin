package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	overviewURI = "/api/overview"
	nodeURI     = "/api/node/"
)

// https://www.rabbitmq.com/monitoring.html
type overview struct {
	objectTotals `json:"object_totals" stm:"object_totals"`
	queueTotals  `json:"queue_totals" stm:"queue_totals"`
	messageStats `json:"message_stats" stm:"message_stats"`
	Node         string
}

type objectTotals struct {
	Consumers   int `stm:"consumers"`
	Queues      int `stm:"queues"`
	Exchanges   int `stm:"exchanges"`
	Connections int `stm:"connections"`
	Channels    int `stm:"channels"`
}

type queueTotals struct {
	MessagesReady          int `json:"messages_ready" stm:"messages_ready"`
	MessagesUnacknowledged int `json:"messages_unacknowledged" stm:"messages_unacknowledged"`
}

// https://rawcdn.githack.com/rabbitmq/rabbitmq-management/master/priv/www/doc/stats.html
type messageStats struct {
	Ack              int `stm:"ack"`
	Publish          int `stm:"publish"`
	PublishIn        int `json:"publish_in" stm:"publish_in"`
	PublishOut       int `json:"publish_out" stm:"publish_out"`
	Confirm          int `stm:"confirm"`
	Deliver          int `stm:"deliver"`
	DeliverNoAck     int `json:"deliver_no_ack" stm:"deliver_no_ack"`
	Get              int `stm:"get"`
	GetNoAck         int `json:"get_no_ack" stm:"get_no_ack"`
	DeliverGet       int `json:"deliver_get" stm:"deliver_get"`
	Redeliver        int `stm:"redeliver"`
	ReturnUnroutable int `json:"return_unroutable" stm:"return_unroutable"`
}

type node struct {
	FDUsed      int `json:"fd_used" stm:"fd_used"`
	MemUsed     int `json:"mem_used" stm:"mem_used"`
	SocketsUsed int `json:"sockets_used" stm:"sockets_used"`
	ProcUsed    int `json:"proc_used" stm:"proc_used"`
	DiskFree    int `json:"disk_free" stm:"disk_free"`
	RunQueue    int `json:"run_queue" stm:"run_queue"`
}

type apiClient struct {
	req        web.Request
	httpClient *http.Client
	nodeName   string
}

func (a *apiClient) getOverview() (overview, error) {
	var overview overview

	req, err := a.createRequest(overviewURI)

	if err != nil {
		return overview, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return overview, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&overview); err != nil {
		return overview, fmt.Errorf("erorr on decode request to %s : %s", req.URL, err)
	}

	a.nodeName = overview.Node

	return overview, nil
}

func (a apiClient) getNodeStats() (node, error) {
	var node node

	req, err := a.createRequest(nodeURI + a.nodeName)

	if err != nil {
		return node, fmt.Errorf("error on creating request : %v", err)
	}

	resp, err := a.doRequestOK(req)

	defer closeBody(resp)

	if err != nil {
		return node, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		return node, fmt.Errorf("erorr on decode request to %s : %s", req.URL, err)
	}

	return node, nil
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
