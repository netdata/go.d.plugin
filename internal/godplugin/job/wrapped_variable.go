package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type wrappedVar struct {
	item charts.Var
}

func (w wrappedVar) set(value int64) string {
	return fmt.Sprintf(formatVarSET, w.item.ID, value)
}
