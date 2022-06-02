package k8s_state

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestKubeState_Init(t *testing.T) {

}

func TestKubeState_Charts(t *testing.T) {

}

func TestKubeState_Cleanup(t *testing.T) {

}

func TestKubeState_Collect(t *testing.T) {
	type (
		testCaseStep func(t *testing.T, ks *KubeState)
		testCase     struct {
			client kubernetes.Interface
			steps  []testCaseStep
		}
	)

	tests := map[string]struct {
		create func(t *testing.T) testCase
	}{
		"only Node": {
			create: func(t *testing.T) testCase {
				client := fake.NewSimpleClientset(
					newNode("node01"),
				)

				step1 := func(t *testing.T, ks *KubeState) {
					mx := ks.Collect()
					expected := map[string]int64{
						"discovery_node_discoverer_state":              1,
						"discovery_pod_discoverer_state":               1,
						"node_node01_age":                              3,
						"node_node01_alloc_cpu_limits_used":            0,
						"node_node01_alloc_cpu_limits_util":            0,
						"node_node01_alloc_cpu_requests_used":          0,
						"node_node01_alloc_cpu_requests_util":          0,
						"node_node01_alloc_mem_limits_used":            0,
						"node_node01_alloc_mem_limits_util":            0,
						"node_node01_alloc_mem_requests_used":          0,
						"node_node01_alloc_mem_requests_util":          0,
						"node_node01_alloc_pods_allocated":             0,
						"node_node01_alloc_pods_available":             110,
						"node_node01_alloc_pods_used":                  0,
						"node_node01_cond_diskpressure":                0,
						"node_node01_cond_memorypressure":              0,
						"node_node01_cond_networkunavailable":          0,
						"node_node01_cond_pidpressure":                 0,
						"node_node01_cond_ready":                       1,
						"node_node01_containers":                       0,
						"node_node01_containers_state_running":         0,
						"node_node01_containers_state_terminated":      0,
						"node_node01_containers_state_waiting":         0,
						"node_node01_init_containers":                  0,
						"node_node01_init_containers_state_running":    0,
						"node_node01_init_containers_state_terminated": 0,
						"node_node01_init_containers_state_waiting":    0,
						"node_node01_pods_cond_containersready":        0,
						"node_node01_pods_cond_podinitialized":         0,
						"node_node01_pods_cond_podready":               0,
						"node_node01_pods_cond_podscheduled":           0,
						"node_node01_pods_phase_failed":                0,
						"node_node01_pods_phase_pending":               0,
						"node_node01_pods_phase_running":               0,
						"node_node01_pods_phase_succeeded":             0,
						"node_node01_pods_readiness":                   0,
						"node_node01_pods_readiness_ready":             0,
						"node_node01_pods_readiness_unready":           0,
					}
					copyAge(expected, mx)
					assert.Equal(t, expected, mx)
					assert.Equal(t,
						len(nodeChartsTmpl)+len(baseCharts),
						len(*ks.Charts()),
					)
				}

				return testCase{
					client: client,
					steps:  []testCaseStep{step1},
				}
			},
		},
		"only Pod": {
			create: func(t *testing.T) testCase {
				pod := newPod("node01", "pod01")
				client := fake.NewSimpleClientset(
					pod,
				)

				step1 := func(t *testing.T, ks *KubeState) {
					mx := ks.Collect()
					expected := map[string]int64{
						"discovery_node_discoverer_state":                         1,
						"discovery_pod_discoverer_state":                          1,
						"pod_default_pod01_age":                                   3,
						"pod_default_pod01_alloc_cpu_limits":                      0,
						"pod_default_pod01_alloc_cpu_limits_used":                 400,
						"pod_default_pod01_alloc_cpu_requests":                    0,
						"pod_default_pod01_alloc_cpu_requests_used":               200,
						"pod_default_pod01_alloc_mem_limits":                      0,
						"pod_default_pod01_alloc_mem_limits_used":                 419430400,
						"pod_default_pod01_alloc_mem_requests":                    0,
						"pod_default_pod01_alloc_mem_requests_used":               209715200,
						"pod_default_pod01_cond_containersready":                  1,
						"pod_default_pod01_cond_podinitialized":                   1,
						"pod_default_pod01_cond_podready":                         1,
						"pod_default_pod01_cond_podscheduled":                     1,
						"pod_default_pod01_container_container1_readiness":        1,
						"pod_default_pod01_container_container1_restarts":         0,
						"pod_default_pod01_container_container1_state_running":    1,
						"pod_default_pod01_container_container1_state_terminated": 0,
						"pod_default_pod01_container_container1_state_waiting":    0,
						"pod_default_pod01_container_container2_readiness":        1,
						"pod_default_pod01_container_container2_restarts":         0,
						"pod_default_pod01_container_container2_state_running":    1,
						"pod_default_pod01_container_container2_state_terminated": 0,
						"pod_default_pod01_container_container2_state_waiting":    0,
						"pod_default_pod01_containers":                            2,
						"pod_default_pod01_containers_state_running":              2,
						"pod_default_pod01_containers_state_terminated":           0,
						"pod_default_pod01_containers_state_waiting":              0,
						"pod_default_pod01_init_containers":                       1,
						"pod_default_pod01_init_containers_state_running":         0,
						"pod_default_pod01_init_containers_state_terminated":      1,
						"pod_default_pod01_init_containers_state_waiting":         0,
						"pod_default_pod01_phase_failed":                          0,
						"pod_default_pod01_phase_pending":                         0,
						"pod_default_pod01_phase_running":                         1,
						"pod_default_pod01_phase_succeeded":                       0,
						"pod_default_pod01_readiness_ready":                       1,
					}
					copyAge(expected, mx)
					assert.Equal(t, expected, mx)
					assert.Equal(t,
						len(podChartsTmpl)+len(containerChartsTmpl)*len(pod.Spec.Containers)+len(baseCharts),
						len(*ks.Charts()),
					)
				}

				return testCase{
					client: client,
					steps:  []testCaseStep{step1},
				}
			},
		},
		"Nodes and Pods": {
			create: func(t *testing.T) testCase {
				node := newNode("node01")
				pod := newPod(node.Name, "pod01")
				client := fake.NewSimpleClientset(
					node,
					pod,
				)

				step1 := func(t *testing.T, ks *KubeState) {
					mx := ks.Collect()
					expected := map[string]int64{
						"discovery_node_discoverer_state":                         1,
						"discovery_pod_discoverer_state":                          1,
						"node_node01_age":                                         3,
						"node_node01_alloc_cpu_limits_used":                       400,
						"node_node01_alloc_cpu_limits_util":                       11428,
						"node_node01_alloc_cpu_requests_used":                     200,
						"node_node01_alloc_cpu_requests_util":                     5714,
						"node_node01_alloc_mem_limits_used":                       419430400,
						"node_node01_alloc_mem_limits_util":                       11428,
						"node_node01_alloc_mem_requests_used":                     209715200,
						"node_node01_alloc_mem_requests_util":                     5714,
						"node_node01_alloc_pods_allocated":                        1,
						"node_node01_alloc_pods_available":                        109,
						"node_node01_alloc_pods_used":                             909,
						"node_node01_cond_diskpressure":                           0,
						"node_node01_cond_memorypressure":                         0,
						"node_node01_cond_networkunavailable":                     0,
						"node_node01_cond_pidpressure":                            0,
						"node_node01_cond_ready":                                  1,
						"node_node01_containers":                                  2,
						"node_node01_containers_state_running":                    2,
						"node_node01_containers_state_terminated":                 0,
						"node_node01_containers_state_waiting":                    0,
						"node_node01_init_containers":                             1,
						"node_node01_init_containers_state_running":               0,
						"node_node01_init_containers_state_terminated":            1,
						"node_node01_init_containers_state_waiting":               0,
						"node_node01_pods_cond_containersready":                   1,
						"node_node01_pods_cond_podinitialized":                    1,
						"node_node01_pods_cond_podready":                          1,
						"node_node01_pods_cond_podscheduled":                      1,
						"node_node01_pods_phase_failed":                           0,
						"node_node01_pods_phase_pending":                          0,
						"node_node01_pods_phase_running":                          1,
						"node_node01_pods_phase_succeeded":                        0,
						"node_node01_pods_readiness":                              100000,
						"node_node01_pods_readiness_ready":                        1,
						"node_node01_pods_readiness_unready":                      0,
						"pod_default_pod01_age":                                   3,
						"pod_default_pod01_alloc_cpu_limits":                      11428,
						"pod_default_pod01_alloc_cpu_limits_used":                 400,
						"pod_default_pod01_alloc_cpu_requests":                    5714,
						"pod_default_pod01_alloc_cpu_requests_used":               200,
						"pod_default_pod01_alloc_mem_limits":                      11428,
						"pod_default_pod01_alloc_mem_limits_used":                 419430400,
						"pod_default_pod01_alloc_mem_requests":                    5714,
						"pod_default_pod01_alloc_mem_requests_used":               209715200,
						"pod_default_pod01_cond_containersready":                  1,
						"pod_default_pod01_cond_podinitialized":                   1,
						"pod_default_pod01_cond_podready":                         1,
						"pod_default_pod01_cond_podscheduled":                     1,
						"pod_default_pod01_container_container1_readiness":        1,
						"pod_default_pod01_container_container1_restarts":         0,
						"pod_default_pod01_container_container1_state_running":    1,
						"pod_default_pod01_container_container1_state_terminated": 0,
						"pod_default_pod01_container_container1_state_waiting":    0,
						"pod_default_pod01_container_container2_readiness":        1,
						"pod_default_pod01_container_container2_restarts":         0,
						"pod_default_pod01_container_container2_state_running":    1,
						"pod_default_pod01_container_container2_state_terminated": 0,
						"pod_default_pod01_container_container2_state_waiting":    0,
						"pod_default_pod01_containers":                            2,
						"pod_default_pod01_containers_state_running":              2,
						"pod_default_pod01_containers_state_terminated":           0,
						"pod_default_pod01_containers_state_waiting":              0,
						"pod_default_pod01_init_containers":                       1,
						"pod_default_pod01_init_containers_state_running":         0,
						"pod_default_pod01_init_containers_state_terminated":      1,
						"pod_default_pod01_init_containers_state_waiting":         0,
						"pod_default_pod01_phase_failed":                          0,
						"pod_default_pod01_phase_pending":                         0,
						"pod_default_pod01_phase_running":                         1,
						"pod_default_pod01_phase_succeeded":                       0,
						"pod_default_pod01_readiness_ready":                       1,
					}
					copyAge(expected, mx)
					assert.Equal(t, expected, mx)
					assert.Equal(t,
						len(nodeChartsTmpl)+len(podChartsTmpl)+len(containerChartsTmpl)*len(pod.Spec.Containers)+len(baseCharts),
						len(*ks.Charts()),
					)
				}

				return testCase{
					client: client,
					steps:  []testCaseStep{step1},
				}
			},
		},
	}

	for name, creator := range tests {
		t.Run(name, func(t *testing.T) {
			test := creator.create(t)

			ks := New()
			ks.newKubeClient = func() (kubernetes.Interface, error) { return test.client, nil }

			require.True(t, ks.Init())
			require.True(t, ks.Check())
			defer ks.Cleanup()

			for i, executeStep := range test.steps {
				if i == 0 {
					_ = ks.Collect()
					time.Sleep(ks.initDelay)
				}
				executeStep(t, ks)
			}
		})
	}
}

func newNode(name string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Status: corev1.NodeStatus{
			Capacity: corev1.ResourceList{
				corev1.ResourceCPU:    mustQuantity("4000m"),
				corev1.ResourceMemory: mustQuantity("4000Mi"),
				"pods":                mustQuantity("110"),
			},
			Allocatable: corev1.ResourceList{
				corev1.ResourceCPU:    mustQuantity("3500m"),
				corev1.ResourceMemory: mustQuantity("3500Mi"),
				"pods":                mustQuantity("110"),
			},
			Conditions: []corev1.NodeCondition{
				{Type: corev1.NodeReady, Status: corev1.ConditionTrue},
				{Type: corev1.NodeMemoryPressure, Status: corev1.ConditionFalse},
				{Type: corev1.NodeDiskPressure, Status: corev1.ConditionFalse},
				{Type: corev1.NodePIDPressure, Status: corev1.ConditionFalse},
				{Type: corev1.NodeNetworkUnavailable, Status: corev1.ConditionFalse},
			},
		},
	}
}

func newPod(nodeName, name string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         corev1.NamespaceDefault,
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Spec: corev1.PodSpec{
			NodeName: nodeName,
			InitContainers: []corev1.Container{
				{
					Name: "init-container1",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("50m"),
							corev1.ResourceMemory: mustQuantity("50Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("10m"),
							corev1.ResourceMemory: mustQuantity("10Mi"),
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name: "container1",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("200m"),
							corev1.ResourceMemory: mustQuantity("200Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("100m"),
							corev1.ResourceMemory: mustQuantity("100Mi"),
						},
					},
				},
				{
					Name: "container2",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("200m"),
							corev1.ResourceMemory: mustQuantity("200Mi")},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    mustQuantity("100m"),
							corev1.ResourceMemory: mustQuantity("100Mi"),
						},
					},
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{Type: corev1.PodReady, Status: corev1.ConditionTrue},
				{Type: corev1.PodScheduled, Status: corev1.ConditionTrue},
				{Type: corev1.PodInitialized, Status: corev1.ConditionTrue},
				{Type: corev1.ContainersReady, Status: corev1.ConditionTrue},
			},
			InitContainerStatuses: []corev1.ContainerStatus{
				{
					Name:  "init-container1",
					State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{}},
				},
			},
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:  "container1",
					Ready: true,
					State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
				},
				{
					Name:  "container2",
					Ready: true,
					State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
				},
			},
		},
	}
}

func mustQuantity(s string) apiresource.Quantity {
	q, err := apiresource.ParseQuantity(s)
	if err != nil {
		panic(fmt.Sprintf("fail to create resource quantity: %v", err))
	}
	return q
}

func copyAge(dst, src map[string]int64) {
	for k, v := range src {
		if !strings.HasSuffix(k, "_age") {
			continue
		}
		if _, ok := dst[k]; ok {
			dst[k] = v
		}
	}
}

/*
	m := ks.Collect()
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	for _, value := range l {
		fmt.Println(fmt.Sprintf("\"%s\": %d,", value, m[value]))
	}
*/
