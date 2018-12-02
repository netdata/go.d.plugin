package main

import (
	"flag"
	"os"
	"path"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/godplugin"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/pkg/multipath"
)

var (
	cd, _       = os.Getwd()
	configPaths = multipath.New(
		os.Getenv("NETDATA_USER_CONFIG_DIR"),
		os.Getenv("NETDATA_STOCK_CONFIG_DIR"),
		path.Join(cd, "/../../../../etc/netdata"),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d"),
	)
)

func main() {
	opt := parseCLI()

	if opt.Debug {
		logger.SetSeverity(logger.DEBUG)
	}

	plugin := createPlugin(opt)

	if !plugin.Setup() {
		return
	}

	plugin.Serve()
}

func createPlugin(opt *cli.Option) *godplugin.Plugin {
	plugin := godplugin.New()

	plugin.Option = opt
	plugin.ConfigPath = configPaths
	plugin.Out = os.Stdout

	if plugin.Option.ConfigDir != "" {
		plugin.ConfigPath = multipath.New(plugin.Option.ConfigDir)
	}

	return plugin
}

func parseCLI() *cli.Option {
	opt, err := cli.Parse(os.Args)
	if err != nil {
		if err != flag.ErrHelp {
			os.Exit(1)
		}
		os.Exit(0)
	}

	return opt
}
