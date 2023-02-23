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
	podTargetGroup struct {
		targets []model.Target
		source  string
	}
	PodTarget struct {
		model.Base `hash:"ignore"`
		hash       uint64
		tuid       string
		Address    string

		Namespace   string
		Name        string
		Annotations map[string]any
		Labels      map[string]any
		NodeName    string
		PodIP       string

		ContName     string
		Image        string
		Env          map[string]any
		Port         string
		PortName     string
		PortProtocol string
	}
)

func (p *PodTarget) Hash() uint64 { return p.hash }
func (p *PodTarget) TUID() string { return p.tuid }

func (p *podTargetGroup) Source() string          { return p.source }
func (p *podTargetGroup) Targets() []model.Target { return p.targets }

func NewPodDiscovery(pod, cmap, secret cache.SharedInformer) *PodDiscovery {
	if cmap == nil || secret == nil {
		panic("nil cmap or secret informer")
	}

	queue := workqueue.NewNamed("pod")

	_, _ = pod.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj any) { enqueue(queue, obj) },
		UpdateFunc: func(_, obj any) { enqueue(queue, obj) },
		DeleteFunc: func(obj any) { enqueue(queue, obj) },
	})

	return &PodDiscovery{
		podInformer:    pod,
		cmapInformer:   cmap,
		secretInformer: secret,
		queue:          queue,
		Logger:         logger.New("k8s discovery", "pod"),
	}
}

type PodDiscovery struct {
	*logger.Logger

	podInformer    cache.SharedInformer
	cmapInformer   cache.SharedInformer
	secretInformer cache.SharedInformer
	queue          *workqueue.Type
}

func (pd *PodDiscovery) String() string {
	return fmt.Sprintf("k8s %s discovery", RolePod)
}

func (pd *PodDiscovery) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	pd.Info("instance is started")
	defer pd.Info("instance is stopped")
	defer pd.queue.ShutDown()

	go pd.podInformer.Run(ctx.Done())
	go pd.cmapInformer.Run(ctx.Done())
	go pd.secretInformer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(),
		pd.podInformer.HasSynced, pd.cmapInformer.HasSynced, pd.secretInformer.HasSynced) {
		pd.Error("failed to sync caches")
		return
	}

	go pd.run(ctx, in)

	<-ctx.Done()
}

func (pd *PodDiscovery) run(ctx context.Context, in chan<- []model.TargetGroup) {
	for {
		item, shutdown := pd.queue.Get()
		if shutdown {
			return
		}

		func() {
			defer pd.queue.Done(item)

			key := item.(string)
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				return
			}

			item, exists, err := pd.podInformer.GetStore().GetByKey(key)
			if err != nil {
				return
			}

			if !exists {
				group := &podTargetGroup{source: podSourceFromNsName(namespace, name)}
				send(ctx, in, group)
				return
			}

			pod, err := toPod(item)
			if err != nil {
				return
			}

			group := pd.buildGroup(pod)
			send(ctx, in, group)
		}()
	}
}

func (pd *PodDiscovery) buildGroup(pod *corev1.Pod) model.TargetGroup {
	if pod.Status.PodIP == "" || len(pod.Spec.Containers) == 0 {
		return &podTargetGroup{
			source: podSource(pod),
		}
	}
	return &podTargetGroup{
		source:  podSource(pod),
		targets: pd.buildTargets(pod),
	}
}

func (pd *PodDiscovery) buildTargets(pod *corev1.Pod) (targets []model.Target) {
	for _, container := range pod.Spec.Containers {
		env := pd.collectEnv(pod.Namespace, container)

		if len(container.Ports) == 0 {
			target := &PodTarget{
				tuid:        podTUID(pod, container),
				Address:     pod.Status.PodIP,
				Namespace:   pod.Namespace,
				Name:        pod.Name,
				Annotations: toMapAny(pod.Annotations),
				Labels:      toMapAny(pod.Labels),
				NodeName:    pod.Spec.NodeName,
				PodIP:       pod.Status.PodIP,
				ContName:    container.Name,
				Image:       container.Image,
				Env:         toMapAny(env),
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
					tuid:         podTUIDWithPort(pod, container, port),
					Address:      net.JoinHostPort(pod.Status.PodIP, portNum),
					Namespace:    pod.Namespace,
					Name:         pod.Name,
					Annotations:  toMapAny(pod.Annotations),
					Labels:       toMapAny(pod.Labels),
					NodeName:     pod.Spec.NodeName,
					PodIP:        pod.Status.PodIP,
					ContName:     container.Name,
					Image:        container.Image,
					Env:          toMapAny(env),
					Port:         portNum,
					PortName:     port.Name,
					PortProtocol: string(port.Protocol),
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

func (pd *PodDiscovery) collectEnv(ns string, container corev1.Container) map[string]string {
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
			pd.envFromConfigMap(vars, ns, src)
		case src.SecretRef != nil:
			pd.envFromSecret(vars, ns, src)
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
			pd.valueFromSecret(vars, ns, env)
		case env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil:
			pd.valueFromConfigMap(vars, ns, env)
		}
	}

	if len(vars) == 0 {
		return nil
	}
	return vars
}

func (pd *PodDiscovery) valueFromConfigMap(vars map[string]string, ns string, env corev1.EnvVar) {
	if env.ValueFrom.ConfigMapKeyRef.Name == "" || env.ValueFrom.ConfigMapKeyRef.Key == "" {
		return
	}

	sr := env.ValueFrom.ConfigMapKeyRef
	key := ns + "/" + sr.Name
	item, exist, err := pd.cmapInformer.GetStore().GetByKey(key)
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

func (pd *PodDiscovery) valueFromSecret(vars map[string]string, ns string, env corev1.EnvVar) {
	if env.ValueFrom.SecretKeyRef.Name == "" || env.ValueFrom.SecretKeyRef.Key == "" {
		return
	}

	keyRef := env.ValueFrom.SecretKeyRef
	key := ns + "/" + keyRef.Name
	item, exist, err := pd.secretInformer.GetStore().GetByKey(key)
	if err != nil || !exist {
		return
	}

	secret, err := toSecret(item)
	if err != nil {
		return
	}

	if v, ok := secret.Data[keyRef.Key]; ok {
		vars[env.Name] = string(v)
	}
}

func (pd *PodDiscovery) envFromConfigMap(vars map[string]string, ns string, src corev1.EnvFromSource) {
	if src.ConfigMapRef.Name == "" {
		return
	}

	key := ns + "/" + src.ConfigMapRef.Name
	item, exist, err := pd.cmapInformer.GetStore().GetByKey(key)
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

func (pd *PodDiscovery) envFromSecret(vars map[string]string, ns string, src corev1.EnvFromSource) {
	if src.SecretRef.Name == "" {
		return
	}

	key := ns + "/" + src.SecretRef.Name
	item, exist, err := pd.secretInformer.GetStore().GetByKey(key)
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

func podTUID(pod *corev1.Pod, container corev1.Container) string {
	return fmt.Sprintf("%s_%s_%s",
		pod.Namespace,
		pod.Name,
		container.Name,
	)
}

func podTUIDWithPort(pod *corev1.Pod, container corev1.Container, port corev1.ContainerPort) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s",
		pod.Namespace,
		pod.Name,
		container.Name,
		strings.ToLower(string(port.Protocol)),
		strconv.FormatInt(int64(port.ContainerPort), 10),
	)
}

func podSourceFromNsName(namespace, name string) string {
	return "k8s/pod/" + namespace + "/" + name
}

func podSource(pod *corev1.Pod) string {
	return podSourceFromNsName(pod.Namespace, pod.Name)
}

func toPod(item any) (*corev1.Pod, error) {
	pod, ok := item.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", item)
	}
	return pod, nil
}

func toConfigMap(item any) (*corev1.ConfigMap, error) {
	cmap, ok := item.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", item)
	}
	return cmap, nil
}

func toSecret(item any) (*corev1.Secret, error) {
	secret, ok := item.(*corev1.Secret)
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

func toMapAny(src map[string]string) map[string]any {
	if src == nil {
		return nil
	}

	m := make(map[string]any, len(src))
	for k, v := range src {
		m[k] = v
	}

	return m
}
