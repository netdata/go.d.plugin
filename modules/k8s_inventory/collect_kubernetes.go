package k8s_inventory

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func (ki *KubernetesInventory) runCollectKubernetes(ctx context.Context, in <-chan resource) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-in:
			switch r.kind() {
			case kindPod:
				ki.collectPod(r)
			case kindNode:
				ki.collectNode(r)
			case kindPersistentVolume:
				ki.collectPV(r)
			case kindPersistentVolumeClaim:
				ki.collectPVC(r)
			}
		}
	}
}

func (ki *KubernetesInventory) collectNode(r resource) {
	node, err := convToNode(r.value())
	if err != nil {
		return
	}
	ki.Infof("GOT SRC: '%s', node: '%s'", r.source(), node.Name)
}

func (ki *KubernetesInventory) collectPod(r resource) {
	pod, err := convToPod(r.value())
	if err != nil {
		return
	}
	ki.Infof("GOT SRC: '%s', pod: '%s'", r.source(), pod.Name)
}

func (ki *KubernetesInventory) collectPV(r resource) {
	pv, err := convToPV(r.value())
	if err != nil {
		return
	}
	ki.Infof("GOT SRC: '%s', name: '%s'", r.source(), pv.Name)
}

func (ki *KubernetesInventory) collectPVC(r resource) {
	pvc, err := convToPVC(r.value())
	if err != nil {
		return
	}
	ki.Infof("GOT SRC: '%s', name: '%s'", r.source(), pvc.Name)
}

func convToNode(v interface{}) (*corev1.Node, error) {
	node, ok := v.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", v)
	}
	return node, nil
}

func convToPod(v interface{}) (*corev1.Pod, error) {
	pod, ok := v.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", v)
	}
	return pod, nil
}

func convToPV(v interface{}) (*corev1.PersistentVolume, error) {
	pv, ok := v.(*corev1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", v)
	}
	return pv, nil
}

func convToPVC(v interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := v.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("received unexpected object type: %T", v)
	}
	return pvc, nil
}
