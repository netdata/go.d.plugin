package whoisquery

import (
	"fmt"
	"regexp"
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

// TODO: do we need this validation at all?
var reValidDomain = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func newProvider(config Config) (provider, error) {
	domain := config.Source
	if valid := reValidDomain.MatchString(domain); !valid {
		return nil, fmt.Errorf("incorrect domain: %s, expected pattern: %s", domain, reValidDomain)
	}

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
