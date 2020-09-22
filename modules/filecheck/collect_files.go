package filecheck

import (
	"fmt"
	"os"
	"time"

	"github.com/netdata/go-orchestrator/module"
)

func (fc *Filecheck) collectFiles(mx map[string]int64) {
	curTime := time.Now()
	for _, filepath := range fc.Files.Include {
		fc.collectFile(mx, filepath, curTime)
	}
}

func (fc *Filecheck) collectFile(mx map[string]int64, filepath string, curTime time.Time) {
	if !fc.collectedFiles[filepath] {
		fc.collectedFiles[filepath] = true
		fc.addFileToCharts(filepath)
	}

	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			mx[fileDimID(filepath, "exists")] = 0
		} else {
			mx[fileDimID(filepath, "exists")] = 1
		}
		return
	}

	if info.IsDir() {
		return
	}

	mx[fileDimID(filepath, "exists")] = 1
	mx[fileDimID(filepath, "size_bytes")] = info.Size()
	mx[fileDimID(filepath, "mtime_ago")] = int64(curTime.Sub(info.ModTime()).Seconds())
}

func (fc *Filecheck) addFileToCharts(filepath string) {
	for _, c := range fileCharts {
		chart := fc.Charts().Get(c.ID)
		if chart == nil {
			continue
		}

		var id string
		switch chart.ID {
		case fileExistenceChart.ID:
			id = fileDimID(filepath, "exists")
		case fileModificationTimeAgoChart.ID:
			id = fileDimID(filepath, "mtime_ago")
		case fileSizeChart.ID:
			id = fileDimID(filepath, "size_bytes")
		default:
			continue
		}

		_ = chart.AddDim(&module.Dim{ID: id, Name: filepath})
		chart.MarkNotCreated()
	}
}

func fileDimID(filepath, metric string) string {
	return fmt.Sprintf("file_%s_%s", metric, filepath)
}
