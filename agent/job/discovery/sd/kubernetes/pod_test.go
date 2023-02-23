package kubernetes

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
)

func TestPodTargetGroup_Source(t *testing.T) {
	tests := map[string]struct {
		sim            func() discoverySim
		expectedSource []string
	}{
		"pods with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
				discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						preparePodGroup(httpd),
						preparePodGroup(nginx),
					},
				}
				return sim
			},
			expectedSource: []string{
				"k8s/pod/default/httpd-dd95c4d68-5bkwl",
				"k8s/pod/default/nginx-7cfd77469b-q6kxj",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sim := test.sim()
			var actual []string
			for _, group := range sim.run(t) {
				actual = append(actual, group.Source())
			}

			assert.Equal(t, test.expectedSource, actual)
		})
	}
}

func TestPodTargetGroup_Targets(t *testing.T) {
	tests := map[string]struct {
		sim                func() discoverySim
		expectedNumTargets int
	}{
		"pods with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
				discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						preparePodGroup(httpd),
						preparePodGroup(nginx),
					},
				}
				return sim
			},
			expectedNumTargets: 4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sim := test.sim()
			var actual int
			for _, group := range sim.run(t) {
				actual += len(group.Targets())
			}

			assert.Equal(t, test.expectedNumTargets, actual)
		})
	}
}

func TestPodTargetGroup_Hash(t *testing.T) {
	tests := map[string]struct {
		sim          func() discoverySim
		expectedHash []uint64
	}{
		"pods with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
				discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						preparePodGroup(httpd),
						preparePodGroup(nginx),
					},
				}
				return sim
			},
			expectedHash: []uint64{
				10984584459239076244,
				10480294460002508451,
				6511065520956605596,
				9821662463142050012,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sim := test.sim()
			var actual []uint64
			for _, group := range sim.run(t) {
				for _, tg := range group.Targets() {
					actual = append(actual, tg.Hash())
				}
			}

			assert.Equal(t, test.expectedHash, actual)
		})
	}
}

func TestPodTargetGroup_TUID(t *testing.T) {
	tests := map[string]struct {
		sim          func() discoverySim
		expectedTUID []string
	}{
		"pods with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
				discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						preparePodGroup(httpd),
						preparePodGroup(nginx),
					},
				}
				return sim
			},
			expectedTUID: []string{
				"default_httpd-dd95c4d68-5bkwl_httpd_tcp_80",
				"default_httpd-dd95c4d68-5bkwl_httpd_tcp_443",
				"default_nginx-7cfd77469b-q6kxj_nginx_tcp_80",
				"default_nginx-7cfd77469b-q6kxj_nginx_tcp_443",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sim := test.sim()
			var actual []string
			for _, group := range sim.run(t) {
				for _, tg := range group.Targets() {
					actual = append(actual, tg.TUID())
				}
			}

			assert.Equal(t, test.expectedTUID, actual)
		})
	}
}

func TestNewPodDiscovery(t *testing.T) {
	tests := map[string]struct {
		podInf    cache.SharedInformer
		cmapInf   cache.SharedInformer
		secretInf cache.SharedInformer
		wantPanic bool
	}{
		"valid informers": {
			podInf:    cache.NewSharedInformer(nil, &corev1.Pod{}, resyncPeriod),
			cmapInf:   cache.NewSharedInformer(nil, &corev1.ConfigMap{}, resyncPeriod),
			secretInf: cache.NewSharedInformer(nil, &corev1.Secret{}, resyncPeriod),
		},
		"nil informers": {wantPanic: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.wantPanic {
				assert.Panics(t, func() { NewPodDiscovery(nil, nil, nil) })
			} else {
				assert.IsType(t, &PodDiscovery{}, NewPodDiscovery(test.podInf, test.cmapInf, test.secretInf))
			}
		})
	}
}

func TestPodDiscovery_String(t *testing.T) {
	p := PodDiscovery{}
	assert.NotEmpty(t, p.String())
}

func TestPodDiscovery_Discover(t *testing.T) {
	tests := map[string]func() discoverySim{
		"ADD: pods exist before run": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroup(httpd),
					preparePodGroup(nginx),
				},
			}
			return sim
		},
		"ADD: pods exist before run and add after sync": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			discovery, clientset := prepareNamespaceAllDiscovery(RolePod, httpd)
			podClient := clientset.CoreV1().Pods("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					_, _ = podClient.Create(ctx, nginx, metav1.CreateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					preparePodGroup(httpd),
					preparePodGroup(nginx),
				},
			}
			return sim
		},
		"DELETE: remove pods after sync": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			discovery, clientset := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)
			podClient := clientset.CoreV1().Pods("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_ = podClient.Delete(ctx, httpd.Name, metav1.DeleteOptions{})
					_ = podClient.Delete(ctx, nginx.Name, metav1.DeleteOptions{})
				},
				expectedGroups: []model.TargetGroup{
					preparePodGroup(httpd),
					preparePodGroup(nginx),
					prepareEmptyPodGroup(httpd),
					prepareEmptyPodGroup(nginx),
				},
			}
			return sim
		},
		"DELETE,ADD: remove and add pods after sync": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			discovery, clientset := prepareNamespaceAllDiscovery(RolePod, httpd)
			podClient := clientset.CoreV1().Pods("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_ = podClient.Delete(ctx, httpd.Name, metav1.DeleteOptions{})
					_, _ = podClient.Create(ctx, nginx, metav1.CreateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					preparePodGroup(httpd),
					prepareEmptyPodGroup(httpd),
					preparePodGroup(nginx),
				},
			}
			return sim
		},
		"ADD: pods with empty PodIP": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			httpd.Status.PodIP = ""
			nginx.Status.PodIP = ""
			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					prepareEmptyPodGroup(httpd),
					prepareEmptyPodGroup(nginx),
				},
			}
			return sim
		},
		"UPDATE: set pods PodIP after sync": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			httpd.Status.PodIP = ""
			nginx.Status.PodIP = ""
			discovery, clientset := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)
			podClient := clientset.CoreV1().Pods("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_, _ = podClient.Update(ctx, prepareHTTPDPod(), metav1.UpdateOptions{})
					_, _ = podClient.Update(ctx, prepareNGINXPod(), metav1.UpdateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					prepareEmptyPodGroup(httpd),
					prepareEmptyPodGroup(nginx),
					preparePodGroup(prepareHTTPDPod()),
					preparePodGroup(prepareNGINXPod()),
				},
			}
			return sim
		},
		"ADD: pods without containers": func() discoverySim {
			httpd, nginx := prepareHTTPDPod(), prepareNGINXPod()
			httpd.Spec.Containers = httpd.Spec.Containers[:0]
			nginx.Spec.Containers = httpd.Spec.Containers[:0]
			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					prepareEmptyPodGroup(httpd),
					prepareEmptyPodGroup(nginx),
				},
			}
			return sim
		},
		"Env: from value": func() discoverySim {
			httpd := prepareHTTPDPod()
			mangle := func(c *corev1.Container) {
				c.Env = []corev1.EnvVar{
					{Name: "key1", Value: "value1"},
				}
			}
			mangleContainers(httpd.Spec.Containers, mangle)
			data := map[string]string{"key1": "value1"}

			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroupWithEnv(httpd, data),
				},
			}
			return sim
		},
		"Env: from Secret": func() discoverySim {
			httpd := prepareHTTPDPod()
			mangle := func(c *corev1.Container) {
				c.Env = []corev1.EnvVar{
					{
						Name: "key1",
						ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"},
							Key:                  "key1",
						}},
					},
				}
			}
			mangleContainers(httpd.Spec.Containers, mangle)
			data := map[string]string{"key1": "value1"}
			secret := prepareSecret("my-secret", data)

			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, secret)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroupWithEnv(httpd, data),
				},
			}
			return sim
		},
		"Env: from ConfigMap": func() discoverySim {
			httpd := prepareHTTPDPod()
			mangle := func(c *corev1.Container) {
				c.Env = []corev1.EnvVar{
					{
						Name: "key1",
						ValueFrom: &corev1.EnvVarSource{ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{Name: "my-cmap"},
							Key:                  "key1",
						}},
					},
				}
			}
			mangleContainers(httpd.Spec.Containers, mangle)
			data := map[string]string{"key1": "value1"}
			cmap := prepareConfigMap("my-cmap", data)

			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, cmap)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroupWithEnv(httpd, data),
				},
			}
			return sim
		},
		"EnvFrom: from ConfigMap": func() discoverySim {
			httpd := prepareHTTPDPod()
			mangle := func(c *corev1.Container) {
				c.EnvFrom = []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "my-cmap"}},
					},
				}
			}
			mangleContainers(httpd.Spec.Containers, mangle)
			data := map[string]string{"key1": "value1", "key2": "value2"}
			cmap := prepareConfigMap("my-cmap", data)

			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, cmap)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroupWithEnv(httpd, data),
				},
			}
			return sim
		},
		"EnvFrom: from Secret": func() discoverySim {
			httpd := prepareHTTPDPod()
			mangle := func(c *corev1.Container) {
				c.EnvFrom = []corev1.EnvFromSource{
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}},
					},
				}
			}
			mangleContainers(httpd.Spec.Containers, mangle)
			data := map[string]string{"key1": "value1", "key2": "value2"}
			secret := prepareSecret("my-secret", data)

			discovery, _ := prepareNamespaceAllDiscovery(RolePod, httpd, secret)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					preparePodGroupWithEnv(httpd, data),
				},
			}
			return sim
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim().run(t) })
	}
}

func mangleContainers(containers []corev1.Container, mangle func(container *corev1.Container)) {
	for i := range containers {
		mangle(&containers[i])
	}
}

func prepareHTTPDPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "httpd-dd95c4d68-5bkwl",
			Namespace:   "default",
			UID:         "1cebb6eb-0c1e-495b-8131-8fa3e6668dc8",
			Annotations: map[string]string{"phase": "prod"},
			Labels:      map[string]string{"app": "httpd", "tier": "frontend"},
		},
		Spec: corev1.PodSpec{
			NodeName: "m01",
			Containers: []corev1.Container{
				{
					Name:  "httpd",
					Image: "httpd",
					Ports: []corev1.ContainerPort{
						{Name: "http", Protocol: corev1.ProtocolTCP, ContainerPort: 80},
						{Name: "https", Protocol: corev1.ProtocolTCP, ContainerPort: 443},
					},
				},
			},
		},
		Status: corev1.PodStatus{
			PodIP: "172.17.0.1",
		},
	}
}

func prepareNGINXPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "nginx-7cfd77469b-q6kxj",
			Namespace:   "default",
			UID:         "09e883f2-d740-4c5f-970d-02cf02876522",
			Annotations: map[string]string{"phase": "prod"},
			Labels:      map[string]string{"app": "nginx", "tier": "frontend"},
		},
		Spec: corev1.PodSpec{
			NodeName: "m01",
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{Name: "http", Protocol: corev1.ProtocolTCP, ContainerPort: 80},
						{Name: "https", Protocol: corev1.ProtocolTCP, ContainerPort: 443},
					},
				},
			},
		},
		Status: corev1.PodStatus{
			PodIP: "172.17.0.2",
		},
	}
}

func prepareConfigMap(name string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			UID:       types.UID("a03b8dc6-dc40-46dc-b571-5030e69d8167" + name),
		},
		Data: data,
	}
}

func prepareSecret(name string, data map[string]string) *corev1.Secret {
	secretData := make(map[string][]byte, len(data))
	for k, v := range data {
		secretData[k] = []byte(v)
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			UID:       types.UID("a03b8dc6-dc40-46dc-b571-5030e69d8161" + name),
		},
		Data: secretData,
	}
}

func prepareEmptyPodGroup(pod *corev1.Pod) *podTargetGroup {
	return &podTargetGroup{source: podSource(pod)}
}

func preparePodGroup(pod *corev1.Pod) *podTargetGroup {
	group := prepareEmptyPodGroup(pod)
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			portNum := strconv.FormatUint(uint64(port.ContainerPort), 10)
			target := &PodTarget{
				tuid:         podTUIDWithPort(pod, container, port),
				Address:      net.JoinHostPort(pod.Status.PodIP, portNum),
				Namespace:    pod.Namespace,
				Name:         pod.Name,
				Annotations:  toMapAny(pod.Annotations),
				Labels:       toMapAny(pod.Labels),
				NodeName:     pod.Spec.NodeName,
				PodIP:        pod.Status.PodIP,
				ContName:     container.Name,
				Image:        container.Image,
				Env:          nil,
				Port:         portNum,
				PortName:     port.Name,
				PortProtocol: string(port.Protocol),
			}
			target.hash = mustCalcHash(target)
			target.Tags().Merge(discoveryTags)
			group.targets = append(group.targets, target)
		}
	}
	return group
}

func preparePodGroupWithEnv(pod *corev1.Pod, env map[string]string) *podTargetGroup {
	group := preparePodGroup(pod)
	for _, target := range group.Targets() {
		target.(*PodTarget).Env = toMapAny(env)
		target.(*PodTarget).hash = mustCalcHash(target)
	}
	return group
}
