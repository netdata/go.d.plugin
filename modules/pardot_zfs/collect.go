package pardot_zfs

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func (p *PardotZFS) collect() map[string]int64 {
	results := make(map[string]int64)

	for _, v := range p.pools {

		stderr := new(bytes.Buffer)
		stdout := new(bytes.Buffer)
		cmd := exec.Command("/usr/sbin/zfs", "list", v, "-Ho", "fragmentation")
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			results[v] = 0
			continue
		}

		be, err := io.ReadAll(stderr)
		if err != nil {
			results[v] = 0
			continue
		}
		if len(be) > 0 {
			results[v] = 0
			continue
		}

		bs, err := io.ReadAll(stdout)
		if err != nil {
			results[v] = 0
			continue
		}

		s := string(bs)
		s = strings.TrimSuffix(s, "%")

		i, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			results[v] = 0
			continue
		}

		results[v] = i
	}

	return results
}
