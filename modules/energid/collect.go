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

var energidRequest = energyRequests{
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

type energidResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `error:"method"`
	Id     string      `id:"method"`
}

type energyResponses []energidResponse

type energyRequest struct {
	JSONRPCversion string   `json:"jsonrpc"`
	ID             string   `json:"id"`
	Method         string   `json:"method"`
	Params         []string `json:"params"`
}

type energyRequests []energyRequest

func (e *Energid) collect() (map[string]int64, error) {
	ms, err := e.scrapeEnergid()
	if err != nil {
		return nil, err
	}

	return stm.ToMap(ms), nil
}

func (e *Energid) scrapeEnergid() (*energidInfo, error) {
	ms := &energidInfo{}

	req, _ := web.NewHTTPRequest(e.Request)
	req.URL.Path = urlPathStat
	req.Method = http.MethodPost
	req.Header.Add("content-type", "application/json")
	eb := energidRequest

	body, err := json.Marshal(eb)
	if err != nil {
		e.Error(err)
		return nil, fmt.Errorf("Cannot marshal JSON %s", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	stats := energyResponses{
		{
			Result: &ms.Blockchain,
		},
		{
			Result: &ms.MemPool,
		},
		{
			Result: &ms.Network,
		},
		{
			Result: &ms.TxOutSet,
		},
	}
	if err := e.doOKDecode(req, &stats); err != nil {
		e.Warning(err)
		return nil, fmt.Errorf("Cannot get response: %s", err)
	}

	return ms, nil
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
