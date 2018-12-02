package modules

import "fmt"

type (
	// Creator is a Job builder
	Creator struct {
		UpdateEvery       int
		DisabledByDefault bool
		Create            func() Module
	}
	// Registry is a collection of Creators
	Registry map[string]Creator
)

// DefaultRegistry DefaultRegistry
var DefaultRegistry = Registry{}

// Register a module
func Register(name string, creator Creator) {
	register(DefaultRegistry, name, creator)
}

func register(registry Registry, name string, creator Creator) {
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("%s is already in registry", name))
	}
	registry[name] = creator
}
