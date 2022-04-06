package k8s_inventory

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newNodeDiscoverer(si cache.SharedInformer) *nodeDiscoverer {
	if si == nil {
		panic("nil node shared informer")
	}

	queue := workqueue.NewNamed("node")
	si.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &nodeDiscoverer{
		informer: si,
		queue:    queue,
	}
}

type nodeResource struct {
	src string
	val interface{}
}

func (r nodeResource) source() string     { return r.src }
func (r nodeResource) kind() resourceKind { return kindNode }
func (r nodeResource) value() interface{} { return r.val }

type nodeDiscoverer struct {
	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (d *nodeDiscoverer) run(ctx context.Context, in chan<- resource) {
	defer func() { fmt.Println("STOP NODE DISCOVERER") }()
	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runProcessQueue(ctx, in)

	<-ctx.Done()
}

func (d *nodeDiscoverer) runProcessQueue(ctx context.Context, in chan<- resource) {
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

			r := &nodeResource{src: nodeSource(name)}
			if exists {
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func nodeSource(name string) string {
	return "k8s/node/" + name
}
