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
		cmp := bytes.TrimSpace(scanner.Bytes())
		switch {
		case l.IP == nil && bytes.HasPrefix(cmp, []byte("lease")):
			line := string(cmp[:len(cmp)-1])
			l.IP = net.ParseIP(line[6:len(line)-1])
		case l.IP == nil && bytes.HasPrefix(cmp, []byte("iaaddr")):
			line := string(cmp[:len(cmp)-1])
			l.IP = net.ParseIP(line[7:len(line)-1])
		case l.EndString == "" && bytes.HasPrefix(cmp, []byte("ends")):	
			line := string(cmp[:len(cmp)-1])
			if bytes.HasPrefix(cmp, []byte("ends epoch")) {
				l.EndString = line[11:]
				val, _ := strconv.ParseInt(l.EndString, 10, 64)
				l.Ends = time.Unix(val, 0 )
			} else {
				l.EndString = line[5:]
				l.Ends, _ = time.Parse("2006/01/02 15:04:05", line[7:])
			}
		case l.State == "" && bytes.HasPrefix(cmp, []byte("binding state")):
			line := string(cmp[:len(cmp)-1])
			l.State = line[14:]
		case bytes.HasPrefix(cmp, []byte("}")):
			if idx, ok := set[l.IP.String()]; ok {
				list[idx] = l
			} else {
				set[l.IP.String()] = len(list)
				list = append(list, l)
			}
			l = LeaseFile{}
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
