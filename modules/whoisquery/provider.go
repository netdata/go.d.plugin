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
	domainAddress    string
	timeout   time.Duration
}

func newProvider(config Config) (provider, error) {
	sourceDomain := string(config.Source)

	return &fromNet{domainAddress: sourceDomain}, nil
}

func (f fromNet) remainingTime() (float64, error) {
	
	raw, err := whois.Whois(f.domainAddress)

	result, err := whoisparser.Parse(raw)
	if err == nil {

		expiryRaw := result.Domain.ExpirationDate

		// The result only has year-month-day
		isExpiryDateOnly, _ := regexp.MatchString(`^\d{4}-\d{1,2}-\d{1,2}$`, expiryRaw)
		if isExpiryDateOnly {
			expiryRaw += "T0:00:00Z"
		}
		expiry, _ := time.Parse(time.RFC3339, expiryRaw)
		remainingToExpire := expiry.Sub(time.Now())
		remainingToExpireSeconds := remainingToExpire.Seconds()
		return remainingToExpireSeconds, nil
	}
	return -1, fmt.Errorf("%v", err)
}