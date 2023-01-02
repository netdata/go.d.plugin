// SPDX-License-Identifier: GPL-3.0-or-later

package rabbitmq

// https://www.rabbitmq.com/monitoring.html#cluster-wide-metrics
type overviewStats struct {
	ObjectTotals struct {
		Consumers   int `stm:"consumers"`
		Queues      int `stm:"queues"`
		Exchanges   int `stm:"exchanges"`
		Connections int `stm:"connections"`
		Channels    int `stm:"channels"`
	} `json:"object_totals" stm:"object_totals"`
	ChurnRates struct {
		ChannelClosed     int `json:"channel_closed" stm:"channel_closed"`
		ChannelCreated    int `json:"channel_created" stm:"channel_created"`
		ConnectionClosed  int `json:"connection_closed" stm:"connection_closed"`
		ConnectionCreated int `json:"connection_created" stm:"connection_created"`
		QueueCreated      int `json:"queue_created" stm:"queue_created"`
		QueueDeclared     int `json:"queue_declared" stm:"queue_declared"`
		QueueDeleted      int `json:"queue_deleted" stm:"queue_deleted"`
	} `json:"churn_rates" stm:"churn_rates"`
	QueueTotals struct {
		Messages               int `json:"messages" stm:"messages"`
		MessagesReady          int `json:"messages_ready" stm:"messages_ready"`
		MessagesUnacknowledged int `json:"messages_unacknowledged" stm:"messages_unacknowledged"`
	} `json:"queue_totals" stm:"queue_totals"`
	MessageStats messageStats `json:"message_stats" stm:"message_stats"`
	Node         string
}

// https://www.rabbitmq.com/monitoring.html#node-metrics
type nodeStats struct {
	FDTotal      int `json:"fd_total" stm:"fd_total"`
	FDUsed       int `json:"fd_used" stm:"fd_used"`
	MemLimit     int `json:"mem_limit" stm:"mem_limit"`
	MemUsed      int `json:"mem_used" stm:"mem_used"`
	SocketsTotal int `json:"sockets_total" stm:"sockets_total"`
	SocketsUsed  int `json:"sockets_used" stm:"sockets_used"`
	ProcTotal    int `json:"proc_total" stm:"proc_total"`
	ProcUsed     int `json:"proc_used" stm:"proc_used"`
	DiskFree     int `json:"disk_free" stm:"disk_free"`
	RunQueue     int `json:"run_queue" stm:"run_queue"`
}

type vhostStats struct {
	Name                   string       `json:"name"`
	Messages               int          `stm:"messages"`
	MessagesReady          int          `stm:"messages_ready"`
	MessagesUnacknowledged int          `stm:"messages_unacknowledged"`
	MessageStats           messageStats `json:"3message_stats" stm:"message_stats"`
}

// https://www.rabbitmq.com/monitoring.html#queue-metrics
type queueStats struct {
	Name                   string       `json:"name"`
	Vhost                  string       `json:"vhost"`
	State                  string       `json:"state"`
	Type                   string       `json:"type"`
	Messages               int          `stm:"messages"`
	MessagesPagedOut       int          `stm:"messages_paged_out"`
	MessagesPersistent     int          `stm:"messages_persistent"`
	MessagesReady          int          `stm:"messages_ready"`
	MessagesUnacknowledged int          `stm:"messages_unacknowledged"`
	MessageStats           messageStats `json:"message_stats" stm:"message_stats"`
}

// https://rawcdn.githack.com/rabbitmq/rabbitmq-server/v3.11.5/deps/rabbitmq_management/priv/www/api/index.html
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
