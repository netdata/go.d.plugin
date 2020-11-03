package example

import "github.com/netdata/go.d.plugin/agent/module"

var chartTemplate = module.Chart{
	ID:    "random_%d",
	Title: "A Random Number",
	Units: "random",
	Fam:   "random",
	Ctx:   "example.random",
}
