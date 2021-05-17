package whoisquery

import (
	"fmt"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
)

type provider interface {
	remainingTime() (float64, error)
}

type fromNet struct {
	domainAddress string
}

func newProvider(config Config) (provider, error) {
	sourceDomain := config.Source
	validDomain, _ := regexp.MatchString(`^[a-zA-Z0-9\-]+\.[a-zA-Z0-9]+$`, sourceDomain)
	if !validDomain {
		return nil, fmt.Errorf("incorrect domain pattern: %v", sourceDomain)
	}
	return &fromNet{domainAddress: sourceDomain}, nil
}

func (f fromNet) remainingTime() (float64, error) {
	raw, err := whois.Whois(f.domainAddress)
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
