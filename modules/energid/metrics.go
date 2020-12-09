package energid

// For more details, take a look at https://github.com/energicryptocurrency/core-api-documentation

type energidStats struct {
	// https://github.com/energicryptocurrency/core-api-documentation#getblockchaininfo
	BlockChain blockchainStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#getmempoolinfo
	MemPool mempoolStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#getnetworkinfo
	Network networkStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#gettxoutsetinfo
	TXout txoutStatistic
}

type energidResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `error:"method"`
	Id     string      `id:"method"`
}

type energyResponses []energidResponse

type energyRequest struct {
	JSONRPCversion string   `json:"jsonrpc"`
	ID             string   `json:"id"`
	Method         string   `json:"method"`
	Params         []string `json:"params"`
}

type energyRequests []energyRequest

type blockchainStatistic struct {
	Blocks     float64 `stm:"blocks" json:"blocks"`
	Headers    float64 `stm:"headers" json:"headers"`
	Difficulty float64 `stm:"difficulty" json:"difficulty"`
}

type mempoolStatistic struct {
	Max    float64 `stm:"maxmempool" json:"maxmempool"`
	Usage  float64 `stm:"usage" json:"usage"`
	TxSize float64 `stm:"bytes" json:"bytes"`
}

type networkStatistic struct {
	Connections float64 `stm:"connections" json:"connections"`
	TimeOffset  float64 `stm:"timeoffset" json:"timeoffset"`
}

type txoutStatistic struct {
	Count float64 `stm:"transactions" json:"transactions"`
	Xfers float64 `stm:"txouts" json:"txouts"`
}
