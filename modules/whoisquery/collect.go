package whoisquery

func (wq *WhoisQuery) collect() (map[string]int64, error) {
	remainingTime, err := wq.prov.remainingTime()
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)
	wq.collectExpiration(mx, remainingTime)
	return mx, nil
}

func (wq WhoisQuery) collectExpiration(mx map[string]int64, remainingTime float64) {
	mx["expiry"] = int64(remainingTime)
	mx["days_until_expiration_warning"] = wq.DaysUntilWarn
	mx["days_until_expiration_critical"] = wq.DaysUntilCrit

}
