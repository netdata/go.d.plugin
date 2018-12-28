package solr

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("solr", creator)
}

var (
	defURL         = "http://127.0.0.1:8983"
	defHTTPTimeout = time.Second

	coreHandlersURI = "group=core&prefix=UPDATE,QUERY&wt=json"
	infoSystemURI   = "/solr/admin/info/system?wt=json"
)

// New creates Solr with default values
func New() *Solr {
	return &Solr{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
	}
}

type InfoSystem struct {
	Lucene struct {
		Version string `json:"solr-spec-version"`
	}
}

type parser interface {
	parse(*http.Response) (map[string]int64, error)
}

// Solr solr module
type Solr struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	reqCoreHandlers *http.Request
	reqInfoSystem   *http.Request
	client          *http.Client

	parser
}

func (s *Solr) doRequest(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

// Cleanup makes cleanup
func (Solr) Cleanup() {}

// Init makes initialization
func (s *Solr) Init() bool {
	if s.URL == "" {
		s.URL = defURL
	}

	var err error

	s.URI = infoSystemURI
	if s.reqInfoSystem, err = web.NewHTTPRequest(s.Request); err != nil {
		s.Errorf("error on creating HTTP request : %s", err)
		return false
	}

	s.URI = coreHandlersURI
	if s.reqCoreHandlers, err = web.NewHTTPRequest(s.Request); err != nil {
		s.Errorf("error on creating HTTP request : %s", err)
		return false
	}

	s.client = web.NewHTTPClient(s.Client)

	s.parser = &v6Parser{parsed: make(map[string]int64)}

	return true
}

// Check makes check
func (Solr) Check() bool {

	return false
}

// Charts creates Charts
func (Solr) Charts() *Charts {
	return nil
}

// Collect collects coresMetrics
func (s *Solr) Collect() map[string]int64 {
	resp, err := s.doRequest(s.reqCoreHandlers)

	if err != nil {
		s.Errorf("error on request to %s : %s", s.reqCoreHandlers.URL, err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		s.Errorf("%s returned HTTP status %d", s.reqCoreHandlers.URL, resp.StatusCode)
		return nil
	}

	metrics, err := s.parse(resp)

	if err != nil {
		s.Error(err)
		return nil
	}

	return metrics
}
