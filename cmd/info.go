package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"

	"github.com/l2isbad/go.d.plugin/modules"
)

var (
	yellowBold = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()
)

func Info() {
	fmt.Println(yellowBold("Available modules:\n"))
	var s []string
	for v := range modules.Registry {
		s = append(s, v)
	}
	sort.Strings(s)
	for idx, n := range s {
		fmt.Printf("%d %s\n", idx+1, n)
	}
	fmt.Printf("\nModule debug mode \t%s\n",
		green("./go.d.plugin -debug -module=MODULE_NAME"))
}
