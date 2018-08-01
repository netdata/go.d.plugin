package main

import (
	"os"
	"path"

	"github.com/l2isbad/go.d.plugin/internal/godplugin"
)

func main() {
	config := godplugin.NewConfig()
	config.Load(pluginConfigFile())
	godplugin.New(os.Args[1:]).Start()

}

// return netdata conf directory
// e.g. /opt/netdata/etc/netdata
func netdataConfigDir() string {
	dir := os.Getenv("NETDATA_CONFIG_DIR")

	if dir == "" {
		cd, _ := os.Getwd()
		dir = path.Join(cd, "/../../../../etc/netdata")
	}
	return configDir
}

// return netdata conf directory
// e.g. /opt/netdata/etc/netdata/go.d.conf
func pluginConfigFile() string {
	return path.Join(netdataConfigDir(), "go.d.conf")
}

func moduleConfigDir() string {
	file := os.Getenv("NETDATA_CONFIG_GO_MODULE")
	if file == "" {
		file = path.Join(netdataConfigDir(), "go.d.conf")
	}
	return file
}
