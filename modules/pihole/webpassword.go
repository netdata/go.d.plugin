package pihole

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (p *Pihole) webPassword() string {
	if p.Password != "" {
		return p.Password
	}

	password, err := findWebPassword()
	if err != nil {
		p.Warningf("error on web password finding : %v", err)
	}
	return password
}

func findWebPassword() (string, error) {
	f, err := os.Open(defaultSetupVarsPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

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
