package powerdns_recursor

// PowerDNS Recursor documentation has no section about statistics objects,
// fortunately authoritative has.
// https://doc.powerdns.com/authoritative/http-api/statistics.html#objects
type (
	statisticMetrics []statisticMetric
	statisticMetric  struct {
		Name  string
		Type  string
		Value interface{}
	}
)
