package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/web"
)

// Client creates new ScaleIO client.
func New(client web.Client, request web.Request) (*Client, error) {
	httpClient, err := web.NewHTTPClient(client)
	if err != nil {
		return nil, err
	}
	if err := request.ParseUserURL(); err != nil {
		return nil, err
	}
	return &Client{
		httpClient: httpClient,
		request:    request,
		token:      newToken(),
	}, nil
}

// Client represents ScaleIO client.
type Client struct {
	httpClient *http.Client
	request    web.Request
	token      *token
}

// IsLoggedIn reports whether the client is logged in.
func (c Client) IsLoggedIn() bool { return c.token.isSet() }

// Login connects to FxFlex Gateway to get the token that is used for later authentication for other requests.
func (c *Client) Login() error {
	if c.IsLoggedIn() {
		_ = c.Logout()
	}

	req := c.request.Copy()
	req.URL.Path = pathLogin

	return c.doWithDecode(c.token, decodeToken, req, false, false)
}

// Logout sends logout request.
func (c *Client) Logout() error {
	if !c.IsLoggedIn() {
		return nil
	}

	req := c.request.Copy()
	req.URL.Path = pathLogout
	req.Password = c.token.get()

	c.token.unset()

	return c.doWithDecode(nil, nil, req, false, false)
}

// APIVersion returns FxFlex Gateway API version.
func (c *Client) APIVersion() (*Version, error) {
	req := c.request.Copy()
	req.URL.Path = pathVersion
	req.Password = c.token.get()

	ver := &Version{}
	err := c.doWithDecode(ver, decodeVersion, req, false, false)

	return ver, err
}

// SelectedStatistics makes the query and decodes response into the passed structure.
func (c *Client) SelectedStatistics(dst interface{}, query string) error {
	req := c.request.Copy()
	req.URL.Path = pathSelectedStatistics
	req.Password = c.token.get()
	req.Method = http.MethodPost
	req.Headers["Content-Type"] = "application/json"
	req.Body = query

	return c.doWithDecode(dst, decodeJson, req, true, true)
}

func (c *Client) doWithDecode(dst interface{}, decode decodeFunc, req web.Request, checkLoggedIn, reAuth bool) error {
	if checkLoggedIn && !c.IsLoggedIn() {
		return errors.New("not logged-in")
	}

	resp, err := c.doOK(req, reAuth)
	defer closeBody(resp)
	if err != nil {
		return err
	}

	if dst == nil || decode == nil {
		return nil
	}

	if err := decode(dst, resp.Body); err != nil {
		return fmt.Errorf("error on parsing response from %s : %v", req.URL, err)
	}

	return nil
}

func (c *Client) doOK(req web.Request, reAuth bool) (*http.Response, error) {
	httpReq, err := web.NewHTTPRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error on creating http request to %s : %v", req.URL, err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if reAuth && resp.StatusCode == http.StatusUnauthorized {
		if err = c.Login(); err != nil {
			return resp, err
		}
		return c.doOK(req, false)
	}

	if err = checkStatusCode(resp); err != nil {
		return resp, fmt.Errorf("%s returned %v", httpReq.URL, err)
	}

	return resp, nil
}

func checkStatusCode(resp *http.Response) error {
	// For all 4xx and 5xx return codes, the body may contain an apiError
	// instance with more specifics about the failure.
	if resp.StatusCode >= 400 {
		e := error(&apiError{})
		if err := decodeJson(e, resp.Body); err != nil {
			e = err
		}
		return fmt.Errorf("HTTP status code %d : %v", resp.StatusCode, e)
	}

	// 200(OK), 201(Created), 202(Accepted), 204 (No Content).
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}

	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

type decodeFunc func(dst interface{}, reader io.Reader) error

var decodeJson decodeFunc = func(dst interface{}, reader io.Reader) error { return json.NewDecoder(reader).Decode(dst) }

var decodeVersion decodeFunc = func(dst interface{}, reader io.Reader) error {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	parts := strings.Split(strings.Trim(string(bs), "\n "), ".")
	if len(parts) != 2 {
		return fmt.Errorf("can't parse")
	}
	ver := dst.(*Version)
	ver.Major, err = strconv.ParseInt(parts[0], 10, 64)
	ver.Minor, err = strconv.ParseInt(parts[1], 10, 64)

	return err
}

var decodeToken decodeFunc = func(dst interface{}, reader io.Reader) error {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	token := dst.(*token)
	token.set(strings.Trim(string(bs), `"`))

	return nil
}
