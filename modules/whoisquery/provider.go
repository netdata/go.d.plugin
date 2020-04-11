package whoisquery

import (
	"fmt"
	"regexp"
	"time"

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

	expiryRaw := result.Domain.ExpirationDate
	// The result only has year-month-day
	isExpiryDateOnly, _ := regexp.MatchString(`^\d{4}-\d{1,2}-\d{1,2}$`, expiryRaw)
	if isExpiryDateOnly {
		expiryRaw += "T0:00:00Z"
	}
	expiry, _ := time.Parse(time.RFC3339, expiryRaw)
	remainingToExpireSeconds := time.Until(expiry).Seconds()
	return remainingToExpireSeconds, nil
}
