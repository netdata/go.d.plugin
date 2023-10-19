// SPDX-License-Identifier: GPL-3.0-or-later

package kubernetes

import "fmt"

type Config struct {
	APIServer  string   `yaml:"api_server"`
	Tags       string   `yaml:"tags"`
	Namespaces []string `yaml:"namespaces"`
	Role       string   `yaml:"role"`
	LocalMode  bool     `yaml:"local_mode"`
	Selector   struct {
		Label string `yaml:"label"`
		Field string `yaml:"field"`
	} `yaml:"selector"`
}

func validateConfig(cfg Config) error {
	if !isRoleValid(cfg.Role) {
		return fmt.Errorf("invalid role '%s', valid roles: '%s', '%s'", cfg.Role, RolePod, RoleService)
	}

	if cfg.Tags == "" {
		return fmt.Errorf("no tags set for '%s' role", cfg.Role)
	}

	return nil
}

func isRoleValid(role string) bool {
	return role == RolePod || role == RoleService
}
