package k8s_kubelet

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"math"

	"github.com/netdata/go-orchestrator/module"
)

func (k *Kubelet) collect() (map[string]int64, error) {
	raw, err := k.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	k.collectToken(raw, mx)
	k.collectRESTClientHTTPRequests(raw, mx)
	k.collectAPIServer(raw, mx)
	k.collectKubelet(raw, mx)
	k.collectVolumeManager(raw, mx)

	return stm.ToMap(mx), nil
}

func (k *Kubelet) collectVolumeManager(raw prometheus.Metrics, mx *metrics) {
	vmPlugins := make(map[string]*volumeManagerPlugin)

	for _, metric := range raw.FindByName("volume_manager_total_volumes") {
		pluginName := metric.Labels.Get("plugin_name")
		state := metric.Labels.Get("state")

		if !k.activeVolumeManagerPlugins[pluginName] {
			_ = k.charts.Add(newVolumeManagerChart(pluginName))
			k.activeVolumeManagerPlugins[pluginName] = true
		}
		if _, ok := vmPlugins[pluginName]; !ok {
			vmPlugins[pluginName] = &volumeManagerPlugin{}
		}

		switch state {
		case "actual_state_of_world":
			vmPlugins[pluginName].State.Actual.Set(metric.Value)
		case "desired_state_of_world":
			vmPlugins[pluginName].State.Desired.Set(metric.Value)
		}
	}

	mx.VolumeManager.Plugins = vmPlugins
}

func (k *Kubelet) collectKubelet(raw prometheus.Metrics, mx *metrics) {
	value := raw.FindByName("kubelet_node_config_error").Max()
	mx.Kubelet.NodeConfigError.Set(value)

	value = raw.FindByName("kubelet_running_container_count").Max()
	mx.Kubelet.RunningContainerCount.Set(value)

	value = raw.FindByName("kubelet_running_pod_count").Max()
	mx.Kubelet.RunningPodCount.Set(value)

	k.collectRuntimeOperations(raw, mx)
	k.collectRuntimeOperationsErrors(raw, mx)
	k.collectDockerOperations(raw, mx)
	k.collectDockerOperationsErrors(raw, mx)
	k.collectPLEGRelisting(raw, mx)
}

func (k *Kubelet) collectAPIServer(raw prometheus.Metrics, mx *metrics) {
	value := raw.FindByName("apiserver_audit_requests_rejected_total").Max()
	mx.APIServer.Audit.Requests.Rejected.Set(value)

	value = raw.FindByName("apiserver_storage_data_key_generation_failures_total").Max()
	mx.APIServer.Storage.DataKeyGeneration.Failures.Set(value)

	value = raw.FindByName("apiserver_storage_envelope_transformation_cache_misses_total").Max()
	mx.APIServer.Storage.EnvelopeTransformation.CacheMisses.Set(value)

	k.collectStorageDataKeyGenerationLatencies(raw, mx)
}

func (k *Kubelet) collectToken(raw prometheus.Metrics, mx *metrics) {
	value := raw.FindByName("get_token_count").Max()
	mx.Token.Count.Set(value)

	value = raw.FindByName("get_token_fail_count").Max()
	mx.Token.FailCount.Set(value)
}

func (k *Kubelet) collectPLEGRelisting(raw prometheus.Metrics, mx *metrics) {
	// Summary
	for _, metric := range raw.FindByName("kubelet_pleg_relist_interval_microseconds") {
		if math.IsNaN(metric.Value) {
			continue
		}
		quantile := metric.Labels.Get("quantile")
		switch quantile {
		case "0.5":
			mx.Kubelet.PLEG.Relist.Interval.Quantile05.Set(metric.Value)
		case "0.9":
			mx.Kubelet.PLEG.Relist.Interval.Quantile09.Set(metric.Value)
		case "0.99":
			mx.Kubelet.PLEG.Relist.Interval.Quantile099.Set(metric.Value)
		}
	}
	for _, metric := range raw.FindByName("kubelet_pleg_relist_latency_microseconds") {
		if math.IsNaN(metric.Value) {
			continue
		}
		quantile := metric.Labels.Get("quantile")
		switch quantile {
		case "0.5":
			mx.Kubelet.PLEG.Relist.Latency.Quantile05.Set(metric.Value)
		case "0.9":
			mx.Kubelet.PLEG.Relist.Latency.Quantile09.Set(metric.Value)
		case "0.99":
			mx.Kubelet.PLEG.Relist.Latency.Quantile099.Set(metric.Value)
		}
	}
}

func (k *Kubelet) collectStorageDataKeyGenerationLatencies(raw prometheus.Metrics, mx *metrics) {
	for _, metric := range raw.FindByName("apiserver_storage_data_key_generation_latencies_microseconds_bucket") {
		bucket := metric.Labels.Get("le")
		switch bucket {
		case "5":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE5.Set(metric.Value)
		case "10":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE10.Set(metric.Value)
		case "20":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE20.Set(metric.Value)
		case "40":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE40.Set(metric.Value)
		case "80":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE80.Set(metric.Value)
		case "160":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE160.Set(metric.Value)
		case "320":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE320.Set(metric.Value)
		case "640":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE640.Set(metric.Value)
		case "1280":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE1280.Set(metric.Value)
		case "2560":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE2560.Set(metric.Value)
		case "5120":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE5120.Set(metric.Value)
		case "10240":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE10240.Set(metric.Value)
		case "20480":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE20480.Set(metric.Value)
		case "40960":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LE40960.Set(metric.Value)
		case "+Inf":
			mx.APIServer.Storage.DataKeyGeneration.Latencies.LEInf.Set(metric.Value)
		}
	}
}

func (k *Kubelet) collectRESTClientHTTPRequests(raw prometheus.Metrics, mx *metrics) {
	metricName := "rest_client_requests_total"
	chart := k.charts.Get("rest_client_requests_by_code")

	for _, metric := range raw.FindByName(metricName) {
		code := metric.Labels.Get("code")
		if code == "" {
			continue
		}
		dimID := "rest_client_requests_" + code
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: code, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.RESTClient.Requests.ByStatusCode[code] = mtx.Gauge(metric.Value)
	}

	chart = k.charts.Get("rest_client_requests_by_method")

	for _, metric := range raw.FindByName(metricName) {
		method := metric.Labels.Get("method")
		if method == "" {
			continue
		}
		dimID := "rest_client_requests_" + method
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: method, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.RESTClient.Requests.ByMethod[method] = mtx.Gauge(metric.Value)
	}
}

func (k *Kubelet) collectRuntimeOperations(raw prometheus.Metrics, mx *metrics) {
	chart := k.charts.Get("kubelet_runtime_operations")

	for _, metric := range raw.FindByName("kubelet_runtime_operations") {
		opType := metric.Labels.Get("operation_type")
		if opType == "" {
			continue
		}
		dimID := "kubelet_runtime_operations_" + opType
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: opType, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.Kubelet.Runtime.Operations[opType] = mtx.Gauge(metric.Value)
	}
}

func (k *Kubelet) collectRuntimeOperationsErrors(raw prometheus.Metrics, mx *metrics) {
	chart := k.charts.Get("kubelet_runtime_operations_errors")

	for _, metric := range raw.FindByName("kubelet_runtime_operations_errors") {
		opType := metric.Labels.Get("operation_type")
		if opType == "" {
			continue
		}
		dimID := "kubelet_runtime_operations_errors_" + opType
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: opType, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.Kubelet.Runtime.OperationsErrors[opType] = mtx.Gauge(metric.Value)
	}
}

func (k *Kubelet) collectDockerOperations(raw prometheus.Metrics, mx *metrics) {
	chart := k.charts.Get("kubelet_docker_operations")

	for _, metric := range raw.FindByName("kubelet_docker_operations") {
		opType := metric.Labels.Get("operation_type")
		if opType == "" {
			continue
		}
		dimID := "kubelet_docker_operations_" + opType
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: opType, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.Kubelet.Docker.Operations[opType] = mtx.Gauge(metric.Value)
	}
}

func (k *Kubelet) collectDockerOperationsErrors(raw prometheus.Metrics, mx *metrics) {
	chart := k.charts.Get("kubelet_docker_operations_errors")

	for _, metric := range raw.FindByName("kubelet_docker_operations_errors") {
		opType := metric.Labels.Get("operation_type")
		if opType == "" {
			continue
		}
		dimID := "kubelet_docker_operations_errors_" + opType
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: opType, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.Kubelet.Docker.OperationsErrors[opType] = mtx.Gauge(metric.Value)
	}
}
