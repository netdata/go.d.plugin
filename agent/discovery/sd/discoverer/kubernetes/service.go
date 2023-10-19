// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type serviceGroup struct {
	targets []model.Target
	source  string
}

func (sg serviceGroup) Provider() string        { return "sd:k8s:service" }
func (sg serviceGroup) Source() string          { return fmt.Sprintf("%s(%s)", sg.Provider(), sg.source) }
func (sg serviceGroup) Targets() []model.Target { return sg.targets }

type ServiceTarget struct {
	model.Base `hash:"ignore"`
	hash       uint64
	tuid       string
	Address    string

	Namespace   string
	Name        string
	Annotations map[string]interface{}
	Labels      map[string]interface{}

	Port         string
	PortName     string
	PortProtocol string
	ClusterIP    string
	ExternalName string
	Type         string
}

func (st ServiceTarget) Hash() uint64 { return st.hash }
func (st ServiceTarget) TUID() string { return st.tuid }

type Service struct {
	informer cache.SharedInformer
	queue    *workqueue.Type
	log      *logger.Logger
}

func NewService(inf cache.SharedInformer) *Service {
	if inf == nil {
		panic("nil service informer")
	}

	queue := workqueue.NewWithConfig(workqueue.QueueConfig{Name: "service"})
	_, _ = inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &Service{
		informer: inf,
		queue:    queue,
		log:      logger.New("k8s service discovery", ""),
	}
}

func (s *Service) String() string {
	return fmt.Sprintf("k8s %s discovery", RoleService)
}

func (s *Service) Discover(ctx context.Context, ch chan<- []model.TargetGroup) {
	s.log.Info("instance is started")
	defer s.log.Info("instance is stopped")
	defer s.queue.ShutDown()

	go s.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), s.informer.HasSynced) {
		s.log.Error("failed to sync caches")
		return
	}

	go s.run(ctx, ch)
	<-ctx.Done()
}

func (s *Service) run(ctx context.Context, ch chan<- []model.TargetGroup) {
	for {
		item, shutdown := s.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer s.queue.Done(item)

			key := item.(string)
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := s.informer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			if !exists {
				group := &serviceGroup{source: serviceSourceFromNsName(namespace, name)}
				send(ctx, ch, group)
				return
			}

			svc, err := toService(item)
			if err != nil {
				return
			}

			group := s.buildGroup(svc)
			send(ctx, ch, group)
		}()
	}
}

func (s *Service) buildGroup(svc *corev1.Service) model.TargetGroup {
	// TODO: headless service?
	if svc.Spec.ClusterIP == "" || len(svc.Spec.Ports) == 0 {
		return &serviceGroup{
			source: serviceSource(svc),
		}
	}
	return &serviceGroup{
		source:  serviceSource(svc),
		targets: s.buildTargets(svc),
	}
}

func (s *Service) buildTargets(svc *corev1.Service) (targets []model.Target) {
	for _, port := range svc.Spec.Ports {
		portNum := strconv.FormatInt(int64(port.Port), 10)
		target := &ServiceTarget{
			tuid:         serviceTUID(svc, port),
			Address:      net.JoinHostPort(svc.Name+"."+svc.Namespace+".svc", portNum),
			Namespace:    svc.Namespace,
			Name:         svc.Name,
			Annotations:  mapAny(svc.Annotations),
			Labels:       mapAny(svc.Labels),
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
	return namespace + "/" + name
}

func serviceSource(svc *corev1.Service) string {
	return serviceSourceFromNsName(svc.Namespace, svc.Name)
}

func toService(o interface{}) (*corev1.Service, error) {
	svc, ok := o.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", o)
	}
	return svc, nil
}
