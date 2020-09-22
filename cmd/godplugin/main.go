package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"github.com/netdata/go-orchestrator/cli"
	"github.com/netdata/go-orchestrator/pkg/logger"
	"github.com/netdata/go-orchestrator/pkg/multipath"
	"github.com/netdata/go-orchestrator/plugin"

	_ "github.com/netdata/go.d.plugin/modules/activemq"
	_ "github.com/netdata/go.d.plugin/modules/apache"
	_ "github.com/netdata/go.d.plugin/modules/bind"
	_ "github.com/netdata/go.d.plugin/modules/cockroachdb"
	_ "github.com/netdata/go.d.plugin/modules/consul"
	_ "github.com/netdata/go.d.plugin/modules/coredns"
	_ "github.com/netdata/go.d.plugin/modules/dnsmasq_dhcp"
	_ "github.com/netdata/go.d.plugin/modules/dnsquery"
	_ "github.com/netdata/go.d.plugin/modules/docker_engine"
	_ "github.com/netdata/go.d.plugin/modules/dockerhub"
	_ "github.com/netdata/go.d.plugin/modules/elasticsearch"
	_ "github.com/netdata/go.d.plugin/modules/example"
	_ "github.com/netdata/go.d.plugin/modules/filecheck"
	_ "github.com/netdata/go.d.plugin/modules/fluentd"
	_ "github.com/netdata/go.d.plugin/modules/freeradius"
	_ "github.com/netdata/go.d.plugin/modules/hdfs"
	_ "github.com/netdata/go.d.plugin/modules/httpcheck"
	_ "github.com/netdata/go.d.plugin/modules/k8s_kubelet"
	_ "github.com/netdata/go.d.plugin/modules/k8s_kubeproxy"
	_ "github.com/netdata/go.d.plugin/modules/lighttpd"
	_ "github.com/netdata/go.d.plugin/modules/lighttpd2"
	_ "github.com/netdata/go.d.plugin/modules/logstash"
	_ "github.com/netdata/go.d.plugin/modules/mysql"
	_ "github.com/netdata/go.d.plugin/modules/nginx"
	_ "github.com/netdata/go.d.plugin/modules/openvpn"
	_ "github.com/netdata/go.d.plugin/modules/phpdaemon"
	_ "github.com/netdata/go.d.plugin/modules/phpfpm"
	_ "github.com/netdata/go.d.plugin/modules/pihole"
	_ "github.com/netdata/go.d.plugin/modules/portcheck"
	_ "github.com/netdata/go.d.plugin/modules/prometheus"
	_ "github.com/netdata/go.d.plugin/modules/pulsar"
	_ "github.com/netdata/go.d.plugin/modules/rabbitmq"
	_ "github.com/netdata/go.d.plugin/modules/scaleio"
	_ "github.com/netdata/go.d.plugin/modules/solr"
	_ "github.com/netdata/go.d.plugin/modules/springboot2"
	_ "github.com/netdata/go.d.plugin/modules/squidlog"
	_ "github.com/netdata/go.d.plugin/modules/tengine"
	_ "github.com/netdata/go.d.plugin/modules/unbound"
	_ "github.com/netdata/go.d.plugin/modules/vcsa"
	_ "github.com/netdata/go.d.plugin/modules/vernemq"
	_ "github.com/netdata/go.d.plugin/modules/vsphere"
	_ "github.com/netdata/go.d.plugin/modules/weblog"
	_ "github.com/netdata/go.d.plugin/modules/whoisquery"
	_ "github.com/netdata/go.d.plugin/modules/wmi"
	_ "github.com/netdata/go.d.plugin/modules/x509check"
	_ "github.com/netdata/go.d.plugin/modules/zookeeper"
)

var (
	cd, _     = os.Getwd()
	name      = "go.d"
	userDir   = os.Getenv("NETDATA_USER_CONFIG_DIR")
	stockDir  = os.Getenv("NETDATA_STOCK_CONFIG_DIR")
	varLibDir = os.Getenv("NETDATA_LIB_DIR")
	lockDir   = os.Getenv("NETDATA_LOCK_DIR")
	watchPath = os.Getenv("NETDATA_PLUGINS_GOD_WATCH_PATH")

	version = "unknown"
)

func confDir(opts *cli.Option) multipath.MultiPath {
	if len(opts.ConfDir) > 0 {
		return opts.ConfDir
	}
	if userDir != "" || stockDir != "" {
		return multipath.New(
			userDir,
			stockDir,
		)
	}
	return multipath.New(
		path.Join(cd, "/../../../../etc/netdata"),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d"),
	)
}

func modulesConfDir(opts *cli.Option) (mpath multipath.MultiPath) {
	if len(opts.ConfDir) > 0 {
		return opts.ConfDir
	}
	if userDir != "" || stockDir != "" {
		if userDir != "" {
			mpath = append(mpath, path.Join(userDir, name))
		}
		if stockDir != "" {
			mpath = append(mpath, path.Join(stockDir, name))
		}
		return multipath.New(mpath...)
	}
	return multipath.New(
		path.Join(cd, "/../../../../etc/netdata", name),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d", name),
	)
}

func watchPaths(opts *cli.Option) []string {
	if watchPath == "" {
		return opts.WatchPath
	}
	return append(opts.WatchPath, watchPath)
}

func stateFile() string {
	if varLibDir == "" {
		return ""
	}
	return path.Join(varLibDir, "god-jobs-statuses.json")
}

func main() {
	opts := parseCLI()

	if opts.Version {
		fmt.Println(fmt.Sprintf("go.d.plugin, version: %s", version))
		return
	}

	if opts.Debug {
		logger.SetSeverity(logger.DEBUG)
	}

	p := plugin.New(plugin.Config{
		Name:              name,
		ConfDir:           confDir(opts),
		ModulesConfDir:    modulesConfDir(opts),
		ModulesSDConfPath: watchPaths(opts),
		StateFile:         stateFile(),
		LockDir:           lockDir,
		RunModule:         opts.Module,
		MinUpdateEvery:    opts.UpdateEvery,
	})

	p.Run()
}

func parseCLI() *cli.Option {
	opt, err := cli.Parse(os.Args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	return opt
}
