package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	apiPath = "/admin/api.php"

	queryVersion                = "version"
	querySummaryRaw             = "summaryRaw"
	queryTopItems               = "topItems"               // AUTH
	queryTopClients             = "topClients"             // AUTH
	queryGetForwardDestinations = "getForwardDestinations" // AUTH
	queryGetQueryTypes          = "getQueryTypes"          // AUTH
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
		Client:      config.Client,
		URL:         config.URL,
		WebPassword: config.WebPassword,
	}
}

// Client represents Pihole client.
type Client struct {
	*http.Client
	URL         string
	WebPassword string
}

// Version does ?version query.
// Returns API version.
func (c *Client) Version() (int, error) {
	u := fmt.Sprintf(
		"%s%s?%s",
		c.URL,
		apiPath,
		queryVersion,
	)

	var v version
	if err := c.getWithDecode(&v, u); err != nil {
		return 0, err
	}

	return v.Version, nil
}

// SummaryRaw does ?summaryRaw query.
// Returns summary statistics.
func (c *Client) SummaryRaw() (*SummaryRaw, error) {
	u := fmt.Sprintf(
		"%s%s?%s",
		c.URL,
		apiPath,
		querySummaryRaw,
	)

	var s SummaryRaw
	if err := c.getWithDecode(&s, u); err != nil {
		return nil, err
	}

	return &s, nil
}

// QueryTypes does ?getQueryTypes query.
// Returns number of queries that the Pi-hole’s DNS server has processed.
func (c *Client) QueryTypes() (*QueryTypes, error) {
	if c.WebPassword == "" {
		return nil, ErrPasswordNotSet
	}

	u := fmt.Sprintf(
		"%s%s?%s&auth=%s",
		c.URL,
		apiPath,
		queryGetQueryTypes,
		c.WebPassword,
	)

	var qt queryTypes
	if err := c.getWithDecode(&qt, u); err != nil {
		return nil, err
	}

	return &qt.Types, nil
}

// ForwardDestinations does ?getForwardDestinations query.
// Returns number of queries that have been forwarded to the targets.
func (c *Client) ForwardDestinations() (*ForwardDestinations, error) {
	if c.WebPassword == "" {
		return nil, ErrPasswordNotSet
	}

	u := fmt.Sprintf(
		"%s%s?%s&auth=%s",
		c.URL,
		apiPath,
		queryGetForwardDestinations,
		c.WebPassword,
	)

	var fd forwardDestinations
	if err := c.getWithDecode(&fd, u); err != nil {
		return nil, err
	}

	return &fd.Destinations, nil
}

// TopItems does ?topClients query.
// Returns top sources.
func (c *Client) TopClients(top int) (*TopClients, error) {
	if c.WebPassword == "" {
		return nil, ErrPasswordNotSet
	}

	u := fmt.Sprintf(
		"%s%s?%s=%d&auth=%s",
		c.URL,
		apiPath,
		queryTopClients,
		top,
		c.WebPassword,
	)

	var tc topClients
	if err := c.getWithDecode(&tc, u); err != nil {
		return nil, err
	}

	return &tc.Clients, nil
}

// TopItems does ?topItems query.
// Returns top domains and top advertisements.
func (c *Client) TopItems(top int) (*TopItems, error) {
	if c.WebPassword == "" {
		return nil, ErrPasswordNotSet
	}

	u := fmt.Sprintf(
		"%s%s?%s=%d&auth=%s",
		c.URL,
		apiPath,
		queryTopItems,
		top,
		c.WebPassword,
	)

	var ti TopItems
	if err := c.getWithDecode(&ti, u); err != nil {
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
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned %d status code", url, resp.StatusCode)
	}

	return resp, nil
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
