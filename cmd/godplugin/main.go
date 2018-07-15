package main

import (
	"os"
	"path"

	"github.com/l2isbad/go.d.plugin/internal/godplugin"
)

func main() {
	godplugin.New(getConfigDir()).Start()
}

func getConfigDir() string {
	pd := os.Getenv("NETDATA_CONFIG_DIR")

	if pd == "" {
		cd, _ := os.Getwd()
		pd = path.Join(cd, "/../../../../etc/netdata")
	}
	return pd
}
