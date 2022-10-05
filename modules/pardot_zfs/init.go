package pardot_zfs

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/netdata/go.d.plugin/logger"
)

func (p *PardotZFS) init() bool {
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	cmd := exec.Command("/usr/sbin/zfs", "list")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		logger.Infof("'/usr/sbin/zfs list' returned error: %v\n", err)
		return false
	}

	if strings.Contains("no datasets available", stderr.String()) {
		logger.Debugln("no datasets available")
		return false
	}

	pools, err := p.getPools(stdout)
	if err != nil {
		return false
	}

	if len(pools) == 0 {
		return false
	}

	p.pools = pools

	return true
}

func (p *PardotZFS) getPools(b *bytes.Buffer) ([]string, error) {

	var pools []string
	s := bufio.NewScanner(b)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "NAME") {
			continue
		}

		fields := strings.Fields(s.Text())
		pools = append(pools, fields[0])
	}

	if err := s.Err(); err != nil {
		logger.Infof("scanner got error: %v\n", err)
		return nil, err
	}

	logger.Debugf("Got pools: %q\n", pools)
	return pools, nil
}
