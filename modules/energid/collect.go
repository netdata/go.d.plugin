package energid

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

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
	// Add functions to parse the output

	return collected, nil
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
	go func() { defer wg.Done(); e.scrapeTxOUT(ms) }()

	wg.Wait()
	return ms
}

func (e *Energid) scrapeBlockChain(ms *energidStats) {
	e.Request.Body = "{\"jsonrpc\": \"1.0\", \"id\":\"1\", \"method\": \"getblockchaininfo\", \"params\": [] }"
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Header.Add("content-type", "application/json;")
	req.Method = "POST"

	var bc blockchainStatistic
	if err := e.doOKDecode(req, &bc); err != nil {
		ms.BlockChain = nil
		return
	}

	ms.BlockChain = &bc
}

func (e *Energid) scrapeMemPool(ms *energidStats) {
	e.Request.Body = "{\"jsonrpc\": \"1.0\", \"id\":\"2\", \"method\": \"getmempoolinfo\", \"params\": [] }"
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Header.Add("content-type", "application/json;")

	var m mempoolStatistic
	if err := e.doOKDecode(req, &m); err != nil {
		ms.BlockChain = nil
		return
	}

	ms.MemPool = &m
}

func (e *Energid) scrapeNetwork(ms *energidStats) {
	e.Request.Body = "{\"jsonrpc\": \"1.0\", \"id\":\"3\", \"method\": \"getnetworkinfo\", \"params\": [] }"
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Header.Add("content-type", "application/json;")

	var ns networkStatistic
	if err := e.doOKDecode(req, &ns); err != nil {
		ms.BlockChain = nil
		return
	}

	ms.Network = &ns
}

func (e *Energid) scrapeTxOUT(ms *energidStats) {
	e.Request.Body = "{\"jsonrpc\": \"1.0\", \"id\":\"4\", \"method\": \"gettxoutsetinfo\", \"params\": [] }"
	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Header.Add("content-type", "application/json;")
	// add body

	var txs txoutStatistic
	if err := e.doOKDecode(req, &txs); err != nil {
		ms.BlockChain = nil
		return
	}

	ms.TxOUT = &txs
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
