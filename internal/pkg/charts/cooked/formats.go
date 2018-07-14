package cooked

import (
	"fmt"
	"os"
)

var (
	formatVarSET        = "VARIABLE CHART '%s' = '%d'\n"
	formatDimSET        = "SET '%s' = '%d'\n"
	formatDimCREATE     = "DIMENSION '%s' '%s' %s %d %d '%s'\n"
	formatChartBEGIN    = "BEGIN %s.%s %d\n"
	formatChartCREATE   = "CHART %s.%s '%s' '%s' '%s' '%s' '%s' %s %d %d go.d %s\n"
	formatChartOBSOLETE = "CHART %s.%s '%s' '%s' '%s' '%s' '%s' %s %d %d go.d %s obsolete\n"
)

// SafePrint prints using fmt.Print and Exit(1) if any write error encountered.
func SafePrint(a ...interface{}) {
	if _, err := fmt.Print(a...); err != nil {
		os.Exit(1)
	}
}
