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

func (e *Energid) collect() (map[string]int64, error) {
	info, err := e.scrapeEnergidInfo()
	if err != nil {
		return nil, err
	}

	return stm.ToMap(info), nil
}

func (e *Energid) scrapeEnergidInfo() (*energidInfo, error) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.Method = http.MethodPost
	req.Header.Set("Content-Type", "application/json")
	request := []rpcRequest{
		{JSONRPC: "1.0", ID: 1, Method: "getblockchaininfo"},
		{JSONRPC: "1.0", ID: 2, Method: "getmempoolinfo"},
		{JSONRPC: "1.0", ID: 3, Method: "getnetworkinfo"},
		{JSONRPC: "1.0", ID: 4, Method: "gettxoutsetinfo"},
		{JSONRPC: "1.0", ID: 5, Method: "getmemoryinfo"},
	}
	body, _ := json.Marshal(request)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	info := energidInfo{
		Blockchain: &blockchainInfo{},
		MemPool:    &memPoolInfo{},
		Network:    &networkInfo{},
		TxOutSet:   &txOutSetInfo{},
		Memory:     &memoryInfo{},
	}
	response := []rpcResponse{
		{Result: &info.Blockchain},
		{Result: &info.MemPool},
		{Result: &info.Network},
		{Result: &info.TxOutSet},
		{Result: &info.Memory},
	}

	if err := e.doOKDecode(req, &response); err != nil {
		return nil, err
	}

	for i, r := range response {
		if r.Error != nil && i < len(request) {
			e.Warningf("error on '%s' method: %v", request[i].Method, r.Error)
		}
	}

	return &info, nil
}

func (e *Energid) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := e.httpClient.Do(req)
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
