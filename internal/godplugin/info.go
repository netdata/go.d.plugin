package godplugin

import (
	"fmt"
	"sort"

	"github.com/l2isbad/go.d.plugin/internal/modules"
)

func info() {
	fmt.Println("Available modules:")
	var s []string
	for v := range modules.Registry {
		s = append(s, v)
	}
	sort.Strings(s)
	for idx, n := range s {
		fmt.Printf("  %d. %s\n", idx+1, n)
	}
}
