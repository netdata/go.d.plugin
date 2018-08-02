package main

import (
	"os"
	"path"
	"syscall"

	"os/signal"

	"flag"

	"github.com/l2isbad/go.d.plugin/internal/godplugin"
	"github.com/l2isbad/go.d.plugin/internal/pkg/cli"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

var log *logger.Logger

func main() {
	opt := parseOptions()

	plugin := createPlugin(opt)
	plugin.Setup()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		for {
			switch <-signalChan {
			case syscall.SIGINT:
				log.Info("SIGINT received. Terminating...")
				plugin.Shutdown()
				return
			}
		}
	}()

	plugin.MainLoop()
}

func createPlugin(opt *cli.Option) *godplugin.Plugin {
	confDir := netdataConfigDir()
	config := godplugin.NewConfig()
	config.Load(confDir)
	plugin := godplugin.NewPlugin()
	plugin.Option = opt
	plugin.Config = config
	plugin.ModuleConfDir = path.Join(confDir, "go.d")
	return plugin
}

func parseOptions() *cli.Option {
	opt, err := cli.Parse(os.Args)
	if err != nil {
		if err != flag.ErrHelp {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if opt.Debug {
		logger.SetLevel(logger.DEBUG)
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
