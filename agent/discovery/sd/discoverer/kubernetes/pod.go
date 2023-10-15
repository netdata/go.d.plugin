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

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type podGroup struct {
	targets []model.Target
	source  string
}

func (pg podGroup) Source() string          { return pg.source }
func (pg podGroup) Targets() []model.Target { return pg.targets }

type PodTarget struct {
	model.Base `hash:"ignore"`
	hash       uint64
	tuid       string
	Address    string

	Namespace   string
	Name        string
	Annotations map[string]interface{}
	Labels      map[string]interface{}
	NodeName    string
	PodIP       string

	ControllerName string
	ControllerKind string

	ContName     string
	Image        string
	Env          map[string]interface{}
	Port         string
	PortName     string
	PortProtocol string
}

func (pt PodTarget) Hash() uint64 { return pt.hash }
func (pt PodTarget) TUID() string { return pt.tuid }

type Pod struct {
	podInformer    cache.SharedInformer
	cmapInformer   cache.SharedInformer
	secretInformer cache.SharedInformer
	queue          *workqueue.Type
	log            *logger.Logger
}

func NewPod(pod, cmap, secret cache.SharedInformer) *Pod {
	if pod == nil || cmap == nil || secret == nil {
		panic("nil cmap or secret informer")
	}

	queue := workqueue.NewWithConfig(workqueue.QueueConfig{Name: "pod"})

	_, _ = pod.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj interface{}) { enqueue(queue, obj) },
		DeleteFunc: func(obj interface{}) { enqueue(queue, obj) },
	})

	return &Pod{
		podInformer:    pod,
		cmapInformer:   cmap,
		secretInformer: secret,
		queue:          queue,
		log:            logger.New("k8s pod discovery", ""),
	}
}

func (p *Pod) String() string {
	return fmt.Sprintf("k8s %s discovery", RolePod)
}

func (p *Pod) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	p.log.Info("instance is started")
	defer p.log.Info("instance is stopped")
	defer p.queue.ShutDown()

	go p.podInformer.Run(ctx.Done())
	go p.cmapInformer.Run(ctx.Done())
	go p.secretInformer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(),
		p.podInformer.HasSynced, p.cmapInformer.HasSynced, p.secretInformer.HasSynced) {
		p.log.Error("failed to sync caches")
		return
	}

	go p.run(ctx, in)
	<-ctx.Done()
}

func (p *Pod) run(ctx context.Context, in chan<- []model.TargetGroup) {
	for {
		item, shutdown := p.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer p.queue.Done(item)

			key := item.(string)
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := p.podInformer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			if !exists {
				group := &podGroup{source: podSourceFromNsName(namespace, name)}
				send(ctx, in, group)
				return
			}

			pod, err := toPod(item)
			if err != nil {
				return
			}

			group := p.buildGroup(pod)
			send(ctx, in, group)
		}()
	}
}

func (p *Pod) buildGroup(pod *apiv1.Pod) model.TargetGroup {
	if pod.Status.PodIP == "" || len(pod.Spec.Containers) == 0 {
		return &podGroup{
			source: podSource(pod),
		}
	}
	return &podGroup{
		source:  podSource(pod),
		targets: p.buildTargets(pod),
	}
}

func (p *Pod) buildTargets(pod *apiv1.Pod) (targets []model.Target) {
	var name, kind string
	for _, ref := range pod.OwnerReferences {
		if ref.Controller != nil && *ref.Controller {
			name = ref.Name
			kind = ref.Kind
			break
		}
	}

	for _, container := range pod.Spec.Containers {
		env := p.collectEnv(pod.Namespace, container)

		if len(container.Ports) == 0 {
			target := &PodTarget{
				tuid:           podTUID(pod, container),
				Address:        pod.Status.PodIP,
				Namespace:      pod.Namespace,
				Name:           pod.Name,
				Annotations:    toMapInterface(pod.Annotations),
				Labels:         toMapInterface(pod.Labels),
				NodeName:       pod.Spec.NodeName,
				PodIP:          pod.Status.PodIP,
				ControllerName: name,
				ControllerKind: kind,
				ContName:       container.Name,
				Image:          container.Image,
				Env:            toMapInterface(env),
			}
			hash, err := calcHash(target)
			if err != nil {
				continue
			}
			target.hash = hash

			targets = append(targets, target)
		} else {
			for _, port := range container.Ports {
				portNum := strconv.FormatUint(uint64(port.ContainerPort), 10)
				target := &PodTarget{
					tuid:           podTUIDWithPort(pod, container, port),
					Address:        net.JoinHostPort(pod.Status.PodIP, portNum),
					Namespace:      pod.Namespace,
					Name:           pod.Name,
					Annotations:    toMapInterface(pod.Annotations),
					Labels:         toMapInterface(pod.Labels),
					NodeName:       pod.Spec.NodeName,
					PodIP:          pod.Status.PodIP,
					ControllerName: name,
					ControllerKind: kind,
					ContName:       container.Name,
					Image:          container.Image,
					Env:            toMapInterface(env),
					Port:           portNum,
					PortName:       port.Name,
					PortProtocol:   string(port.Protocol),
				}
				hash, err := calcHash(target)
				if err != nil {
					continue
				}
				target.hash = hash

				targets = append(targets, target)
			}
		}
	}
	return targets
}

func (p *Pod) collectEnv(ns string, container apiv1.Container) map[string]string {
	vars := make(map[string]string)

	// When a key exists in multiple sources,
	// the value associated with the last source will take precedence.
	// Values defined by an Env with a duplicate key will take precedence.
	//
	// Order (https://github.com/kubernetes/kubectl/blob/master/pkg/describe/describe.go)
	// - envFrom: configMapRef, secretRef
	// - env: value || valueFrom: fieldRef, resourceFieldRef, secretRef, configMap

	for _, src := range container.EnvFrom {
		switch {
		case src.ConfigMapRef != nil:
			p.envFromConfigMap(vars, ns, src)
		case src.SecretRef != nil:
			p.envFromSecret(vars, ns, src)
		}
	}

	for _, env := range container.Env {
		if env.Name == "" || isVar(env.Name) {
			continue
		}
		switch {
		case env.Value != "":
			vars[env.Name] = env.Value
		case env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil:
			p.valueFromSecret(vars, ns, env)
		case env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil:
			p.valueFromConfigMap(vars, ns, env)
		}
	}
	if len(vars) == 0 {
		return nil
	}
	return vars
}

func (p *Pod) valueFromConfigMap(vars map[string]string, ns string, env apiv1.EnvVar) {
	if env.ValueFrom.ConfigMapKeyRef.Name == "" || env.ValueFrom.ConfigMapKeyRef.Key == "" {
		return
	}

	sr := env.ValueFrom.ConfigMapKeyRef
	key := ns + "/" + sr.Name
	item, exist, err := p.cmapInformer.GetStore().GetByKey(key)
	if err != nil || !exist {
		return
	}
	cmap, err := toConfigMap(item)
	if err != nil {
		return
	}
	if v, ok := cmap.Data[sr.Key]; ok {
		vars[env.Name] = v
	}
}

func (p *Pod) valueFromSecret(vars map[string]string, ns string, env apiv1.EnvVar) {
	if env.ValueFrom.SecretKeyRef.Name == "" || env.ValueFrom.SecretKeyRef.Key == "" {
		return
	}

	secretKey := env.ValueFrom.SecretKeyRef
	key := ns + "/" + secretKey.Name

	item, exist, err := p.secretInformer.GetStore().GetByKey(key)
	if err != nil || !exist {
		return
	}

	secret, err := toSecret(item)
	if err != nil {
		return
	}

	if v, ok := secret.Data[secretKey.Key]; ok {
		vars[env.Name] = string(v)
	}
}

func (p *Pod) envFromConfigMap(vars map[string]string, ns string, src apiv1.EnvFromSource) {
	if src.ConfigMapRef.Name == "" {
		return
	}

	key := ns + "/" + src.ConfigMapRef.Name
	item, exist, err := p.cmapInformer.GetStore().GetByKey(key)
	if err != nil || !exist {
		return
	}

	cmap, err := toConfigMap(item)
	if err != nil {
		return
	}

	for k, v := range cmap.Data {
		vars[src.Prefix+k] = v
	}
}

func (p *Pod) envFromSecret(vars map[string]string, ns string, src apiv1.EnvFromSource) {
	if src.SecretRef.Name == "" {
		return
	}

	key := ns + "/" + src.SecretRef.Name
	item, exist, err := p.secretInformer.GetStore().GetByKey(key)
	if err != nil || !exist {
		return
	}

	secret, err := toSecret(item)
	if err != nil {
		return
	}

	for k, v := range secret.Data {
		vars[src.Prefix+k] = string(v)
	}
}

func podTUID(pod *apiv1.Pod, container apiv1.Container) string {
	return fmt.Sprintf("%s_%s_%s",
		pod.Namespace,
		pod.Name,
		container.Name,
	)
}

func podTUIDWithPort(pod *apiv1.Pod, container apiv1.Container, port apiv1.ContainerPort) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s",
		pod.Namespace,
		pod.Name,
		container.Name,
		strings.ToLower(string(port.Protocol)),
		strconv.FormatUint(uint64(port.ContainerPort), 10),
	)
}

func podSourceFromNsName(namespace, name string) string {
	return "k8s/pod/" + namespace + "/" + name
}

func podSource(pod *apiv1.Pod) string {
	return podSourceFromNsName(pod.Namespace, pod.Name)
}

func toPod(item interface{}) (*apiv1.Pod, error) {
	pod, ok := item.(*apiv1.Pod)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", item)
	}
	return pod, nil
}

func toConfigMap(item interface{}) (*apiv1.ConfigMap, error) {
	cmap, ok := item.(*apiv1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", item)
	}
	return cmap, nil
}

func toSecret(item interface{}) (*apiv1.Secret, error) {
	secret, ok := item.(*apiv1.Secret)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", item)
	}
	return secret, nil
}

func isVar(name string) bool {
	// Variable references $(VAR_NAME) are expanded using the previous defined
	// environment variables in the container and any service environment
	// variables.
	return strings.IndexByte(name, '$') != -1
}

func toMapInterface(src map[string]string) map[string]interface{} {
	if src == nil {
		return nil
	}
	m := make(map[string]interface{}, len(src))
	for k, v := range src {
		m[k] = v
	}
	return m
}
