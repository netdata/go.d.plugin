package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type (
	serviceTargetGroup struct {
		targets []model.Target
		source  string
	}
	ServiceTarget struct {
		model.Base `hash:"ignore"`
		hash       uint64
		tuid       string
		Address    string

		Namespace   string
		Name        string
		Annotations map[string]any
		Labels      map[string]any

		Port         string
		PortName     string
		PortProtocol string
		ClusterIP    string
		ExternalName string
		Type         string
	}
)

func (s ServiceTarget) Hash() uint64 { return s.hash }
func (s ServiceTarget) TUID() string { return s.tuid }

func (s serviceTargetGroup) Source() string          { return s.source }
func (s serviceTargetGroup) Targets() []model.Target { return s.targets }

func NewServiceDiscovery(inf cache.SharedInformer) *ServiceDiscovery {
	queue := workqueue.NewNamed("service")
	_, _ = inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj any) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj any) { enqueue(queue, obj) },
		DeleteFunc: func(obj any) { enqueue(queue, obj) },
	})

	return &ServiceDiscovery{
		informer: inf,
		queue:    queue,
		Logger:   logger.New("k8s discovery", "service"),
	}
}

type ServiceDiscovery struct {
	*logger.Logger

	informer cache.SharedInformer
	queue    *workqueue.Type
}

func (sd *ServiceDiscovery) String() string {
	return fmt.Sprintf("k8s %sd discovery", RoleService)
}

func (sd *ServiceDiscovery) Discover(ctx context.Context, ch chan<- []model.TargetGroup) {
	sd.Info("instance is started")
	defer sd.Info("instance is stopped")

	defer sd.queue.ShutDown()

	go sd.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), sd.informer.HasSynced) {
		sd.Error("failed to sync caches")
		return
	}

	go sd.run(ctx, ch)
	<-ctx.Done()
}

func (sd *ServiceDiscovery) run(ctx context.Context, ch chan<- []model.TargetGroup) {
	for {
		item, shutdown := sd.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer sd.queue.Done(item)

			key := item.(string)
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := sd.informer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			if !exists {
				group := &serviceTargetGroup{source: serviceSourceFromNsName(namespace, name)}
				send(ctx, ch, group)
				return
			}

			svc, err := toService(item)
			if err != nil {
				return
			}

			group := sd.buildGroup(svc)
			send(ctx, ch, group)
		}()
	}
}

func (sd *ServiceDiscovery) buildGroup(svc *corev1.Service) model.TargetGroup {
	// TODO: headless service?
	if svc.Spec.ClusterIP == "" || len(svc.Spec.Ports) == 0 {
		return &serviceTargetGroup{
			source: serviceSource(svc),
		}
	}
	return &serviceTargetGroup{
		source:  serviceSource(svc),
		targets: sd.buildTargets(svc),
	}
}

func (sd *ServiceDiscovery) buildTargets(svc *corev1.Service) (targets []model.Target) {
	for _, port := range svc.Spec.Ports {
		portNum := strconv.FormatInt(int64(port.Port), 10)
		target := &ServiceTarget{
			tuid:         serviceTUID(svc, port),
			Address:      net.JoinHostPort(svc.Name+"."+svc.Namespace+".svc", portNum),
			Namespace:    svc.Namespace,
			Name:         svc.Name,
			Annotations:  toMapAny(svc.Annotations),
			Labels:       toMapAny(svc.Labels),
			Port:         portNum,
			PortName:     port.Name,
			PortProtocol: string(port.Protocol),
			ClusterIP:    svc.Spec.ClusterIP,
			ExternalName: svc.Spec.ExternalName,
			Type:         string(svc.Spec.Type),
		}
		hash, err := calcHash(target)
		if err != nil {
			continue
		}
		target.hash = hash

		targets = append(targets, target)
	}
	return targets
}

func serviceTUID(svc *corev1.Service, port corev1.ServicePort) string {
	return fmt.Sprintf("%s_%s_%s_%s",
		svc.Namespace,
		svc.Name,
		strings.ToLower(string(port.Protocol)),
		strconv.FormatInt(int64(port.Port), 10),
	)
}

func serviceSourceFromNsName(namespace, name string) string {
	return "k8s/service/" + namespace + "/" + name
}

func serviceSource(svc *corev1.Service) string {
	return serviceSourceFromNsName(svc.Namespace, svc.Name)
}

func toService(obj any) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", obj)
	}
	return svc, nil
}
