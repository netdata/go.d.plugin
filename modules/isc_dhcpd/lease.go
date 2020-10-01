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

type LeaseFile struct {
	IP    net.IP
	Ends  time.Time
	EndString string
	State string
}

func ParseDHCPd(r io.Reader) ([]LeaseFile) {
	var list []LeaseFile
	set := make(map[string]int)

	scanner := bufio.NewScanner(r)
	l := LeaseFile{}
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		switch {
		case l.IP == nil && bytes.HasPrefix(line, []byte("lease")):
			str := string(line[:len(line)-1])
			l.IP = net.ParseIP(str[6:len(str)-1])
		case l.IP == nil && bytes.HasPrefix(line, []byte("iaaddr")):
			str := string(line[:len(line)-1])
			l.IP = net.ParseIP(str[7:len(str)-1])
		case l.EndString == "" && bytes.HasPrefix(line, []byte("ends")):	
			str := string(line[:len(line)-1])
			if bytes.HasPrefix(line, []byte("ends epoch")) {
				l.EndString = str[11:]
				val, _ := strconv.ParseInt(l.EndString, 10, 64)
				l.Ends = time.Unix(val, 0 )
			} else {
				l.EndString = str[5:]
				l.Ends, _ = time.Parse("2006/01/02 15:04:05", str[7:])
			}
		case l.State == "" && bytes.HasPrefix(line, []byte("binding state")):
			str := string(line[:len(line)-1])
			l.State = str[14:]
		case bytes.HasPrefix(line, []byte("}")):
			//This test was added to parse IPV6 lease correctly
			if l.IP != nil {
				if idx, ok := set[l.IP.String()]; ok {
					list[idx] = l
				} else {
					set[l.IP.String()] = len(list)
					list = append(list, l)
				}
				l = LeaseFile{}
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
