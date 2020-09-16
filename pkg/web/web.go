package web

// HTTP is a struct with embedded Request and Client.
// This structure intended to be part of the module configuration.
// Supported configuration file formats: YAML.
type HTTP struct {
	Request `yaml:",inline"`
	Client  `yaml:",inline"`
}
