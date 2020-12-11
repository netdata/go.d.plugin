package nginxvts

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const vtsURLPath = "/status/format/json"

func (vts *NginxVTS) collect() (map[string]int64, error) {
	ms, err := vts.scapeTotalMetrics()
	if err != nil {
		return nil, nil
	}

	collected := make(map[string]interface{})
	vts.collectMainMetrics(collected, ms)
	vts.collectSharedZonesMetrics(collected, ms)
	vts.collectServerZonesMetrics(collected, ms)

	return stm.ToMap(collected), nil
}

func (vts *NginxVTS) collectMainMetrics(collected map[string]interface{}, ms *vtsMetrics) {
	collected["loadmsec"] = ms.LoadMsec
	collected["nowmsec"] = ms.NowMsec
	collected["connections"] = ms.Connections
}

func (vts *NginxVTS) collectSharedZonesMetrics(collected map[string]interface{}, ms *vtsMetrics) {
	collected["sharedzones"] = ms.SharedZones
}

func (vts *NginxVTS) collectServerZonesMetrics(collected map[string]interface{}, ms *vtsMetrics) {
	if !ms.hasServerZones() {
		return
	}

	// "*" means all servers
	collected["total"] = ms.ServerZones["*"]
}

func (vts *NginxVTS) scapeTotalMetrics() (*vtsMetrics, error) {
	req, _ := web.NewHTTPRequest(vts.Request)
	req.URL.Path = vtsURLPath

	var total vtsMetrics

	if err := vts.doOKDecode(req, &total); err != nil {
		vts.Warning(err)
		return nil, err
	}
	return &total, nil
}

func (vts *NginxVTS) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := vts.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on HTTP request '%s': %v", req.URL, err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("'%s' returned HTTP status code: %d", req.URL, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(in); err != nil {
		return fmt.Errorf("error on decoding response from '%s': %v", req.URL, err)
	}
	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
