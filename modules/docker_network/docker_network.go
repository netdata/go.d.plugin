package docker_network

import (
	"context"
	_ "embed"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
	"time"
)

//go:embed "config_schema.json"
var configSchema string

// Our init function registers the module with the agent.
func init() {
	module.Register("docker_network", module.Creator{
		JobConfigSchema: configSchema,
		Create:          func() module.Module { return New() },
	})
}

// New creates a new instance of our module.
func New() *DockerNetwork {
	return &DockerNetwork{
		// This config is only overridden by the config file.
		Config: Config{
			Address: docker.DefaultDockerHost,
			Timeout: web.Duration{Duration: time.Second * 5},
		},
		charts: summaryCharts.Copy(), // TODO: Implement summaryCharts
		newClient: func(cfg Config) (dockerClient, error) {
			return docker.NewClientWithOpts(docker.WithHost(cfg.Address))
		},
		containers: make(map[string]bool),
	}
}

type Config struct {
	Timeout web.Duration `yaml:"timeout"`
	Address string       `yaml:"address"`
}

type (
	DockerNetwork struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		newClient     func(Config) (dockerClient, error)
		client        dockerClient
		verNegotiated bool

		containers map[string]bool
	}
	// For our docker client, we use the official docker client library.
	// We embed the client in our module, so we can easily mock it in our tests.
	dockerClient interface {
		NegotiateAPIVersion(ctx context.Context)
		Info(ctx context.Context) (types.Info, error)
		// We just need a list of containers and the stats about it
		ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
		// TODO: Implement container stats
		Close() error
	}
)

// These are all boilerplate functions that we need to implement for our module.

// Init will initialize our module.
func (d *DockerNetwork) Init() bool {
	return true
}

// Check will check if the module is able to collect metrics.
func (d *DockerNetwork) Check() bool {
	return len(d.Collect()) > 0
}

// Charts returns the charts that we want to expose to the agent.
func (d *DockerNetwork) Charts() *module.Charts {
	return d.charts
}

// Collect will collect the metrics from the docker client.
func (d *DockerNetwork) Collect() map[string]int64 {
	// All we'll be collecting is the current bit rate of the network interface.
	// This does mean we need to store a previous value, so we can calculate the difference.
	mx, err := d.collect() // TODO: Implement collect
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

// Cleanup will close our docker client.s
func (d *DockerNetwork) Cleanup() {
	if d.client == nil {
		return
	}
	if err := d.client.Close(); err != nil {
		d.Warningf("error on closing docker client: %v", err)
	}
	d.client = nil
}
