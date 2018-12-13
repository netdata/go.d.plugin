package rabbitmq

import (
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("rabbitmq", creator)
}

// https://www.rabbitmq.com/monitoring.html
type (
	apiOverview struct {
		objectTotals `json:"object_totals" stm:"object_totals"`
		queueTotals  `json:"queue_totals" stm:"queue_totals"`
		messageStats `json:"message_stats" stm:"message_stats"`
	}
	apiNode []node
)

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

// New creates Rabbitmq with default values
func New() *Rabbitmq {
	return &Rabbitmq{}
}

// Rabbitmq rabbitmq module
type Rabbitmq struct {
	modules.Base // should be embedded by every module

}

// Cleanup makes cleanup
func (Rabbitmq) Cleanup() {}

// Init makes initialization
func (Rabbitmq) Init() bool {
	return false
}

// Check makes check
func (Rabbitmq) Check() bool {
	return false
}

// Charts creates Charts
func (Rabbitmq) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (Rabbitmq) GatherMetrics() map[string]int64 {
	return nil
}
