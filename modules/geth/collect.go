package geth

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (v *Geth) collect() (map[string]int64, error) {
	pms, err := v.prom.Scrape()
	if err != nil {
		return nil, err
	}
	mx := v.collectGeth(pms)

	return stm.ToMap(mx), nil
}

func (g *Geth) collectGeth(pms prometheus.Metrics) map[string]float64 {
	mx := make(map[string]float64)
	g.collectChainData(mx, pms)
	g.collectP2P(mx, pms)
	g.collectTxPool(mx, pms)
	g.collectRpc(mx, pms)
	return mx
}

func (v *Geth) collectChainData(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		chainValidation,
		chainWrite,
		ethDbChainDataAncientRead,
		ethDbChainDataAncientWrite,
		ethDbChaindataDiskRead,
		ethDbChainDataDiskWrite,
		chainHeadBlock,
		chainHeadHeader,
		chainHeadReceipt,
		ethDbChainDataAncientSize,
		ethDbChainDataDiskSize,
		reorgsAdd,
		reorgsDropped,
		reorgsExecuted,
		goRoutines,
	)
	v.collectEth(mx, pms)

}

func (v *Geth) collectRpc(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		rpcRequests,
		rpcSuccess,
		rpcFailure,
	)
	v.collectEth(mx, pms)
}

func (v *Geth) collectTxPool(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		txPoolInvalid,
		txPoolPending,
		txPoolLocal,
		txPoolPendingDiscard,
		txPoolNofunds,
		txPoolPendingRatelimit,
		txPoolPendingReplace,
		txPoolQueuedDiscard,
		txPoolQueuedEviction,
		txPoolQueuedEviction,
		txPoolQueuedRatelimit,
	)
	v.collectEth(mx, pms)
}

func (v *Geth) collectP2P(mx map[string]float64, pms prometheus.Metrics) {
	pms = pms.FindByNames(
		p2pDials,
		p2pEgress,
		p2pIngress,
		p2pPeers,
		p2pServes,
	)
	v.collectEth(mx, pms)
}

func (v *Geth) collectEth(mx map[string]float64, pms prometheus.Metrics) {
	for _, pm := range pms {
		mx[pm.Name()] += pm.Value
	}
}
