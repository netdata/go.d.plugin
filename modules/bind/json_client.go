package bind

type jsonStats struct {
	OpCodes   map[string]int
	QTypes    map[string]int
	NSStats   map[string]int
	SockStats map[string]int
	Views     map[string]jsonView
}

type jsonView struct {
	Resolver map[string]map[string]int
}

type jsonClient struct{}
