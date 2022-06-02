package k8s_state

import (
	"context"

	"github.com/netdata/go.d.plugin/logger"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newNodeDiscoverer(si cache.SharedInformer, l *logger.Logger) *nodeDiscoverer {
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
		Logger:   l,
		informer: si,
		queue:    queue,
		started:  make(chan struct{}),
	}
}

type nodeResource struct {
	src string
	val interface{}
}

func (r nodeResource) source() string         { return r.src }
func (r nodeResource) kind() kubeResourceKind { return kubeResourceNode }
func (r nodeResource) value() interface{}     { return r.val }

type nodeDiscoverer struct {
	*logger.Logger
	informer cache.SharedInformer
	queue    *workqueue.Type
	started  chan struct{}
}

func (d *nodeDiscoverer) run(ctx context.Context, in chan<- resource) {
	d.Info("node_discoverer is started")
	defer func() { d.Info("node_discoverer is stopped") }()

	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runDiscover(ctx, in)
	close(d.started)

	<-ctx.Done()
}

func (d *nodeDiscoverer) ready() bool {
	select {
	case <-d.started:
		return true
	default:
		return false
	}
}

func (d *nodeDiscoverer) runDiscover(ctx context.Context, in chan<- resource) {
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

			r := &nodeResource{src: nodeSource(replaceDots(name))}
			if exists {
				if n, err := toNode(item); err == nil {
					n.Name = replaceDots(name)
				}
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func nodeSource(name string) string {
	return "k8s/node/" + name
}
