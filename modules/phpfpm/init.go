package phpfpm

import (
	"errors"
	"fmt"
	"os"

	"github.com/netdata/go.d.plugin/pkg/web"
)

func (p Phpfpm) initClient() (client, error) {
	if p.UserURL != "" {
		return p.initHTTPClient()
	}
	if p.Socket != "" {
		return p.initSocketClient()
	}
	return nil, errors.New("neither 'socket' nor 'url' set")
}

func (p Phpfpm) initHTTPClient() (*httpClient, error) {
	if err := p.Request.ParseUserURL(); err != nil {
		return nil, fmt.Errorf("parse URL: %v", err)
	}
	c, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return nil, fmt.Errorf("create HTTP client: %v", err)
	}
	p.Debugf("using HTTP client, URL: %s", p.URL)
	p.Debugf("using timeout: %s", p.Timeout.Duration)
	return newHTTPClient(c, p.Request), nil
}

func (p Phpfpm) initSocketClient() (*socketClient, error) {
	if _, err := os.Stat(p.Socket); err != nil {
		return nil, fmt.Errorf("the socket '%s' does not exist: %v", p.Socket, err)
	}
	p.Debugf("using socket client : %s", p.Socket)
	p.Debugf("using timeout: %s", p.Timeout.Duration)
	return newSocketClient(p.Socket, p.Timeout.Duration), nil
}
