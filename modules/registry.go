package modules

type (
	// Creator is a builder to create Job instance
	Creator struct {
		UpdateEvery       int
		DisabledByDefault bool
		Create            func() Module
	}
	Registry map[string]Creator
)

// DefaultRegistry DefaultRegistry
var DefaultRegistry = Registry{}

// Register a module
func Register(name string, creator Creator) {
	DefaultRegistry[name] = creator
}
