// SPDX-License-Identifier: GPL-3.0-or-later

package wireguard

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestWireGuard_Init(t *testing.T) {

}

func TestWireGuard_Charts(t *testing.T) {

}

func TestWireGuard_Cleanup(t *testing.T) {

}

func TestWireGuard_Check(t *testing.T) {

}

func TestWireGuard_Collect(t *testing.T) {
	type testCaseStep struct {
		prepareMock func(m *mockClient)
		check       func(t *testing.T, w *WireGuard)
	}
	tests := map[string][]testCaseStep{
		"several devices no peers": {
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(1))
					m.devices = append(m.devices, prepareDevice(2))
				},
				check: func(t *testing.T, w *WireGuard) {
					mx := w.Collect()

					expected := map[string]int64{
						"device_wg1_peers":    0,
						"device_wg1_receive":  0,
						"device_wg1_transmit": 0,
						"device_wg2_peers":    0,
						"device_wg2_receive":  0,
						"device_wg2_transmit": 0,
					}

					copyLatestHandshake(mx, expected)
					assert.Equal(t, expected, mx)
					assert.Equal(t, len(deviceChartsTmpl)*2, len(*w.Charts()))
				},
			},
		},
		"several devices several peers each": {
			{
				prepareMock: func(m *mockClient) {
					d1 := prepareDevice(1)
					d1.Peers = append(d1.Peers, preparePeer("11"))
					d1.Peers = append(d1.Peers, preparePeer("12"))
					m.devices = append(m.devices, d1)

					d2 := prepareDevice(2)
					d2.Peers = append(d2.Peers, preparePeer("21"))
					d2.Peers = append(d2.Peers, preparePeer("22"))
					m.devices = append(m.devices, d2)
				},
				check: func(t *testing.T, w *WireGuard) {
					mx := w.Collect()

					expected := map[string]int64{
						"device_wg1_peers":    2,
						"device_wg1_receive":  0,
						"device_wg1_transmit": 0,
						"device_wg2_peers":    2,
						"device_wg2_receive":  0,
						"device_wg2_transmit": 0,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_latest_handshake_ago": 60,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_receive":              0,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_transmit":             0,
						"peer_wg1_cGVlcjEyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_latest_handshake_ago": 60,
						"peer_wg1_cGVlcjEyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_receive":              0,
						"peer_wg1_cGVlcjEyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_transmit":             0,
						"peer_wg2_cGVlcjIxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_latest_handshake_ago": 60,
						"peer_wg2_cGVlcjIxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_receive":              0,
						"peer_wg2_cGVlcjIxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_transmit":             0,
						"peer_wg2_cGVlcjIyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_latest_handshake_ago": 60,
						"peer_wg2_cGVlcjIyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_receive":              0,
						"peer_wg2_cGVlcjIyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_transmit":             0,
					}

					copyLatestHandshake(mx, expected)
					assert.Equal(t, expected, mx)
					assert.Equal(t, len(deviceChartsTmpl)*2+len(peerChartsTmpl)*4, len(*w.Charts()))
				},
			},
		},
		"device added at runtime": {
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(1))
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*1, len(*w.Charts()))
				},
			},
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(2))
				},
				check: func(t *testing.T, w *WireGuard) {
					mx := w.Collect()

					expected := map[string]int64{
						"device_wg1_peers":    0,
						"device_wg1_receive":  0,
						"device_wg1_transmit": 0,
						"device_wg2_peers":    0,
						"device_wg2_receive":  0,
						"device_wg2_transmit": 0,
					}
					copyLatestHandshake(mx, expected)
					assert.Equal(t, expected, mx)
					assert.Equal(t, len(deviceChartsTmpl)*2, len(*w.Charts()))

				},
			},
		},
		"device removed at run time, no cleanup occurred": {
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(1))
					m.devices = append(m.devices, prepareDevice(2))
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
				},
			},
			{
				prepareMock: func(m *mockClient) {
					m.devices = m.devices[:len(m.devices)-1]
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*2, len(*w.Charts()))
					assert.Equal(t, 0, calcObsoleteCharts(w.Charts()))
				},
			},
		},
		"device removed at run time, cleanup occurred": {
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(1))
					m.devices = append(m.devices, prepareDevice(2))
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
				},
			},
			{
				prepareMock: func(m *mockClient) {
					m.devices = m.devices[:len(m.devices)-1]
				},
				check: func(t *testing.T, w *WireGuard) {
					w.cleanupEvery = time.Second
					time.Sleep(time.Second)
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*2, len(*w.Charts()))
					assert.Equal(t, len(deviceChartsTmpl)*1, calcObsoleteCharts(w.Charts()))
				},
			},
		},
		"peer added at runtime": {
			{
				prepareMock: func(m *mockClient) {
					m.devices = append(m.devices, prepareDevice(1))
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*1, len(*w.Charts()))
				},
			},
			{
				prepareMock: func(m *mockClient) {
					d1 := m.devices[0]
					d1.Peers = append(d1.Peers, preparePeer("11"))
				},
				check: func(t *testing.T, w *WireGuard) {
					mx := w.Collect()

					expected := map[string]int64{
						"device_wg1_peers":    1,
						"device_wg1_receive":  0,
						"device_wg1_transmit": 0,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_latest_handshake_ago": 60,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_receive":              0,
						"peer_wg1_cGVlcjExAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=_transmit":             0,
					}
					copyLatestHandshake(mx, expected)
					assert.Equal(t, expected, mx)
					assert.Equal(t, len(deviceChartsTmpl)*1+len(peerChartsTmpl)*1, len(*w.Charts()))

				},
			},
		},
		"peer removed at run time, no cleanup occurred": {
			{
				prepareMock: func(m *mockClient) {
					d1 := prepareDevice(1)
					d1.Peers = append(d1.Peers, preparePeer("11"))
					d1.Peers = append(d1.Peers, preparePeer("12"))
					m.devices = append(m.devices, d1)
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
				},
			},
			{
				prepareMock: func(m *mockClient) {
					d1 := m.devices[0]
					d1.Peers = d1.Peers[:len(d1.Peers)-1]
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*1+len(peerChartsTmpl)*2, len(*w.Charts()))
					assert.Equal(t, 0, calcObsoleteCharts(w.Charts()))
				},
			},
		},
		"peer removed at run time, cleanup occurred": {
			{
				prepareMock: func(m *mockClient) {
					d1 := prepareDevice(1)
					d1.Peers = append(d1.Peers, preparePeer("11"))
					d1.Peers = append(d1.Peers, preparePeer("12"))
					m.devices = append(m.devices, d1)
				},
				check: func(t *testing.T, w *WireGuard) {
					_ = w.Collect()
				},
			},
			{
				prepareMock: func(m *mockClient) {
					d1 := m.devices[0]
					d1.Peers = d1.Peers[:len(d1.Peers)-1]
				},
				check: func(t *testing.T, w *WireGuard) {
					w.cleanupEvery = time.Second
					time.Sleep(time.Second)
					_ = w.Collect()
					assert.Equal(t, len(deviceChartsTmpl)*1+len(peerChartsTmpl)*2, len(*w.Charts()))
					assert.Equal(t, len(peerChartsTmpl)*1, calcObsoleteCharts(w.Charts()))
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := New()
			require.True(t, w.Init())
			m := &mockClient{}
			w.client = m

			for i, step := range test {
				t.Run(fmt.Sprintf("step[%d]", i), func(t *testing.T) {
					step.prepareMock(m)
					step.check(t, w)
				})
			}
		})
	}
}

type mockClient struct {
	devices      []*wgtypes.Device
	errOnDevices bool
}

func (m *mockClient) Devices() ([]*wgtypes.Device, error) {
	if m.errOnDevices {
		return nil, errors.New("mock.Devices() error")
	}
	return m.devices, nil
}

func (m *mockClient) Close() error {
	return nil
}

func prepareDevice(num uint8) *wgtypes.Device {
	return &wgtypes.Device{
		Name: fmt.Sprintf("wg%d", num),
	}
}

func preparePeer(s string) wgtypes.Peer {
	b := make([]byte, 32)
	b = append(b[:0], fmt.Sprintf("peer%s", s)...)
	k, _ := wgtypes.NewKey(b[:32])

	return wgtypes.Peer{
		PublicKey:         k,
		LastHandshakeTime: time.Now().Add(-time.Minute),
		ReceiveBytes:      0,
		TransmitBytes:     0,
	}
}

func copyLatestHandshake(dst, src map[string]int64) {
	for k, v := range src {
		if strings.HasSuffix(k, "latest_handshake_ago") {
			if _, ok := dst[k]; ok {
				dst[k] = v
			}
		}
	}
}

func calcObsoleteCharts(charts *module.Charts) int {
	var num int
	for _, c := range *charts {
		if c.Obsolete {
			num++
		}
	}
	return num
}
