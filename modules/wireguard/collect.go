package wireguard

import (
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func (w *WireGuard) collect() (map[string]int64, error) {
	if w.client == nil {
		client, err := wgctrl.New()
		if err != nil {
			return nil, err
		}
		w.client = client
	}

	ds, err := w.client.Devices()
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)
	now := time.Now()

	for _, d := range ds {
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
				w.addNewPeerCharts(d.Name, pubKey)
			}

			mx["peer_"+id+"_receive"] = p.ReceiveBytes
			mx["peer_"+id+"_transmit"] = p.TransmitBytes
			mx["peer_"+id+"_last_handshake_ago"] = int64(now.Sub(p.LastHandshakeTime).Seconds())
		}
	}

	return mx, nil
}

func peerID(device, peerPublicKey string) string {
	return device + "_" + peerPublicKey
}
