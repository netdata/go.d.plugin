// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"
	"github.com/netdata/go.d.plugin/pkg/k8sclient"

	"github.com/ilyam8/hashstructure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	RolePod     = "pod"
	RoleService = "service"
)

const (
	envNodeName = "MY_NODE_NAME"
)

func NewTargetDiscoverer(cfg Config) (*TargetDiscoverer, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("k8s td config validation: %v", err)
	}

	client, err := k8sclient.New("Netdata/service-td")
	if err != nil {
		return nil, fmt.Errorf("create clientset: %v", err)
	}

	namespaces := cfg.Namespaces
	if len(namespaces) == 0 {
		namespaces = []string{corev1.NamespaceAll}
	}

	if cfg.LocalMode && cfg.Role == RolePod {
		name := os.Getenv(envNodeName)
		if name == "" {
			return nil, fmt.Errorf("local_mode is enabled, but env '%s' not set", envNodeName)
		}
		cfg.Selector.Field = joinSelectors(cfg.Selector.Field, "spec.nodeName="+name)
	}

	d := &TargetDiscoverer{
		Logger:        logger.New("k8s td manager", ""),
		namespaces:    namespaces,
		role:          cfg.Role,
		selectorLabel: cfg.Selector.Label,
		selectorField: cfg.Selector.Field,
		client:        client,
		discoverers:   make([]model.Discoverer, 0, len(namespaces)),
		started:       make(chan struct{}),
	}

	return d, nil
}

type TargetDiscoverer struct {
	*logger.Logger

	namespaces    []string
	role          string
	selectorLabel string
	selectorField string
	client        kubernetes.Interface
	discoverers   []model.Discoverer
	started       chan struct{}
}

func (d *TargetDiscoverer) String() string {
	return "k8s td manager"
}

const resyncPeriod = 10 * time.Minute

func (d *TargetDiscoverer) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	for _, namespace := range d.namespaces {
		var disc model.Discoverer
		switch d.role {
		case RolePod:
			disc = d.setupPodDiscoverer(ctx, namespace)
		case RoleService:
			disc = d.setupServiceDiscoverer(ctx, namespace)
		default:
			panic(fmt.Sprintf("unknown k8 td role: '%s'", d.role))
		}
		d.discoverers = append(d.discoverers, disc)
	}
	if len(d.discoverers) == 0 {
		panic("k8s cant run td: zero discoverers")
	}

	d.Infof("registered: %v", d.discoverers)

	var wg sync.WaitGroup
	updates := make(chan []model.TargetGroup)

	for _, disc := range d.discoverers {
		wg.Add(1)
		go func(disc model.Discoverer) { defer wg.Done(); disc.Discover(ctx, updates) }(disc)
	}

	wg.Add(1)
	go func() { defer wg.Done(); d.run(ctx, updates, in) }()

	close(d.started)

	wg.Wait()
	<-ctx.Done()
}

func (d *TargetDiscoverer) run(ctx context.Context, updates chan []model.TargetGroup, in chan<- []model.TargetGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case tggs := <-updates:
			select {
			case <-ctx.Done():
				return
			case in <- tggs:
			}
		}
	}
}

func (d *TargetDiscoverer) setupPodDiscoverer(ctx context.Context, namespace string) *PodTargetDiscoverer {
	pod := d.client.CoreV1().Pods(namespace)
	podLW := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = d.selectorField
			options.LabelSelector = d.selectorLabel
			return pod.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = d.selectorField
			options.LabelSelector = d.selectorLabel
			return pod.Watch(ctx, options)
		},
	}

	cmap := d.client.CoreV1().ConfigMaps(namespace)
	cmapLW := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return cmap.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return cmap.Watch(ctx, options)
		},
	}

	secret := d.client.CoreV1().Secrets(namespace)
	secretLW := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return secret.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return secret.Watch(ctx, options)
		},
	}

	return NewPodTargetDiscoverer(
		cache.NewSharedInformer(podLW, &corev1.Pod{}, resyncPeriod),
		cache.NewSharedInformer(cmapLW, &corev1.ConfigMap{}, resyncPeriod),
		cache.NewSharedInformer(secretLW, &corev1.Secret{}, resyncPeriod),
	)
}

func (d *TargetDiscoverer) setupServiceDiscoverer(ctx context.Context, namespace string) *ServiceTargetDiscoverer {
	svc := d.client.CoreV1().Services(namespace)

	svcLW := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = d.selectorField
			options.LabelSelector = d.selectorLabel
			return svc.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = d.selectorField
			options.LabelSelector = d.selectorLabel
			return svc.Watch(ctx, options)
		},
	}

	inf := cache.NewSharedInformer(svcLW, &corev1.Service{}, resyncPeriod)

	return NewServiceTargetDiscoverer(inf)
}

func enqueue(queue *workqueue.Type, obj any) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	queue.Add(key)
}

func send(ctx context.Context, in chan<- []model.TargetGroup, tgg model.TargetGroup) {
	if tgg == nil {
		return
	}
	select {
	case <-ctx.Done():
	case in <- []model.TargetGroup{tgg}:
	}
}

func calcHash(obj any) (uint64, error) {
	return hashstructure.Hash(obj, nil)
}

func joinSelectors(srs ...string) string {
	var i int
	for _, v := range srs {
		if v != "" {
			srs[i] = v
			i++
		}
	}
	return strings.Join(srs[:i], ",")
}
