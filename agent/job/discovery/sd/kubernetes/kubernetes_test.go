package kubernetes

import (
	"fmt"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMain(m *testing.M) {
	_ = os.Setenv(envNodeName, "m01")
	_ = os.Setenv(envFakeClient, "true")
	code := m.Run()
	_ = os.Unsetenv(envNodeName)
	_ = os.Unsetenv(envFakeClient)
	os.Exit(code)
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		cfg     DiscoveryConfig
		wantErr bool
	}{
		"role pod and local mode": {
			wantErr: false,
			cfg:     DiscoveryConfig{Role: RolePod, Tags: "k8s", LocalMode: true},
		},
		"role service and local mode": {
			wantErr: false,
			cfg:     DiscoveryConfig{Role: RoleService, Tags: "k8s", LocalMode: true},
		},
		"empty config": {
			wantErr: true,
		},
		"invalid role": {
			wantErr: true,
			cfg:     DiscoveryConfig{Role: "invalid"},
		},
		"lack of tags": {
			wantErr: true,
			cfg:     DiscoveryConfig{Role: RolePod},
		},
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
	prodNamespace := prepareNamespace(prod)
	devNamespace := prepareNamespace(dev)

	tests := map[string]func() discoverySim{
		"multiple namespaces pod discovery": func() discoverySim {
			httpdProd, nginxProd := prepareHTTPDPod(), prepareNGINXPod()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := prepareHTTPDPod(), prepareNGINXPod()
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
			httpdProd, nginxProd := prepareHTTPDClusterIPService(), prepareNGINXClusterIPService()
			httpdProd.Namespace = prod
			nginxProd.Namespace = prod

			httpdDev, nginxDev := prepareHTTPDClusterIPService(), prepareNGINXClusterIPService()
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

func prepareNamespaceAllDiscovery(role string, objects ...runtime.Object) (*Discovery, kubernetes.Interface) {
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

func prepareNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

func mustCalcHash(target any) uint64 {
	hash, err := calcHash(target)
	if err != nil {
		panic(fmt.Sprintf("hash calculation: %v", err))
	}
	return hash
}
