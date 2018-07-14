package main

import (
	"os"
	"path"

	"github.com/l2isbad/go.d.plugin/internal/godplugin"
)

func main() {
	// plugin and modules dirs
	var pd, md string

	pd = os.Getenv("NETDATA_CONFIG_DIR")

	if pd == "" {
		cd, _ := os.Getwd()
		pd = path.Join(cd, "/../../../../etc/netdata")
	}

	md = path.Join(pd, "go.d/")

	godplugin.New(pd, md).Start()
}
