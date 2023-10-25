// SPDX-License-Identifier: GPL-3.0-or-later

package hostnetsocket

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
	"github.com/netdata/go.d.plugin/logger"

	"github.com/ilyam8/hashstructure"
)

func NewTargetDiscoverer(path string) (*TargetDiscoverer, error) {
	d := &TargetDiscoverer{
		Logger:   logger.New("qq", "qq"),
		interval: time.Second * 60,
		ll: &localListenersExec{
			binPath: path,
			timeout: time.Second * 10,
		},
		started: make(chan struct{}),
	}

	return d, nil
}

type (
	TargetDiscoverer struct {
		*logger.Logger

		interval time.Duration
		ll       localListeners

		started chan struct{}
	}
	localListeners interface {
		discover(ctx context.Context) ([]byte, error)
	}
)

func (d *TargetDiscoverer) Discover(ctx context.Context, in chan<- []model.TargetGroup) {
	close(d.started)

	if err := d.discoverLocalListeners(ctx, in); err != nil {
		d.Error(err)
		return
	}

	tk := time.NewTicker(d.interval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			if err := d.discoverLocalListeners(ctx, in); err != nil {
				d.Error(err)
				return
			}
		}
	}
}

func (d *TargetDiscoverer) discoverLocalListeners(ctx context.Context, in chan<- []model.TargetGroup) error {
	bs, err := d.ll.discover(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	tggs, err := d.parseLocalListeners(bs)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	case in <- tggs:
	}
	return nil
}

func (d *TargetDiscoverer) parseLocalListeners(bs []byte) ([]model.TargetGroup, error) {
	var tgts []model.Target

	sc := bufio.NewScanner(bytes.NewReader(bs))
	for sc.Scan() {
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}

		// Protocol|Address|Port|Cmdline
		parts := strings.SplitN(text, "|", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("unexpected data: '%s'", text)
		}

		tgt := listenerTarget{
			Protocol: parts[0],
			Address:  parts[1],
			Port:     parts[2],
			Comm:     extractComm(parts[3]),
			Cmdline:  parts[3],
		}

		hash, err := calcHash(tgt)
		if err != nil {
			continue
		}

		tgt.hash = hash

		tgts = append(tgts, &tgt)
	}

	tgg := &listenerTargetGroup{
		provider: "hostsocket",
		source:   "local_listeners",
		targets:  tgts,
	}

	return []model.TargetGroup{tgg}, nil
}

func extractComm(s string) string {
	i := strings.IndexByte(s, ' ')
	if i <= 0 {
		return ""
	}
	_, comm := filepath.Split(s[:i])
	return comm
}

func calcHash(obj any) (uint64, error) {
	return hashstructure.Hash(obj, nil)
}
