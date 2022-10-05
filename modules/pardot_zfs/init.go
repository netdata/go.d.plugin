package pardot_zfs

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func (z *zfsMetric) init() bool {
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	cmd := exec.Command("/usr/sbin/zfs", "list")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return false
	}

	if strings.Contains("no datasets available", stderr.String()) {
		return false
	}

	pools, err := z.getPools(stdout)
	if err != nil {
		return false
	}

	if len(pools) == 0 {
		return false
	}

	z.pools = pools

	return true
}

func (z *zfsMetric) getPools(b *bytes.Buffer) ([]string, error) {

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
		return nil, err
	}

	return pools, nil
}
