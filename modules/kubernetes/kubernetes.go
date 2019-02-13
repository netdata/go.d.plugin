package kubernetes

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("kubernetes", creator)
}

const (
	defaultHTTPTimeout = time.Second * 2
	defaultURL         = "http://192.168.99.106:10255"
	// defaultURL         = "http://127.0.0.1:10255"
)

type Config struct {
	web.HTTP `yaml:",inline"`
}

// New creates Kubernetes with default values.
func New() *Kubernetes {
	return &Kubernetes{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: defaultURL},
				Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
			},
		},
		activePods: make(map[string]bool),
		charts:     &Charts{},
	}
}

// Kubernetes Kubernetes module.
type Kubernetes struct {
	module.Base
	Config `yaml:",inline"`

	charts    *Charts
	apiClient *apiClient
	// TODO: likely wrong
	activePods map[string]bool
}

// Cleanup makes cleanup.
func (Kubernetes) Cleanup() {}

// Init makes initialization.
func (k *Kubernetes) Init() bool {
	if k.URL == "" {
		k.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(k.Client)

	if err != nil {
		k.Errorf("error on creating http client : %v", err)
		return false
	}

	k.apiClient = newAPIClient(client, k.Request)

	return true
}

// Check makes check.
func (k *Kubernetes) Check() bool { return len(k.Collect()) > 0 }

// Charts creates Charts.
func (k Kubernetes) Charts() *Charts { return k.charts }

// Collect collects metrics.
func (k *Kubernetes) Collect() map[string]int64 {
	stats, err := k.apiClient.getStatsSummary()

	if err != nil {
		k.Error(err)
		return nil
	}

	var (
		metrics     = make(map[string]int64)
		updatedPods = make(map[string]bool)
	)

	for _, pod := range stats.Pods {
		if !k.activePods[pod.PodRef.UID] {
			k.activePods[pod.PodRef.UID] = true
			k.addPodToCharts(pod)
		}
		k.collectPodStats(metrics, pod)
		updatedPods[pod.PodRef.UID] = true
	}

	// TODO: remove immediately?
	for podIUD := range updatedPods {
		if k.activePods[podIUD] {
			continue
		}
		delete(k.activePods, podIUD)
		k.removePodFromCharts(podIUD)
	}

	return metrics
}

func (k *Kubernetes) removePodFromCharts(podIUD string) {}

func (k *Kubernetes) addPodToCharts(pod PodStats) {}

func (k *Kubernetes) collectPodStats(metrics map[string]int64, pod PodStats) {}
