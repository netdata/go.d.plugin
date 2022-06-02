package k8s_state

import (
	"context"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/logger"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newKubeDiscovery(client kubernetes.Interface, l *logger.Logger) *kubeDiscovery {
	return &kubeDiscovery{
		client:  client,
		Logger:  l,
		started: make(chan struct{}),
	}
}

type kubeDiscovery struct {
	*logger.Logger
	client      kubernetes.Interface
	discoverers []discoverer
	started     chan struct{}
}

func (d *kubeDiscovery) run(ctx context.Context, in chan<- resource) {
	d.Info("kube_discoverer is started")
	defer func() { d.Info("kube_discoverer is stopped") }()

	d.discoverers = d.setupDiscoverers(ctx)

	var wg sync.WaitGroup
	updates := make(chan resource)

	for _, dd := range d.discoverers {
		wg.Add(1)
		go func(dd discoverer) { defer wg.Done(); dd.run(ctx, updates) }(dd)
	}

	wg.Add(1)
	go func() { defer wg.Done(); d.runDiscover(ctx, updates, in) }()

	close(d.started)
	wg.Wait()
	<-ctx.Done()
}

func (d *kubeDiscovery) ready() bool {
	select {
	case <-d.started:
		for _, dd := range d.discoverers {
			if !dd.ready() {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (d *kubeDiscovery) runDiscover(ctx context.Context, updates chan resource, in chan<- resource) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-updates:
			select {
			case <-ctx.Done():
				return
			case in <- r:
			}
		}
	}
}

const resyncPeriod = 10 * time.Minute

var (
	myNodeName = os.Getenv("MY_NODE_NAME")
)

func (d *kubeDiscovery) setupDiscoverers(ctx context.Context) []discoverer {
	node := d.client.CoreV1().Nodes()
	nodeWatcher := &cache.ListWatch{
		ListFunc:  func(options metav1.ListOptions) (runtime.Object, error) { return node.List(ctx, options) },
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) { return node.Watch(ctx, options) },
	}

	pod := d.client.CoreV1().Pods(corev1.NamespaceAll)
	podWatcher := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			if myNodeName != "" {
				options.FieldSelector = "spec.nodeName=" + myNodeName
			}
			return pod.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			if myNodeName != "" {
				options.FieldSelector = "spec.nodeName=" + myNodeName
			}
			return pod.Watch(ctx, options)
		},
	}

	return []discoverer{
		newNodeDiscoverer(cache.NewSharedInformer(nodeWatcher, &corev1.Node{}, resyncPeriod), d.Logger),
		newPodDiscoverer(cache.NewSharedInformer(podWatcher, &corev1.Pod{}, resyncPeriod), d.Logger),
	}
}

func enqueue(queue *workqueue.Type, obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	queue.Add(key)
}

func send(ctx context.Context, in chan<- resource, r resource) {
	if r == nil {
		return
	}
	select {
	case <-ctx.Done():
	case in <- r:
	}
}

var reDots = regexp.MustCompile(`\.`)

func replaceDots(v string) string {
	return reDots.ReplaceAllString(v, "-")
}
