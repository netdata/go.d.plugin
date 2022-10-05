package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
)

func (c *Consul) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	if err := c.collectLocalChecks(mx); err != nil {
		return nil, err
	}

	return mx, nil
}

func (c *Consul) doOKDecode(urlPath string, in interface{}) error {
	req, err := web.NewHTTPRequest(c.Request.Copy())
	if err != nil {
		return fmt.Errorf("error on creating request: %v", err)
	}

	req.URL.Path = urlPath
	if c.ACLToken != "" {
		req.Header.Set("X-Consul-Token", c.ACLToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on request to %s : %v", req.URL, err)
	}

	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&in); err != nil {
		return fmt.Errorf("error on decoding response from %s : %v", req.URL, err)
	}

	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
