package cli

import (
	"os"
	"strconv"

	"github.com/l2isbad/go.d.plugin/internal/pkg/cli/flags"
)

type ParsedCMD struct {
	Debug       bool
	Module      string
	UpdateEvery int
}

// Parse returns parsed command-line flags in ParsedCMD struct
// Available flags:
func Parse() ParsedCMD {
	var (
		d bool
		m string
		u = 1
	)

	f := flags.New()
	f.BoolVar(&d, "debug", "d", false, "debug mode")
	f.StringVar(&m, "module", "m", "all", "module name")
	f.Parse()

	// override update every should be the last elem in os.Args
	last := os.Args[len(os.Args)-1]
	if v, _ := strconv.Atoi(last); v > 0 {
		u = v
	}

	return ParsedCMD{
		Debug:       d,
		Module:      m,
		UpdateEvery: u,
	}
}
