package modules

import (
	"path"
	"runtime"
)

type Creator func() Module

func (c Creator) MakeModule() Module {
	return c()
}

type Creators map[string]Creator

func (c Creators) Remove(v string) {
	if _, ok := c[v]; ok {
		c[v] = nil
	}
}

var Registry = make(Creators)

func Add(c Creator) {
	name := getFileName(2)

	if name != "" {
		Registry[name] = c
	}
}

func getFileName(skip int) string {
	var name string
	_, n, _, _ := runtime.Caller(skip)

	if n != "" {
		_, name = path.Split(n[:len(n)-3])
	}
	return name
}
