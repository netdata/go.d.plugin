package client

import "sync"

func newToken() *token { return &token{mux: &sync.RWMutex{}} }

type token struct {
	mux   *sync.RWMutex
	value string
}

func (t *token) set(v string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.value = v
}

func (t token) get() string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

func (t *token) unset() { t.set("") }

func (t token) isSet() bool { return t.get() != "" }
