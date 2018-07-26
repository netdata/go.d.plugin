package job

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type variable struct {
	item charts.Var
}

func (w variable) set(value int64) string {
	return fmt.Sprintf(formatVarSET, w.item.ID, value)
}
