// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import (
	"fmt"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/pkg/k8sclient"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMain(m *testing.M) {
	_ = os.Setenv(envNodeName, "m01")
	_ = os.Setenv(k8sclient.EnvFakeClient, "true")
	code := m.Run()
	_ = os.Unsetenv(envNodeName)
	_ = os.Unsetenv(k8sclient.EnvFakeClient)
	os.Exit(code)
}

func TestNewDiscovery(t *testing.T) {
	tests := map[string]struct {
		cfg     Config
		wantErr bool
	}{
		"role pod and local mode":     {cfg: Config{Role: RolePod, Tags: "k8s", LocalMode: true}},
		"role service and local mode": {cfg: Config{Role: RoleService, Tags: "k8s", LocalMode: true}},
		"empty config":                {wantErr: true},
		"invalid role":                {cfg: Config{Role: "invalid"}, wantErr: true},
		"lack of tags":                {cfg: Config{Role: RolePod}, wantErr: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			discovery, err := NewDiscovery(test.cfg)

			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, discovery)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, discovery)
				if test.cfg.LocalMode && test.cfg.Role == RolePod {
					assert.Contains(t, discovery.selectorField, "spec.nodeName=m01")
				}
				if test.cfg.LocalMode && test.cfg.Role != RolePod {
					assert.Empty(t, discovery.selectorField)
				}
			}
		})
	}
}

func TestDiscovery_Discover(t *testing.T) {
	const prod = "prod"
	const dev = "dev"
	prodNamespace := newNamespace(prod)
	devNamespace := newNamespace(dev)

	tests := map[string]func() discoverySim{
		"multiple namespaces pod discovery": func() discoverySim {
			httpdProd, nginxProd := newHTTPDPod(), newNGINXPod()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := newHTTPDPod(), newNGINXPod()
			httpdDev.Namespace = dev
			nginxDev.Namespace = dev

			discovery, _ := prepareDiscovery(
				RolePod,
				[]string{prod, dev},
				prodNamespace, devNamespace, httpdProd, nginxProd, httpdDev, nginxDev)

			sim := discoverySim{
				discovery:        discovery,
				sortBeforeVerify: true,
				expectedGroups: []model.TargetGroup{
					preparePodGroup(httpdDev),
					preparePodGroup(nginxDev),
					preparePodGroup(httpdProd),
					preparePodGroup(nginxProd),
				},
			}
			return sim
		},
		"multiple namespaces ClusterIP service discovery": func() discoverySim {
			httpdProd, nginxProd := newHTTPDClusterIPService(), newNGINXClusterIPService()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := newHTTPDClusterIPService(), newNGINXClusterIPService()
			httpdDev.Namespace = dev
			nginxDev.Namespace = dev

			discovery, _ := prepareDiscovery(
				RoleService,
				[]string{prod, dev},
				prodNamespace, devNamespace, httpdProd, nginxProd, httpdDev, nginxDev)

			sim := discoverySim{
				discovery:        discovery,
				sortBeforeVerify: true,
				expectedGroups: []model.TargetGroup{
					prepareSvcGroup(httpdDev),
					prepareSvcGroup(nginxDev),
					prepareSvcGroup(httpdProd),
					prepareSvcGroup(nginxProd),
				},
			}
			return sim
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim().run(t) })
	}
}

var discoveryTags model.Tags = map[string]struct{}{"k8s": {}}

func prepareAllNsDiscovery(role string, objects ...runtime.Object) (*Discovery, kubernetes.Interface) {
	return prepareDiscovery(role, []string{corev1.NamespaceAll}, objects...)
}

func prepareDiscovery(role string, namespaces []string, objects ...runtime.Object) (*Discovery, kubernetes.Interface) {
	clientset := fake.NewSimpleClientset(objects...)
	discovery := &Discovery{
		tags:          discoveryTags,
		namespaces:    namespaces,
		role:          role,
		selectorLabel: "",
		selectorField: "",
		client:        clientset,
		discoverers:   nil,
		started:       make(chan struct{}),
	}
	return discovery, clientset
}

func newNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

func mustCalcHash(target interface{}) uint64 {
	hash, err := calcHash(target)
	if err != nil {
		panic(fmt.Sprintf("hash calculation: %v", err))
	}
	return hash
}
