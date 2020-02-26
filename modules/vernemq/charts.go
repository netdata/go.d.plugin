package vernemq

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
)

var charts = Charts{
	chartOpenSockets.Copy(),
	chartSocketEvents.Copy(),
	chartClientKeepaliveExpired.Copy(),
	chartSocketErrors.Copy(),
	chartSocketCloseTimeout.Copy(),

	chartQueueProcesses.Copy(),
	chartQueueProcessesEvents.Copy(),
	chartQueueProcessesOfflineStorage.Copy(),
	chartQueueMessagesInQueues.Copy(),
	chartQueueMessages.Copy(),
	chartQueueUndeliveredMessages.Copy(),

	chartRouterSubscriptions.Copy(),
	chartRouterMatchedSubscriptions.Copy(),
	chartRouterMemory.Copy(),

	chartSystemUtilization.Copy(),
	chartSystemProcesses.Copy(),
	chartSystemReductions.Copy(),
	chartSystemContextSwitches.Copy(),
	chartSystemIO.Copy(),
	chartSystemRunQueue.Copy(),
	chartSystemGCCount.Copy(),
	chartSystemGCWordsReclaimed.Copy(),
	chartSystemMemoryAllocated.Copy(),

	chartBandwidth.Copy(),

	chartRetainMessages.Copy(),
	chartRetainMemoryUsage.Copy(),

	chartClusterCommunicationBandwidth.Copy(),
	chartClusterCommunicationDropped.Copy(),
	chartNetSplitUnresolved.Copy(),
	chartNetSplitEvents.Copy(),

	chartMQTTv5AUTH.Copy(),
	chartMQTTv5AUTHReceivedReason.Copy(),
	chartMQTTv5AUTHSentReason.Copy(),

	chartMQTTv4v5CONNECT.Copy(),
	chartMQTTv4v5CONNECTSentReason.Copy(),

	chartMQTTv4v5DISCONNECT.Copy(),
	chartMQTTv5DISCONNECTReceivedReason.Copy(),
	chartMQTTv5DISCONNECTSentReason.Copy(),

	chartMQTTv4v5SUBSCRIBE.Copy(),
	chartMQTTv4v5SUBSCRIBEError.Copy(),
	chartMQTTv4v5SUBSCRIBEAuthError.Copy(),

	chartMQTTv4v5UNSUBSCRIBE.Copy(),
	chartMQTTv4v5UNSUBSCRIBEError.Copy(),

	chartMQTTv4v5PUBLISH.Copy(),
	chartMQTTv4v5PUBLISHErrors.Copy(),
	chartMQTTv4v5PUBLISHAuthErrors.Copy(),
	chartMQTTv4v5PUBACK.Copy(),
	chartMQTTv5PUBACKReceivedReason.Copy(),
	chartMQTTv5PUBACKSentReason.Copy(),
	chartMQTTv4v5PUBACKUnexpected.Copy(),
	chartMQTTv4v5PUBREC.Copy(),
	chartMQTTv5PUBRECReceivedReason.Copy(),
	chartMQTTv5PUBRECSentReason.Copy(),
	chartMQTTv4PUBRECUnexpected.Copy(),
	chartMQTTv4v5PUBREL.Copy(),
	chartMQTTv5PUBRELReceivedReason.Copy(),
	chartMQTTv5PUBRELSentReason.Copy(),
	chartMQTTv4v5PUBCOMP.Copy(),
	chartMQTTv5PUBCOMReceivedReason.Copy(),
	chartMQTTv5PUBCOMSentReason.Copy(),
	chartMQTTv4v5PUBCOMPUnexpected.Copy(),

	chartMQTTv4v5PING.Copy(),

	chartUptime.Copy(),
}

// Sockets
var (
	chartOpenSockets = Chart{
		ID:    "sockets",
		Title: "Open Sockets",
		Units: "sockets",
		Fam:   "sockets",
		Ctx:   "vernemq.sockets",
		Dims: Dims{
			{ID: "open_sockets", Name: "open"},
		},
	}
	chartSocketEvents = Chart{
		ID:    "socket_events",
		Title: "Socket Open and Close Events",
		Units: "events/s",
		Fam:   "sockets",
		Ctx:   "vernemq.socket_operations",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricSocketOpen, Name: "open", Algo: module.Incremental},
			{ID: metricSocketClose, Name: "close", Algo: module.Incremental},
		},
	}
	chartClientKeepaliveExpired = Chart{
		ID:    "client_keepalive_expired",
		Title: "Closed Sockets due to Keepalive Time Expired",
		Units: "sockets/s",
		Fam:   "sockets",
		Ctx:   "vernemq.client_keepalive_expired",
		Dims: Dims{
			{ID: metricClientKeepaliveExpired, Name: "closed", Algo: module.Incremental},
		},
	}
	chartSocketCloseTimeout = Chart{
		ID:    "socket_close_timeout",
		Title: "Closed Sockets due to CONNECT Frame Hasn't Been Received On Time",
		Units: "sockets/s",
		Fam:   "sockets",
		Ctx:   "vernemq.socket_close_timeout",
		Dims: Dims{
			{ID: metricSocketCloseTimeout, Name: "closed", Algo: module.Incremental},
		},
	}
	chartSocketErrors = Chart{
		ID:    "socket_errors",
		Title: "Socket Errors",
		Units: "errors/s",
		Fam:   "sockets",
		Ctx:   "vernemq.socket_errors",
		Dims: Dims{
			{ID: metricSocketError, Name: "errors", Algo: module.Incremental},
		},
	}
)

// Queues
var (
	chartQueueProcesses = Chart{
		ID:    "queue_processes",
		Title: "Living Queues in an Online or an Offline State",
		Units: "queue processes",
		Fam:   "queues",
		Ctx:   "vernemq.queue_processes",
		Dims: Dims{
			{ID: metricQueueProcesses, Name: "queue_processes"},
		},
	}
	chartQueueProcessesEvents = Chart{
		ID:    "queue_processes_events",
		Title: "Queue Processes Setup and Teardown Events",
		Units: "events/s",
		Fam:   "queues",
		Ctx:   "vernemq.queue_processes_operations",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricQueueSetup, Name: "setup", Algo: module.Incremental},
			{ID: metricQueueTeardown, Name: "teardown", Algo: module.Incremental},
		},
	}
	chartQueueProcessesOfflineStorage = Chart{
		ID:    "queue_process_init_from_storage",
		Title: "Queue Processes Initialized from Offline Storage",
		Units: "queue processes/s",
		Fam:   "queues",
		Ctx:   "vernemq.queue_process_init_from_storage",
		Dims: Dims{
			{ID: metricQueueInitializedFromStorage, Name: "queue processes", Algo: module.Incremental},
		},
	}
	chartQueueMessagesInQueues = Chart{
		ID:    "queue_messages_in_queues",
		Title: "PUBLISH Messages that Currently in the Queues",
		Units: "messages",
		Fam:   "queues",
		Ctx:   "vernemq.queue_messages_in_queues",
		Dims: Dims{
			{ID: "queue_messages_current", Name: "in"},
		},
	}
	chartQueueMessages = Chart{
		ID:    "queue_messages",
		Title: "Received and Sent PUBLISH Messages",
		Units: "messages/s",
		Fam:   "queues",
		Ctx:   "vernemq.queue_messages",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricQueueMessageIn, Name: "received", Algo: module.Incremental},
			{ID: metricQueueMessageOut, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartQueueUndeliveredMessages = Chart{
		ID:    "queue_undelivered_messages",
		Title: "Undelivered PUBLISH Messages",
		Units: "messages/s",
		Fam:   "queues",
		Ctx:   "vernemq.queue_undelivered_messages",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricQueueMessageDrop, Name: "dropped", Algo: module.Incremental},
			{ID: metricQueueMessageExpired, Name: "expired", Algo: module.Incremental},
			{ID: metricQueueMessageUnhandled, Name: "unhandled", Algo: module.Incremental},
		},
	}
)

// Subscriptions
var (
	chartRouterSubscriptions = Chart{
		ID:    "router_subscriptions",
		Title: "Subscriptions in the Routing Table",
		Units: "subscriptions",
		Fam:   "subscriptions",
		Ctx:   "vernemq.router_subscriptions",
		Dims: Dims{
			{ID: metricRouterSubscriptions, Name: "subscriptions"},
		},
	}
	chartRouterMatchedSubscriptions = Chart{
		ID:    "router_matched_subscriptions",
		Title: "Matched Subscriptions",
		Units: "subscriptions/s",
		Fam:   "subscriptions",
		Ctx:   "vernemq.router_matched_subscriptions",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricRouterMatchesLocal, Name: "local", Algo: module.Incremental},
			{ID: metricRouterMatchesRemote, Name: "remote", Algo: module.Incremental},
		},
	}
	chartRouterMemory = Chart{
		ID:    "router_memory",
		Title: "Routing Table Memory Usage",
		Units: "KiB",
		Fam:   "subscriptions",
		Ctx:   "vernemq.router_memory",
		Dims: Dims{
			{ID: metricRouterMemory, Name: "used", Div: 1024},
		},
	}
)

// Erlang VM
var (
	chartSystemUtilization = Chart{
		ID:    "system_utilization",
		Title: "Average Scheduler Utilization",
		Units: "percentage",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_utilization",
		Dims: Dims{
			{ID: metricSystemUtilization, Name: "utilization"},
		},
	}
	chartSystemProcesses = Chart{
		ID:    "system_processes",
		Title: "Erlang Processes",
		Units: "processes",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_processes",
		Dims: Dims{
			{ID: metricSystemProcessCount, Name: "processes"},
		},
	}
	chartSystemReductions = Chart{
		ID:    "system_reductions",
		Title: "Reductions",
		Units: "ops/s",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_reductions",
		Dims: Dims{
			{ID: metricSystemReductions, Name: "reductions", Algo: module.Incremental},
		},
	}
	chartSystemContextSwitches = Chart{
		ID:    "system_context_switches",
		Title: "Context Switches",
		Units: "ops/s",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_context_switches",
		Dims: Dims{
			{ID: metricSystemContextSwitches, Name: "context switches", Algo: module.Incremental},
		},
	}
	chartSystemIO = Chart{
		ID:    "system_io",
		Title: "Received and Sent Traffic through Ports",
		Units: "KiB/s",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_io",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricSystemIOIn, Name: "received", Algo: module.Incremental, Div: 1024},
			{ID: metricSystemIOOut, Name: "sent", Algo: module.Incremental, Div: -1024},
		},
	}
	chartSystemRunQueue = Chart{
		ID:    "system_run_queue",
		Title: "Processes that are Ready to Run on All Run-Queue",
		Units: "processes",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_run_queue",
		Dims: Dims{
			{ID: metricSystemRunQueue, Name: "ready"},
		},
	}
	chartSystemGCCount = Chart{
		ID:    "system_gc_count",
		Title: "GC Count",
		Units: "ops/s",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_gc_count",
		Dims: Dims{
			{ID: metricSystemGCCount, Name: "gc", Algo: module.Incremental},
		},
	}
	chartSystemGCWordsReclaimed = Chart{
		ID:    "system_gc_words_reclaimed",
		Title: "GC Words Reclaimed",
		Units: "ops/s",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_gc_words_reclaimed",
		Dims: Dims{
			{ID: metricSystemWordsReclaimedByGC, Name: "words reclaimed", Algo: module.Incremental},
		},
	}
	chartSystemMemoryAllocated = Chart{
		ID:    "system_allocated_memory",
		Title: "Memory Allocated by the Erlang Processes and by the Emulator",
		Units: "KiB",
		Fam:   "erlang vm",
		Ctx:   "vernemq.system_allocated_memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricVMMemoryProcesses, Name: "processes", Div: 1024},
			{ID: metricVMMemorySystem, Name: "system", Div: 1024},
		},
	}
)

// Bandwidth
var (
	chartBandwidth = Chart{
		ID:    "bandwidth",
		Title: "Bandwidth",
		Units: "KiB/s",
		Fam:   "bandwidth",
		Ctx:   "vernemq.bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricBytesReceived, Name: "received", Algo: module.Incremental, Div: 1024},
			{ID: metricBytesSent, Name: "sent", Algo: module.Incremental, Div: -1024},
		},
	}
)

// Retain
var (
	chartRetainMessages = Chart{
		ID:    "retain_messages",
		Title: "Stored Retained Messages",
		Units: "messages",
		Fam:   "retain",
		Ctx:   "vernemq.retain_messages",
		Dims: Dims{
			{ID: metricRetainMessages, Name: "messages"},
		},
	}
	chartRetainMemoryUsage = Chart{
		ID:    "retain_memory",
		Title: "Stored Retained Messages Memory Usage",
		Units: "KiB",
		Fam:   "retain",
		Ctx:   "vernemq.retain_memory",
		Dims: Dims{
			{ID: metricRetainMemory, Name: "used", Div: 1024},
		},
	}
)

// Cluster
var (
	chartClusterCommunicationBandwidth = Chart{
		ID:    "cluster_bandwidth",
		Title: "Communication with Other Nodes from the Cluster",
		Units: "KiB/s",
		Fam:   "cluster",
		Ctx:   "vernemq.cluster_bandwidth",
		Dims: Dims{
			{ID: metricClusterBytesReceived, Name: "received", Algo: module.Incremental, Div: 1024},
			{ID: metricClusterBytesSent, Name: "sent", Algo: module.Incremental, Div: -1024},
		},
	}
	chartClusterCommunicationDropped = Chart{
		ID:    "cluster_dropped",
		Title: "Dropped Traffic During Communication",
		Units: "KiB/s",
		Fam:   "cluster",
		Ctx:   "vernemq.cluster_dropped",
		Dims: Dims{
			{ID: metricClusterBytesDropped, Name: "dropped", Algo: module.Incremental, Div: 1024},
		},
	}
	chartNetSplitUnresolved = Chart{
		ID:    "netsplit_unresolved",
		Title: "Unresolved Netsplits",
		Units: "netsplits",
		Fam:   "cluster",
		Ctx:   "vernemq.netsplit_unresolved",
		Dims: Dims{
			{ID: "netsplit_unresolved", Name: "unresolved"},
		},
	}
	chartNetSplitEvents = Chart{
		ID:    "netsplit_events",
		Title: "Netsplit Events",
		Units: "netsplits/s",
		Fam:   "cluster",
		Ctx:   "vernemq.netsplits",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: metricNetSplitResolved, Name: "resolved"},
			{ID: metricNetSplitDetected, Name: "detected"},
		},
	}
)

// AUTH
var (
	chartMQTTv5AUTH = Chart{
		ID:    "mqtt_auth",
		Title: "MQTTv5 AUTH",
		Units: "packets/s",
		Fam:   "mqtt auth",
		Ctx:   "vernemq.mqtt_auth",
		Dims: Dims{
			{ID: metricAUTHReceived, Name: "received", Algo: module.Incremental},
			{ID: metricAUTHSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5AUTHReceivedReason = Chart{
		ID:    "mqtt_auth_received_reason",
		Title: "MQTTv5 AUTH Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt auth",
		Ctx:   "vernemq.mqtt_auth_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricAUTHReceived, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricAUTHReceived, "continue_authentication"), Name: "continue", Algo: module.Incremental},
			{ID: join(metricAUTHReceived, "reauthenticate"), Name: "reauthenticate", Algo: module.Incremental},
		},
	}
	chartMQTTv5AUTHSentReason = Chart{
		ID:    "mqtt_auth_sent_reason",
		Title: "MQTTv5 AUTH Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt auth",
		Ctx:   "vernemq.mqtt_auth_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricAUTHSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricAUTHSent, "continue_authentication"), Name: "continue", Algo: module.Incremental},
			{ID: join(metricAUTHSent, "reauthenticate"), Name: "reauthenticate", Algo: module.Incremental},
		},
	}
)

// CONNECT
var (
	chartMQTTv4v5CONNECT = Chart{
		ID:    "mqtt_connect",
		Title: "MQTTv4/v5 CONNECT and CONNACK",
		Units: "packets/s",
		Fam:   "mqtt connect",
		Ctx:   "vernemq.mqtt_connect",
		Dims: Dims{
			{ID: metricCONNECTReceived, Name: "CONNECT", Algo: module.Incremental},
			{ID: metricCONNACKSent, Name: "CONNACK", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv4v5CONNECTSentReason = Chart{
		ID:    "mqtt_connect_sent_reason",
		Title: "MQTTv4/v5 CONNACK Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt connect",
		Ctx:   "vernemq.mqtt_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricCONNACKSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "unsupported_protocol_version"), Name: "unsupported_protocol_version", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "client_identifier_not_valid"), Name: "client_identifier_not_valid", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "server_unavailable"), Name: "server_unavailable", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "bad_username_or_password"), Name: "bad_username_or_password", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},

			{ID: join(metricCONNACKSent, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "malformed_packet"), Name: "malformed_packet", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "protocol_error"), Name: "protocol_error", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "server_busy"), Name: "server_busy", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "banned"), Name: "banned", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "bad_authentication_method"), Name: "bad_authentication_method", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "packet_too_large"), Name: "packet_too_large", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "retain_not_supported"), Name: "retain_not_supported", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "qos_not_supported"), Name: "qos_not_supported", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "use_another_server"), Name: "use_another_server", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "server_moved"), Name: "server_moved", Algo: module.Incremental},
			{ID: join(metricCONNACKSent, "connection_rate_exceeded"), Name: "connection_rate_exceeded", Algo: module.Incremental},
		},
	}
)

// DISCONNECT
var (
	chartMQTTv4v5DISCONNECT = Chart{
		ID:    "mqtt_disconnect",
		Title: "MQTTv4/v5 DISCONNECT",
		Units: "packets/s",
		Fam:   "mqtt disconnect",
		Ctx:   "vernemq.mqtt_disconnect",
		Dims: Dims{
			{ID: metricDISCONNECTReceived, Name: "received", Algo: module.Incremental},
			{ID: metricDISCONNECTSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5DISCONNECTReceivedReason = Chart{
		ID:    "mqtt_disconnect_received_reason",
		Title: "MQTTv5 DISCONNECT Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt disconnect",
		Ctx:   "vernemq.mqtt_disconnect_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricDISCONNECTReceived, "normal_disconnect"), Name: "normal_disconnect", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "disconnect_with_will_msg"), Name: "disconnect_with_will_msg", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "malformed_packet"), Name: "malformed_packet", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "protocol_error"), Name: "protocol_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "receive_max_exceeded"), Name: "receive_max_exceeded", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "topic_alias_invalid"), Name: "topic_alias_invalid", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "packet_too_large"), Name: "packet_too_large", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "message_rate_too_high"), Name: "message_rate_too_high", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "administrative_action"), Name: "administrative_action", Algo: module.Incremental},
			{ID: join(metricDISCONNECTReceived, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
		},
	}
	chartMQTTv5DISCONNECTSentReason = Chart{
		ID:    "mqtt_disconnect_sent_reason",
		Title: "MQTTv5 DISCONNECT Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt disconnect",
		Ctx:   "vernemq.mqtt_disconnect_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricDISCONNECTSent, "normal_disconnect"), Name: "normal_disconnect", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "malformed_packet"), Name: "malformed_packet", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "protocol_error"), Name: "protocol_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "server_busy"), Name: "server_busy", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "server_shutting_down"), Name: "server_shutting_down", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "keep_alive_timeout"), Name: "keep_alive_timeout", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "session_taken_over"), Name: "session_taken_over", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "topic_filter_invalid"), Name: "topic_filter_invalid", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "receive_max_exceeded"), Name: "receive_max_exceeded", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "topic_alias_invalid"), Name: "topic_alias_invalid", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "packet_too_large"), Name: "packet_too_large", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "message_rate_too_high"), Name: "message_rate_too_high", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "administrative_action"), Name: "administrative_action", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "retain_not_supported"), Name: "retain_not_supported", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "qos_not_supported"), Name: "qos_not_supported", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "use_another_server"), Name: "use_another_server", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "server_moved"), Name: "server_moved", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "shared_subs_not_supported"), Name: "shared_subs_not_supported", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "connection_rate_exceeded"), Name: "connection_rate_exceeded", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "max_connect_time"), Name: "max_connect_time", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "subscription_ids_not_supported"), Name: "subscription_ids_not_supported", Algo: module.Incremental},
			{ID: join(metricDISCONNECTSent, "wildcard_subs_not_supported"), Name: "wildcard_subs_not_supported", Algo: module.Incremental},
		},
	}
)

// SUBSCRIBE
var (
	chartMQTTv4v5SUBSCRIBE = Chart{
		ID:    "mqtt_subscribe",
		Title: "MQTTv4/v5 SUBSCRIBE and SUBACK",
		Units: "packets/s",
		Fam:   "mqtt subscribe",
		Ctx:   "vernemq.mqtt_subscribe",
		Dims: Dims{
			{ID: metricSUBSCRIBEReceived, Name: "SUBSCRIBE", Algo: module.Incremental},
			{ID: metricSUBACKSent, Name: "SUBACK", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv4v5SUBSCRIBEError = Chart{
		ID:    "mqtt_subscribe_error",
		Title: "MQTTv4/v5 Failed SUBSCRIBE Operations due to a Netsplit",
		Units: "ops/s",
		Fam:   "mqtt subscribe",
		Ctx:   "vernemq.mqtt_subscribe_error",
		Dims: Dims{
			{ID: metricSUBSCRIBEError, Name: "failed", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5SUBSCRIBEAuthError = Chart{
		ID:    "mqtt_subscribe_auth_error",
		Title: "MQTTv4/v5 Unauthorized SUBSCRIBE Attempts",
		Units: "attempts/s",
		Fam:   "mqtt subscribe",
		Ctx:   "vernemq.mqtt_subscribe_auth_error",
		Dims: Dims{
			{ID: metricSUBSCRIBEAuthError, Name: "unauth", Algo: module.Incremental},
		},
	}
)

// UNSUBSCRIBE
var (
	chartMQTTv4v5UNSUBSCRIBE = Chart{
		ID:    "mqtt_unsubscribe",
		Title: "MQTTv4/v5 UNSUBSCRIBE and UNSUBACK",
		Units: "packets/s",
		Fam:   "mqtt unsubscribe",
		Ctx:   "vernemq.mqtt_unsubscribe",
		Dims: Dims{
			{ID: metricUNSUBSCRIBEReceived, Name: "UNSUBSCRIBE", Algo: module.Incremental},
			{ID: metricUNSUBACKSent, Name: "UNSUBACK", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv4v5UNSUBSCRIBEError = Chart{
		ID:    "mqtt_unsubscribe_error",
		Title: "MQTTv4/v5 Failed UNSUBSCRIBE Operations due to a Netsplit",
		Units: "ops/s",
		Fam:   "mqtt unsubscribe",
		Ctx:   "vernemq.mqtt_unsubscribe_error",
		Dims: Dims{
			{ID: metricUNSUBSCRIBEError, Name: "failed", Algo: module.Incremental},
		},
	}
)

// PUBLISH
var (
	chartMQTTv4v5PUBLISH = Chart{
		ID:    "mqtt_publish",
		Title: "MQTTv4/v5 QOS 0,1,2 PUBLISH",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_publish",
		Dims: Dims{
			{ID: metricPUBSLISHReceived, Name: "received", Algo: module.Incremental},
			{ID: metricPUBSLIHSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv4v5PUBLISHErrors = Chart{
		ID:    "mqtt_publish_errors",
		Title: "MQTTv4/v5 Failed PUBLISH Operations due to a Netsplit",
		Units: "ops/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_publish_errors",
		Dims: Dims{
			{ID: metricPUBLISHError, Name: "failed", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBLISHAuthErrors = Chart{
		ID:    "mqtt_publish_auth_errors",
		Title: "MQTTv4/v5 Unauthorized PUBLISH Attempts",
		Units: "attempts/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_publish_auth_errors",
		Type:  module.Area,
		Dims: Dims{
			{ID: metricPUBLISHAuthError, Name: "unauth", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBACK = Chart{
		ID:    "mqtt_puback",
		Title: "MQTTv4/v5 QOS 1 PUBACK",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_puback",
		Dims: Dims{
			{ID: metricPUBACKReceived, Name: "received", Algo: module.Incremental},
			{ID: metricPUBACKSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5PUBACKReceivedReason = Chart{
		ID:    "mqtt_puback_received_reason",
		Title: "MQTTv5 PUBACK QOS 1 Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_puback_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBACKReceived, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "no_matching_subscribers"), Name: "no_matching_subscribers", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "packet_id_in_use"), Name: "packet_id_in_use", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricPUBACKReceived, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
		},
	}
	chartMQTTv5PUBACKSentReason = Chart{
		ID:    "mqtt_puback_sent_reason",
		Title: "MQTTv5 PUBACK QOS 1 Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_puback_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBACKSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "no_matching_subscribers"), Name: "no_matching_subscribers", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "packet_id_in_use"), Name: "packet_id_in_use", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricPUBACKSent, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBACKUnexpected = Chart{
		ID:    "mqtt_puback_unexpected",
		Title: "MQTTv4/v5 PUBACK QOS 1 Received Unexpected Messages",
		Units: "messages/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_puback_invalid_error",
		Dims: Dims{
			{ID: metricPUBACKInvalid, Name: "unexpected", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBREC = Chart{
		ID:    "mqtt_pubrec",
		Title: "MQTTv4/v5 PUBREC QOS 2",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrec",
		Dims: Dims{
			{ID: metricPUBRECReceived, Name: "received", Algo: module.Incremental},
			{ID: metricPUBRECSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5PUBRECReceivedReason = Chart{
		ID:    "mqtt_pubrec_received_reason",
		Title: "MQTTv5 PUBREC QOS 2 Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrec_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBRECReceived, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "no_matching_subscribers"), Name: "no_matching_subscribers", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "packet_id_in_use"), Name: "packet_id_in_use", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricPUBRECReceived, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
		},
	}
	chartMQTTv5PUBRECSentReason = Chart{
		ID:    "mqtt_pubrec_sent_reason",
		Title: "MQTTv5 PUBREC QOS 2 Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrec_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBRECSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "no_matching_subscribers"), Name: "no_matching_subscribers", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "unspecified_error"), Name: "unspecified_error", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "impl_specific_error"), Name: "impl_specific_error", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "not_authorized"), Name: "not_authorized", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "topic_name_invalid"), Name: "topic_name_invalid", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "packet_id_in_use"), Name: "packet_id_in_use", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "quota_exceeded"), Name: "quota_exceeded", Algo: module.Incremental},
			{ID: join(metricPUBRECSent, "payload_format_invalid"), Name: "payload_format_invalid", Algo: module.Incremental},
		},
	}
	chartMQTTv4PUBRECUnexpected = Chart{
		ID:    "mqtt_pubrec_unexpected",
		Title: "MQTTv4 PUBREC QOS 2 Received Unexpected Messages",
		Units: "messages/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrec_invalid_error",
		Dims: Dims{
			{ID: metricPUBRECInvalid, Name: "unexpected", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBREL = Chart{
		ID:    "mqtt_pubrel",
		Title: "MQTTv4/v5 PUBREL QOS 2",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrel",
		Dims: Dims{
			{ID: metricPUBRELReceived, Name: "received", Algo: module.Incremental},
			{ID: metricPUBRELSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5PUBRELReceivedReason = Chart{
		ID:    "mqtt_pubrel_received_reason",
		Title: "MQTTv5 PUBREL QOS 2 Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrel_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBRELReceived, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBRELReceived, "packet_id_not_found"), Name: "packet_id_not_found", Algo: module.Incremental},
		},
	}
	chartMQTTv5PUBRELSentReason = Chart{
		ID:    "mqtt_pubrel_sent_reason",
		Title: "MQTTv5 PUBREL QOS 2 Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubrel_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBRELSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBRELSent, "packet_id_not_found"), Name: "packet_id_not_found", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBCOMP = Chart{
		ID:    "mqtt_pubcomp",
		Title: "MQTTv4/v5 PUBCOMP QOS 2",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubcom",
		Dims: Dims{
			{ID: metricPUBCOMPReceived, Name: "received", Algo: module.Incremental},
			{ID: metricPUBCOMPSent, Name: "sent", Algo: module.Incremental, Mul: -1},
		},
	}
	chartMQTTv5PUBCOMReceivedReason = Chart{
		ID:    "mqtt_pubcomp_received_reason",
		Title: "MQTTv5 PUBCOMP QOS 2 Received by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubcomp_received_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBCOMPReceived, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBCOMPReceived, "packet_id_not_found"), Name: "packet_id_not_found", Algo: module.Incremental},
		},
	}
	chartMQTTv5PUBCOMSentReason = Chart{
		ID:    "mqtt_pubcomp_sent_reason",
		Title: "MQTTv5 PUBCOMP QOS 2 Sent by Reason",
		Units: "packets/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubcomp_sent_reason",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: join(metricPUBCOMPSent, "success"), Name: "success", Algo: module.Incremental},
			{ID: join(metricPUBCOMPSent, "packet_id_not_found"), Name: "packet_id_not_found", Algo: module.Incremental},
		},
	}
	chartMQTTv4v5PUBCOMPUnexpected = Chart{
		ID:    "mqtt_pubcomp_unexpected",
		Title: "MQTTv4/v5 PUBCOMP QOS 2 Received Unexpected Messages",
		Units: "messages/s",
		Fam:   "mqtt publish",
		Ctx:   "vernemq.mqtt_pubcomp_invalid_error",
		Dims: Dims{
			{ID: metricPUNCOMPInvalid, Name: "unexpected", Algo: module.Incremental},
		},
	}
)

// PING
var (
	chartMQTTv4v5PING = Chart{
		ID:    "mqtt_ping",
		Title: "MQTTv4/v5 PING",
		Units: "packets/s",
		Fam:   "mqtt ping",
		Ctx:   "vernemq.mqtt_ping",
		Dims: Dims{
			{ID: metricPINGREQReceived, Name: "PINGREQ", Algo: module.Incremental},
			{ID: metricPINGRESPSent, Name: "PINGRESP", Algo: module.Incremental, Mul: -1},
		},
	}
)

var (
	chartUptime = Chart{
		ID:    "node_uptime",
		Title: "Node Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "vernemq.node_uptime",
		Dims: Dims{
			{ID: metricSystemWallClock, Name: "time", Div: 1000},
		},
	}
)
