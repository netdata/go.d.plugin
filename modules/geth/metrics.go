package geth

// summary
const (
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