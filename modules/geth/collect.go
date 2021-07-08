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
		ethDbChainDataAncientRread,
		ethDbChainDataAncientWrite,
		ethDbChaindataDiskRead,
		ethDbChainDataDiskWrite,
		chainHeadBlock,
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
		p2pIngressEth650x00,
		p2pIngressEth650x00Packets,
		p2pIngressEth650x01,
		p2pIngressEth650x01Packets,
		p2pIngressEth650x03,
		p2pIngressEth650x03Packets,
		p2pIngressEth650x04,
		p2pIngressEth650x04Packets,
		p2pIngressEth650x05,
		p2pIngressEth650x05Packets,
		p2pIngressEth650x06,
		p2pIngressEth650x06Packets,
		p2pIngressEth650x08,
		p2pIngressEth650x08Packets,
		p2pEgressEth650x00,
		p2pEgressEth650x00Packets,
		p2pEgressEth650x01,
		p2pEgressEth650x01Packets,
		p2pEgressEth650x03,
		p2pEgressEth650x03Packets,
		p2pEgressEth650x04,
		p2pEgressEth650x04Packets,
		p2pEgressEth650x05,
		p2pEgressEth650x05Packets,
		p2pEgressEth650x06,
		p2pEgressEth650x06Packets,
		p2pEgressEth650x08,
		p2pEgressEth650x08Packets,
		p2pIngressEth660x00,
		p2pIngressEth660x00Packets,
		p2pIngressEth660x01,
		p2pIngressEth660x01Packets,
		p2pIngressEth660x03,
		p2pIngressEth660x03Packets,
		p2pIngressEth660x04,
		p2pIngressEth660x04Packets,
		p2pIngressEth660x06,
		p2pIngressEth660x06Packets,
		p2pIngressEth660x08,
		p2pIngressEth660x08Packets,
		p2pTrackedEth660x03,
		p2pTrackedEth660x05,
	)
	v.collectEth(mx, pms)
}

func (v *Geth) collectEth(mx map[string]float64, pms prometheus.Metrics) {
	for _, pm := range pms {
		mx[pm.Name()] += pm.Value
	}
}
