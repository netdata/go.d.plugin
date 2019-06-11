package dnsmasq_dhcp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
)

func (d *DnsmasqDHCP) findDHCPRanges() ([]string, error) {
	cs, err := findConfigurations(d.ConfDir)
	if err != nil {
		d.Warningf("error during configuration dir scanning : %v", err)
	}

	configs := []string{d.ConfPath}
	configs = append(configs, cs...)
	d.Infof("configuration files to read : %v", configs)

	seen := make(map[string]bool)
	var ranges []string

	for _, config := range configs {
		d.Debugf("reading from %s", config)
		rs, err := findDHCPRanges(config)
		if err != nil {
			return nil, err
		}

		for _, r := range rs {
			d.Debugf("found '%s'", r)
			if r = parseDHCPRangeLine(r); r == "" || seen[r] {
				continue
			}

			d.Debugf("adding '%s'", r)
			seen[r] = true
			ranges = append(ranges, r)
		}
	}

	return ranges, nil
}

func findConfigurations(confDir string) ([]string, error) {
	fis, err := ioutil.ReadDir(confDir)
	if err != nil {
		return nil, err
	}

	var configs []string
	for _, fi := range fis {
		if !fi.Mode().IsRegular() || !strings.HasSuffix(fi.Name(), ".conf") {
			continue
		}
		configs = append(configs, path.Join(confDir, fi.Name()))
	}

	return configs, nil
}

func findDHCPRanges(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("'%s' is not a regular file", filePath)
	}

	var ranges []string
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "dhcp-range") {
			continue
		}
		ranges = append(ranges, line)
	}

	return ranges, nil
}

/*
Examples:
  - dhcp-range=192.168.0.50,192.168.0.150,12h
  - dhcp-range=192.168.0.50,192.168.0.150,255.255.255.0,12h
  - dhcp-range=set:red,1.1.1.50,1.1.2.150, 255.255.252.0
  - dhcp-range=192.168.0.0,static
  - dhcp-range=1234::2, 1234::500, 64, 12h
  - dhcp-range=1234::2, 1234::500
  - dhcp-range=1234::2, 1234::500, slaac
  - dhcp-range=1234::, ra-only
  - dhcp-range=1234::, ra-names
  - dhcp-range=1234::, ra-stateless
*/
var reDHCPRange = regexp.MustCompile(`(?:[=,])([0-9a-f.:]+),([0-9a-f.:]+)`)

// parseDHCPRangeLine expects lines that starts with 'dhcp-range='
func parseDHCPRangeLine(s string) (r string) {
	if strings.Contains(s, "ra-stateless") {
		return
	}

	s = strings.ReplaceAll(s, " ", "")

	match := reDHCPRange.FindStringSubmatch(s)
	if match == nil {
		return
	}

	start, end := net.ParseIP(match[1]), net.ParseIP(match[2])
	if start == nil || end == nil {
		return
	}

	return fmt.Sprintf("%s-%s", start, end)
}
