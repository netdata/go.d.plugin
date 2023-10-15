// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"errors"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/discoverer/kubernetes"
)

type Config struct {
	K8S []kubernetes.Config `yaml:"k8s"`
}

func validateConfig(cfg Config) error {
	if len(cfg.K8S) == 0 {
		return errors.New("empty config")
	}
	return nil
}
