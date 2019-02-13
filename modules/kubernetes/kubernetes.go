package kubernetes

import (
	"fmt"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
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
	defaultURL         = "http://192.168.99.111:10255"
	// defaultURL         = "http://127.0.0.1:10255"
)

type Config struct {
	web.HTTP    `yaml:",inline"`
	UpdateEvery int `yaml:"update_every"`
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
	Config     `yaml:",inline"`
	PermitPods string

	apiClient *apiClient
	// TODO: likely wrong
	activePods map[string]bool
	permitPods matcher.Matcher

	charts *Charts
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

	if k.PermitPods != "" {
		m, err := matcher.NewSimplePatternsMatcher(k.PermitPods)
		if err != nil {
			k.Errorf("error on creating permit_pods matcher : %v", err)
			return false
		}
		// k.permitPods = matcher.WithCache(m)
		k.permitPods = m
	}

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
		// TODO: match on what?
		if k.permitPods != nil && !k.permitPods.MatchString(pod.PodRef.Name) {
			continue
		}
		if !k.activePods[pod.PodRef.UID] {
			k.activePods[pod.PodRef.UID] = true
			k.addPodToCharts(&pod)
		}

		for k, v := range podStatsToMap(&pod) {
			metrics[k] = v
		}

		updatedPods[pod.PodRef.UID] = true
	}

	// TODO: remove immediately?
	for podIUD := range updatedPods {
		if !k.activePods[podIUD] {
			delete(k.activePods, podIUD)
			k.removePodFromCharts(podIUD)
		}
	}

	return metrics
}

func (k *Kubernetes) removePodFromCharts(podIUD string) {
	for _, chart := range *k.charts {
		if strings.Contains(chart.ID, podIUD) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (k *Kubernetes) addPodToCharts(pod *PodStats) {
	for _, c := range []Chart{
		chartCPUStats,
		chartMemoryStatsUsage,
		chartMemoryStatsPageFaults,
	} {
		chart := c.Copy()
		chart.ID = fmt.Sprintf(chart.ID, fmt.Sprintf("%s_%s_%s", pod.PodRef.Name, pod.PodRef.UID, pod.PodRef.Namespace))
		chart.Fam = pod.PodRef.Name
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, pod.PodRef.UID)
		}
		_ = k.charts.Add(chart)
	}

	// TODO:
	// 1. dim.Div = k.UpdateEvery * dim.Div in CPU
	// 2. remove AvailableBytes if nil from MemStats

}

func podStatsToMap(pod *PodStats) map[string]int64 {
	rv := make(map[string]int64)
	if has(pod.CPU.UsageCoreNanoSeconds) {
		rv[pod.PodRef.UID+"_cpu_stats_usage_core_nano_seconds"] = *pod.CPU.UsageCoreNanoSeconds
	}
	if has(pod.Memory.AvailableBytes) {
		rv[pod.PodRef.UID+"_memory_stats_available_bytes"] = *pod.Memory.AvailableBytes
	}
	if has(pod.Memory.UsageBytes) {
		rv[pod.PodRef.UID+"_memory_stats_usage_bytes"] = *pod.Memory.UsageBytes
	}
	if has(pod.Memory.WorkingSetBytes) {
		rv[pod.PodRef.UID+"_memory_stats_working_set_bytes"] = *pod.Memory.WorkingSetBytes
	}
	if has(pod.Memory.RSSBytes) {
		rv[pod.PodRef.UID+"_memory_stats_rss_bytes"] = *pod.Memory.RSSBytes
	}
	if has(pod.Memory.PageFaults) {
		rv[pod.PodRef.UID+"_memory_stats_page_faults"] = *pod.Memory.PageFaults
	}
	if has(pod.Memory.MajorPageFaults) {
		rv[pod.PodRef.UID+"_memory_stats_major_page_faults"] = *pod.Memory.MajorPageFaults
	}

	return rv
}

func has(v *int64) bool { return v != nil }
