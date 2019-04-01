package coredns

func (cd *CoreDNS) collect() (map[string]int64, error) {
	_, err := cd.prom.Scrape()

	if err != nil {
		return nil, err
	}
	return nil, nil
}
