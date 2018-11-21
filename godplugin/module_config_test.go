package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	moduleConfNotExist = "tests/module.conf-not-exist.yml"
	moduleConf         = "tests/module.conf.yml"
)

func TestModuleConfigNone(t *testing.T) {
	config := moduleConfig{}
	assert.Error(t, config.load(moduleConfNotExist))

}

func TestModuleConfigDefault(t *testing.T) {
	config := moduleConfig{}

	require.NoError(t, config.load(moduleConf))
	require.Len(t, config.Jobs, 2)

	config.updateJobs(0, 1)

	assert.Equal(
		t, config.Jobs,
		[]rawConfig{
			{
				"name":                "job1",
				"update_every":        1,  // default
				"autodetection_retry": 10, // global
				"param":               10, // global
			},
			{
				"name":                "job2",
				"update_every":        22, // job specific
				"autodetection_retry": 10, // global
				"param":               22, // job specific
			},
		},
	)
}

func TestModuleConfigWithModUpdateEvery(t *testing.T) {
	config := moduleConfig{}

	require.NoError(t, config.load(moduleConf))
	require.Len(t, config.Jobs, 2)

	// updateEvery default 1 => 99
	config.updateJobs(99, 1)

	assert.Equal(
		t, config.Jobs,
		[]rawConfig{
			{
				"name":                "job1",
				"update_every":        99, // overridden default
				"autodetection_retry": 10, // global
				"param":               10, // global
			},
			{
				"name":                "job2",
				"update_every":        22, // job specific
				"autodetection_retry": 10, // global
				"param":               22, // job specific
			},
		},
	)
}

func TestModuleConfigWithGlobalUpdateEvery(t *testing.T) {
	config := moduleConfig{}

	require.NoError(t, config.load(moduleConf))
	require.Len(t, config.Jobs, 2)

	// minimum updateEvery 1 => 15
	config.updateJobs(0, 15)

	assert.Equal(
		t, config.Jobs,
		[]rawConfig{
			{
				"name":                "job1",
				"update_every":        15, // overridden default
				"autodetection_retry": 10, // global
				"param":               10, // global
			},
			{
				"name":                "job2",
				"update_every":        22, // job specific
				"autodetection_retry": 10, // global
				"param":               22, // job specific
			},
		},
	)
}

func TestModuleConfigWithGlobalAndModuleUpdateEvery(t *testing.T) {
	config := moduleConfig{}

	require.NoError(t, config.load(moduleConf))
	require.Len(t, config.Jobs, 2)

	// updateEvery default 1 => 5, minimum updateEvery 1 => 15
	config.updateJobs(5, 15)

	assert.Equal(
		t, config.Jobs,
		[]rawConfig{
			{
				"name":                "job1",
				"update_every":        15, // overridden default
				"autodetection_retry": 10, // global
				"param":               10, // global
			},
			{
				"name":                "job2",
				"update_every":        22, // job specific
				"autodetection_retry": 10, // global
				"param":               22, // job specific
			},
		},
	)
}
