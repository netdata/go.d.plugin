package web

// HTTP is a struct with embedded RawRequest and RawClient.
type HTTP struct {
	RawRequest `yaml:",inline"`
	RawClient  `yaml:",inline"`
}
