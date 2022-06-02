package k8s_state

import (
	"context"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"k8s.io/client-go/kubernetes"
)

func init() {
	module.Register("k8s_state", module.Creator{
		Defaults: module.Defaults{
			//Disabled: true,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *KubeState {
	return &KubeState{
		initDelay:     time.Second * 3,
		newKubeClient: newKubeClient,
		charts:        baseCharts.Copy(),
		once:          &sync.Once{},
		wg:            &sync.WaitGroup{},
		state:         newKubeState(),
	}
}

type (
	discoverer interface {
		run(ctx context.Context, in chan<- resource)
		hasSynced() bool
	}

	KubeState struct {
		module.Base

		newKubeClient func() (kubernetes.Interface, error)

		startTime time.Time
		initDelay time.Duration

		charts *module.Charts

		client     kubernetes.Interface
		once       *sync.Once
		wg         *sync.WaitGroup
		discoverer discoverer
		ctx        context.Context
		ctxCancel  context.CancelFunc
		state      *kubeState
	}
)

func (ks *KubeState) Init() bool {
	client, err := ks.initClient()
	if err != nil {
		ks.Errorf("client initialization: %v", err)
		return false
	}
	ks.client = client

	ks.ctx, ks.ctxCancel = context.WithCancel(context.Background())

	ks.discoverer = ks.initDiscoverer(ks.client)

	return true
}

func (ks *KubeState) Check() bool {
	if ks.client == nil || ks.discoverer == nil {
		ks.Error("not initialized job")
		return false
	}

	ver, err := ks.client.Discovery().ServerVersion()
	if err != nil {
		ks.Errorf("failed to connect to the Kubernetes API server: %v", err)
		return false
	}

	ks.Infof("successfully connected to the Kubernetes API server '%s'", ver)
	return true
}

func (ks *KubeState) Charts() *module.Charts {
	return ks.charts
}

func (ks *KubeState) Collect() map[string]int64 {
	ms, err := ks.collect()
	if err != nil {
		ks.Error(err)
	}

	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (ks *KubeState) Cleanup() {
	if ks.ctxCancel == nil {
		return
	}
	ks.ctxCancel()

	c := make(chan struct{})
	go func() { defer close(c); ks.wg.Wait() }()

	t := time.NewTimer(time.Second * 3)
	defer t.Stop()

	select {
	case <-c:
		return
	case <-t.C:
		return
	}
}
