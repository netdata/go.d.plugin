package k8s_inventory

import (
	"k8s.io/client-go/kubernetes"
)

func (ki KubernetesInventory) initClient() (kubernetes.Interface, error) {
	return newKubeClient()
}

func (ki KubernetesInventory) initDiscoverer(client kubernetes.Interface) discoverer {
	return &kubeDiscovery{
		client: client,
	}
}
