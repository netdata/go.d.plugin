package prometheus

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/textparse"

	"github.com/pkg/errors"
)

type (
	// Scraper Scraper
	Scraper struct {
		URL    string
		Client *http.Client

		buf     *bytes.Buffer
		gzipr   *gzip.Reader
		bodybuf *bufio.Reader
	}
)

const (
	acceptHeader    = `text/plain;version=0.0.4;q=1,*/*;q=0.1`
	userAgentHeader = `netdata/go.d.plugin`
)

// Scrape Scrape
func (p *Scraper) Scrape(metrics *Metrics) error {
	if p.buf == nil {
		p.buf = bytes.NewBuffer(make([]byte, 0, 16000))
	}
	p.buf.Reset()
	err := p.fetch(p.buf)
	if err != nil {
		return err
	}
	return p.parse(p.buf.Bytes(), metrics)
}

func (p *Scraper) parse(prometheusText []byte, metrics *Metrics) error {
	parser := textparse.New(prometheusText)

	for {
		entry, err := parser.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if entry != textparse.EntrySeries {
			continue
		}
		_, _, val := parser.Series()
		var lbls labels.Labels
		parser.Metric(&lbls)
		*metrics = append(*metrics, Metric{lbls, val})
	}
	return nil
}

func (p *Scraper) fetch(w io.Writer) error {
	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgentHeader)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
		p.gzipr.Reset(p.bodybuf)
	}
	_, err = io.Copy(w, p.gzipr)
	p.gzipr.Close()
	return err
}
