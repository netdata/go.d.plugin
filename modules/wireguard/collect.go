// SPDX-License-Identifier: GPL-3.0-or-later

package wireguard

import (
	"fmt"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (w *WireGuard) collect() (map[string]int64, error) {
	if w.client == nil {
		client, err := w.newWGClient()
		if err != nil {
			return nil, fmt.Errorf("creating wireguard client: %v", err)
		}
		w.client = client
	}

	devices, err := w.client.Devices()
	if err != nil {
		return nil, fmt.Errorf("retrieving WireGuard devices: %v", err)
	}

	now := time.Now()
	if w.cleanupLastTime.IsZero() {
		w.cleanupLastTime = now
	}

	mx := make(map[string]int64)

	w.collectDevicesPeers(mx, devices, now)

	if now.Sub(w.cleanupLastTime) > w.cleanupEvery {
		w.cleanupLastTime = now
		w.cleanupDevicesPeers(devices)
	}

	return mx, nil
}

func (w *WireGuard) collectDevicesPeers(mx map[string]int64, devices []*wgtypes.Device, now time.Time) {
	for _, d := range devices {
		if !w.devices[d.Name] {
			w.devices[d.Name] = true
			w.addNewDeviceCharts(d.Name)
		}

		mx["device_"+d.Name+"_peers"] = int64(len(d.Peers))

		for _, p := range d.Peers {
			pubKey := p.PublicKey.String()
			id := peerID(d.Name, pubKey)

			if !w.peers[id] {
				w.peers[id] = true
				w.addNewPeerCharts(id, d.Name, pubKey)
			}

			mx["device_"+d.Name+"_receive"] += p.ReceiveBytes
			mx["device_"+d.Name+"_transmit"] += p.TransmitBytes
			mx["peer_"+id+"_receive"] = p.ReceiveBytes
			mx["peer_"+id+"_transmit"] = p.TransmitBytes
			mx["peer_"+id+"_last_handshake_ago"] = int64(now.Sub(p.LastHandshakeTime).Seconds())
		}
	}
}

func (w *WireGuard) cleanupDevicesPeers(devices []*wgtypes.Device) {
	seenDevices, seenPeers := make(map[string]bool), make(map[string]bool)
	for _, d := range devices {
		seenDevices[d.Name] = true
		for _, p := range d.Peers {
			seenPeers[peerID(d.Name, p.PublicKey.String())] = true
		}
	}
	for d := range w.devices {
		if !seenDevices[d] {
			delete(w.devices, d)
			w.removeDeviceCharts(d)
		}
	}
	for p := range w.peers {
		if !seenPeers[p] {
			delete(w.peers, p)
			w.removePeerCharts(p)
		}
	}
}

func peerID(device, peerPublicKey string) string {
	return device + "_" + peerPublicKey
}
