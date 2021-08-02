package mongo

import (
	"time"
)

const (
	defaultTimeout = 20
	defaultUri     = "mongodb://localhost:27017"
)

type Config struct {
	Uri     string        `yaml:"uri"`
	Timeout time.Duration `yaml:"timeout"`
}
