package web

// Web is a struct with embedded RawRequest and RawClient.
type Web struct {
	RawRequest `yaml:",inline"`
	RawClient  `yaml:",inline"`
}
