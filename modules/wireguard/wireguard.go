// SPDX-License-Identifier: GPL-3.0-or-later

package wireguard

import (
	"github.com/netdata/go.d.plugin/agent/module"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func init() {
	module.Register("wireguard", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *WireGuard {
	return &WireGuard{
		charts:  &module.Charts{},
		devices: make(map[string]bool),
		peers:   make(map[string]bool),
	}
}

type WireGuard struct {
	module.Base

	charts *module.Charts

	client *wgctrl.Client

	devices map[string]bool
	peers   map[string]bool
}

func (w *WireGuard) Init() bool {
	return true
}

func (w *WireGuard) Check() bool {
	return len(w.Collect()) > 0
}

func (w *WireGuard) Charts() *module.Charts {
	return w.charts
}

func (w *WireGuard) Collect() map[string]int64 {
	mx, err := w.collect()
	if err != nil {
		w.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (w *WireGuard) Cleanup() {}
