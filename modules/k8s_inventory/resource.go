package k8s_inventory

type resource interface {
	source() string
	kind() resourceKind
	value() interface{}
}

type resourceKind string

const (
	kindNode                  resourceKind = "node"
	kindPod                   resourceKind = "pod"
	kindPersistentVolume      resourceKind = "pv"
	kindPersistentVolumeClaim resourceKind = "pvc"
)
