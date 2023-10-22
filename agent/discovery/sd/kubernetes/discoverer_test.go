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

func TestNewTargetDiscoverer(t *testing.T) {
	tests := map[string]struct {
		cfg     Config
		wantErr bool
	}{
		"role pod and local mode": {
			wantErr: false,
			cfg:     Config{Role: RolePod, LocalMode: true},
		},
		"role service and local mode": {
			wantErr: false,
			cfg:     Config{Role: RoleService, LocalMode: true},
		},
		"empty config": {
			wantErr: true,
			cfg:     Config{},
		},
		"invalid role": {
			wantErr: true,
			cfg:     Config{Role: "invalid"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			disc, err := NewTargetDiscoverer(test.cfg)

			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, disc)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, disc)
			if test.cfg.LocalMode {
				if test.cfg.Role == RolePod {
					assert.Contains(t, disc.selectorField, "spec.nodeName=m01")
				} else {
					assert.Empty(t, disc.selectorField)
				}
			}
		})
	}
}

func TestTargetDiscoverer_Discover(t *testing.T) {
	const prod = "prod"
	const dev = "dev"
	prodNamespace := newNamespace(prod)
	devNamespace := newNamespace(dev)

	tests := map[string]func() discoverySim{
		"multiple namespaces pod td": func() discoverySim {
			httpdProd, nginxProd := newHTTPDPod(), newNGINXPod()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := newHTTPDPod(), newNGINXPod()
			httpdDev.Namespace = dev
			nginxDev.Namespace = dev

			disc, _ := prepareDiscoverer(
				RolePod,
				[]string{prod, dev},
				prodNamespace, devNamespace, httpdProd, nginxProd, httpdDev, nginxDev)

			return discoverySim{
				td:               disc,
				sortBeforeVerify: true,
				wantTargetGroups: []model.TargetGroup{
					preparePodTargetGroup(httpdDev),
					preparePodTargetGroup(nginxDev),
					preparePodTargetGroup(httpdProd),
					preparePodTargetGroup(nginxProd),
				},
			}
		},
		"multiple namespaces ClusterIP service td": func() discoverySim {
			httpdProd, nginxProd := newHTTPDClusterIPService(), newNGINXClusterIPService()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := newHTTPDClusterIPService(), newNGINXClusterIPService()
			httpdDev.Namespace = dev
			nginxDev.Namespace = dev

			disc, _ := prepareDiscoverer(
				RoleService,
				[]string{prod, dev},
				prodNamespace, devNamespace, httpdProd, nginxProd, httpdDev, nginxDev)

			return discoverySim{
				td:               disc,
				sortBeforeVerify: true,
				wantTargetGroups: []model.TargetGroup{
					prepareSvcTargetGroup(httpdDev),
					prepareSvcTargetGroup(nginxDev),
					prepareSvcTargetGroup(httpdProd),
					prepareSvcTargetGroup(nginxProd),
				},
			}
		},
	}

	for name, createSim := range tests {
		t.Run(name, func(t *testing.T) {
			sim := createSim()
			sim.run(t)
		})
	}
}

func prepareAllNsDiscoverer(role string, objects ...runtime.Object) (*TargetDiscoverer, kubernetes.Interface) {
	return prepareDiscoverer(role, []string{corev1.NamespaceAll}, objects...)
}

func prepareDiscoverer(role string, namespaces []string, objects ...runtime.Object) (*TargetDiscoverer, kubernetes.Interface) {
	client := fake.NewSimpleClientset(objects...)
	disc := &TargetDiscoverer{
		namespaces:    namespaces,
		role:          role,
		selectorLabel: "",
		selectorField: "",
		client:        client,
		discoverers:   nil,
		started:       make(chan struct{}),
	}
	return disc, client
}

func newNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

func mustCalcHash(obj any) uint64 {
	hash, err := calcHash(obj)
	if err != nil {
		panic(fmt.Sprintf("hash calculation: %v", err))
	}
	return hash
}
