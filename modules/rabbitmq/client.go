package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	overviewURLPath = "/api/overview"
	nodeURLPath     = "/api/nodes/"
	vhostsURLPath   = "/api/vhosts"
)

// https://www.rabbitmq.com/monitoring.html
type (
	overviewStats struct {
		ObjectTotals *objectTotals `json:"object_totals" stm:"object_totals"`
		QueueTotals  *queueTotals  `json:"queue_totals" stm:"queue_totals"`
		MessageStats *messageStats `json:"message_stats" stm:"message_stats"`
		Node         string
	}

	objectTotals struct {
		Consumers   int `stm:"consumers"`
		Queues      int `stm:"queues"`
		Exchanges   int `stm:"exchanges"`
		Connections int `stm:"connections"`
		Channels    int `stm:"channels"`
	}

	queueTotals struct {
		MessagesReady          int `json:"messages_ready" stm:"messages_ready"`
		MessagesUnacknowledged int `json:"messages_unacknowledged" stm:"messages_unacknowledged"`
	}

	// https://rawcdn.githack.com/rabbitmq/rabbitmq-management/master/priv/www/doc/stats.html
	messageStats struct {
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
)

type nodeStats struct {
	FDUsed      *int `json:"fd_used" stm:"fd_used"`
	MemUsed     *int `json:"mem_used" stm:"mem_used"`
	SocketsUsed *int `json:"sockets_used" stm:"sockets_used"`
	ProcUsed    *int `json:"proc_used" stm:"proc_used"`
	DiskFree    *int `json:"disk_free" stm:"disk_free"`
	RunQueue    *int `json:"run_queue" stm:"run_queue"`
}

type (
	vhostStats struct {
		Name         string
		MessageStats *messageStats `json:"message_stats" stm:"message_stats"`
	}

	vhostsStats []vhostStats
)

func newClient(httpClient *http.Client, request web.Request) *client {
	return &client{httpClient: httpClient, request: request}
}

type client struct {
	request    web.Request
	httpClient *http.Client
	nodeName   string
}

func (c *client) findNodeName() error {
	stats, err := c.scrapeOverview()
	if err != nil {
		return err
	}
	c.nodeName = stats.Node
	return nil
}

func (c *client) scrapeOverview() (*overviewStats, error) {
	var stats overviewStats

	req := c.request.Copy()
	req.URL.Path = overviewURLPath

	err := c.doOKWithDecode(req, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (c *client) scrapeNodeStats() (*nodeStats, error) {
	if c.nodeName == "" {
		return nil, errors.New("node name not set")
	}
	var stats nodeStats

	req := c.request.Copy()
	req.URL.Path = nodeURLPath + c.nodeName

	err := c.doOKWithDecode(req, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (c *client) scrapeVhostsStats() (vhostsStats, error) {
	var stats vhostsStats

	req := c.request.Copy()
	req.URL.Path = vhostsURLPath

	err := c.doOKWithDecode(req, &stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (c *client) do(req web.Request) (*http.Response, error) {
	httpReq, err := web.NewHTTPRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error on creating http request to %s : %v", req.URL, err)
	}
	return c.httpClient.Do(httpReq)
}

func (c *client) doOK(req web.Request) (*http.Response, error) {
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned %d", req.URL, resp.StatusCode)
	}
	return resp, nil
}

func (c *client) doOKWithDecode(req web.Request, dst interface{}) error {
	resp, err := c.doOK(req)
	defer closeBody(resp)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(dst)
	if err != nil {
		return fmt.Errorf("error on decoding response from %s : %v", req.URL, err)
	}
	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
