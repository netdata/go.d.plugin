package unbound

import (
	"sort"
	"strings"
)

func (u *Unbound) collect() (map[string]int64, error) {
	return nil, nil
}

func findByPrefix(ms []string, prefix string) []string {
	from := sort.Search(len(ms), func(i int) bool { return ms[i] >= prefix })
	if from == len(ms) || !strings.HasPrefix(ms[from], prefix) {
		return nil
	}
	until := from + 1
	for until < len(ms) && strings.HasPrefix(ms[until], prefix) {
		until++
	}
	return ms[from:until]
}
