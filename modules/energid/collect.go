package energid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathStat = "/"
)

func (e *Energid) collect() (map[string]int64, error) {
	ms := e.scrapeEnergid()

	collected := make(map[string]int64)
	e.collectBlockChain(collected, ms)
	e.collectMemPool(collected, ms)
	e.collectNetwork(collected, ms)
	e.collectTXout(collected, ms)

	return collected, nil
}

func (Energid) collectBlockChain(collected map[string]int64, ms *energidStats) {
	for metric, value := range stm.ToMap(ms.BlockChain) {
		collected["blockchain_"+metric] = int64(value)
	}
}

func (Energid) collectMemPool(collected map[string]int64, ms *energidStats) {
	for metric, value := range stm.ToMap(ms.MemPool) {
		switch metric {
		case "maxmempool":
			collected["mempool_max"] = int64(value)
		case "usage":
			collected["mempool_current"] = int64(value)
		case "bytes":
			collected["mempool_txsize"] = int64(value)
		}
	}
}

func (Energid) collectNetwork(collected map[string]int64, ms *energidStats) {
	for metric, value := range stm.ToMap(ms.Network) {
		collected["network_"+metric] = int64(value)
	}
}

func (Energid) collectTXout(collected map[string]int64, ms *energidStats) {
	for metric, value := range stm.ToMap(ms.TXout) {
		switch metric {
		case "transactions":
			collected["utxo_xfers"] = int64(value)
		case "txouts":
			collected["utxo_count"] = int64(value)
		}
	}
}

func (e *Energid) scrapeEnergid() *energidStats {
	ms := &energidStats{}

	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := e.energidMakeBody()

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return nil
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	stats := energyResponses {
		{
			Result: &ms.BlockChain,
		},
		{
			Result: &ms.MemPool,
		},
		{
			Result: &ms.Network,
		},
		{
			Result: &ms.TXout,
		},
	}
	if err := e.doOKDecode(req, &stats); err != nil {
		e.Warning(err)
		return nil
	}

	return ms
}

func (e *Energid) energidMakeBody() *energyBodies {
	return &energyBodies{ 
		{
			JSONRPCversion: "1.0",
			ID:             "1",
			Method:         "getblockchaininfo",
			Params:         make([]string, 0),
		},
		{
			JSONRPCversion: "1.0",
			ID:             "2",
			Method:         "getmempoolinfo",
			Params:         make([]string, 0),
		},
		{
			JSONRPCversion: "1.0",
			ID:             "3",
			Method:         "getnetworkinfo",
			Params:         make([]string, 0),
		},
		{
			JSONRPCversion: "1.0",
			ID:             "4",
			Method:         "gettxoutsetinfo",
			Params:         make([]string, 0),
		},
	}
}

func (d *Energid) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := d.httpClient.Do(req)
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
