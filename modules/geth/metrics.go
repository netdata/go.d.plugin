package geth

// Source Code Metrics:
//  - https://github.com/vernemq/vernemq/blob/master/apps/vmq_server/src/vmq_metrics.erl
//  - https://github.com/vernemq/vernemq/blob/master/apps/vmq_server/src/vmq_metrics.hrl

// Source Code FSM:
//  - https://github.com/vernemq/vernemq/blob/master/apps/vmq_server/src/vmq_mqtt_fsm.erl
//  - https://github.com/vernemq/vernemq/blob/master/apps/vmq_server/src/vmq_mqtt5_fsm.erl

// MQTT Packet Types:
//  - v4: http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180834
//  - v5: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901019

// Erlang VM:
//  - http://erlang.org/documentation/doc-5.7.1/erts-5.7.1/doc/html/erlang.html

// Not used metrics (https://docs.vernemq.com/monitoring/introduction):
// - "mqtt_connack_accepted_sent"              // v4, not populated,  "mqtt_connack_sent" used instead
// - "mqtt_connack_unacceptable_protocol_sent" // v4, not populated,  "mqtt_connack_sent" used instead
// - "mqtt_connack_identifier_rejected_sent"   // v4, not populated,  "mqtt_connack_sent" used instead
// - "mqtt_connack_server_unavailable_sent"    // v4, not populated,  "mqtt_connack_sent" used instead
// - "mqtt_connack_bad_credentials_sent"       // v4, not populated,  "mqtt_connack_sent" used instead
// - "mqtt_connack_not_authorized_sent"        // v4, not populated,  "mqtt_connack_sent" used instead
// - "system_exact_reductions"
// - "system_runtime"
// - "vm_memory_atom"
// - "vm_memory_atom_used"
// - "vm_memory_binary"
// - "vm_memory_code"
// - "vm_memory_ets"
// - "vm_memory_processes_used"
// - "vm_memory_total"

// -----------------------------------------------MQTT------------------------------------------------------------------
const (
	// AUTH
	
)

// summary
const (
	// Sockets
	chainValidation = "chain_validation"
	chainWrite = "chain_write"
	chainHeadBlock = "chain_head_block"
)

// + rate
const (
	ethDbChainDataAncientRread = "eth_db_chaindata_ancient_read"
	ethDbChainDataAncientWrite = "eth_db_chaindata_ancient_write"
	ethDbChaindataDiskRead = "eth_db_chaindata_disk_read"
	ethDbChainDataDiskWrite = "eth_db_chaindata_disk_write"
	txPoolInvalid = "txpool_invalid"
	txPoolPending="txpool_pending"
	txPoolLocal="txpool_local"
	txPoolPendingDiscard="txpool_pending_discard"
	txPoolNofunds="txpool_pending_nofunds"
	txPoolPendingRatelimit="txpool_pending_ratelimit"
	txPoolPendingReplace="txpool_pending_replace"
	txPoolQueuedDiscard="txpool_queued_discard"
	txPoolQueuedEviction="txpool_queued_eviction"
	txPoolQueuedNofunds="txpool_queued_nofunds"
	txPoolQueuedRatelimit="txpool_queued_ratelimit"



)

const (
// gauge
p2pDials="p2p_dials"
p2pEgress="p2p_egress"
p2pIngress="p2p_ingress"
p2pIngressEth650x00="p2p_ingress_eth_65_0x00"
p2pIngressEth650x00Packets="p2p_ingress_eth_65_0x00_packets"
p2pIngressEth650x01="p2p_ingress_eth_65_0x01"
p2pIngressEth650x01Packets="p2p_ingress_eth_65_0x01_packets"
p2pIngressEth650x03="p2p_ingress_eth_65_0x03"
p2pIngressEth650x03Packets="p2p_ingress_eth_65_0x03_packets"
p2pIngressEth650x04="p2p_ingress_eth_65_0x04"
p2pIngressEth650x04Packets="p2p_ingress_eth_65_0x04_packets"
p2pIngressEth650x05="p2p_ingress_eth_65_0x05"
p2pIngressEth650x05Packets="p2p_ingress_eth_65_0x05_packets"
p2pIngressEth650x06="p2p_ingress_eth_65_0x06"
p2pIngressEth650x06Packets="p2p_ingress_eth_65_0x06_packets"
p2pIngressEth650x08="p2p_ingress_eth_65_0x08"
p2pIngressEth650x08Packets="p2p_ingress_eth_65_0x08_packets"

p2pEgressEth650x00="p2p_egress_eth_65_0x00"
p2pEgressEth650x00Packets="p2p_egress_eth_65_0x00_packets"
p2pEgressEth650x01="p2p_egress_eth_65_0x01"
p2pEgressEth650x01Packets="p2p_egress_eth_65_0x01_packets"
p2pEgressEth650x03="p2p_egress_eth_65_0x03"
p2pEgressEth650x03Packets="p2p_egress_eth_65_0x03_packets"
p2pEgressEth650x04="p2p_egress_eth_65_0x04"
p2pEgressEth650x04Packets="p2p_egress_eth_65_0x04_packets"
p2pEgressEth650x05="p2p_egress_eth_65_0x05"
p2pEgressEth650x05Packets="p2p_egress_eth_65_0x05_packets"
p2pEgressEth650x06="p2p_egress_eth_65_0x06"
p2pEgressEth650x06Packets="p2p_egress_eth_65_0x06_packets"
p2pEgressEth650x08="p2p_egress_eth_65_0x08"
p2pEgressEth650x08Packets="p2p_egress_eth_65_0x08_packets"

p2pIngressEth660x00="p2p_ingress_eth_66_0x00"
p2pIngressEth660x00Packets="p2p_ingress_eth_66_0x00_packets"
p2pIngressEth660x03="p2p_ingress_eth_66_0x03"
p2pIngressEth660x03Packets="p2p_ingress_eth_66_0x03_packets"
p2pIngressEth660x04="p2p_ingress_eth_66_0x04"
p2pIngressEth660x04Packets="p2p_ingress_eth_66_0x04_packets"
p2pIngressEth660x06="p2p_ingress_eth_66_0x06"
p2pIngressEth660x06Packets="p2p_ingress_eth_66_0x06_packets"
p2pIngressEth660x08="p2p_ingress_eth_66_0x08"
p2pIngressEth660x08Packets="p2p_ingress_eth_66_0x08_packets"
p2pIngressEth660x01="p2p_ingress_snap_1_0x01"
p2pIngressEth660x01Packets="p2p_ingress_snap_1_0x01_packets"

p2pPeers="p2p_peers"
p2pServes="p2p_serves"
p2pTrackedEth660x03="p2p_tracked_eth_66_0x03"
p2pTrackedEth660x05="p2p_tracked_eth_66_0x05"

rpcRequests = "rpc_requests"
rpcSuccess = "rpc_success"
rpcFailure = "rpcFailure"
)