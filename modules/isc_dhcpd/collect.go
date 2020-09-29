package isc_dhcpd

func (d *DHCPd) collect() (map[string]int64, error) {
	cm := make(map[string]int64)

	d.parseLease(cm)

	return cm, nil
}
