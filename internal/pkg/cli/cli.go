package cli

import (
	"strconv"

	"github.com/jessevdk/go-flags"
)

// Option defines command line options.
type Option struct {
	Debug       bool   `short:"d" long:"debug" description:"debug mode"`
	Module      string `short:"m" long:"module" description:"module name" default:"all"`
	UpdateEvery int
}

// Parse returns parsed command-line flags in Option struct
func Parse(args []string) (*Option, error) {
	opt := &Option{UpdateEvery: 1}
	parser := flags.NewParser(opt, flags.Default)
	parser.Name = "go.d.plugin"
	parser.Usage = "[OPTIONS] [update every]"
	rest, err := parser.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	if len(rest) > 1 {
		opt.UpdateEvery, err = strconv.Atoi(rest[1])
		if err != nil {
			return nil, err
		}
	}
	return opt, nil
}
