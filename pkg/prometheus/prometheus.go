package prometheus

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/textparse"
)

type (
	// Prometheus is a helper for scrape and parse prometheus format metrics.
	Prometheus interface {
		// Scrape and parse prometheus format metrics
		Scrape() (Metrics, error)
	}

	prometheus struct {
		client  *http.Client
		request web.Request
		metrics Metrics

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
		client:  client,
		request: request,
		buf:     bytes.NewBuffer(make([]byte, 0, 16000)),
	}
}

// Scrape scrapes metrics, parses and sorts
func (p *prometheus) Scrape() (Metrics, error) {
	p.metrics.Reset()
	if err := p.scrape(&p.metrics); err != nil {
		return nil, err
	}
	p.metrics.Sort()
	return p.metrics, nil
}

func (p *prometheus) scrape(metrics *Metrics) error {
	p.buf.Reset()
	err := p.fetch(p.buf)
	if err != nil {
		return err
	}
	return parse(p.buf.Bytes(), metrics)
}

func parse(prometheusText []byte, metrics *Metrics) error {
	var parser = textparse.NewPromParser(prometheusText)

	for {
		et, err := parser.Next()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch et {
		case textparse.EntrySeries:
			var lbs labels.Labels
			_, _, val := parser.Series()
			parser.Metric(&lbs)

			metrics.Add(Metric{lbs, val})
		case textparse.EntryType:
		case textparse.EntryHelp:
		case textparse.EntryComment:
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
