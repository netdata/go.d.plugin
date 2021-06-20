package geth

import "github.com/netdata/go.d.plugin/agent/module"

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
	Dim    = module.Dim
)

var charts = Charts{
	chartAncientChainData.Copy(),
	chartChaidataDisk.Copy(),
	chartTxPoolPending.Copy(),
	chartBlockProcessingTime.Copy(),
	chartTxPoolQueued.Copy(),
	chartP2PNetwork.Copy(),
	chartP2PNetworkDetails.Copy(),
	chartNumberOfPeers.Copy(),
	chartRpcInformation.Copy(),
}

var (

	chartAncientChainData = Chart{
		ID:    "chaindata_ancient",
		Title: "Ancient Chaindata",
		Units: "bytes",
		Fam:   "chaindata",
		Ctx:   "geth.chaindata",
		Dims: Dims{
			{ID:ethDbChainDataAncientRread , Name: "Ancient chaindata reads"},
			{ID: ethDbChainDataAncientWrite, Name: "Ancient chaindata writes", Mul: -1},
		},
		
	}
	chartChaidataDisk = Chart{
		ID:    "chaindata_disk",
		Title: "Chaindata on disk",
		Units: "bytes",
		Fam:   "chaindata",
		Ctx:   "geth.chaindata_disk",
		Dims: Dims{
			{ID:ethDbChaindataDiskRead , Name: "Disk Chaindata reads"},
			{ID: ethDbChainDataDiskWrite, Name: "Disk Chaindata writes", Mul: -1},
		},
	}
	chartBlockProcessingTime = Chart{
		ID: "blockProcessing_time", 
		Title: "Block processing time",
		Units: "seconds",
		Fam: "block_processing", 
		Ctx: "geth.block_processing",
		Dims: Dims{
			{ID: blockProcessing, Name: "Block processing Time"},
		},
	}
	chartTxPoolPending = Chart{
		ID: "txpoolpending",
		Title: "Pending Transaction Pool",
		Units: "transactions", 
		Fam: "tx_pool",
		Ctx: "geth.tx_pool_pending", 
		Dims: Dims{
			{ID: txPoolInvalid, Name: "Invalid transaction pool"},
			{ID: txPoolPending, Name: "Pending transaction pool"},
			{ID: txPoolLocal, Name: "Local transaction pool"},
			{ID: txPoolPendingDiscard, Name: "Pending discard transaction pool"},
			{ID: txPoolNofunds, Name: "Pool of transactions with no funds" },
			{ID: txPoolPendingRatelimit, Name: "Pending transaction pool ratelimit"},
			{ID: txPoolPendingReplace, Name: "Pending transaction pool to replace"},
		},
	}
	chartTxPoolQueued = Chart{
		ID: "txpoolqueued",
		Title: "Queued Transaction Pool",
		Units: "transactions", 
		Fam: "tx_pool",
		Ctx: "geth.tx_pool_queued", 
		Dims: Dims{
			{ID: txPoolQueuedDiscard, Name: "Transaction pool queued for discard"},
			{ID: txPoolQueuedEviction, Name: "Transaction pool queued for eviction"},
			{ID:txPoolQueuedNofunds, Name: "Transaction pool queued with no funds" },
			{ID: txPoolQueuedRatelimit, Name: "Transaction pool queued rate limit"},
		},
	}
	chartP2PNetwork = Chart{
		ID: "p2p_network", 
		Title: "P2P bandwidth",
		Units: "bytes", 
		Fam: "p2p_bandwidth", 
		Ctx: "geth.p2p_bandwidth",
		Dims: Dims{
			{ID: p2pEgress, Name: "P2P Ingress network"},
			{ID: p2pEgress, Name: "P2P Egress network", Mul: -1},
		},
	}
	chartP2PNetworkDetails = Chart{
		ID: "p2p_eth_65", 
		Title: "Eth/65 Network utilization",
		Units: "bytes", 
		Fam: "p2p_eth_65",
		Ctx: "geth.p2p_eth_65",
		Dims: Dims{
			{ID: p2pIngressEth650x00, Name: "Eth/65 handshake ingress"},
			{ID: p2pIngressEth650x01, Name: "Eth/65 new block hash ingress"},
			{ID: p2pIngressEth650x03, Name: "Eth/65 block header request ingress"},
			{ID: p2pIngressEth650x04, Name: "Eth/65 block header response ingress"},
			{ID: p2pIngressEth650x05, Name: "Eth/65 block body request ingress"},
			{ID: p2pIngressEth650x06, Name: "Eth/65 block body response ingress"},
			{ID: p2pIngressEth650x08, Name: "Eth/65 transactions announcement ingress"},
			{ID: p2pEgressEth650x00, Name: "Eth/65 handshake egress", Mul: -1},
			{ID: p2pEgressEth650x01, Name: "Eth/65 new block hash egress", Mul: -1},
			{ID: p2pEgressEth650x03, Name: "Eth/65 block header request egress", Mul: -1},
			{ID: p2pEgressEth650x04, Name: "Eth/65 block header response egress", Mul: -1},
			{ID: p2pEgressEth650x05, Name: "Eth/65 block body request egress", Mul: -1},
			{ID: p2pEgressEth650x06, Name: "Eth/65 block body response egress", Mul: -1},
			{ID: p2pEgressEth650x08, Title: "Eth/65 transactions announcement egress", Mul: -1},
		},
	}
	chartNumberOfPeers = Chart{
		ID: "p2p_general_info",
		Title: "Number of Peers",
		Units: "peers",
		Fam: "p2p_peers",
		CTX: "geth.p2p_peers",
		Dims: Dims{
			{ID: p2pPeers, Name: "Node Peers"},
		},
	}
	chartRpcInformation = Chart{
		ID: "rpc_metrics", 
		Title: "RPC information", 
		Units: "rpc calls",
		Fam: "rpc_metrics",
		CTX: "geth.rpc_metrics",
		Dims: Dims{
			{ID: rpcFailure, Name: "Failed RPC requests", Algo: "percentage-of-absolute-row"},
			{ID: rpcSuccess, Name: "Successful RPC requests", Algo: "percentage-of-absolute-row"},
		},
	}
)