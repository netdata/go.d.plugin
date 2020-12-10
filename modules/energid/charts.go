package energid

import "github.com/netdata/go.d.plugin/agent/module"

type (
	Charts = module.Charts
	Dims   = module.Dims
)

var charts = Charts{
	{
		ID:    "blockindex",
		Title: "Blockchain Index",
		Units: "count",
		Fam:   "blockchain",
		Ctx:   "energid.blockindex",
		Type:  module.Area,
		Dims: Dims{
			{ID: "blockchain_blocks", Name: "Blocks", Algo: module.Absolute},
			{ID: "blockchain_headers", Name: "Headers", Algo: module.Absolute},
		},
	},
	{
		ID:    "difficulty",
		Title: "Blockchain Difficulty",
		Units: "difficulty",
		Fam:   "blockchain",
		Ctx:   "energid.difficulty",
		Type:  module.Line,
		Dims: Dims{
			{ID: "blockchain_difficulty", Name: "Diff", Algo: module.Absolute},
		},
	},
	{
		ID:    "mempool",
		Title: "MemPool",
		Units: "MiB",
		Fam:   "memory",
		Ctx:   "energid.mempool",
		Type:  module.Area,
		Dims: Dims{
			{ID: "mempool_max", Name: "Max", Algo: module.Absolute, Div: 1024 * 1024},
			{ID: "mempool_current", Name: "Usage", Algo: module.Absolute, Div: 1024 * 1024},
			{ID: "mempool_txsize", Name: "TX Size", Algo: module.Absolute, Div: 1024 * 1024},
		},
	},
	{
		ID:    "secmem",
		Title: "Secure Memory",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "energid.secmem",
		Type:  module.Area,
		Dims: Dims{
			{ID: "secmem_total", Name: "total"},
			{ID: "secmem_used", Name: "used"},
			{ID: "secmem_free", Name: "free"},
			{ID: "secmem_locked", Name: "locked"},
		},
	},
	{
		ID:    "network",
		Title: "Network",
		Units: "count",
		Fam:   "network",
		Ctx:   "energid.network",
		Type:  module.Line,
		Dims: Dims{
			{ID: "network_connections", Name: "Connections", Algo: module.Absolute},
		},
	},
	{
		ID:    "timeoffset",
		Title: "Network",
		Units: "seconds",
		Fam:   "network",
		Ctx:   "energid.timeoffset",
		Type:  module.Line,
		Dims: Dims{
			{ID: "network_timeoffset", Name: "Offseet", Algo: module.Absolute},
		},
	},
	{
		ID:    "utxo",
		Title: "UTXO",
		Units: "count",
		Fam:   "UTXO",
		Ctx:   "energid.utxo",
		Type:  module.Line,
		Dims: Dims{
			{ID: "utxo_count", Name: "UTXO", Algo: module.Absolute},
		},
	},
	{
		ID:    "xfers",
		Title: "UTXO",
		Units: "count",
		Fam:   "UTXO",
		Ctx:   "energid.xfers",
		Type:  module.Line,
		Dims: Dims{
			{ID: "utxo_xfers", Name: "Xfers", Algo: module.Absolute},
		},
	},
}
