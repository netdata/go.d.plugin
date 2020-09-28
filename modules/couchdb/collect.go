package couchdb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathOverviewStats = "/_node/_stats"
	urlPathSystemStats   = "/_node/_system"
	urlPathActiveTasks   = "/_active_tasks"
)

func (cdb *CouchDB) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	return collected, nil
}

func (cdb CouchDB) pingCouchDB() error {
	req, _ := web.NewHTTPRequest(cdb.Request)

	var info struct{ Name string }
	return cdb.doOKDecode(req, &info)
}

func (cdb CouchDB) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := cdb.httpClient.Do(req)
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
