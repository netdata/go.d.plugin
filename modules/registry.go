package modules

type (
	// Creator is a builder to create job instance
	Creator struct {
		UpdateEvery       *int
		DisabledByDefault bool
		Create            func() Module
	}
	Registry map[string]Creator
)

// DefaultRegistry DefaultRegistry
var DefaultRegistry = Registry{}

// Register a job
func Register(name string, creator Creator) {
	DefaultRegistry[name] = creator
}
