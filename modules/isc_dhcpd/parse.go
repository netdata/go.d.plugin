package isc_dhcpd

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

type leaseEntry struct {
	ip    string
	ends string
	bindingState string
}

func parseDHCPdLeases(leases []leaseEntry, r io.Reader) ([]leaseEntry) {
	set := make(map[string]int)

	scanner := bufio.NewScanner(r)
	l := leaseEntry{}
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		switch {
		case l.ip == "" && bytes.HasPrefix(line, []byte("lease")):
			text := string(line)
			// "lease 192.168.0.1 {" => "192.168.0.1"
			l.ip = text[6:len(text)-2]
		case l.ip == "" && bytes.HasPrefix(line, []byte("iaaddr")):
			text := string(line)
			// "iaaddr 1985:470:1f0b:c9a::001 {" =>  "1985:470:1f0b:c9a::001"
			l.ip = text[7:len(text)-2]
		case l.ends == "" && bytes.HasPrefix(line, []byte("ends")):	
			text := string(line)
			if bytes.HasPrefix(line, []byte("ends epoch")) {
				// "ends epoch 1600137140;" => "epoch 1600137140"
				l.ends = text[5:len(text)-1]
			} else if bytes.HasPrefix(line, []byte("ends never")) {
				// "ends never;" => "never"
				l.ends = text[5:len(text)-1]
			} else {
				// "ends 5 2020/09/15 05:49:01;" => "2020/09/15 05:49:01"
				l.ends = text[5:len(text) - 1]
			}
		case l.bindingState == "" && bytes.HasPrefix(line, []byte("binding state")):
			// "binding state active;" => "active"
			text := string(line)
			l.bindingState = text[14:len(text) - 1]
		case bytes.HasPrefix(line, []byte("}")):
			//This test was added to parse IPV6 lease correctly
			if l.ip != "" && l.bindingState != "" && l.ends != "" {
				if idx, ok := set[l.ip]; ok {
					leases[idx] = l
				} else {
					set[l.ip] = len(leases)
					leases = append(leases, l)
				}
			}
			l = leaseEntry{}
		}
	}

	return leases
}

func (d *DHCPd) parseDHCPLease() error {
	info, err := os.Stat(d.Config.LeaseFile)
	if err != nil {
		return errors.New("Cannot get file information")
	}

	if info.ModTime().Unix() == d.LastModification {
		return nil
	}

	d.LastModification = info.ModTime().Unix()

	f, err := os.Open(d.Config.LeaseFile)
	if err != nil {
		d.leases = nil
		return errors.New("Cannot open file")
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	d.leases = parseDHCPdLeases(d.leases, buf)

	return nil
}

func (d *DHCPd) parseLease() {
	if !d.collectedLeases {
		d.collectedLeases = true
		d.addPoolsToCharts()
	}

	err := d.parseDHCPLease()
	if err != nil {
		return
	}
}
