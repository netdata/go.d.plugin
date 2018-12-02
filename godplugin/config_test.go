package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginConfig_isModuleEnabled(t *testing.T) {
	modName1 := "modName1"
	modName2 := "modName2"
	modName3 := "modName3"

	conf := config{
		DefaultRun: true,
		Modules: map[string]bool{
			modName1: true,
			modName2: false,
		},
	}

	assert.True(t, conf.isModuleEnabled(modName1, false))
	assert.False(t, conf.isModuleEnabled(modName2, false))
	assert.Equal(
		t,
		conf.DefaultRun,
		conf.isModuleEnabled(modName3, false),
	)

	assert.True(t, conf.isModuleEnabled(modName1, true))
	assert.False(t, conf.isModuleEnabled(modName2, true))
	assert.Equal(
		t,
		!conf.DefaultRun,
		conf.isModuleEnabled(modName3, true),
	)

	conf.DefaultRun = false

	assert.True(t, conf.isModuleEnabled(modName1, false))
	assert.False(t, conf.isModuleEnabled(modName2, false))
	assert.Equal(
		t,
		conf.DefaultRun,
		conf.isModuleEnabled(modName3, false),
	)

	assert.True(t, conf.isModuleEnabled(modName1, true))
	assert.False(t, conf.isModuleEnabled(modName2, true))
	assert.Equal(
		t,
		conf.DefaultRun,
		conf.isModuleEnabled(modName3, true),
	)

}

func TestModuleConfig_updateJobs(t *testing.T) {
	conf := moduleConfig{
		Global: &moduleGlobal{
			UpdateEvery:        10,
			AutoDetectionRetry: 10,
		},
		Jobs: []map[string]interface{}{
			{"name": "job1"},
			{"name": "job2", "update_every": 1},
		},
	}

	conf.updateJobs(0, 0)

	assert.Equal(
		t,
		[]map[string]interface{}{
			{"name": "job1", "update_every": 10, "autodetection_retry": 10},
			{"name": "job2", "update_every": 1, "autodetection_retry": 10},
		},
		conf.Jobs,
	)
}

func TestModuleConfig_UpdateJobsRewriteModuleUpdateEvery(t *testing.T) {
	conf := moduleConfig{
		Global: &moduleGlobal{
			UpdateEvery:        10,
			AutoDetectionRetry: 10,
		},
		Jobs: []map[string]interface{}{
			{"name": "job1"},
			{"name": "job2", "update_every": 1},
		},
	}

	conf.updateJobs(20, 0)

	assert.Equal(
		t,
		[]map[string]interface{}{
			{"name": "job1", "update_every": 20, "autodetection_retry": 10},
			{"name": "job2", "update_every": 1, "autodetection_retry": 10},
		},
		conf.Jobs,
	)
}

func TestModuleConfig_UpdateJobsRewritePluginUpdateEvery(t *testing.T) {
	conf := moduleConfig{
		Global: &moduleGlobal{
			UpdateEvery:        10,
			AutoDetectionRetry: 10,
		},
		Jobs: []map[string]interface{}{
			{"name": "job1"},
			{"name": "job2", "update_every": 1},
		},
	}

	conf.updateJobs(0, 5)

	assert.Equal(
		t,
		[]map[string]interface{}{
			{"name": "job1", "update_every": 10, "autodetection_retry": 10},
			{"name": "job2", "update_every": 5, "autodetection_retry": 10},
		},
		conf.Jobs,
	)
}
