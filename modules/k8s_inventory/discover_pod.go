package k8s_inventory

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newPodDiscoverer(si cache.SharedInformer) *podDiscoverer {
	if si == nil {
		panic("nil pod shared informer")
	}

	queue := workqueue.NewNamed("pod")
	si.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &podDiscoverer{
		informer: si,
		queue:    queue,
	}
}

type podResource struct {
	src string
	val interface{}
}

func (r podResource) source() string     { return r.src }
func (r podResource) kind() resourceKind { return kindPod }
func (r podResource) value() interface{} { return r.val }

type podDiscoverer struct {
	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (d *podDiscoverer) run(ctx context.Context, in chan<- resource) {
	defer func() { fmt.Println("STOP POD DISCOVERER") }()
	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runProcessQueue(ctx, in)

	<-ctx.Done()
}

func (d *podDiscoverer) runProcessQueue(ctx context.Context, in chan<- resource) {
	for {
		item, shutdown := d.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer d.queue.Done(item)

			key := item.(string)
			ns, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := d.informer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			r := &podResource{src: podSource(ns, name)}
			if exists {
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func podSource(namespace, name string) string {
	return "k8s/pod/" + namespace + "/" + name
}
