package wireguard

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
	"golang.zx2c4.com/wireguard/wgctrl"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("wireguard", creator)
}

// New creates Wireguard with default values
func New() *Wireguard {
	return &Wireguard{
		metrics: make(map[string]int64),
	}
}

// Config is the Wireguard module configuration file.
type Config struct {
	Interface string `yaml:"interface"`
}

// Wireguard example module
type Wireguard struct {
	module.Base // should be embedded by every module

	Config      `yaml:",inline"`
	UpdateEvery int64 `yaml:"update_every"`
	connection  *wgctrl.Client
	metrics     map[string]int64
}

// Cleanup makes cleanup
func (w *Wireguard) Cleanup() {
	if w.connection != nil {
		if err := w.connection.Close(); err != nil {
			w.Errorf("Error when try to close wg connection: %v", err)
		}
	}
}

// Init makes initialization
func (w *Wireguard) Init() bool {
	connection, err := wgctrl.New()
	if err != nil {
		w.Errorf("Failed to open wgctl(wireguad): %v", err)
		return false
	}
	w.connection = connection

	if w.Interface == "" {
		w.Infof("You do not define a wireguard network interface. It will try to use wg0")
		w.Interface = "wg0"
	}

	_, err = connection.Device(w.Interface)
	if err != nil {
		w.Errorf("failed to get device: %v", err)
		return false
	}
	return true
}

// Check makes check
func (Wireguard) Check() bool {
	return true
}

// Charts creates Charts
func (w *Wireguard) Charts() *Charts {
	wc := charts.Copy()

	for _, chart := range *wc {
		chart.ID = fmt.Sprintf(chart.ID, w.Interface)
		chart.Title = fmt.Sprintf(chart.Title, w.Interface)
		chart.Fam = fmt.Sprintf(chart.Fam, w.Interface)
	}

	device, _ := w.connection.Device(w.Interface)
	for id, peer := range device.Peers {
		wbc := bandwitchChart.Copy()

		for _, chart := range *wbc {
			chart.ID = fmt.Sprintf(chart.ID, id)
			chart.Title = fmt.Sprintf(chart.Title, peer.PublicKey.String())
			chart.Fam = fmt.Sprintf(chart.Fam, id)

			for _, dim := range chart.Dims {
				dim.ID = fmt.Sprintf(dim.ID, id)
			}
		}
		_ = wc.Add(*wbc...)
	}
	return wc
}

// Collect collects metrics
func (w *Wireguard) Collect() map[string]int64 {
	w.metrics["received_total"] = 0
	w.metrics["sent_total"] = 0

	device, _ := w.connection.Device(w.Interface)

	for id, peer := range device.Peers {
		receivedKey := fmt.Sprintf("received_%v", id)
		sentKey := fmt.Sprintf("sent_%v", id)

		w.metrics[receivedKey] = peer.ReceiveBytes * 8 / w.UpdateEvery
		w.metrics[sentKey] = peer.TransmitBytes * 8 / w.UpdateEvery
		w.metrics["received_total"] += peer.ReceiveBytes
		w.metrics["sent_total"] += peer.TransmitBytes
	}
	return w.metrics
}
