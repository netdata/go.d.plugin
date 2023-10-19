// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func TestServiceGroup_Source(t *testing.T) {
	tests := map[string]struct {
		sim            func() discoverySim
		expectedSource []string
	}{
		"ClusterIP svc with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
				discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						prepareSvcGroup(httpd),
						prepareSvcGroup(nginx),
					},
				}
				return sim
			},
			expectedSource: []string{
				"sd:k8s:service(default/httpd-cluster-ip-service)",
				"sd:k8s:service(default/nginx-cluster-ip-service)",
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

func TestServiceGroup_Targets(t *testing.T) {
	tests := map[string]struct {
		sim                func() discoverySim
		expectedNumTargets int
	}{
		"ClusterIP svc with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
				discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						prepareSvcGroup(httpd),
						prepareSvcGroup(nginx),
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

func TestServiceTarget_Hash(t *testing.T) {
	tests := map[string]struct {
		sim          func() discoverySim
		expectedHash []uint64
	}{
		"ClusterIP svc with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
				discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						prepareSvcGroup(httpd),
						prepareSvcGroup(nginx),
					},
				}
				return sim
			},
			expectedHash: []uint64{
				17611803477081780974,
				6019985892433421258,
				4151907287549842238,
				5757608926096186119,
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

func TestServiceTarget_TUID(t *testing.T) {
	tests := map[string]struct {
		sim          func() discoverySim
		expectedTUID []string
	}{
		"ClusterIP svc with multiple ports": {
			sim: func() discoverySim {
				httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
				discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

				sim := discoverySim{
					discovery: discovery,
					expectedGroups: []model.TargetGroup{
						prepareSvcGroup(httpd),
						prepareSvcGroup(nginx),
					},
				}
				return sim
			},
			expectedTUID: []string{
				"default_httpd-cluster-ip-service_tcp_80",
				"default_httpd-cluster-ip-service_tcp_443",
				"default_nginx-cluster-ip-service_tcp_80",
				"default_nginx-cluster-ip-service_tcp_443",
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

func TestNewService(t *testing.T) {
	tests := map[string]struct {
		informer  cache.SharedInformer
		wantPanic bool
	}{
		"valid informer": {informer: cache.NewSharedInformer(nil, &corev1.Service{}, resyncPeriod)},
		"nil informer":   {wantPanic: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.wantPanic {
				assert.Panics(t, func() { NewService(nil) })
			} else {
				assert.IsType(t, &Service{}, NewService(test.informer))
			}
		})
	}
}

func TestService_String(t *testing.T) {
	var s Service
	assert.NotEmpty(t, s.String())
}

func TestService_Discover(t *testing.T) {
	tests := map[string]func() discoverySim{
		"ADD: ClusterIP svc exist before run": func() discoverySim {
			httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
			discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					prepareSvcGroup(httpd),
					prepareSvcGroup(nginx),
				},
			}
			return sim
		},
		"ADD: ClusterIP svc exist before run and add after sync": func() discoverySim {
			httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
			discovery, clientset := prepareAllNsDiscovery(RoleService, httpd)
			svcClient := clientset.CoreV1().Services("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					_, _ = svcClient.Create(ctx, nginx, metav1.CreateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					prepareSvcGroup(httpd),
					prepareSvcGroup(nginx),
				},
			}
			return sim
		},
		"DELETE: ClusterIP svc remove after sync": func() discoverySim {
			httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
			discovery, clientset := prepareAllNsDiscovery(RoleService, httpd, nginx)
			svcClient := clientset.CoreV1().Services("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_ = svcClient.Delete(ctx, httpd.Name, metav1.DeleteOptions{})
					_ = svcClient.Delete(ctx, nginx.Name, metav1.DeleteOptions{})
				},
				expectedGroups: []model.TargetGroup{
					prepareSvcGroup(httpd),
					prepareSvcGroup(nginx),
					prepareEmptySvcGroup(httpd),
					prepareEmptySvcGroup(nginx),
				},
			}
			return sim
		},
		"ADD,DELETE: ClusterIP svc remove and add after sync": func() discoverySim {
			httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
			discovery, clientset := prepareAllNsDiscovery(RoleService, httpd)
			svcClient := clientset.CoreV1().Services("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_ = svcClient.Delete(ctx, httpd.Name, metav1.DeleteOptions{})
					_, _ = svcClient.Create(ctx, nginx, metav1.CreateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					prepareSvcGroup(httpd),
					prepareEmptySvcGroup(httpd),
					prepareSvcGroup(nginx),
				},
			}
			return sim
		},
		"ADD: Headless svc exist before run": func() discoverySim {
			httpd, nginx := newHTTPDHeadlessService(), newNGINXHeadlessService()
			discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					prepareEmptySvcGroup(httpd),
					prepareEmptySvcGroup(nginx),
				},
			}
			return sim
		},
		"UPDATE: Headless => ClusterIP svc after sync": func() discoverySim {
			httpd, nginx := newHTTPDHeadlessService(), newNGINXHeadlessService()
			httpdUpd, nginxUpd := *httpd, *nginx
			httpdUpd.Spec.ClusterIP = "10.100.0.1"
			nginxUpd.Spec.ClusterIP = "10.100.0.2"
			discovery, clientset := prepareAllNsDiscovery(RoleService, httpd, nginx)
			svcClient := clientset.CoreV1().Services("default")

			sim := discoverySim{
				discovery: discovery,
				runAfterSync: func(ctx context.Context) {
					time.Sleep(time.Millisecond * 50)
					_, _ = svcClient.Update(ctx, &httpdUpd, metav1.UpdateOptions{})
					_, _ = svcClient.Update(ctx, &nginxUpd, metav1.UpdateOptions{})
				},
				expectedGroups: []model.TargetGroup{
					prepareEmptySvcGroup(httpd),
					prepareEmptySvcGroup(nginx),
					prepareSvcGroup(&httpdUpd),
					prepareSvcGroup(&nginxUpd),
				},
			}
			return sim
		},
		"ADD: ClusterIP svc with zero exposed ports": func() discoverySim {
			httpd, nginx := newHTTPDClusterIPService(), newNGINXClusterIPService()
			httpd.Spec.Ports = httpd.Spec.Ports[:0]
			nginx.Spec.Ports = httpd.Spec.Ports[:0]
			discovery, _ := prepareAllNsDiscovery(RoleService, httpd, nginx)

			sim := discoverySim{
				discovery: discovery,
				expectedGroups: []model.TargetGroup{
					prepareEmptySvcGroup(httpd),
					prepareEmptySvcGroup(nginx),
				},
			}
			return sim
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim().run(t) })
	}

}

func newHTTPDClusterIPService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "httpd-cluster-ip-service",
			Namespace:   "default",
			Annotations: map[string]string{"phase": "prod"},
			Labels:      map[string]string{"app": "httpd", "tier": "frontend"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http", Protocol: corev1.ProtocolTCP, Port: 80},
				{Name: "https", Protocol: corev1.ProtocolTCP, Port: 443},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "10.100.0.1",
			Selector:  map[string]string{"app": "httpd", "tier": "frontend"},
		},
	}
}

func newNGINXClusterIPService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "nginx-cluster-ip-service",
			Namespace:   "default",
			Annotations: map[string]string{"phase": "prod"},
			Labels:      map[string]string{"app": "nginx", "tier": "frontend"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http", Protocol: corev1.ProtocolTCP, Port: 80},
				{Name: "https", Protocol: corev1.ProtocolTCP, Port: 443},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "10.100.0.2",
			Selector:  map[string]string{"app": "nginx", "tier": "frontend"},
		},
	}
}

func newHTTPDHeadlessService() *corev1.Service {
	svc := newHTTPDClusterIPService()
	svc.Name = "httpd-headless-service"
	svc.Spec.ClusterIP = ""
	return svc
}

func newNGINXHeadlessService() *corev1.Service {
	svc := newNGINXClusterIPService()
	svc.Name = "nginx-headless-service"
	svc.Spec.ClusterIP = ""
	return svc
}

func prepareEmptySvcGroup(svc *corev1.Service) *serviceGroup {
	return &serviceGroup{source: serviceSource(svc)}
}

func prepareSvcGroup(svc *corev1.Service) *serviceGroup {
	group := prepareEmptySvcGroup(svc)
	for _, port := range svc.Spec.Ports {
		portNum := strconv.FormatInt(int64(port.Port), 10)
		target := &ServiceTarget{
			tuid:         serviceTUID(svc, port),
			Address:      net.JoinHostPort(svc.Name+"."+svc.Namespace+".svc", portNum),
			Namespace:    svc.Namespace,
			Name:         svc.Name,
			Annotations:  mapAny(svc.Annotations),
			Labels:       mapAny(svc.Labels),
			Port:         portNum,
			PortName:     port.Name,
			PortProtocol: string(port.Protocol),
			ClusterIP:    svc.Spec.ClusterIP,
			ExternalName: svc.Spec.ExternalName,
			Type:         string(svc.Spec.Type),
		}
		target.hash = mustCalcHash(target)
		target.Tags().Merge(discoveryTags)
		group.targets = append(group.targets, target)
	}
	return group
}
