package main

import (
	"flag"
	"os"
	"path"

	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/godplugin"
	"github.com/netdata/go.d.plugin/logger"
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
	confDir := netdataConfigDir()
	config := godplugin.NewConfig()
	_ = config.Load(confDir)

	plugin := godplugin.New()
	plugin.Option = opt
	plugin.Config = config
	plugin.ModuleConfDir = path.Join(confDir, "go.d")
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

// return netdata conf directory (e.g. /opt/netdata/etc/netdata)
func netdataConfigDir() string {
	dir := os.Getenv("NETDATA_CONFIG_DIR")

	if dir == "" {
		cd, _ := os.Getwd()
		dir = path.Join(cd, "/../../../../etc/netdata")
	}
	return dir
}
