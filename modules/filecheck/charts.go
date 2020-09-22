package filecheck

import "github.com/netdata/go-orchestrator/module"

var (
	fileCharts = module.Charts{
		fileExistenceChart.Copy(),
		fileModificationTimeAgoChart.Copy(),
		fileSizeChart.Copy(),
	}

	fileExistenceChart = module.Chart{
		ID:    "file_existence",
		Title: "File Existence",
		Units: "boolean",
		Fam:   "files",
		Ctx:   "filecheck.file_existence",
	}
	fileModificationTimeAgoChart = module.Chart{
		ID:    "file_mtime_ago",
		Title: "File Time Since the Last Modification",
		Units: "seconds",
		Fam:   "files",
		Ctx:   "filecheck.file_mtime_ago",
	}
	fileSizeChart = module.Chart{
		ID:    "file_size",
		Title: "File Size",
		Units: "bytes",
		Fam:   "files",
		Ctx:   "filecheck.file_size",
	}
)

var (
	dirCharts = module.Charts{
		dirExistenceChart.Copy(),
		dirModificationTimeChart.Copy(),
		dirNumOfFilesChart.Copy(),
	}

	dirExistenceChart = module.Chart{
		ID:    "dir_existence",
		Title: "Dir Existence",
		Units: "boolean",
		Fam:   "dirs",
		Ctx:   "filecheck.dir_existence",
	}
	dirModificationTimeChart = module.Chart{
		ID:    "dir_mtime_ago",
		Title: "Dir Time Since the Last Modification",
		Units: "seconds",
		Fam:   "dirs",
		Ctx:   "filecheck.dir_mtime_ago",
	}
	dirNumOfFilesChart = module.Chart{
		ID:    "dir_num_of_files",
		Title: "Dir Number of Files",
		Units: "files",
		Fam:   "dirs",
		Ctx:   "filecheck.dir_num_of_files",
	}
)
