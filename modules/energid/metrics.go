package energid

// API docs https://github.com/energicryptocurrency/core-api-documentation

type energidInfo struct {
	Blockchain blockchainInfo `stm:"blockchain"`
	MemPool    memPoolInfo    `stm:"mempool"`
	Network    networkInfo    `stm:"network"`
	TxOutSet   txOutSetInfo   `stm:"utxo"`
}

// https://github.com/energicryptocurrency/core-api-documentation#getblockchaininfo
type blockchainInfo struct {
	Blocks     float64 `stm:"blocks" json:"blocks"`
	Headers    float64 `stm:"headers" json:"headers"`
	Difficulty float64 `stm:"difficulty" json:"difficulty"`
}

// https://github.com/energicryptocurrency/core-api-documentation#getmempoolinfo
type memPoolInfo struct {
	Size       float64 `stm:"txcount" json:"size"`
	Bytes      float64 `stm:"txsize" json:"bytes"`
	Usage      float64 `stm:"current" json:"usage"`
	MaxMemPool float64 `stm:"max" json:"maxmempool"`
}

// https://github.com/energicryptocurrency/core-api-documentation#getnetworkinfo
type networkInfo struct {
	TimeOffset  float64 `stm:"timeoffset" json:"timeoffset"`
	Connections float64 `stm:"connections" json:"connections"`
}

// https://github.com/energicryptocurrency/core-api-documentation#gettxoutsetinfo
type txOutSetInfo struct {
	Transactions float64 `stm:"xfers" json:"transactions"`
	TxOuts       float64 `stm:"count" json:"txouts"`
}
