// SPDX-License-Identifier: GPL-3.0-or-later

package wireguard

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioDevicePeers = module.Priority + iota
	prioDeviceTraffic
	prioPeerTraffic
	prioPeerLastHandShake
)

var (
	deviceChartsTmpl = module.Charts{
		devicePeersChartTmpl.Copy(),
		deviceTrafficChartTmpl.Copy(),
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
	deviceTrafficChartTmpl = module.Chart{
		ID:       "device_%s_traffic",
		Title:    "Device traffic",
		Units:    "B/s",
		Fam:      "device traffic",
		Ctx:      "wireguard.device_traffic",
		Type:     module.Area,
		Priority: prioDeviceTraffic,
		Dims: module.Dims{
			{ID: "device_%s_receive", Name: "receive", Algo: module.Incremental},
			{ID: "device_%s_transmit", Name: "transmit", Algo: module.Incremental, Mul: -1},
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

func newPeerCharts(id, device, pubKey string) *module.Charts {
	charts := peerChartsTmpl.Copy()

	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, id)
		c.Labels = []module.Label{
			{Key: "device", Value: device},
			{Key: "public_key", Value: pubKey},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, id)
		}
	}

	return charts
}

func (w *WireGuard) addNewPeerCharts(id, device, pubKey string) {
	charts := newPeerCharts(id, device, pubKey)

	if err := w.Charts().Add(*charts...); err != nil {
		w.Warning(err)
	}
}

func (w *WireGuard) removePeerCharts(id string) {
	prefix := fmt.Sprintf("peer_%s", id)

	for _, c := range *w.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}
