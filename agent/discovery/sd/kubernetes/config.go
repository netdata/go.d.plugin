// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import "fmt"

type Config struct {
	APIServer  string   `yaml:"api_server"` // TODO: implement?
	Namespaces []string `yaml:"namespaces"`
	Role       string   `yaml:"role"`
	LocalMode  bool     `yaml:"local_mode"`
	Selector   struct {
		Label string `yaml:"label"`
		Field string `yaml:"field"`
	} `yaml:"selector"`
}

func validateConfig(cfg Config) error {
	if !(cfg.Role == RolePod || cfg.Role == RoleService) {
		return fmt.Errorf("invalid role '%s', valid roles: '%s', '%s'", cfg.Role, RolePod, RoleService)
	}

	return nil
}
