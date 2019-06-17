package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Query string

func (q Query) WithValue(value string) Query {
	if len(value) == 0 {
		return q
	}

	parts := strings.Split(string(q), "=")
	switch len(parts) {
	case 0:
		return q
	case 1:
		return Query(fmt.Sprintf("%s=%s", string(q), value))
	default:
		return Query(fmt.Sprintf("%s&%s=%s", string(q), parts[0], value))
	}
}

const (
	apiPath = "/admin/api.php"

	queryAuth                   Query = "auth"
	QueryVersion                Query = "version"
	QuerySummaryRaw             Query = "summaryRaw"
	QueryTopItems               Query = "topItems"               // AUTH
	QueryTopClients             Query = "topClients"             // AUTH
	QueryGetForwardDestinations Query = "getForwardDestinations" // AUTH
	QueryGetQueryTypes          Query = "getQueryTypes"          // AUTH
)

var ErrPasswordNotSet = errors.New("password not set")

type Configuration struct {
	Client      *http.Client
	URL         string
	WebPassword string
}

func New(config Configuration) *Client {
	if config.Client == nil {
		config.Client = http.DefaultClient
	}

	return &Client{
		HTTPClient:  config.Client,
		URL:         config.URL,
		WebPassword: config.WebPassword,
	}
}

// Client represents Pihole client.
type Client struct {
	HTTPClient  *http.Client
	URL         string
	WebPassword string
}

func (c *Client) Collect(dst interface{}, query Query) error {
	if needAuth(query) && c.WebPassword == "" {
		return ErrPasswordNotSet
	}
	u, err := makeURL(c.URL, query, c.WebPassword)
	if err != nil {
		return err
	}
	return c.getWithDecode(dst, u)
}

// Version does ?version query.
// Returns API version.
func (c *Client) Version() (int, error) {
	var v version
	err := c.Collect(&v, QueryVersion)
	if err != nil {
		return 0, err
	}

	return v.Version, nil
}

// SummaryRaw does ?summaryRaw query.
// Returns summary statistics.
func (c *Client) SummaryRaw() (*SummaryRaw, error) {
	var s SummaryRaw
	err := c.Collect(&s, QuerySummaryRaw)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// QueryTypes does ?getQueryTypes query.
// Returns number of queries that the Pi-hole’s DNS server has processed.
func (c *Client) QueryTypes() (*QueryTypes, error) {
	var qt queryTypes
	err := c.Collect(&qt, QueryGetQueryTypes)
	if err != nil {
		return nil, err
	}

	return &qt.Types, nil
}

// ForwardDestinations does ?getForwardDestinations query.
// Returns number of queries that have been forwarded to the targets.
func (c *Client) ForwardDestinations() (*ForwardDestinations, error) {
	var fd forwardDestinations
	err := c.Collect(&fd, QueryGetForwardDestinations)
	if err != nil {
		return nil, err
	}

	return &fd.Destinations, nil
}

// TopItems does ?topClients query.
// Returns top sources.
func (c *Client) TopClients(top int) (*TopClients, error) {
	var tc topClients
	err := c.Collect(&tc, QueryTopClients)
	if err != nil {
		return nil, err
	}

	return &tc.Clients, nil
}

// TopItems does ?topItems query.
// Returns top domains and top advertisements.
func (c *Client) TopItems(top int) (*TopItems, error) {
	var ti TopItems
	err := c.Collect(&ti, QueryTopItems)
	if err != nil {
		return nil, err
	}

	return &ti, nil
}

func (c *Client) getWithDecode(dst interface{}, url string) error {
	resp, err := c.getOK(url)
	defer closeBody(resp)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error on reading response from %s : %v", url, err)
	}

	// empty array if:
	// - unauthorized query
	// - wrong query
	if isEmptyArray(b) {
		return fmt.Errorf("unauthorized access to %s", url)
	}

	if err = json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("error on parsing response from %s : %v", url, err)
	}

	return nil
}

func (c *Client) getOK(url string) (*http.Response, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned %d status code", url, resp.StatusCode)
	}

	return resp, nil
}

func makeURL(host string, query Query, password string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	u.Path = apiPath
	u.RawQuery = string(query)
	if password != "" {
		u.RawQuery += "&" + string(queryAuth.WithValue(password))
	}

	return u.String(), nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

func isEmptyArray(data []byte) bool {
	empty := "[]"
	return len(data) == len(empty) && string(data) == empty
}

func needAuth(q Query) bool {
	switch q {
	case QueryGetQueryTypes, QueryGetForwardDestinations, QueryTopItems, QueryTopClients:
		return true
	}
	return false
}
