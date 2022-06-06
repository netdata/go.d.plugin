package k8s_state

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

// NETDATA_CHART_PRIO_CGROUPS_CONTAINERS 40000
const prioDiscoveryDiscovererState = 39999

const (
	prioNodeAllocatableCPUUtil = 38100 + iota
	prioNodeAllocatableCPUUsed
	prioNodeAllocatableMemUtil
	prioNodeAllocatableMemUsed
	prioNodeAllocatablePodsUtil
	prioNodeAllocatablePodsUsage
	prioNodeConditions
	prioNodePodsReadiness
	prioNodePodsReadinessState
	prioNodePodsCondition
	prioNodePodsPhase
	prioNodeContainersCount
	prioNodeContainersState
	prioNodeInitContainersState
	prioNodeAge
)

const (
	prioPodAllocatedCPU = 38200 + iota
	prioPodAllocatedCPUUsed
	prioPodAllocatedMem
	prioPodAllocatedMemUsed
	prioPodCondition
	prioPodPhase
	prioPodAge
	prioPodContainersCount
	prioPodContainersState
	prioPodInitContainersState
	prioPodContainerReadinessState
	prioPodContainerRestarts
	prioPodContainerState
	prioPodContainerWaitingStateReason
	prioPodContainerTerminatedStateReason
)

const (
	labelKeyPrefix         = "k8s_"
	labelKeyClusterID      = labelKeyPrefix + "cluster_id"
	labelKeyClusterName    = labelKeyPrefix + "cluster_name"
	labelKeyNamespace      = labelKeyPrefix + "namespace"
	labelKeyKind           = labelKeyPrefix + "kind"
	labelKeyPodName        = labelKeyPrefix + "pod_name"
	labelKeyNodeName       = labelKeyPrefix + "node_name"
	labelKeyPodUID         = labelKeyPrefix + "pod_uid"
	labelKeyControllerKind = labelKeyPrefix + "controller_kind"
	labelKeyControllerName = labelKeyPrefix + "controller_name"
	labelKeyContainerName  = labelKeyPrefix + "container_name"
	labelKeyContainerID    = labelKeyPrefix + "container_id"
	labelKeyQoSClass       = labelKeyPrefix + "qos_class"
)

var baseCharts = module.Charts{
	discoveryStatusChart.Copy(),
}

var nodeChartsTmpl = module.Charts{
	nodeAllocatableCPUUtilizationChartTmpl.Copy(),
	nodeAllocatableCPUUsedChartTmpl.Copy(),
	nodeAllocatableMemUtilizationChartTmpl.Copy(),
	nodeAllocMemUsedChartTmpl.Copy(),
	nodeAllocatablePodsUtilizationChartTmpl.Copy(),
	nodeAllocatablePodsUsageChartTmpl.Copy(),
	nodeConditionsChartTmpl.Copy(),
	nodePodsReadinessChartTmpl.Copy(),
	nodePodsReadinessStateChartTmpl.Copy(),
	nodePodsConditionChartTmpl.Copy(),
	nodePodsPhaseChartTmpl.Copy(),
	nodeContainersChartTmpl.Copy(),
	nodeContainersStateChartTmpl.Copy(),
	nodeInitContainersStateChartTmpl.Copy(),
	nodeAgeChartTmpl.Copy(),
}

var podChartsTmpl = module.Charts{
	podAllocatedCPUChartTmpl.Copy(),
	podAllocatedCPUUsedChartTmpl.Copy(),
	podAllocatedMemChartTmpl.Copy(),
	podAllocatedMemUsedChartTmpl.Copy(),
	podConditionChartTmpl.Copy(),
	podPhaseChartTmpl.Copy(),
	podAgeChartTmpl.Copy(),
	podContainersCountChartTmpl.Copy(),
	podContainersStateChartTmpl.Copy(),
	podInitContainersStateChartTmpl.Copy(),
}

var containerChartsTmpl = module.Charts{
	containerReadinessStateChartTmpl.Copy(),
	containerRestartsChartTmpl.Copy(),
	containersStateChartTmpl.Copy(),
	containersStateWaitingChartTmpl.Copy(),
	containersStateTerminatedChartTmpl.Copy(),
}

var (
	// CPU resource
	nodeAllocatableCPUUtilizationChartTmpl = module.Chart{
		ID:       "node_%s.allocatable_cpu_utilization",
		Title:    "CPU resource utilization",
		Units:    "%",
		Fam:      "node cpu resource",
		Ctx:      "k8s_state.node_allocatable_cpu_utilization",
		Priority: prioNodeAllocatableCPUUtil,
		Dims: module.Dims{
			{ID: "node_%s_alloc_cpu_requests_util", Name: "requests", Div: precision},
			{ID: "node_%s_alloc_cpu_limits_util", Name: "limits", Div: precision},
		},
	}
	nodeAllocatableCPUUsedChartTmpl = module.Chart{
		ID:       "node_%s.allocatable_cpu_used",
		Title:    "CPU resource used",
		Units:    "millicpu",
		Fam:      "node cpu resource",
		Ctx:      "k8s_state.node_allocatable_cpu_used",
		Priority: prioNodeAllocatableCPUUsed,
		Dims: module.Dims{
			{ID: "node_%s_alloc_cpu_requests_used", Name: "requests"},
			{ID: "node_%s_alloc_cpu_limits_used", Name: "limits"},
		},
	}
	// memory resource
	nodeAllocatableMemUtilizationChartTmpl = module.Chart{
		ID:       "node_%s.allocatable_mem_utilization",
		Title:    "Memory resource utilization",
		Units:    "%",
		Fam:      "node mem resource",
		Ctx:      "k8s_state.node_allocatable_mem_utilization",
		Priority: prioNodeAllocatableMemUtil,
		Dims: module.Dims{
			{ID: "node_%s_alloc_mem_requests_util", Name: "requests", Div: precision},
			{ID: "node_%s_alloc_mem_limits_util", Name: "limits", Div: precision},
		},
	}
	nodeAllocMemUsedChartTmpl = module.Chart{
		ID:       "node_%s.allocatable_mem_used",
		Title:    "Memory resource used",
		Units:    "bytes",
		Fam:      "node mem resource",
		Ctx:      "k8s_state.node_allocatable_mem_used",
		Priority: prioNodeAllocatableMemUsed,
		Dims: module.Dims{
			{ID: "node_%s_alloc_mem_requests_used", Name: "requests"},
			{ID: "node_%s_alloc_mem_limits_used", Name: "limits"},
		},
	}
	// pods resource
	nodeAllocatablePodsUtilizationChartTmpl = module.Chart{
		ID:       "node_%s.allocatable_pods_utilization",
		Title:    "Pods resource utilization",
		Units:    "%",
		Fam:      "node pods resource",
		Ctx:      "k8s_state.node_allocatable_pods_utilization",
		Priority: prioNodeAllocatablePodsUtil,
		Dims: module.Dims{
			{ID: "node_%s_alloc_pods_used", Name: "allocated", Div: precision},
		},
	}
	nodeAllocatablePodsUsageChartTmpl = module.Chart{
		ID:       "node_%s.allocated_pods_usage",
		Title:    "Pods resource usage",
		Units:    "pods",
		Fam:      "node pods resource",
		Ctx:      "k8s_state.node_allocatable_pods_usage",
		Type:     module.Stacked,
		Priority: prioNodeAllocatablePodsUsage,
		Dims: module.Dims{
			{ID: "node_%s_alloc_pods_available", Name: "available"},
			{ID: "node_%s_alloc_pods_allocated", Name: "allocated"},
		},
	}
	// condition
	nodeConditionsChartTmpl = module.Chart{
		ID:       "node_%s.condition_status",
		Title:    "Condition status",
		Units:    "status",
		Fam:      "node condition",
		Ctx:      "k8s_state.node_condition",
		Priority: prioNodeConditions,
	}
	// pods readiness
	nodePodsReadinessChartTmpl = module.Chart{
		ID:       "node_%s.pods_readiness",
		Title:    "Pods readiness",
		Units:    "%",
		Fam:      "node pods readiness",
		Ctx:      "k8s_state.node_pods_readiness",
		Priority: prioNodePodsReadiness,
		Dims: module.Dims{
			{ID: "node_%s_pods_readiness", Name: "ready", Div: precision},
		},
	}
	nodePodsReadinessStateChartTmpl = module.Chart{
		ID:       "node_%s.pods_readiness_state",
		Title:    "Pods readiness state",
		Units:    "pods",
		Fam:      "node pods readiness",
		Ctx:      "k8s_state.node_pods_readiness_state",
		Type:     module.Stacked,
		Priority: prioNodePodsReadinessState,
		Dims: module.Dims{
			{ID: "node_%s_pods_readiness_ready", Name: "ready"},
			{ID: "node_%s_pods_readiness_unready", Name: "unready"},
		},
	}
	// pods condition
	nodePodsConditionChartTmpl = module.Chart{
		ID:       "node_%s.pods_condition",
		Title:    "Pods condition",
		Units:    "pods",
		Fam:      "node pods condition",
		Ctx:      "k8s_state.node_pods_condition",
		Priority: prioNodePodsCondition,
		Dims: module.Dims{
			{ID: "node_%s_pods_cond_podready", Name: "PodReady"},
			{ID: "node_%s_pods_cond_podscheduled", Name: "PodScheduled"},
			{ID: "node_%s_pods_cond_podinitialized", Name: "PodInitialized"},
			{ID: "node_%s_pods_cond_containersready", Name: "ContainersReady"},
		},
	}
	// pods phase
	nodePodsPhaseChartTmpl = module.Chart{
		ID:       "node_%s.pods_phase",
		Title:    "Pods phase",
		Units:    "pods",
		Fam:      "node pods phase",
		Ctx:      "k8s_state.node_pods_phase",
		Type:     module.Stacked,
		Priority: prioNodePodsPhase,
		Dims: module.Dims{
			{ID: "node_%s_pods_phase_running", Name: "Running"},
			{ID: "node_%s_pods_phase_failed", Name: "Failed"},
			{ID: "node_%s_pods_phase_succeeded", Name: "Succeeded"},
			{ID: "node_%s_pods_phase_pending", Name: "Pending"},
		},
	}
	// containers
	nodeContainersChartTmpl = module.Chart{
		ID:       "node_%s.containers",
		Title:    "Containers",
		Units:    "containers",
		Fam:      "node containers",
		Ctx:      "k8s_state.node_containers",
		Priority: prioNodeContainersCount,
		Dims: module.Dims{
			{ID: "node_%s_containers", Name: "containers"},
			{ID: "node_%s_init_containers", Name: "init_containers"},
		},
	}
	nodeContainersStateChartTmpl = module.Chart{
		ID:       "node_%s.containers_state",
		Title:    "Containers state",
		Units:    "containers",
		Fam:      "node containers",
		Ctx:      "k8s_state.node_containers_state",
		Type:     module.Stacked,
		Priority: prioNodeContainersState,
		Dims: module.Dims{
			{ID: "node_%s_containers_state_running", Name: "Running"},
			{ID: "node_%s_containers_state_waiting", Name: "Waiting"},
			{ID: "node_%s_containers_state_terminated", Name: "Terminated"},
		},
	}
	nodeInitContainersStateChartTmpl = module.Chart{
		ID:       "node_%s.init_containers_state",
		Title:    "Init containers state",
		Units:    "containers",
		Fam:      "node containers",
		Ctx:      "k8s_state.node_init_containers_state",
		Type:     module.Stacked,
		Priority: prioNodeInitContainersState,
		Dims: module.Dims{
			{ID: "node_%s_init_containers_state_running", Name: "Running"},
			{ID: "node_%s_init_containers_state_waiting", Name: "Waiting"},
			{ID: "node_%s_init_containers_state_terminated", Name: "Terminated"},
		},
	}
	// age
	nodeAgeChartTmpl = module.Chart{
		ID:       "node_%s.age",
		Title:    "Age",
		Units:    "seconds",
		Fam:      "node age",
		Ctx:      "k8s_state.node_age",
		Priority: prioNodeAge,
		Dims: module.Dims{
			{ID: "node_%s_age", Name: "age"},
		},
	}
)

func (ks *KubeState) newNodeCharts(ns *nodeState) *module.Charts {
	cs := nodeChartsTmpl.Copy()
	for _, c := range *cs {
		c.ID = fmt.Sprintf(c.ID, replaceDots(ns.id()))
		c.Labels = []module.Label{
			{Key: labelKeyKind, Value: "node", Source: module.LabelSourceK8s},
			{Key: labelKeyClusterID, Value: ks.kubeClusterID, Source: module.LabelSourceK8s},
			{Key: labelKeyClusterName, Value: ks.kubeClusterName, Source: module.LabelSourceK8s},
		}
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, ns.id())
		}
	}
	return cs
}

func (ks *KubeState) addNodeCharts(ns *nodeState) {
	cs := ks.newNodeCharts(ns)
	if err := ks.Charts().Add(*cs...); err != nil {
		ks.Warning(err)
	}
}

func (ks *KubeState) removeNodeCharts(ns *nodeState) {
	prefix := fmt.Sprintf("node_%s", replaceDots(ns.id()))
	for _, c := range *ks.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

func (ks *KubeState) addNodeConditionToCharts(ns *nodeState, cond string) {
	id := fmt.Sprintf(nodeConditionsChartTmpl.ID, replaceDots(ns.id()))
	c := ks.Charts().Get(id)
	if c == nil {
		ks.Warningf("chart '%s' does not exist", id)
		return
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("node_%s_cond_%s", ns.id(), strings.ToLower(cond)),
		Name: cond,
	}
	if err := c.AddDim(dim); err != nil {
		ks.Warning(err)
		return
	}
	c.MarkNotCreated()
}

var (
	podAllocatedCPUChartTmpl = module.Chart{
		ID:       "pod_%s.allocated_cpu",
		Title:    "Allocated CPU",
		Units:    "%",
		Fam:      "pod allocated cpu",
		Ctx:      "k8s_state.pod_allocated_cpu",
		Priority: prioPodAllocatedCPU,
		Dims: module.Dims{
			{ID: "pod_%s_alloc_cpu_requests", Name: "requests", Div: precision},
			{ID: "pod_%s_alloc_cpu_limits", Name: "limits", Div: precision},
		},
	}
	podAllocatedCPUUsedChartTmpl = module.Chart{
		ID:       "pod_%s.allocated_cpu_used",
		Title:    "Allocated CPU used",
		Units:    "millicpu",
		Fam:      "pod  allocated cpu",
		Ctx:      "k8s_state.pod_allocated_cpu_used",
		Priority: prioPodAllocatedCPUUsed,
		Dims: module.Dims{
			{ID: "pod_%s_alloc_cpu_requests_used", Name: "requests"},
			{ID: "pod_%s_alloc_cpu_limits_used", Name: "limits"},
		},
	}

	podAllocatedMemChartTmpl = module.Chart{
		ID:       "pod_%s.allocated_mem",
		Title:    "Allocated memory",
		Units:    "%",
		Fam:      "pod allocated mem",
		Ctx:      "k8s_state.pod_allocated_mem",
		Priority: prioPodAllocatedMem,
		Dims: module.Dims{
			{ID: "pod_%s_alloc_mem_requests", Name: "requests", Div: precision},
			{ID: "pod_%s_alloc_mem_limits", Name: "limits", Div: precision},
		},
	}
	podAllocatedMemUsedChartTmpl = module.Chart{
		ID:       "pod_%s.allocated_mem_used",
		Title:    "Allocated memory used",
		Units:    "bytes",
		Fam:      "pod allocated mem",
		Ctx:      "k8s_state.pod_allocated_mem_used",
		Priority: prioPodAllocatedMemUsed,
		Dims: module.Dims{
			{ID: "pod_%s_alloc_mem_requests_used", Name: "requests"},
			{ID: "pod_%s_alloc_mem_limits_used", Name: "limits"},
		},
	}
	podConditionChartTmpl = module.Chart{
		ID:       "pod_%s.condition",
		Title:    "Condition",
		Units:    "state",
		Fam:      "pod condition",
		Ctx:      "k8s_state.pod_condition",
		Priority: prioPodCondition,
		Dims: module.Dims{
			{ID: "pod_%s_cond_podready", Name: "PodReady"},
			{ID: "pod_%s_cond_podscheduled", Name: "PodScheduled"},
			{ID: "pod_%s_cond_podinitialized", Name: "PodInitialized"},
			{ID: "pod_%s_cond_containersready", Name: "ContainersReady"},
		},
	}
	podPhaseChartTmpl = module.Chart{
		ID:       "pod_%s.phase",
		Title:    "Phase",
		Units:    "state",
		Fam:      "pod phase",
		Ctx:      "k8s_state.pod_phase",
		Priority: prioPodPhase,
		Dims: module.Dims{
			{ID: "pod_%s_phase_running", Name: "Running"},
			{ID: "pod_%s_phase_failed", Name: "Failed"},
			{ID: "pod_%s_phase_succeeded", Name: "Succeeded"},
			{ID: "pod_%s_phase_pending", Name: "Pending"},
		},
	}
	podAgeChartTmpl = module.Chart{
		ID:       "pod_%s.age",
		Title:    "Age",
		Units:    "seconds",
		Fam:      "pod age",
		Ctx:      "k8s_state.pod_age",
		Priority: prioPodAge,
		Dims: module.Dims{
			{ID: "pod_%s_age", Name: "age"},
		},
	}
	podContainersCountChartTmpl = module.Chart{
		ID:       "pod_%s.containers_count",
		Title:    "Containers",
		Units:    "containers",
		Fam:      "pod containers",
		Ctx:      "k8s_state.pod_containers",
		Priority: prioPodContainersCount,
		Dims: module.Dims{
			{ID: "pod_%s_containers", Name: "containers"},
			{ID: "pod_%s_init_containers", Name: "init_containers"},
		},
	}
	podContainersStateChartTmpl = module.Chart{
		ID:       "pod_%s.containers_state",
		Title:    "Containers state",
		Units:    "containers",
		Fam:      "pod containers",
		Ctx:      "k8s_state.pod_containers_state",
		Type:     module.Stacked,
		Priority: prioPodContainersState,
		Dims: module.Dims{
			{ID: "pod_%s_containers_state_running", Name: "Running"},
			{ID: "pod_%s_containers_state_waiting", Name: "Waiting"},
			{ID: "pod_%s_containers_state_terminated", Name: "Terminated"},
		},
	}
	podInitContainersStateChartTmpl = module.Chart{
		ID:       "pod_%s.init_containers_state",
		Title:    "Init containers state",
		Units:    "containers",
		Fam:      "pod containers",
		Ctx:      "k8s_state.pod_init_containers_state",
		Type:     module.Stacked,
		Priority: prioPodInitContainersState,
		Dims: module.Dims{
			{ID: "pod_%s_init_containers_state_running", Name: "Running"},
			{ID: "pod_%s_init_containers_state_waiting", Name: "Waiting"},
			{ID: "pod_%s_init_containers_state_terminated", Name: "Terminated"},
		},
	}
)

func (ks *KubeState) newPodCharts(ps *podState) *module.Charts {
	charts := podChartsTmpl.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, replaceDots(ps.id()))
		c.Labels = ks.newPodChartLabels(ps)
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, ps.id())
		}
	}
	return charts
}

func (ks *KubeState) newPodChartLabels(ps *podState) []module.Label {
	labels := []module.Label{
		{Key: labelKeyNamespace, Value: ps.namespace, Source: module.LabelSourceK8s},
		{Key: labelKeyPodName, Value: ps.name, Source: module.LabelSourceK8s},
		{Key: labelKeyNodeName, Value: ps.nodeName, Source: module.LabelSourceK8s},
		{Key: labelKeyKind, Value: "pod", Source: module.LabelSourceK8s},
		{Key: labelKeyPodUID, Value: ps.uid, Source: module.LabelSourceK8s},
		{Key: labelKeyQoSClass, Value: ps.qosClass, Source: module.LabelSourceK8s},
		{Key: labelKeyControllerKind, Value: ps.controllerKind, Source: module.LabelSourceK8s},
		{Key: labelKeyControllerName, Value: ps.controllerName, Source: module.LabelSourceK8s},
		{Key: labelKeyClusterID, Value: ks.kubeClusterID, Source: module.LabelSourceK8s},
		{Key: labelKeyClusterName, Value: ks.kubeClusterName, Source: module.LabelSourceK8s},
	}
	return labels
}

func (ks *KubeState) addPodCharts(ps *podState) {
	charts := ks.newPodCharts(ps)
	if err := ks.Charts().Add(*charts...); err != nil {
		ks.Warning(err)
	}
}

func (ks *KubeState) removePodCharts(ps *podState) {
	prefix := fmt.Sprintf("pod_%s", replaceDots(ps.id()))
	for _, c := range *ks.Charts() {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

var (
	containerReadinessStateChartTmpl = module.Chart{
		ID:       "pod_%s.container_%s_readiness_state",
		Title:    "Readiness state",
		Units:    "state",
		Fam:      "%s",
		Ctx:      "k8s_state.pod_container_readiness_state",
		Priority: prioPodContainerReadinessState,
		Dims: module.Dims{
			{ID: "pod_%s_container_%s_readiness", Name: "ready"},
		},
	}
	containerRestartsChartTmpl = module.Chart{
		ID:       "pod_%s.container_%s_restarts",
		Title:    "Restarts",
		Units:    "restarts",
		Fam:      "%s",
		Ctx:      "k8s_state.pod_container_restarts",
		Priority: prioPodContainerRestarts,
		Dims: module.Dims{
			{ID: "pod_%s_container_%s_restarts", Name: "restarts"},
		},
	}
	containersStateChartTmpl = module.Chart{
		ID:       "pod_%s.container_%s_state",
		Title:    "Container state",
		Units:    "state",
		Fam:      "%s",
		Ctx:      "k8s_state.pod_container_state",
		Priority: prioPodContainerState,
		Dims: module.Dims{
			{ID: "pod_%s_container_%s_state_running", Name: "Running"},
			{ID: "pod_%s_container_%s_state_waiting", Name: "Waiting"},
			{ID: "pod_%s_container_%s_state_terminated", Name: "Terminated"},
		},
	}
	containersStateWaitingChartTmpl = module.Chart{
		ID:       "pod_%s.container_%s_state_waiting_reason",
		Title:    "Container waiting state reason",
		Units:    "state",
		Fam:      "%s",
		Ctx:      "k8s_state.pod_container_waiting_state_reason",
		Priority: prioPodContainerWaitingStateReason,
	}
	containersStateTerminatedChartTmpl = module.Chart{
		ID:       "pod_%s.container_%s_state_terminated_reason",
		Title:    "Container terminated state reason",
		Units:    "state",
		Fam:      "%s",
		Ctx:      "k8s_state.pod_container_terminated_state_reason",
		Priority: prioPodContainerTerminatedStateReason,
	}
)

func (ks *KubeState) newContainerCharts(ps *podState, cs *containerState) *module.Charts {
	charts := containerChartsTmpl.Copy()
	for _, c := range *charts {
		c.ID = fmt.Sprintf(c.ID, replaceDots(ps.id()), cs.name)
		c.Fam = fmt.Sprintf(c.Fam, cs.name)
		c.Labels = ks.newContainerChartLabels(ps, cs)
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, ps.id(), cs.name)
		}
	}
	return charts
}

func (ks *KubeState) newContainerChartLabels(ps *podState, cs *containerState) []module.Label {
	labels := ks.newPodChartLabels(ps)
	for i, v := range labels {
		if v.Key == labelKeyKind {
			labels[i].Value = "container"
			break
		}
	}
	labels = append(labels, []module.Label{
		{Key: labelKeyContainerName, Value: cs.name, Source: module.LabelSourceK8s},
		{Key: labelKeyContainerID, Value: cs.uid, Source: module.LabelSourceK8s},
	}...)
	return labels
}

func (ks *KubeState) addContainerCharts(ps *podState, cs *containerState) {
	charts := ks.newContainerCharts(ps, cs)
	if err := ks.Charts().Add(*charts...); err != nil {
		ks.Warning(err)
	}
}

func (ks *KubeState) addContainerWaitingStateReasonToChart(ps *podState, cs *containerState, reason string) {
	id := fmt.Sprintf(containersStateWaitingChartTmpl.ID, replaceDots(ps.id()), cs.name)
	c := ks.Charts().Get(id)
	if c == nil {
		ks.Warningf("chart '%s' does not exist", id)
		return
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("pod_%s_container_%s_state_waiting_reason_%s", ps.id(), cs.name, reason),
		Name: reason,
	}
	if err := c.AddDim(dim); err != nil {
		ks.Warning(err)
		return
	}
	c.MarkNotCreated()
}

func (ks *KubeState) addContainerTerminatedStateReasonToChart(ps *podState, cs *containerState, reason string) {
	id := fmt.Sprintf(containersStateTerminatedChartTmpl.ID, replaceDots(ps.id()), cs.name)
	c := ks.Charts().Get(id)
	if c == nil {
		ks.Warningf("chart '%s' does not exist", id)
		return
	}
	dim := &module.Dim{
		ID:   fmt.Sprintf("pod_%s_container_%s_state_terminated_reason_%s", ps.id(), cs.name, reason),
		Name: reason,
	}
	if err := c.AddDim(dim); err != nil {
		ks.Warning(err)
		return
	}
	c.MarkNotCreated()
}

var discoveryStatusChart = module.Chart{
	ID:       "discovery_discoverers_state",
	Title:    "Running discoverers state",
	Units:    "state",
	Fam:      "discovery",
	Ctx:      "k8s_state.discovery_discoverers_state",
	Priority: prioDiscoveryDiscovererState,
	Opts:     module.Opts{Hidden: true},
	Dims: module.Dims{
		{ID: "discovery_node_discoverer_state", Name: "node"},
		{ID: "discovery_pod_discoverer_state", Name: "pod"},
	},
}

var reDots = regexp.MustCompile(`\.`)

func replaceDots(v string) string {
	return reDots.ReplaceAllString(v, "-")
}
