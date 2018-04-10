package goplugin

import (
	"bytes"
)

type myBool struct {
	Bool bool
}

func (b *myBool) UnmarshalTOML(input []byte) error {
	s := string(bytes.Trim(input, "'\""))
	switch s {
	case "yes", "true", "on", "+":
		b.Bool = true
	default:
		b.Bool = false
	}
	return nil
}

func NewConf() *Conf {
	return &Conf{
		Modules:    make(map[string]*myBool),
		DefaultRun: myBool{true},
		Enabled:    myBool{true},
	}
}

type Conf struct {
	Enabled    myBool             `toml:"enabled"`
	DefaultRun myBool             `toml:"default_run"`
	MaxProcs   int                `toml:"max_procs,    range:[1:]"`
	Modules    map[string]*myBool `toml:"modules"`
}
