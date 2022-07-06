// SPDX-License-Identifier: GPL-3.0-or-later

package registry

import (
	"path/filepath"

	"github.com/gofrs/flock"
)

type FileLockRegistry struct {
	Dir   string
	locks map[string]*flock.Flock
}

func NewFileLockRegistry(dir string) *FileLockRegistry {
	return &FileLockRegistry{
		Dir:   dir,
		locks: make(map[string]*flock.Flock),
	}
}

const suffix = ".collector.lock"

func (r *FileLockRegistry) Register(name string) (bool, error) {
	name = filepath.Join(r.Dir, name+suffix)
	if _, ok := r.locks[name]; ok {
		return true, nil
	}

	locker := flock.New(name)
	ok, err := locker.TryLock()
	if ok {
		r.locks[name] = locker
	} else {
		_ = locker.Close()
	}
	return ok, err
}

func (r *FileLockRegistry) Unregister(name string) error {
	name = filepath.Join(r.Dir, name+suffix)
	locker, ok := r.locks[name]
	if !ok {
		return nil
	}
	delete(r.locks, name)
	return locker.Close()
}
