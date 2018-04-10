package job

import (
	"time"

	"github.com/l2isbad/go.d.plugin/shared"
)

type (
	timers struct {
		sinceLast  shared.Duration
		spentOnRun shared.Duration
		penalty    time.Duration
		curRun     time.Time
		lastRun    time.Time
	}
)
