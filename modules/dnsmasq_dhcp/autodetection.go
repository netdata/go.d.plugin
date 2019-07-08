package dnsmasq_dhcp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// configDir represents conf-dir directive
//
// #conf-dir=/etc/dnsmasq.d
//
// # Include all the files in a directory except those ending in .bak
// #conf-dir=/etc/dnsmasq.d,.bak
//
// # Include all files in a directory which end in .conf
// #conf-dir=/etc/dnsmasq.d/,*.conf
type configDir struct {
	path          string
	includeSuffix []string
	excludeSuffix []string
}

func (d *DnsmasqDHCP) findDHCPRanges() ([]string, error) {
	configs := d.findConfigs(d.ConfPath)

	configs = append([]string{d.ConfPath}, configs...)

	configs = unique(configs)

	d.Infof("configuration files to read: %v", configs)

	seen := make(map[string]bool)
	var ranges []string

	for _, config := range configs {
		d.Debugf("reading %s", config)
		rs, err := findDHCPRanges(config)
		if err != nil && !os.IsNotExist(err) {
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

// findConfigs recursively finds and reads configuration files respecting
// conf-file and conf-dir directives.
// findConfigs is tolerant to IO errors and finds as maximum config
// files as possible, therefore if an error occurrs during scanning process,
// it will be just logged with warning severity
func (d *DnsmasqDHCP) findConfigs(confPath string) []string {
	config, err := os.Open(confPath)
	if err != nil {
		d.Warningf("error during configuration file %q reading: %v", confPath, err)
		return nil
	}

	defer config.Close()

	var (
		includeFiles []string
		includeDirs  []configDir
	)

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		line := scanner.Text()

		if path, ok := getConfValue(line, "conf-file"); ok {
			includeFiles = append(includeFiles, path)
			continue
		}

		if path, ok := getConfValue(line, "conf-dir"); ok {
			args := strings.Split(path, ",")

			dir := configDir{
				path: args[0],
			}

			for _, arg := range args[1:] {
				arg := strings.TrimSpace(arg)
				// dnsmasq treats suffixes with asterisk as "to include" and without
				// asterisk as "to exclude"
				if strings.HasPrefix(arg, "*") {
					dir.includeSuffix = append(dir.includeSuffix, arg[1:])
				} else {
					dir.excludeSuffix = append(dir.excludeSuffix, arg)
				}
			}

			includeDirs = append(includeDirs, dir)

			continue
		}
	}

	for _, dir := range includeDirs {
		dirFiles, err := dir.findConfigs()
		if err != nil {
			d.Warningf("error during configuration dir %q scanning: %v", dir.path, err)
		}

		includeFiles = append(includeFiles, dirFiles...)
	}

	for _, file := range includeFiles {
		files := d.findConfigs(file)

		includeFiles = append(includeFiles, files...)
	}

	return includeFiles
}

func getConfValue(line, prefix string) (value string, ok bool) {
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}

	value = strings.TrimPrefix(line, prefix)

	value = strings.TrimSpace(value)

	if !strings.HasPrefix(value, "=") {
		// got some unexpected line, has prefix like conf-file but there is no
		// assign sign
		return "", false
	}

	value = strings.TrimPrefix(value, "=")

	value = strings.TrimSpace(value)

	return value, true
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

func (dir configDir) findConfigs() ([]string, error) {
	fis, err := ioutil.ReadDir(dir.path)
	if err != nil {
		return nil, err
	}

	var configs []string

	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}

		name := fi.Name()
		if !dir.isValidFileName(name) {
			continue
		}

		if !dir.matchFileName(name) {
			continue
		}

		configs = append(configs, filepath.Join(dir.path, fi.Name()))
	}

	return configs, nil
}

func (dir configDir) isValidFileName(name string) bool {
	// We copy the dnsmasq's logic
	//
	// /* ignore emacs backups and dotfiles */
	// if (len == 0 ||
	//     ent->d_name[len - 1] == '~' ||
	//     (ent->d_name[0] == '#' && ent->d_name[len - 1] == '#') ||
	//     ent->d_name[0] == '.')
	//    continue;

	if strings.HasSuffix(name, "~") ||
		(strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#")) ||
		strings.HasPrefix(name, ".") {
		return false
	}

	return true
}

func (dir configDir) matchFileName(name string) bool {
	if len(dir.includeSuffix) > 0 {
		including := false
		for _, suffix := range dir.includeSuffix {
			if strings.HasSuffix(name, suffix) {
				including = true
				break
			}
		}

		if !including {
			return false
		}
	}

	for _, suffix := range dir.excludeSuffix {
		if strings.HasSuffix(name, suffix) {
			return false
		}
	}

	return true
}

func unique(slice []string) []string {
	result := []string{}

	for _, item := range slice {
		if !contains(result, item) {
			result = append(result, item)
		}
	}
	return result
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}
