package parser

type GroupMap map[string]string

func (gm GroupMap) Has(key string) bool {
	_, ok := gm[key]
	return ok
}

func (gm GroupMap) Get(key string) string {
	return gm[key]
}

func (gm GroupMap) Lookup(key string) (string, bool) {
	v, ok := gm[key]
	return v, ok
}
