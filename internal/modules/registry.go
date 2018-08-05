package modules

type (
	// Creator is a builder to create module instance
	Creator struct {
		UpdateEvery       *int
		ChartCleanup      *int
		DisabledByDefault bool
		NoConfig          bool
		Create            func() Module
	}
)

// Registry Registry
var Registry = map[string]Creator{}

// Register a module
func Register(name string, creator Creator) {
	Registry[name] = creator
}
