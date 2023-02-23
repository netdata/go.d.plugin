package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"

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

func NewDiscovery(cfg DiscoveryConfig) (*Discovery, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("k8s discovery config validation: %v", err)
	}

	d, err := initDiscovery(cfg)
	if err != nil {
		return nil, fmt.Errorf("k8s discovery initialization ('%s'): %v", cfg.Role, err)
	}

	return d, nil
}

type DiscoveryConfig struct {
	APIServer  string   `yaml:"api_server"`
	Tags       string   `yaml:"tags"`
	Namespaces []string `yaml:"namespaces"`
	Role       string   `yaml:"role"`
	LocalMode  bool     `yaml:"local_mode"`
	Selector   struct {
		Label string `yaml:"label"`
		Field string `yaml:"field"`
	} `yaml:"selector"`
}

type (
	Discovery struct {
		*logger.Logger

		tags          model.Tags
		namespaces    []string
		role          string
		selectorLabel string
		selectorField string
		client        kubernetes.Interface
		discoverers   []discoverer
		started       chan struct{}
	}

	discoverer interface {
		Discover(ctx context.Context, ch chan<- []model.TargetGroup)
	}
)

func (d *Discovery) String() string {
	return "k8s discovery manager"
}

const resyncPeriod = 10 * time.Minute

func (d *Discovery) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	if err := d.setupDiscoverers(ctx); err != nil {
		// TODO: no panic pls
		panic(err)
	}

	var wg sync.WaitGroup
	updates := make(chan []model.TargetGroup)

	for _, dd := range d.discoverers {
		wg.Add(1)
		go func(dd discoverer) { defer wg.Done(); dd.Discover(ctx, updates) }(dd)
	}

	wg.Add(1)
	go func() { defer wg.Done(); d.run(ctx, updates, in) }()

	close(d.started)

	wg.Wait()

	<-ctx.Done()
}

func (d *Discovery) run(ctx context.Context, updates chan []model.TargetGroup, in chan<- []model.TargetGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case groups := <-updates:
			for _, group := range groups {
				for _, target := range group.Targets() {
					target.Tags().Merge(d.tags)
				}
			}
			select {
			case <-ctx.Done():
				return
			case in <- groups:
			}
		}
	}
}

func (d *Discovery) setupDiscoverers(ctx context.Context) error {
	for _, namespace := range d.namespaces {
		var dd discoverer
		switch d.role {
		case RolePod:
			dd = d.setupPodDiscovery(ctx, namespace)
		case RoleService:
			dd = d.setupServiceDiscovery(ctx, namespace)
		default:
			return fmt.Errorf("unknown k8 discovery role: '%s'", d.role)
		}
		d.discoverers = append(d.discoverers, dd)
	}

	if len(d.discoverers) == 0 {
		return errors.New("k8s cant run discovery: zero discoverers")
	}

	d.Infof("registered: %v", d.discoverers)
	return nil
}

func (d *Discovery) setupPodDiscovery(ctx context.Context, namespace string) *PodDiscovery {
	pod := d.client.CoreV1().Pods(namespace)
	podLW := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			opts.FieldSelector = d.selectorField
			opts.LabelSelector = d.selectorLabel
			return pod.List(ctx, opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			opts.FieldSelector = d.selectorField
			opts.LabelSelector = d.selectorLabel
			return pod.Watch(ctx, opts)
		},
	}

	cmap := d.client.CoreV1().ConfigMaps(namespace)
	cmapLW := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return cmap.List(ctx, opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return cmap.Watch(ctx, opts)
		},
	}

	secret := d.client.CoreV1().Secrets(namespace)
	secretLW := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return secret.List(ctx, opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return secret.Watch(ctx, opts)
		},
	}

	return NewPodDiscovery(
		cache.NewSharedInformer(podLW, &corev1.Pod{}, resyncPeriod),
		cache.NewSharedInformer(cmapLW, &corev1.ConfigMap{}, resyncPeriod),
		cache.NewSharedInformer(secretLW, &corev1.Secret{}, resyncPeriod),
	)
}

func (d *Discovery) setupServiceDiscovery(ctx context.Context, namespace string) *ServiceDiscovery {
	svc := d.client.CoreV1().Services(namespace)
	clw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			opts.FieldSelector = d.selectorField
			opts.LabelSelector = d.selectorLabel
			return svc.List(ctx, opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			opts.FieldSelector = d.selectorField
			opts.LabelSelector = d.selectorLabel
			return svc.Watch(ctx, opts)
		},
	}
	inf := cache.NewSharedInformer(clw, &corev1.Service{}, resyncPeriod)

	return NewServiceDiscovery(inf)
}

func validateConfig(cfg DiscoveryConfig) error {
	if !(cfg.Role == RolePod || cfg.Role == RoleService) {
		return fmt.Errorf("invalid role '%s', valid roles: '%s', '%s'", cfg.Role, RolePod, RoleService)
	}
	if cfg.Tags == "" {
		return fmt.Errorf("no tags set for '%s' role", cfg.Role)
	}
	return nil
}

func initDiscovery(cfg DiscoveryConfig) (*Discovery, error) {
	tags, err := model.ParseTags(cfg.Tags)
	if err != nil {
		return nil, fmt.Errorf("parse config->tags: %v", err)
	}

	client, err := newClientset()
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

	d := &Discovery{
		tags:          tags,
		namespaces:    namespaces,
		role:          cfg.Role,
		selectorLabel: cfg.Selector.Label,
		selectorField: cfg.Selector.Field,
		client:        client,
		discoverers:   make([]discoverer, 0, len(namespaces)),
		started:       make(chan struct{}),
		Logger:        logger.New("k8s discovery", "kube"),
	}

	return d, nil
}

func enqueue(queue *workqueue.Type, obj any) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	queue.Add(key)
}

func send(ctx context.Context, in chan<- []model.TargetGroup, group model.TargetGroup) {
	if group == nil {
		return
	}
	select {
	case <-ctx.Done():
	case in <- []model.TargetGroup{group}:
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
