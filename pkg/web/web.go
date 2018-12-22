package web

// HTTP is a struct with embedded Request and Client.
type HTTP struct {
	Request `yaml:",inline"`
	Client  `yaml:",inline"`
}
