// SPDX-License-Identifier: GPL-3.0-or-later

package web

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// HTTP is a struct with embedded Request and Client.
// This structure intended to be part of the module configuration.
// Supported configuration file formats: YAML.
type HTTP struct {
	Request `yaml:",inline"`
	Client  `yaml:",inline"`
}

var hostname = ""

func init() {
	path, err := exec.LookPath("hostname")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	bs, err := exec.CommandContext(ctx, path).Output()
	if err != nil {
		return
	}

	hostname = string(bytes.TrimSpace(bs))
}
