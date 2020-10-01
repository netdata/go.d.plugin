package isc_dhcpd

import (
	"errors"
	"net"
	"os"
	"time"
	"io"
	"bufio"
	"bytes"
	"strconv"
)

type LeaseEntry struct {
	ip    net.IP
	ends  time.Time
	endstring string
	state string
}

func ParseDHCPd(r io.Reader) ([]LeaseEntry) {
	var list []LeaseEntry
	set := make(map[string]int)

	scanner := bufio.NewScanner(r)
	l := LeaseEntry{}
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		switch {
		case l.ip == nil && bytes.HasPrefix(line, []byte("lease")):
			str := string(line[:len(line)-1])
			l.ip = net.ParseIP(str[6:len(str)-1])
		case l.ip == nil && bytes.HasPrefix(line, []byte("iaaddr")):
			str := string(line[:len(line)-1])
			l.ip = net.ParseIP(str[7:len(str)-1])
		case l.endstring == "" && bytes.HasPrefix(line, []byte("ends")):	
			str := string(line[:len(line)-1])
			if bytes.HasPrefix(line, []byte("ends epoch")) {
				l.endstring = str[11:]
				val, _ := strconv.ParseInt(l.endstring, 10, 64)
				l.ends = time.Unix(val, 0 )
			} else {
				l.endstring = str[5:]
				l.ends, _ = time.Parse("2006/01/02 15:04:05", str[7:])
			}
		case l.state == "" && bytes.HasPrefix(line, []byte("binding state")):
			str := string(line[:len(line)-1])
			l.state = str[14:]
		case bytes.HasPrefix(line, []byte("}")):
			//This test was added to parse IPV6 lease correctly
			if l.ip != nil {
				if idx, ok := set[l.ip.String()]; ok {
					list[idx] = l
				} else {
					set[l.ip.String()] = len(list)
					list = append(list, l)
				}
				l = LeaseEntry{}
			}
		}
	}

	return list
}

func (d *DHCPd) parseDHCPLease() error {
	info, err := os.Stat(d.Config.LeaseFile)
	if err != nil {
		return errors.New("Cannot get file information")
	}

	if info.ModTime().Unix() == d.Config.LastModification {
		return nil
	}

	d.Config.LastModification = info.ModTime().Unix()

	f, err := os.Open(d.LeaseFile)
	if err != nil {
		d.Config.data = nil
		return errors.New("Cannot open file")
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	l := ParseDHCPd(buf)

	if len(l) > 0 {
		d.Config.data = l
	}

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
