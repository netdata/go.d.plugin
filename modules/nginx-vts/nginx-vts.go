package nginxvts

import (
	"encoding/json"

	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {

	modules.Register("nginx-vts", modules.Creator{
		Create: func() modules.Module { return New() },
	})
}

type (
	// NginxVTS Nginx-vts module
	NginxVTS struct {
		modules.Base // should be embedded by every module
		web.HTTP     `yaml:",inline"`

		httpClient web.Client
	}

	vts struct {
		HostName     string      `json:"hostName"`
		NginxVersion string      `json:"nginxVersion"`
		SharedZones  sharedZones `json:"sharedZones"`
		ServerZones  serverZones `json:"serverZones"`
	}

	sharedZones struct {
		Name     string `json:"name"`
		MaxSize  int    `json:"maxSize" stm:"max_size"`
		UsedSize int    `json:"usedSize" stm:"used_size"`
	}

	serverZones = map[string]serverZone

	serverZone struct {
		RequestCounter int       `json:"requestCounter" stm:"req"`
		InBytes        int       `json:"inBytes" stm:"received"`
		OutBytes       int       `json:"outBytes" stm:"sent"`
		Responses      responses `json:"responses" stm:"resp"`
	}

	responses struct {
		Status1xx int `json:"1xx" stm:"1xx"`
		Status2xx int `json:"2xx" stm:"2xx"`
		Status3xx int `json:"3xx" stm:"3xx"`
		Status4xx int `json:"4xx" stm:"4xx"`
		Status5xx int `json:"5xx" stm:"5xx"`
	}

	data struct {
		SharedZones sharedZones           `stm:"shared"`
		ServerZones map[string]serverZone `stm:"server"`
	}
)

// New creates NginxVTS module with default values
func New() modules.Module {
	return &NginxVTS{}
}

// Init makes initialization
func (n *NginxVTS) Init() bool {
	n.httpClient = n.CreateHTTPClient()
	return true
}

// Check makes check
func (n *NginxVTS) Check() bool {
	req, err := n.CreateHTTPRequest()
	if err != nil {
		n.Warning("create http request error: ", err)
		return false
	}
	resp, err := n.httpClient.Do(req)
	if err != nil {
		n.Info("skip job due to http request error: ", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		n.Info("skip job due to %v status %d", req, resp.StatusCode)
		return false
	}

	var data vts
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		n.Info("skip job due to decode json error: ", err)
		return false
	}

	if data.HostName == "" || data.NginxVersion == "" || data.SharedZones.Name == "" {
		n.Info("skip job due to invalid JSON format")
		return false
	}
	return true
}

// Charts creates Charts
func (n *NginxVTS) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (n *NginxVTS) GatherMetrics() map[string]int64 {
	req, err := n.CreateHTTPRequest()
	if err != nil {
		n.Error("create http request error: ", err)
		return nil
	}
	resp, err := n.httpClient.Do(req)
	if err != nil {
		n.Error("request error: ", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		n.Errorf("%v status %d", req, resp.StatusCode)
		return nil
	}

	var raw vts
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		n.Error("decode json error: ", err)
		return nil
	}

	var metrics data

	metrics.SharedZones = raw.SharedZones
	metrics.ServerZones = raw.ServerZones

	if len(raw.ServerZones) > 0 {

		for name := range raw.ServerZones {
			// TODO: add charts
			_ = name
		}
	}

	return stm.ToMap(metrics)
}
