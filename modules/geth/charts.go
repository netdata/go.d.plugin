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
	chartNewBlocksCreation.Copy(),
	chartTxPoolQueued.Copy(),
	chartP2PNetwork.Copy(),
	chartP2PNetworkDetails.Copy(),
	chartNumberOfPeers.Copy(),
	chartRpcInformation.Copy(),
	chartAncientChainDataRate.Copy(),
	chartChaidataDiskRate.Copy()
}

var (

	chartAncientChainDataRate = Chart{
		ID:    "chaindata_ancient_rate",
		Title: "Ancient Chaindata rate",
		Units: "bytes/s",
		Fam:   "chaindata",
		Ctx:   "geth.eth_db_chaindata_ancient_io_rate",
		Dims: Dims{
			{ID:ethDbChainDataAncientRread , Name: "reads", Algo: "incremental"},
			{ID: ethDbChainDataAncientWrite, Name: "writes", Mul: -1, Algo: "incremental",},
		},
		
	}

	chartAncientChainData = Chart{
		ID:    "chaindata_ancient",
		Title: "Ancient Chaindata",
		Units: "bytes",
		Fam:   "chaindata",
		Ctx:   "geth.eth_db_chaindata_ancient_io",
		Dims: Dims{
			{ID:ethDbChainDataAncientRread , Name: "reads"},
			{ID: ethDbChainDataAncientWrite, Name: "writes", Mul: -1},
		},
		
	}
	chartChaidataDiskRate = Chart{
		ID:    "chaindata_disk_date",
		Title: "On disk Chaindata rate",
		Units: "bytes/s",
		Fam:   "chaindata",
		Ctx:   "geth.eth_db_chaindata_disk_io_rate",
		Dims: Dims{
			{ID:ethDbChaindataDiskRead , Name: "reads", Algo: "incremental"},
			{ID: ethDbChainDataDiskWrite, Name: "writes", Mul: -1, Algo: "incremental"},
		},
	}
	chartChaidataDisk = Chart{
		ID:    "chaindata_disk",
		Title: "Chaindata on disk",
		Units: "bytes",
		Fam:   "chaindata",
		Ctx:   "geth.eth_db_chaindata_disk_io",
		Dims: Dims{
			{ID:ethDbChaindataDiskRead , Name: "reads"},
			{ID: ethDbChainDataDiskWrite, Name: "writes", Mul: -1},
		},
	}
	chartNewBlocksCreation = Chart{
		ID: "chaindata_block_rate", 
		Title: "Block creation rate",
		Units: "blocks/s",
		Fam: "block_creation", 
		Ctx: "geth.block_creation",
		Dims: Dims{
			{ID: chain_head_block, Name: "new blocks", Algo: "incremental"},
		},
	}
	chartTxPoolPending = Chart{
		ID: "txpoolpending",
		Title: "Pending Transaction Pool",
		Units: "transactions", 
		Fam: "tx_pool",
		Ctx: "geth.tx_pool_pending", 
		Dims: Dims{
			{ID: txPoolInvalid, Name: "invalid"},
			{ID: txPoolPending, Name: "pending"},
			{ID: txPoolLocal, Name: "local"},
			{ID: txPoolPendingDiscard, Name: "pending discard"},
			{ID: txPoolNofunds, Name: "no funds" },
			{ID: txPoolPendingRatelimit, Name: "ratelimit"},
			{ID: txPoolPendingReplace, Name: "replace"},
		},
	}
	chartTxPoolCurrent = Chart{
		ID: "txpoolcurrent",
		Title: "Transaction Pool",
		Units: "transactions", 
		Fam: "tx_pool",
		Ctx: "geth.tx_pool_current", 
		Dims: Dims{
			{ID: txPoolInvalid, Name: "invalid"},
			{ID: txPoolPending, Name: "pending"},
			{ID: txPoolLocal, Name: "local"},
			{ID: txPoolNofunds, Name: "pool" },
		},
	}
	chartTxPoolQueued = Chart{
		ID: "txpoolqueued",
		Title: "Queued Transaction Pool",
		Units: "transactions", 
		Fam: "tx_pool",
		Ctx: "geth.tx_pool_queued", 
		Dims: Dims{
			{ID: txPoolQueuedDiscard, Name: "discard"},
			{ID: txPoolQueuedEviction, Name: "eviction"},
			{ID:txPoolQueuedNofunds, Name: "no_funds" },
			{ID: txPoolQueuedRatelimit, Name: "ratelimit"},
		},
	}
	chartP2PNetwork = Chart{
		ID: "p2p_network", 
		Title: "P2P bandwidth",
		Units: "bytes", 
		Fam: "p2p_bandwidth", 
		Ctx: "geth.p2p_bandwidth",
		Dims: Dims{
			{ID: p2pIngress, Name: "ingress"},
			{ID: p2pEgress, Name: "egress", Mul: -1},
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
			{ID: p2pEgressEth650x08, Name: "Eth/65 transactions announcement egress", Mul: -1},
		},
	}
	chartNumberOfPeers = Chart{
		ID: "p2p_general_info",
		Title: "Number of Peers",
		Units: "peers",
		Fam: "p2p_peers",
		Ctx: "geth.p2p_peers",
		Dims: Dims{
			{ID: p2pPeers, Name: "Peers"},
		},
	}
	chartRpcInformation = Chart{
		ID: "rpc_metrics", 
		Title: "RPC information", 
		Units: "rpc calls",
		Fam: "rpc_metrics",
		Ctx: "geth.rpc_metrics",
		Dims: Dims{
			{ID: rpcFailure, Name: "Failed RPC requests", Algo: "percentage-of-absolute-row"},
			{ID: rpcSuccess, Name: "Successful RPC requests", Algo: "percentage-of-absolute-row"},
		},
	}
)