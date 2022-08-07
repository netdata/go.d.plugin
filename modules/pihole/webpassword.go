// SPDX-License-Identifier: GPL-3.0-or-later

package pihole

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (p *Pihole) webPassword() string {
	// do no read setupVarsPath is password is net in the configuration file
	if p.Password != "" {
		return p.Password
	}
	if !isLocalHost(p.URL) {
		p.Info("abort web password auto detection, host is not localhost")
		return ""
	}

	p.Infof("starting web password auto detection, reading : %s", p.SetupVarsPath)
	pass, err := findWebPassword(p.SetupVarsPath)
	if err != nil {
		p.Warningf("error during reading '%s' : %v", p.SetupVarsPath, err)
	}

	return pass
}

func findWebPassword(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	s := bufio.NewScanner(f)
	var password string

	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "WEBPASSWORD") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return "", fmt.Errorf("unparsable line : %s", line)
		}

		password = parts[1]
		break
	}

	return password, nil
}

func isLocalHost(u string) bool {
	if strings.Contains(u, "127.0.0.1") {
		return true
	}
	if strings.Contains(u, "localhost") {
		return true
	}

	return false
}
