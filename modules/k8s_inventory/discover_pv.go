package k8s_inventory

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newPVDiscoverer(si cache.SharedInformer) *pvDiscoverer {
	if si == nil {
		panic("nil persistent volume shared informer")
	}

	queue := workqueue.NewNamed("pv")
	si.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &pvDiscoverer{
		informer: si,
		queue:    queue,
	}
}

type pvResource struct {
	src string
	val interface{}
}

func (r pvResource) source() string     { return r.src }
func (r pvResource) kind() resourceKind { return kindPersistentVolume }
func (r pvResource) value() interface{} { return r.val }

type pvDiscoverer struct {
	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (d *pvDiscoverer) run(ctx context.Context, in chan<- resource) {
	defer func() { fmt.Println("STOP PV DISCOVERER") }()
	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runProcessQueue(ctx, in)

	<-ctx.Done()
}

func (d *pvDiscoverer) runProcessQueue(ctx context.Context, in chan<- resource) {
	for {
		item, shutdown := d.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer d.queue.Done(item)

			key := item.(string)
			_, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := d.informer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			r := &pvResource{src: pvSource(name)}
			if exists {
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func pvSource(name string) string {
	return "k8s/pv/" + name
}
