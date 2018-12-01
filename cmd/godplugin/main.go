package main

import (
	"flag"
	"os"
	"path"

	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/godplugin"
	"github.com/netdata/go.d.plugin/logger"
)

var log = logger.New("main", "")

var (
	cd, _       = os.Getwd()
	configPaths = multipath.New(
		os.Getenv("NETDATA_CONFIG_DIR"),
		os.Getenv("NETDATA_USER_CONFIG_DIR"),
		os.Getenv("NETDATA_STOCK_CONFIG_DIR"),
		path.Join(cd, "/../../../../etc/netdata"),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d"),
	)
)

func main() {
	opt := parseCLI()
	plugin := createPlugin(opt)

	if !plugin.Setup() {
		return
	}

	plugin.Serve()
}

func createPlugin(opt *cli.Option) *godplugin.Plugin {
	config := godplugin.NewConfig()
	config.Load(configPaths.MustFind("go.d.conf"))

	plugin := godplugin.New()
	plugin.Option = opt
	plugin.Config = config
	plugin.ConfigPath = configPaths
	plugin.Out = os.Stdout
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

	if opt.Debug {
		logger.SetSeverity(logger.DEBUG)
	}
	return opt
}
