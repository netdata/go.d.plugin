package x509check

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
		Title: "Time Until Certificate Expiration",
		Units: "seconds",
		Fam:   "expiration time",
		Ctx:   "x509check.time_until_expiration",
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

var withRevocationCharts = Charts{
	{
		ID:    "time_until_expiration",
		Title: "Time Until Certificate Expiration",
		Units: "seconds",
		Fam:   "expiration time",
		Ctx:   "x509check.time_until_expiration",
		Opts:  Opts{StoreFirst: true},
		Dims: Dims{
			{ID: "expiry"},
		},
		Vars: Vars{
			{ID: "days_until_expiration_warning"},
			{ID: "days_until_expiration_critical"},
		},
	},
	{
		ID:    "revocation_status",
		Title: "Revocation Status",
		Units: "boolean",
		Fam:   "revocation",
		Ctx:   "x509check.revocation_status",
		Opts:  Opts{StoreFirst: true},
		Dims: Dims{
			{ID: "revoked"},
		},
	},
}
