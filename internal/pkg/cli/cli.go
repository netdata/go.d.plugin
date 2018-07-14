package cli

import (
	"os"
	"strconv"

	"github.com/l2isbad/go.d.plugin/internal/pkg/cli/flags"
)

type ParsedCLI struct {
	Debug       bool
	Module      string
	UpdateEvery int
}

// Parse returns parsed command-line flags in ParsedCLI struct
// Available flags:
func Parse() ParsedCLI {
	var (
		d bool
		m string
		u = 1
	)

	f := flags.New()
	f.BoolVar(&d, "debug", "d", false, "true or false")
	f.StringVar(&m, "module", "m", "all", "module name")
	f.Parse()

	// override update every should be the last elem in os.Args
	last := os.Args[len(os.Args)-1]
	if v, _ := strconv.Atoi(last); v > 0 {
		u = v
	}

	return ParsedCLI{
		Debug:       d,
		Module:      m,
		UpdateEvery: u,
	}
}
