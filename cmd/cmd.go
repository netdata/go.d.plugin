package cmd

import (
	"flag"
	"os"
	"strconv"
)

type parsedCmd struct {
	Debug    bool
	UpdEvery int
	ModRun   string
}

// Parse returns parsed command-line flags in parsedCmd struct
// Available flags:
// -debug   bool   default false
// -module  string default "all"
//
// UpdEvery must be the only parameter that is passed if the flags are not used, otherwise the last one
// Examples:
// ./go.d.plugin 5
// ./go.d.plugin -debug -module=example 1
func Parse() *parsedCmd {
	var debug bool
	var module string
	flag.BoolVar(&debug, "debug", false, "true or false")
	flag.StringVar(&module, "module", "all", "module name")
	flag.Parse()
	if v, e := strconv.Atoi(os.Args[len(os.Args)-1]); e == nil && v > 0 {
		return &parsedCmd{debug, v, module}
	}
	return &parsedCmd{debug, 1, module}
}
