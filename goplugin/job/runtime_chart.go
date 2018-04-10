package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/charts/cooked"
)

var (
	formatCREATE = "CHART netdata.runtime_%s '' 'Execution Time' ms 'go.d' 'netdata.god_runtime' line 146000 %d\n" +
		"DIMENSION run_time 'run time' absolute\n\n"
	formatUPDATE = "BEGIN netdata.runtime_%s %d\nSET run_time = %d\nEND\n"
)

type runtimeChart struct {
	updated bool
}

func (r *runtimeChart) create(name string, upd int) {
	cooked.SafePrint(fmt.Sprintf(formatCREATE, name, upd))
}

func (r *runtimeChart) update(name string, sinceLast, elapsed int) {
	if !r.updated {
		sinceLast = 0
	}
	r.updated = true
	cooked.SafePrint(fmt.Sprintf(formatUPDATE, name, sinceLast, elapsed))
}
