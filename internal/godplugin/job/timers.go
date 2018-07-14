package job

import (
	"time"

	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type timers struct {
	sinceLast  utils.Duration
	spentOnRun utils.Duration
	penalty    time.Duration
	curRun     time.Time
	lastRun    time.Time
}
