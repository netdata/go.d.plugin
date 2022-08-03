package wireguard

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioDevicePeers = module.Priority + iota
	prioPeerTraffic
	prioPeerLastHandShake
)

var (
	deviceChartsTmpl = module.Charts{
		devicePeersChartTmpl.Copy(),
	}

	devicePeersChartTmpl = module.Chart{
		ID:       "device_%s_peers",
		Title:    "Device peers",
		Units:    "peers",
		Fam:      "device peers",
		Ctx:      "wireguard.device_peers",
		Priority: prioDevicePeers,
		Dims: module.Dims{
			{ID: "device_%s_peers", Name: "peers"},
		},
	}
)

var (
	peerChartsTmpl = module.Charts{
		peerTrafficChartTmpl.Copy(),
		peerLastHandShakeChartTmpl.Copy(),
	}

	peerTrafficChartTmpl = module.Chart{
		ID:       "peer_%s_traffic",
		Title:    "Peer traffic",
		Units:    "B/s",
		Fam:      "peer traffic",
		Ctx:      "wireguard.peer_traffic",
		Type:     module.Area,
		Priority: prioPeerTraffic,
		Dims: module.Dims{
			{ID: "peer_%s_receive", Name: "receive", Algo: module.Incremental},
			{ID: "peer_%s_transmit", Name: "transmit", Algo: module.Incremental, Mul: -1},
		},
	}
	peerLastHandShakeChartTmpl = module.Chart{
		ID:       "peer_%s_last_handshake_ago",
		Title:    "Peer time elapsed sine the latest handshake",
		Units:    "seconds",
		Fam:      "peer latest handshake",
		Ctx:      "wireguard.peer_last_handshake_ago",
		Priority: prioPeerLastHandShake,
		Dims: module.Dims{
			{ID: "peer_%s_last_handshake_ago", Name: "time"},
		},
	}
)

func newDeviceCharts(device string) *module.Charts {
	charts := deviceChartsTmpl.Copy()

	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, device)
		c.Labels = []module.Label{
			{Key: "device", Value: device},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, device)
		}
	}

	return charts
}

func (w *WireGuard) addNewDeviceCharts(device string) {
	charts := newDeviceCharts(device)

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WireGuard) removeDeviceCharts(device string) {
	prefix := fmt.Sprintf("device_%s", device)

	for _, c := range *w.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

func newPeerCharts(device, peerPublicKey string) *module.Charts {
	charts := peerChartsTmpl.Copy()

	id := peerID(device, peerPublicKey)

	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, id)
		c.Labels = []module.Label{
			{Key: "device", Value: device},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, id)
		}
	}

	return charts
}

func (w *WireGuard) addNewPeerCharts(device, peerPublicKey string) {
	charts := newPeerCharts(device, peerPublicKey)

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WireGuard) removePeerCharts(device, peerPublicKey string) {
	id := peerID(device, peerPublicKey)
	prefix := fmt.Sprintf("peer_%s", id)

	for _, c := range *w.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}
