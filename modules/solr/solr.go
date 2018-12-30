package solr

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
)

var (
	minSupportedVersion = 6.4
	coresHandlersURI    = "/solr/admin/metrics?group=core&prefix=UPDATE,QUERY&wt=json"
	infoSystemURI       = "/solr/admin/info/system?wt=json"
)

// New creates Solr with default values
func New() *Solr {
	return &Solr{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		cores: make(map[string]bool),
	}
}

// Solr solr module
type Solr struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	cores map[string]bool

	reqInfoSystem   *http.Request
	reqCoreHandlers *http.Request
	client          *http.Client

	version float64
	charts  *Charts
}

func (s *Solr) doRequest(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

// Cleanup makes cleanup
func (Solr) Cleanup() {}

// Init makes initialization
func (s *Solr) Init() bool {
	if s.URL == "" {
		s.Error("URL not specified")
		return false
	}

	if err := s.createRequests(); err != nil {
		s.Error(err)
		return false
	}

	s.client = web.NewHTTPClient(s.Client)

	return true
}

// Check makes check
func (s *Solr) Check() bool {
	if err := s.getVersion(); err != nil {
		s.Error(err)
		return false
	}

	if s.version < minSupportedVersion {
		s.Errorf("unsupported Solr version : %.1f", s.version)
		return false
	}

	return true
}

// Charts creates Charts
func (s *Solr) Charts() *Charts {
	s.charts = &Charts{}

	return s.charts
}

// Collect collects metrics
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
		s.Errorf("error on parse response from %s : %s", s.reqCoreHandlers.URL, err)
		return nil
	}

	return metrics
}

func (s *Solr) createRequests() error {
	var err error

	s.URI = infoSystemURI
	if s.reqInfoSystem, err = web.NewHTTPRequest(s.Request); err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}
	s.URI = coresHandlersURI
	if s.reqCoreHandlers, err = web.NewHTTPRequest(s.Request); err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	return nil
}

func (s *Solr) getVersion() error {
	resp, err := s.doRequest(s.reqInfoSystem)

	if err != nil {
		return fmt.Errorf("error on request to %s : %s", s.reqInfoSystem.URL, err)
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned HTTP status %d", s.reqInfoSystem.URL, resp.StatusCode)
	}

	var info infoSystem

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return fmt.Errorf("error on decoding %s response : %s", s.reqInfoSystem.URL, err)
	}

	var idx int

	if idx = strings.LastIndex(info.Lucene.Version, "."); idx == -1 {
		return fmt.Errorf("error on parsing version '%s': bad format", info.Lucene.Version)
	}

	if s.version, err = strconv.ParseFloat(info.Lucene.Version[:idx], 10); err != nil {
		return fmt.Errorf("error on parsing version '%s' :  %s", info.Lucene.Version, err)
	}

	return nil
}

type infoSystem struct {
	Lucene struct {
		Version string `json:"solr-spec-version"`
	}
}
