package freeradius

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "authentication_access",
		Title: "Authentication Access",
		Units: "packets/s",
		Fam:   "authentication",
		Ctx:   "freeradius.authentication",
		Dims: Dims{
			{ID: "access-accepts", Name: "accepts", Algo: modules.Incremental},
			{ID: "access-rejects", Name: "rejects", Algo: modules.Incremental},
		},
	},
	{
		ID:    "bad_authentication_requests",
		Title: "Bad Authentication Requests",
		Units: "packets/s",
		Fam:   "authentication",
		Ctx:   "freeradius.bad_authentication",
		Dims: Dims{
			{ID: "auth-dropped-requests", Name: "dropped", Algo: modules.Incremental},
			{ID: "auth-duplicate-requests", Name: "duplicate", Algo: modules.Incremental},
			{ID: "auth-invalid-requests", Name: "invalid", Algo: modules.Incremental},
			{ID: "auth-malformed-requests", Name: "malformed", Algo: modules.Incremental},
			{ID: "auth-unknown-types", Name: "unknown-type", Algo: modules.Incremental},
		},
	},
	{
		ID:    "proxy_authentication_access",
		Title: "Proxy Authentication Access",
		Units: "packets/s",
		Fam:   "authentication",
		Ctx:   "freeradius.authentication",
		Dims: Dims{
			{ID: "proxy-access-accepts", Name: "accepts", Algo: modules.Incremental},
			{ID: "proxy-access-rejects", Name: "rejects", Algo: modules.Incremental},
		},
	},
	{
		ID:    "bad_proxy_authentication_requests",
		Title: "Bad Proxy Authentication Requests",
		Units: "packets/s",
		Fam:   "authentication",
		Ctx:   "freeradius.bad_authentication",
		Dims: Dims{
			{ID: "proxy-auth-dropped-requests", Name: "dropped", Algo: modules.Incremental},
			{ID: "proxy-auth-duplicate-requests", Name: "duplicate", Algo: modules.Incremental},
			{ID: "proxy-auth-invalid-requests", Name: "invalid", Algo: modules.Incremental},
			{ID: "proxy-auth-malformed-requests", Name: "malformed", Algo: modules.Incremental},
			{ID: "proxy-auth-unknown-types", Name: "unknown-types", Algo: modules.Incremental},
		},
	},
	{
		ID:    "accounting",
		Title: "Accounting",
		Units: "packets/s",
		Fam:   "accounting",
		Ctx:   "freeradius.accounting",
		Dims: Dims{
			{ID: "accounting-requests", Name: "requests", Algo: modules.Incremental},
			{ID: "accounting-responses", Name: "responses", Algo: modules.Incremental},
		},
	},
	{
		ID:    "bad_accounting_requests",
		Title: "Bad Accounting Requests",
		Units: "packets/s",
		Fam:   "accounting",
		Ctx:   "freeradius.bad_accounting",
		Dims: Dims{
			{ID: "acct-dropped-requests", Name: "dropped", Algo: modules.Incremental},
			{ID: "acct-duplicate-requests", Name: "duplicate", Algo: modules.Incremental},
			{ID: "acct-invalid-requests", Name: "invalid", Algo: modules.Incremental},
			{ID: "acct-malformed-requests", Name: "malformed", Algo: modules.Incremental},
			{ID: "acct-unknown-types", Name: "unknown-types", Algo: modules.Incremental},
		},
	},
	{
		ID:    "proxy_accounting",
		Title: "Proxy Accounting",
		Units: "packets/s",
		Fam:   "accounting",
		Ctx:   "freeradius.accounting",
		Dims: Dims{
			{ID: "proxy-accounting-requests", Name: "requests", Algo: modules.Incremental},
			{ID: "proxy-accounting-responses", Name: "responses", Algo: modules.Incremental},
		},
	},
	{
		ID:    "bad_proxy_accounting_requests",
		Title: "Bad Proxy Accounting Requests",
		Units: "packets/s",
		Fam:   "accounting",
		Ctx:   "freeradius.bad_accounting",
		Dims: Dims{
			{ID: "proxy-acct-dropped-requests", Name: "dropped", Algo: modules.Incremental},
			{ID: "proxy-acct-duplicate-requests", Name: "duplicate", Algo: modules.Incremental},
			{ID: "proxy-acct-invalid-requests", Name: "invalid", Algo: modules.Incremental},
			{ID: "proxy-acct-malformed-requests", Name: "malformed", Algo: modules.Incremental},
			{ID: "proxy-acct-unknown-types", Name: "unknown-types", Algo: modules.Incremental},
		},
	},
}
