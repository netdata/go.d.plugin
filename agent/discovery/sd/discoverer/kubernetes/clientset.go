// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import (
	"errors"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const envFakeClient = "KUBERNETES_FAKE_CLIENTSET"

func newClientset() (kubernetes.Interface, error) {
	switch {
	case os.Getenv(envFakeClient) != "":
		return fake.NewSimpleClientset(), nil
	case os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "":
		return newClientsetInCluster()
	default:
		return newClientsetOutOfCluster()
	}
}

func newClientsetInCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.UserAgent = "Netdata/auto-discovery"

	return kubernetes.NewForConfig(config)
}

func newClientsetOutOfCluster() (*kubernetes.Clientset, error) {
	home := homeDir()
	if home == "" {
		return nil, errors.New("couldn't find home directory")
	}

	path := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}

	config.UserAgent = "Netdata/auto-discovery"

	return kubernetes.NewForConfig(config)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
