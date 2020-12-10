package filecheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (fc *Filecheck) collectFiles(mx map[string]int64) {
	curTime := time.Now()
	for _, filespath := range fc.Files.Include {
		filepaths, _ := filepath.Glob(filespath)
		if filepaths == nil {
			fc.collectFile(mx, filespath, curTime)
		} else {
			for _, filepathpart := range filepaths {
				fc.collectFile(mx, filepathpart, curTime)
			}
		}
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
		fc.Debug(err)
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
	for _, chart := range *fc.Charts() {
		if !strings.HasPrefix(chart.ID, "file_") {
			continue
		}

		var id string
		switch chart.ID {
		case fileExistenceChart.ID:
			id = fileDimID(filepath, "exists")
		case fileModTimeAgoChart.ID:
			id = fileDimID(filepath, "mtime_ago")
		case fileSizeChart.ID:
			id = fileDimID(filepath, "size_bytes")
		default:
			fc.Warningf("add dimension: couldn't dim id for '%s' chart (file '%s')", chart.ID, filepath)
			continue
		}

		dim := &module.Dim{ID: id, Name: filepath}

		if err := chart.AddDim(dim); err != nil {
			fc.Warning(err)
			continue
		}
		chart.MarkNotCreated()
	}
}

func fileDimID(filepath, metric string) string {
	return fmt.Sprintf("file_%s_%s", filepath, metric)
}
