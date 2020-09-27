package isc_dhcpd

func (d *DHCPD) collect() (map[string]int64, error) {
	cm := make(map[string]int64)

	d.parseLease(cm)

	return cm, nil
}
