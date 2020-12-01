package redis

import (
	"context"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/blang/semver/v4"
	"github.com/go-redis/redis/v8"
)

func init() {
	module.Register("redis", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Redis {
	return &Redis{
		Config: Config{
			Address: "redis://@localhost:6379",
			Timeout: web.Duration{Duration: time.Second},
		},

		collectedCommands: make(map[string]bool),
		collectedDbs:      make(map[string]bool),
	}
}

type Config struct {
	Address          string       `yaml:"address"`
	Timeout          web.Duration `yaml:"timeout"`
	tlscfg.TLSConfig `yaml:",inline"`
}

type (
	Redis struct {
		module.Base
		Config `yaml:",inline"`

		rdb redisClient

		server  string
		version *semver.Version

		collectedCommands map[string]bool
		collectedDbs      map[string]bool

		charts *module.Charts
	}
	redisClient interface {
		Info(ctx context.Context, section ...string) *redis.StringCmd
		Close() error
	}
)

func (r *Redis) Init() bool {
	err := r.validateConfig()
	if err != nil {
		r.Errorf("config validation: %v", err)
		return false
	}

	rdb, err := r.initRedisClient()
	if err != nil {
		r.Errorf("init redis client: %v", err)
		return false
	}
	r.rdb = rdb

	charts, err := r.initCharts()
	if err != nil {
		r.Errorf("init charts: %v", err)
		return false
	}
	r.charts = charts

	return true
}

func (r *Redis) Check() bool {
	return len(r.Collect()) > 0
}

func (r *Redis) Charts() *module.Charts {
	return r.charts
}

func (r *Redis) Collect() map[string]int64 {
	ms, err := r.collect()
	if err != nil {
		r.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (r *Redis) Cleanup() {
	if r.rdb == nil {
		return
	}
	err := r.rdb.Close()
	if err != nil {
		r.Warningf("cleanup: error on closing redis client [%s]: %v", r.Address, err)
	}
	r.rdb = nil
}
