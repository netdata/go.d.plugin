package whoisquery

func (x *WhoisQuery) collect() (map[string]int64, error) {
	remainingTime, err := x.prov.remainingTime()
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)
	x.collectExpiration(mx, remainingTime)
	return mx, nil
}

func (x WhoisQuery) collectExpiration(mx map[string]int64, remainingTime float64) {
	mx["expiry"] = int64(remainingTime)
	mx["days_until_expiration_warning"] = x.DaysUntilWarn
	mx["days_until_expiration_critical"] = x.DaysUntilCrit

}
