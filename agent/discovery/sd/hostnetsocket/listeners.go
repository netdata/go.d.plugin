// SPDX-License-Identifier: GPL-3.0-or-later

package hostnetsocket

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
)

type listenerTargetGroup struct {
	provider string
	source   string
	targets  []model.Target
}

func (g *listenerTargetGroup) Provider() string        { return g.provider }
func (g *listenerTargetGroup) Source() string          { return g.source }
func (g *listenerTargetGroup) Targets() []model.Target { return g.targets }

type listenerTarget struct {
	model.Base

	hash uint64

	Protocol string
	Address  string
	Port     string
	Comm     string
	Cmdline  string
}

func (t *listenerTarget) TUID() string { return t.tuid() }
func (t *listenerTarget) Hash() uint64 { return t.hash }
func (t *listenerTarget) tuid() string {
	return fmt.Sprintf("%s_%s_%d", strings.ToLower(t.Protocol), t.Port, t.hash)
}

type localListenersExec struct {
	binPath string
	timeout time.Duration
}

func (e *localListenersExec) discover(ctx context.Context) ([]byte, error) {
	execCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(execCtx, e.binPath, "tcp") // TODO: tcp6?

	bs, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error on '%s': %v", cmd, err)
	}

	return bs, nil
}
