package k8s_inventory

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newPVCDiscoverer(si cache.SharedInformer) *pvcDiscoverer {
	if si == nil {
		panic("nil persistent volume claim shared informer")
	}

	queue := workqueue.NewNamed("pvc")
	si.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &pvcDiscoverer{
		informer: si,
		queue:    queue,
	}
}

type pvcResource struct {
	src string
	val interface{}
}

func (r pvcResource) source() string     { return r.src }
func (r pvcResource) kind() resourceKind { return kindPersistentVolumeClaim }
func (r pvcResource) value() interface{} { return r.val }

type pvcDiscoverer struct {
	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (d *pvcDiscoverer) run(ctx context.Context, in chan<- resource) {
	defer func() { fmt.Println("STOP PVC DISCOVERER") }()
	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runProcessQueue(ctx, in)

	<-ctx.Done()
}

func (d *pvcDiscoverer) runProcessQueue(ctx context.Context, in chan<- resource) {
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

			r := &pvcResource{src: pvcSource(ns, name)}
			if exists {
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func pvcSource(namespace, name string) string {
	return "k8s/pvc/" + namespace + "/" + name
}
