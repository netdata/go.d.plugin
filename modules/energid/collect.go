package energid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathStat = "/"
)

func (e *Energid) collect() (map[string]int64, error) {
	ms := e.scrapeEnergid()
	if ms.empty() {
		return nil, nil
	}

	collected := make(map[string]int64)
	e.collectBlockChain(collected, ms)
	e.collectMemPool(collected, ms)
	e.collectNetwork(collected, ms)
	e.collectTXout(collected, ms)

	return collected, nil
}

func (Energid) collectBlockChain(collected map[string]int64, ms *energidStats) {
	if !ms.hasBlockChain() {
		return
	}

	for metric, value := range stm.ToMap(ms.BlockChain) {
		collected["blockchain_"+metric] = int64(value)
	}
}

func (Energid) collectMemPool(collected map[string]int64, ms *energidStats) {
	if !ms.hasMemPool() {
		return
	}

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
	if !ms.hasNetwork() {
		return
	}

	for metric, value := range stm.ToMap(ms.Network) {
		collected["network_"+metric] = int64(value)
	}
}

func (Energid) collectTXout(collected map[string]int64, ms *energidStats) {
	if !ms.hasTXout() {
		return
	}

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
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() { defer wg.Done(); e.scrapeBlockChain(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); e.scrapeMemPool(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); e.scrapeNetwork(ms) }()

	wg.Add(1)
	go func() { defer wg.Done(); e.scrapeTXout(ms) }()

	wg.Wait()
	return ms
}

func (e *Energid) scrapeBlockChain(ms *energidStats) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := e.energidMakeBody("1", "getblockchaininfo")

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	var bc blockchainStatistic
	if err := e.doOKDecode(req, &bc); err != nil {
		e.Warning(err)
		return
	}

	ms.BlockChain = &bc
}

func (e *Energid) scrapeMemPool(ms *energidStats) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := e.energidMakeBody("2", "getmempoolinfo")

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	var mem mempoolStatistic
	if err := e.doOKDecode(req, &mem); err != nil {
		e.Warning(err)
		return
	}

	ms.MemPool = &mem
}

func (e *Energid) scrapeNetwork(ms *energidStats) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := e.energidMakeBody("3", "getnetworkinfo")

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	var ns networkStatistic
	if err := e.doOKDecode(req, &ns); err != nil {
		e.Warning(err)
		return
	}

	ms.Network = &ns
}

func (e *Energid) scrapeTXout(ms *energidStats) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := e.energidMakeBody("4", "gettxoutsetinfo")

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	var txs txoutStatistic
	if err := e.doOKDecode(req, &txs); err != nil {
		e.Warning(err)
		return
	}

	ms.TXout = &txs
}

func (e *Energid) energidMakeBody(id string, method string) *energyBody {
	return &energyBody{
		JSONRPCversion: "1.0",
		ID:             id,
		Method:         method,
		Params:         make([]string, 0),
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
