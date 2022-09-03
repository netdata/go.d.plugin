// SPDX-License-Identifier: GPL-3.0-or-later

package whoisquery

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type provider interface {
	remainingTime() (float64, error)
}

type fromNet struct {
	domainAddress string
	client        *whois.Client
}

func newProvider(config Config) (provider, error) {
	domain := config.Source
	client := whois.NewClient()
	client.SetTimeout(config.Timeout.Duration)

	return &fromNet{
		domainAddress: domain,
		client:        client,
	}, nil
}

func (f *fromNet) remainingTime() (float64, error) {
	raw, err := f.client.Whois(f.domainAddress)
	if err != nil {
		return 0, err
	}

	result, err := whoisparser.Parse(raw)
	if err != nil {
		return 0, err
	}

	expire, err := dateparse.ParseAny(result.Domain.ExpirationDate)
	if err != nil {
		return 0, err
	}

	return time.Until(expire).Seconds(), nil
}
