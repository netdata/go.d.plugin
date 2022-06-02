package k8s_state

import (
	"context"

	"github.com/netdata/go.d.plugin/logger"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newPodDiscoverer(si cache.SharedInformer, l *logger.Logger) *podDiscoverer {
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
		Logger:   l,
		informer: si,
		queue:    queue,
	}
}

type podResource struct {
	src string
	val interface{}
}

func (r podResource) source() string         { return r.src }
func (r podResource) kind() kubeResourceKind { return kubeResourcePod }
func (r podResource) value() interface{}     { return r.val }

type podDiscoverer struct {
	*logger.Logger
	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (d *podDiscoverer) run(ctx context.Context, in chan<- resource) {
	d.Info("pod_discoverer is started")
	defer func() { d.Info("pod_discoverer is stopped") }()

	defer d.queue.ShutDown()

	go d.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), d.informer.HasSynced) {
		return
	}

	go d.runDiscover(ctx, in)

	<-ctx.Done()
}

func (d *podDiscoverer) hasSynced() bool { return d.informer.HasSynced() }

func (d *podDiscoverer) runDiscover(ctx context.Context, in chan<- resource) {
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

			r := &podResource{src: podSource(ns, replaceDots(name))}
			if exists {
				if p, err := toPod(item); err == nil {
					p.Name = replaceDots(p.Name)
					p.Spec.NodeName = replaceDots(p.Spec.NodeName)
				}
				r.val = item
			}
			send(ctx, in, r)
		}()
	}
}

func podSource(namespace, name string) string {
	return "k8s/pod/" + namespace + "/" + name
}
