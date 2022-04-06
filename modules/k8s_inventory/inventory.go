package k8s_inventory

import (
	"context"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"k8s.io/client-go/kubernetes"
)

func init() {
	module.Register("ki", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *KubernetesInventory {
	return &KubernetesInventory{
		once: &sync.Once{},
		wg:   &sync.WaitGroup{},
	}
}

type Config struct {
}

type (
	discoverer interface {
		run(ctx context.Context, in chan<- resource)
	}

	KubernetesInventory struct {
		module.Base
		Config `yaml:",inline"`

		client     kubernetes.Interface
		once       *sync.Once
		wg         *sync.WaitGroup
		discoverer discoverer
		ctx        context.Context
		ctxCancel  context.CancelFunc
	}
)

func (ki *KubernetesInventory) Init() bool {
	client, err := ki.initClient()
	if err != nil {
		ki.Errorf("client initialization: %v", err)
		return false
	}
	ki.client = client

	ki.ctx, ki.ctxCancel = context.WithCancel(context.Background())

	ki.discoverer = ki.initDiscoverer(ki.client)

	return true
}

func (ki *KubernetesInventory) Check() bool {
	return len(ki.Collect()) > 0
	if ki.client == nil {
		ki.Error("not initialized client")
		return false
	}

	ver, err := ki.client.Discovery().ServerVersion()
	if err != nil {
		ki.Errorf("failed to connect to the Kuberneter API server: %v", err)
		return false
	}

	ki.Infof("successfully connected to the Kuberneter API server '%s'", ver)
	return true
}

func (ki *KubernetesInventory) Charts() *module.Charts {
	return nil
}

func (ki *KubernetesInventory) Collect() map[string]int64 {
	ms, err := ki.collect()
	if err != nil {
		ki.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (ki *KubernetesInventory) Cleanup() {
	if ki.ctxCancel != nil {
		ki.ctxCancel()
	}
	t := time.NewTimer(time.Second * 5)
	defer t.Stop()
	select {}
}
