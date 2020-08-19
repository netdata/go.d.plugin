package prometheus

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/textparse"
)

type (
	// Prometheus is a helper for scrape and parse prometheus format metrics.
	Prometheus interface {
		// Scrape and parse prometheus format metrics
		Scrape() (Metrics, error)
		// Metadata returns last scrape metrics metadata
		Metadata() Metadata
	}

	prometheus struct {
		client   *http.Client
		request  web.Request
		metrics  Metrics
		metadata Metadata
		sr       selector.Selector

		// internal use
		buf     *bytes.Buffer
		gzipr   *gzip.Reader
		bodybuf *bufio.Reader
	}
)

const (
	acceptHeader    = `text/plain;version=0.0.4;q=1,*/*;q=0.1`
	userAgentHeader = `netdata/go.d.plugin`
)

// New creates a Prometheus instance.
func New(client *http.Client, request web.Request) Prometheus {
	return &prometheus{
		client:   client,
		request:  request,
		metadata: make(Metadata),
		buf:      bytes.NewBuffer(make([]byte, 0, 16000)),
	}
}

// New creates a Prometheus instance with the selector.
func NewWithSelector(client *http.Client, request web.Request, sr selector.Selector) Prometheus {
	return &prometheus{
		client:   client,
		request:  request,
		metadata: make(Metadata),
		sr:       sr,
		buf:      bytes.NewBuffer(make([]byte, 0, 16000)),
	}
}

// Scrape scrapes metrics, parses and sorts
func (p *prometheus) Scrape() (Metrics, error) {
	p.metrics.Reset()
	p.metadata.reset()
	if err := p.scrape(&p.metrics, p.metadata); err != nil {
		return nil, err
	}
	p.metrics.Sort()
	return p.metrics, nil
}

func (p prometheus) Metadata() Metadata {
	return p.metadata
}

func (p *prometheus) scrape(metrics *Metrics, meta Metadata) error {
	p.buf.Reset()
	if err := p.fetch(p.buf); err != nil {
		return err
	}
	return p.parse(p.buf.Bytes(), metrics, meta)
}

func (p *prometheus) parse(prometheusText []byte, metrics *Metrics, meta Metadata) error {
	parser := textparse.NewPromParser(prometheusText)
	for {
		entry, err := parser.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch entry {
		case textparse.EntrySeries:
			var lbs labels.Labels
			_, _, val := parser.Series()
			parser.Metric(&lbs)
			if p.sr != nil && !p.sr.Matches(lbs) {
				continue
			}
			metrics.Add(Metric{lbs, val})
		case textparse.EntryType:
			meta.setType(parser.Type())
		case textparse.EntryHelp:
			meta.setHelp(parser.Help())
		}
	}
	return nil
}

func (p *prometheus) fetch(w io.Writer) error {
	req, err := web.NewHTTPRequest(p.request)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgentHeader)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned HTTP status %s", resp.Status)
	}

	if resp.Header.Get("Content-Encoding") != "gzip" {
		_, err = io.Copy(w, resp.Body)
		return err
	}

	if p.gzipr == nil {
		p.bodybuf = bufio.NewReader(resp.Body)
		p.gzipr, err = gzip.NewReader(p.bodybuf)
		if err != nil {
			return err
		}
	} else {
		p.bodybuf.Reset(resp.Body)
		_ = p.gzipr.Reset(p.bodybuf)
	}
	_, err = io.Copy(w, p.gzipr)
	_ = p.gzipr.Close()
	return err
}
