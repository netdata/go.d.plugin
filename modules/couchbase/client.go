package couchbase

import (
	"fmt"

	couchbase "github.com/couchbase/go-couchbase"
)

type cbClient interface {
	connect(base, username, password string) (cbConnection, error)
}
type cbConnection interface {
	GetPool(name string) (couchbase.Pool, error)
}

type client struct{}

func (client) connect(base, username, password string) (cbConnection, error) {
	c, err := couchbase.ConnectWithAuthCreds(base, username, password)
	if err != nil {
		return nil, fmt.Errorf("error on creating a connection: %v", err)
	}
	return &c, nil
}

func newCouchbaseClient() *client {
	return &client{}
}
