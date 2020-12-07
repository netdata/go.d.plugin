package energid

// For more details, take a look at https://github.com/energicryptocurrency/core-api-documentation

type energidStats struct {
	// https://github.com/energicryptocurrency/core-api-documentation#getblockchaininfo
	BlockChain *blockchainStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#getmempoolinfo
	MemPool *mempoolStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#getnetworkinfo
	Network *networkStatistic

	// https://github.com/energicryptocurrency/core-api-documentation#gettxoutsetinfo
	TxOUT *txoutStatistic
}

func (e energidStats) empty() bool {
	switch {
	case e.hasBlockChain(), e.hasMemPool(), e.hasNetwork(), e.hasTxOUT():
			return false
	}
	return true
}

func (e energidStats) hasBlockChain() bool { return e.BlockChain != nil }
func (e energidStats) hasMemPool() bool   { return e.MemPool != nil }
func (e energidStats) hasNetwork() bool  { return e.Network != nil }
func (e energidStats) hasTxOUT() bool     { return e.TxOUT != nil }

type blockchainStatistic struct {
	Blocks float64 `stm:"blocks" json:"blocks"`
	Headers float64 `stm:"headers" json:"headers"`
	Difficulty float64 `stm:"difficulty" json:"difficulty"`
}

type mempoolStatistic struct {
	Max float64 `stm:"maxmempool" json:"maxmempool"`
	Usage float64 `stm:"usage" json:"usage"`
	TxSize float64 `stm:"bytes" json:"bytes"`
}

type networkStatistic struct {
	Connections float64 `stm:"connections" json:"connections"`
	TimeOffset float64 `stm:"timeoffset" json:"timeoffset"`
}

type txoutStatistic struct {
	Count float64 `stm:"transactions" json:"transactions"`
	Xfers float64 `stm:"txouts" json:"txouts"`
}