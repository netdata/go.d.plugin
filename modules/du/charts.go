package du

import "github.com/netdata/go.d.plugin/agent/module"

var chartTemplate = module.Chart{
	ID:    "filesize",
	Title: "File/Folder size",
	Units: "bytes",
	Fam:   "diskusage",
	Ctx:   "du.filesize",
}
