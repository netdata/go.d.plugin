package k8s_inventory

import (
	"context"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type kubeDiscovery struct {
	client      kubernetes.Interface
	discoverers []discoverer
}

func (d *kubeDiscovery) run(ctx context.Context, in chan<- resource) {
	d.discoverers = d.setupDiscoverers(ctx)

	var wg sync.WaitGroup
	updates := make(chan resource)

	for _, dd := range d.discoverers {
		wg.Add(1)
		go func(dd discoverer) { defer wg.Done(); dd.run(ctx, updates) }(dd)
	}

	wg.Add(1)
	go func() { defer wg.Done(); d.runProcessUpdates(ctx, updates, in) }()

	wg.Wait()
	<-ctx.Done()
}

func (d *kubeDiscovery) runProcessUpdates(ctx context.Context, updates chan resource, in chan<- resource) {
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

func (d *kubeDiscovery) setupDiscoverers(ctx context.Context) []discoverer {
	node := d.client.CoreV1().Nodes()
	nodeWatcher := &cache.ListWatch{
		ListFunc:  func(options metav1.ListOptions) (runtime.Object, error) { return node.List(ctx, options) },
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) { return node.Watch(ctx, options) },
	}

	pod := d.client.CoreV1().Pods(corev1.NamespaceAll)
	podWatcher := &cache.ListWatch{
		ListFunc:  func(options metav1.ListOptions) (runtime.Object, error) { return pod.List(ctx, options) },
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) { return pod.Watch(ctx, options) },
	}

	pv := d.client.CoreV1().PersistentVolumes()
	pvWatcher := &cache.ListWatch{
		ListFunc:  func(options metav1.ListOptions) (runtime.Object, error) { return pv.List(ctx, options) },
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) { return pv.Watch(ctx, options) },
	}

	pvc := d.client.CoreV1().PersistentVolumeClaims(corev1.NamespaceAll)
	pvcWatcher := &cache.ListWatch{
		ListFunc:  func(options metav1.ListOptions) (runtime.Object, error) { return pvc.List(ctx, options) },
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) { return pvc.Watch(ctx, options) },
	}

	return []discoverer{
		newNodeDiscoverer(cache.NewSharedInformer(nodeWatcher, &corev1.Node{}, resyncPeriod)),
		newPodDiscoverer(cache.NewSharedInformer(podWatcher, &corev1.Pod{}, resyncPeriod)),
		newPVCDiscoverer(cache.NewSharedInformer(pvcWatcher, &corev1.PersistentVolumeClaim{}, resyncPeriod)),
		newPVDiscoverer(cache.NewSharedInformer(pvWatcher, &corev1.PersistentVolume{}, resyncPeriod)),
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
