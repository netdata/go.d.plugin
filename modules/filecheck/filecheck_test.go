package filecheck

import (
	"strings"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestFilecheck_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestFilecheck_Init(t *testing.T) {
	tests := map[string]struct {
		config          Config
		wantNumOfCharts int
		wantFail        bool
	}{
		"default": {
			config:   New().Config,
			wantFail: true,
		},
		"empty files->include and dirs->include": {
			config: Config{
				Files: filesConfig{},
				Dirs:  dirsConfig{},
			},
			wantFail: true,
		},
		"files->include and dirs->include": {
			config: Config{
				Files: filesConfig{
					Include: []string{
						"/path/to/file1",
						"/path/to/file2",
					},
				},
				Dirs: dirsConfig{
					Include: []string{
						"/path/to/dir1",
						"/path/to/dir2",
					},
				},
			},
			wantNumOfCharts: len(fileCharts) + len(dirCharts),
		},
		"only files->include": {
			config: Config{
				Files: filesConfig{
					Include: []string{
						"/path/to/file1",
						"/path/to/file2",
					},
				},
			},
			wantNumOfCharts: len(fileCharts),
		},
		"only dirs->include": {
			config: Config{
				Dirs: dirsConfig{
					Include: []string{
						"/path/to/dir1",
						"/path/to/dir2",
					},
				},
			},
			wantNumOfCharts: len(dirCharts),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es := New()
			es.Config = test.config

			if test.wantFail {
				assert.False(t, es.Init())
			} else {
				require.True(t, es.Init())
				assert.Equal(t, test.wantNumOfCharts, len(*es.Charts()))
			}
		})
	}
}

func TestFilecheck_Check(t *testing.T) {
	tests := map[string]struct {
		prepare func() *Filecheck
	}{
		"collect files":                   {prepare: prepareFilecheckFiles},
		"collect only non existent files": {prepare: prepareFilecheckNonExistentFiles},
		"collect dirs":                    {prepare: prepareFilecheckDirs},
		"collect only non existent dirs":  {prepare: prepareFilecheckNonExistentDirs},
		"collect files and dirs":          {prepare: prepareFilecheckDirs},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fc := test.prepare()
			require.True(t, fc.Init())

			assert.True(t, fc.Check())
		})
	}
}

func TestFilecheck_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *Filecheck
		wantCollected map[string]int64
	}{
		"collect files": {
			prepare: prepareFilecheckFiles,
			wantCollected: map[string]int64{
				"file_exists_testdata/empty_file.log":        1,
				"file_exists_testdata/file.log":              1,
				"file_exists_testdata/non_existent_file.log": 0,
				"file_mtime_ago_testdata/empty_file.log":     1621,
				"file_mtime_ago_testdata/file.log":           701,
				"file_size_bytes_testdata/empty_file.log":    0,
				"file_size_bytes_testdata/file.log":          5707,
			},
		},
		"collect only non existent files": {
			prepare: prepareFilecheckNonExistentFiles,
			wantCollected: map[string]int64{
				"file_exists_testdata/non_existent_file.log": 0,
			},
		},
		"collect dirs": {
			prepare: prepareFilecheckDirs,
			wantCollected: map[string]int64{
				"dir_exists_testdata/dir":              1,
				"dir_exists_testdata/empty_dir":        1,
				"dir_exists_testdata/non_existent_dir": 0,
				"dir_mtime_ago_testdata/dir":           657,
				"dir_mtime_ago_testdata/empty_dir":     662,
				"dir_num_of_files_testdata/dir":        3,
				"dir_num_of_files_testdata/empty_dir":  0,
			},
		},
		"collect only non existent dirs": {
			prepare: prepareFilecheckNonExistentDirs,
			wantCollected: map[string]int64{
				"dir_exists_testdata/non_existent_dir": 0,
			},
		},
		"collect files and dirs": {
			prepare: prepareFilecheckFilesDirs,
			wantCollected: map[string]int64{
				"dir_exists_testdata/dir":                    1,
				"dir_exists_testdata/empty_dir":              1,
				"dir_exists_testdata/non_existent_dir":       0,
				"dir_mtime_ago_testdata/dir":                 657,
				"dir_mtime_ago_testdata/empty_dir":           662,
				"dir_num_of_files_testdata/dir":              3,
				"dir_num_of_files_testdata/empty_dir":        0,
				"file_exists_testdata/empty_file.log":        1,
				"file_exists_testdata/file.log":              1,
				"file_exists_testdata/non_existent_file.log": 0,
				"file_mtime_ago_testdata/empty_file.log":     1621,
				"file_mtime_ago_testdata/file.log":           701,
				"file_size_bytes_testdata/empty_file.log":    0,
				"file_size_bytes_testdata/file.log":          5707,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fc := test.prepare()
			require.True(t, fc.Init())

			collected := fc.Collect()

			copyModTime(test.wantCollected, collected)
			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, fc, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, fc *Filecheck, collected map[string]int64) {
	// TODO: check other charts
	for _, chart := range *fc.Charts() {
		if chart.Obsolete {
			continue
		}
		switch chart.ID {
		case fileExistenceChart.ID, dirExistenceChart.ID:
			for _, dim := range chart.Dims {
				_, ok := collected[dim.ID]
				assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
			}
			for _, v := range chart.Vars {
				_, ok := collected[v.ID]
				assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
			}
		}
	}
}

func prepareFilecheckFiles() *Filecheck {
	fc := New()
	fc.Config.Files.Include = []string{
		"testdata/empty_file.log",
		"testdata/file.log",
		"testdata/non_existent_file.log",
	}
	return fc
}

func prepareFilecheckNonExistentFiles() *Filecheck {
	fc := New()
	fc.Config.Files.Include = []string{
		"testdata/non_existent_file.log",
	}
	return fc
}

func prepareFilecheckDirs() *Filecheck {
	fc := New()
	fc.Config.Dirs.Include = []string{
		"testdata/empty_dir",
		"testdata/dir",
		"testdata/non_existent_dir",
	}
	return fc
}

func prepareFilecheckNonExistentDirs() *Filecheck {
	fc := New()
	fc.Config.Dirs.Include = []string{
		"testdata/non_existent_dir",
	}
	return fc
}

func prepareFilecheckFilesDirs() *Filecheck {
	fc := New()
	fc.Config.Files.Include = []string{
		"testdata/empty_file.log",
		"testdata/file.log",
		"testdata/non_existent_file.log",
	}
	fc.Config.Dirs.Include = []string{
		"testdata/empty_dir",
		"testdata/dir",
		"testdata/non_existent_dir",
	}
	return fc
}

func copyModTime(dst, src map[string]int64) {
	if src == nil || dst == nil {
		return
	}
	for key := range src {
		if strings.Contains(key, "mtime") {
			dst[key] = src[key]
		}
	}
}
