// SPDX-License-Identifier: GPL-3.0-or-later

package jobmgr

import (
	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/vnodes"
)

type (
	noopFileLocker  struct{}
	noopStatusSaver struct{}
	noopStatusStore struct{}
	noopVnodes      struct{}
	noopDyncfg      struct{}
)

func (n noopFileLocker) Lock(string) (bool, error) { return true, nil }
func (n noopFileLocker) Unlock(string) error       { return nil }

func (n noopStatusSaver) Save(confgroup.Config, string) {}
func (n noopStatusSaver) Remove(confgroup.Config)       {}

func (n noopStatusStore) Contains(confgroup.Config, ...string) bool { return false }

func (n noopVnodes) Lookup(string) (*vnodes.VirtualNode, bool) { return nil, false }

func (n noopDyncfg) Register(confgroup.Config)                     { return }
func (n noopDyncfg) Unregister(confgroup.Config)                   { return }
func (n noopDyncfg) UpdateStatus(confgroup.Config, string, string) { return }
