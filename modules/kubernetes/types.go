package kubernetes

// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/apis/stats/v1alpha1/types.go

type statsSummary = Summary

// Summary is a top-level container for holding NodeStats and PodStats.
type Summary struct {
	Node NodeStats  `json:"node"`
	Pods []PodStats `json:"pods"`
}

// NodeStats holds node-level unprocessed sample stats.
type NodeStats struct {
	NodeName         string           `json:"nodeName"`
	SystemContainers []ContainerStats `json:"systemContainers"`
	CPU              *CPUStats        `json:"cpu"`
	Memory           *MemoryStats     `json:"memory"`
	Network          *NetworkStats    `json:"network"`
	Fs               *FsStats         `json:"fs"`
	Runtime          *RuntimeStats    `json:"runtime"`
	Rlimit           *RlimitStats     `json:"rlimit"`
}

// RlimitStats are stats rlimit of OS.
type RlimitStats struct {
	MaxPID                *int64 `json:"maxpid"`
	NumOfRunningProcesses *int64 `json:"curproc"`
}

type RuntimeStats struct {
	ImageFs *FsStats `json:"imageFs"`
}

// PodStats holds pod-level unprocessed sample stats.
type PodStats struct {
	PodRef           PodReference     `json:"podRef"`
	Containers       []ContainerStats `json:"containers"`
	CPU              *CPUStats        `json:"cpu"`
	Memory           *MemoryStats     `json:"memory"`
	Network          *NetworkStats    `json:"network"`
	VolumeStats      []VolumeStats    `json:"volume"`
	EphemeralStorage *FsStats         `json:"ephemeral-storage"`
}

// ContainerStats holds container-level unprocessed sample stats.
type ContainerStats struct {
	Name         string             `json:"name"`
	CPU          *CPUStats          `json:"cpu"`
	Memory       *MemoryStats       `json:"memory"`
	Accelerators []AcceleratorStats `json:"accelerator"`
	Rootfs       *FsStats           `json:"rootfs"`
	Logs         *FsStats           `json:"logs"`
}

// PodReference contains enough information to locate the referenced pod.
type PodReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

// InterfaceStats contains resource value data about interface.
type InterfaceStats struct {
	Name     string `json:"name"`
	RxBytes  *int64 `json:"rxBytes"`
	RxErrors *int64 `json:"rxErrors"`
	TxBytes  *int64 `json:"txBytes"`
	TxErrors *int64 `json:"txErrors"`
}

// NetworkStats contains data about network resources.
type NetworkStats struct {
	// Stats for the default interface, if found
	InterfaceStats `json:",inline"`
	Interfaces     []InterfaceStats `json:"interfaces"`
}

// CPUStats contains data about CPU usage.
type CPUStats struct {
	UsageNanoCores       *int64 `json:"usageNanoCores"`
	UsageCoreNanoSeconds *int64 `json:"usageCoreNanoSeconds"`
}

// MemoryStats contains data about memory usage.
type MemoryStats struct {
	AvailableBytes  *int64 `json:"availableBytes"` // if memory limit is undefined, the available bytes is omitted.
	UsageBytes      *int64 `json:"usageBytes"`
	WorkingSetBytes *int64 `json:"workingSetBytes"`
	RSSBytes        *int64 `json:"rssBytes"`
	PageFaults      *int64 `json:"pageFaults"`
	MajorPageFaults *int64 `json:"majorPageFaults"`
}

// AcceleratorStats contains stats for accelerators attached to the container.
type AcceleratorStats struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	ID          string `json:"id"`
	MemoryTotal int64  `json:"memoryTotal"`
	MemoryUsed  int64  `json:"memoryUsed"`
	DutyCycle   int64  `json:"dutyCycle"`
}

// VolumeStats contains data about Volume filesystem usage.
type VolumeStats struct {
	FsStats
	Name   string        `json:"name"`
	PVCRef *PVCReference `json:"pvcRef"`
}

// PVCReference contains enough information to describe the referenced PVC.
type PVCReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// FsStats contains data about filesystem usage.
type FsStats struct {
	AvailableBytes *int64 `json:"availableBytes"`
	CapacityBytes  *int64 `json:"capacityBytes"`
	UsedBytes      *int64 `json:"usedBytes"`
	InodesFree     *int64 `json:"inodesFree"`
	Inodes         *int64 `json:"inodes"`
	InodesUsed     *int64 `json:"inodesUsed"`
}
