package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/netdata/go.d.plugin/agent"
	"github.com/netdata/go.d.plugin/cli"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/pkg/multipath"

	"github.com/jessevdk/go-flags"

	_ "github.com/netdata/go.d.plugin/modules"
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

func init() {
	// https://github.com/netdata/netdata/issues/8949#issuecomment-638294959
	if v := os.Getenv("TZ"); strings.HasPrefix(v, ":") {
		_ = os.Unsetenv("TZ")
	}
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

	a := agent.New(agent.Config{
		Name:              name,
		ConfDir:           confDir(opts),
		ModulesConfDir:    modulesConfDir(opts),
		ModulesSDConfPath: watchPaths(opts),
		StateFile:         stateFile(),
		LockDir:           lockDir,
		RunModule:         opts.Module,
		MinUpdateEvery:    opts.UpdateEvery,
	})

	a.Debugf("plugin: name=%s, version=%s", a.Name, version)
	if u, err := user.Current(); err == nil {
		a.Debugf("current user: name=%s, uid=%s", u.Username, u.Uid)
	}

	a.Run()
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
