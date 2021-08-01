package mongo

import (
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultHost, job.Local.Host, "wrong host")
	assert.Equal(t, time.Duration(defaultTimeout), job.Timeout, "wrong timeout")
	assert.Equal(t, defaultAuthDb, job.Authdb, "wrong auth db")
}

func TestMongo_Init(t *testing.T) {
	tests := map[string]struct {
		config  Config
		success bool
	}{
		"success on default config": {
			success: true,
			config:  New().Config,
		},
		"fails on unset 'address'": {
			success: false,
			config: Config{
				Local: Local{Host: "", Port: 0},
				Auth:  Auth{Host: "", Port: 0},
			},
		},
		"fails on invalid port": {
			success: false,
			config: Config{
				Local: Local{Host: "", Port: 999999},
				Auth:  Auth{Host: "", Port: 999999},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := New()
			m.Config = test.config
			assert.Equal(t, test.success, m.Init())
		})
	}
}

func TestMongo_Charts(t *testing.T) {
	m := New()
	require.True(t, m.Init())
	assert.NotNil(t, m.Charts())
}

func TestMongo_Cleanup(t *testing.T) {
	m := New()
	assert.NotPanics(t, m.Cleanup)

	require.True(t, m.Init())
	m.Cleanup()
	assert.Nil(t, m.client)

	require.NoError(t, m.initMongoClient())
	m.Cleanup()
	assert.Nil(t, m.client)
}

func TestMongo_initMongoClient_default(t *testing.T) {
	m := New()
	require.NoError(t, m.initMongoClient())
}
func TestMongo_initMongoClient_local(t *testing.T) {

	m := New()
	m.Local.Host = "localhost"
	m.Local.Port = 27017
	require.NoError(t, m.initMongoClient())
}
func TestMongo_initMongoClient_auth(t *testing.T) {
	m := New()
	m.Auth.Host = "localhost"
	m.Auth.Port = 27017
	m.Auth.User = "user"
	m.Auth.Pass = "pass"
	require.NoError(t, m.initMongoClient())
}

func TestMongo_initMongoClient_connectionStr(t *testing.T) {
	m := New()
	m.Config.ConnectionStr = "mongodb://user:pass@localhost:27017"
	require.NoError(t, m.initMongoClient())
}

func TestMongo_Check(t *testing.T) {
	m := New()
	m.Config.Local.Host = "not_resolvable_host"
	m.Config.Timeout = 0
	assert.False(t, m.Check())
}
