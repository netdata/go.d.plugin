package whoisquery

import "github.com/netdata/go.d.plugin/agent/module"

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
	Opts   = module.Opts
)

var charts = Charts{
	{
		ID:    "time_until_expiration",
		Title: "Time Until Domain Expiration",
		Units: "seconds",
		Fam:   "expiration time",
		Ctx:   "whoisquery.time_until_expiration",
		Opts:  Opts{StoreFirst: true},
		Dims: Dims{
			{ID: "expiry"},
		},
		Vars: Vars{
			{ID: "days_until_expiration_warning"},
			{ID: "days_until_expiration_critical"},
		},
	},
}
