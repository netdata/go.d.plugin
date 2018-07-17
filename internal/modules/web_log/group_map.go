package web_log

type groupMap map[string]string

func (gm groupMap) has(s string) bool {
	_, ok := gm[s]
	return ok
}

func (gm groupMap) get(s string) string {
	return gm[s]
}
func (gm groupMap) lookup(s string) (string, bool) {
	v, ok := gm[s]
	return v, ok
}

func (gm *groupMap) update(keys, values []string) {
	m := make(groupMap)
	for idx, v := range keys[1:] {
		m[v] = values[idx+1]
	}
	*gm = m
}
